package websocket

import (
	"context"
	"errors"
	"fmt"
	"mcpbinance/internal/websocket/enum"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	ErrClosedConn      = errors.New("connection closed")
	ErrReconnectNeeded = errors.New("reconnect needed")
)

type WebsocketConfig struct {
	name      string
	url       string
	msgCap    int
	readLimit int64

	reconnPeriod time.Duration
	retryBackoff time.Duration
	pingWait     time.Duration
	pongPeriod   time.Duration
}

type WebsocketClient struct {
	wc *WebsocketConfig

	conn   *websocket.Conn
	mu     sync.RWMutex
	cancel context.CancelFunc
	msgCh  chan []byte

	subs   map[string]struct{}
	subsMu sync.Mutex

	wg   sync.WaitGroup
	once sync.Once

	hooks clientHooks
}

type clientHooks interface {
	buildSubs(symbols []enum.Symbol, streamTypes []enum.StreamType) []string
	makeSubMsg(streams []string) map[string]any
}

func (c *WebsocketClient) Run(ctx context.Context, symbols []enum.Symbol, streamTypes []enum.StreamType) error {
	// log.Printf("[%sClient] starting client with URL: %s", c.wc.name, c.wc.url)
	defer func() {
		close(c.msgCh)
		// log.Printf("[%sClient] message channel closed", c.wc.name)
	}()

	c.initSubscription(symbols, streamTypes)
	for {
		// log.Printf("[%sClient] starting session", c.wc.name)
		err := c.runSession(ctx)
		if err != nil && !errors.Is(err, ErrReconnectNeeded) {
			// log.Printf("[%sClient] session ended with error: %v", c.wc.name, err)
		}

		select {
		case <-ctx.Done():
			// log.Printf("[%sClient] context canceled, stopping client", c.wc.name)
			c.stop()
			return nil
		case <-time.After(c.wc.retryBackoff):
			// log.Printf("[%sClient] retrying connection after %s", c.wc.name, c.wc.retryBackoff)
		}
	}
}

func (c *WebsocketClient) Receive() <-chan []byte {
	return c.msgCh
}

func (c *WebsocketClient) Subscribe(symbols, streamTypes []string) error {
	return nil
}

func (c *WebsocketClient) Unsubscribe(symbols, streamTypes []string) error {
	return nil
}

func (c *WebsocketClient) connect(ctx context.Context) error {
	// log.Printf("[%sClient] connecting to %s", c.wc.name, c.wc.url)

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, c.wc.url, nil)
	if err != nil {
		return fmt.Errorf(
			"failed to establish connection with %s: %w", c.wc.url, err,
		)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.conn = conn

	return nil
}

func (c *WebsocketClient) initSubscription(symbols []enum.Symbol, streamTypes []enum.StreamType) {
	streams := c.hooks.buildSubs(symbols, streamTypes)

	c.subsMu.Lock()
	defer c.subsMu.Unlock()

	for _, s := range streams {
		c.subs[s] = struct{}{}
	}

	// log.Printf("[%sClient] initialized subscriptions: %v", c.wc.name, streams)
}

func (c *WebsocketClient) resubscribe() error {
	if len(c.subs) == 0 {
		// log.Printf("[%sClient] no subscription to resubscribe", c.wc.name)
		return nil
	}

	c.subsMu.Lock()
	streams := make([]string, 0, len(c.subs))
	for s := range c.subs {
		streams = append(streams, s)
	}
	c.subsMu.Unlock()

	msg := c.hooks.makeSubMsg(streams)

	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.conn == nil {
		return ErrClosedConn
	}

	if err := c.conn.WriteJSON(msg); err != nil {
		return fmt.Errorf("failed to resubscribe: %w", err)
	}

	// log.Printf("[%sClient] successfully resubscribe", c.wc.name)
	return nil
}

func (c *WebsocketClient) readMessages(ctx context.Context) error {
	// log.Printf("[%sClient] starting message reader", c.wc.name)

	c.conn.SetReadLimit(c.wc.readLimit)
	c.conn.SetReadDeadline(time.Now().Add(c.wc.pongPeriod))
	c.conn.SetPongHandler(func(appData string) error {
		c.conn.SetReadDeadline(time.Now().Add(c.wc.pongPeriod))
		// log.Printf("[BinanceClient] received pong")
		return nil
	})

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseNormalClosure,
			) {
				// log.Printf("[%sClient] connection closed by server: %v", c.wc.name, err)
			} else {
				// log.Printf("[%sClient] read error: %v", c.wc.name, err)
			}
			return err
		}

		select {
		case <-ctx.Done():
			// log.Printf("[%sClient] context canceled, stopping client", c.wc.name)
			return nil
		case c.msgCh <- msg:
		default:
			// log.Printf("[%sClient] message channel is full, skipping message", c.wc.name)
		}
	}
}

func (c *WebsocketClient) pingLoop(ctx context.Context) error {
	// log.Printf("[%sClient] starting ping loop with interval %s", c.wc.name, c.wc.pingWait)

	ticker := time.NewTicker(c.wc.pingWait)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// log.Printf("[%sClient] stoping ping loop", c.wc.name)
			return nil
		case <-ticker.C:
			c.mu.RLock()
			if c.conn == nil {
				c.mu.RUnlock()
				// log.Printf("[%sClient] connection is nil, stopping ping loop", c.wc.name)
				return nil
			}
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			c.mu.RUnlock()
			if err != nil {
				return fmt.Errorf("failed to write ping message: %w", err)
			}
			// log.Printf("[%sClient] ping sent", c.wc.name)
		}
	}
}

func (c *WebsocketClient) closeConnection(safeClose bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		// log.Printf("[%sClient] connection already closed", c.wc.name)
		return
	}

	// log.Printf("[%sClient] closing connection (safeClose=%t)", c.wc.name, safeClose)

	if safeClose {
		if err := c.conn.WriteMessage(
			websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		); err != nil {
			// log.Printf("[%sClient] failed to write close message: %v", c.wc.name, err)
		}
	}

	if err := c.conn.Close(); err != nil {
		// log.Printf("[%sClient] failed to close connection with: %v", c.wc.name, err)
	}

	c.conn = nil
	// log.Printf("[%sClient] connection closed", c.wc.name)
}

func (c *WebsocketClient) runSession(ctx context.Context) error {
	// log.Printf("[%sClient] establishing new connection...", c.wc.name)
	connCtx, cancel := context.WithCancel(ctx)
	c.cancel = cancel

	if err := c.connect(connCtx); err != nil {
		return fmt.Errorf("connect error: %w", err)
	}

	if err := c.resubscribe(); err != nil {
		return fmt.Errorf("resubscribe error: %w", err)
	}

	c.wg.Add(2)

	errCh := make(chan error, 2)

	go func() {
		defer c.wg.Done()
		errCh <- c.readMessages(connCtx)
	}()

	go func() {
		defer c.wg.Done()
		errCh <- c.pingLoop(connCtx)
	}()

	reconnTimer := time.NewTimer(c.wc.reconnPeriod)
	defer reconnTimer.Stop()

	defer func() {
		// log.Printf("[%sClient] cleaning up session", c.wc.name)
		cancel()
		c.closeConnection(false)
		c.wg.Wait()
	}()

	select {
	case <-ctx.Done():
		// log.Printf("[%sClient] context canceled during session", c.wc.name)
		return nil
	case <-errCh:
		// log.Printf("[%sClient] error occurred, reconnecting", c.wc.name)
		return ErrReconnectNeeded
	case <-reconnTimer.C:
		// log.Printf("[%sClient] reconnect timer expired", c.wc.name)
		return ErrReconnectNeeded
	}
}

func (c *WebsocketClient) stop() {
	c.once.Do(func() {
		if c.cancel != nil {
			c.cancel()
		}
		c.closeConnection(true)
	})
}
