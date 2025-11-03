package config

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
)

type EnvConfig struct {
	WsURL              string   `envconfig:"WS_URL"`
	WsMsgCap           int      `envconfig:"WS_MSG_CAP"`
	WsReadLimit        int64    `envconfig:"WS_READ_LIMIT"`
	WsReconnPeriodHr   int      `envconfig:"WS_RECONN_PERIOD_HR"`
	WsRetryBackoffSec  int      `envconfig:"WS_RETRY_BACKOFF_SEC"`
	WsPingWaitMin      int      `envconfig:"WS_PING_WAIT_MIN"`
	WsPongPeriodMin    int      `envconfig:"WS_PONG_PERIOD_MIN"`
	StartupSymbols     []string `envconfig:"STARTUP_SYMBOLS"`
	StartupStreamTypes []string `envconfig:"STARTUP_STREAM_TYPES"`
	ServerName         string   `envconfig:"SERVER_NAME"`
	ServerVersion      string   `envconfig:"SERVER_VERSION"`
	DbURL              string   `envconfig:"DB_URL"`
}

func MustConfig() (*EnvConfig, error) {
	var cfg EnvConfig

	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
