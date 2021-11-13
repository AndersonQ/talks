package config

import (
	"encoding/json"
	"os"

	"github.com/caarlos0/env"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//revive:disable:exported Obvious config struct.
type Config struct {
	ServiceName string `env:"SERVICE_NAME" envDefault:"otel-example"`
	Port        string `env:"PORT" envDefault:"4242"`

	Debug bool `env:"DEBUG" envDefault:"false"`

	KafkaBootstrapServer string `env:"KAFKA_BOOTSTRAP_SERVERS" envDefault:"localhost:9092"`
	KafkaGroupID         string `env:"KAFKA_GROUP_ID" envDefault:"otel-example"`
	Topic                string `env:"KAFKA_TOPIC" envDefault:"otel-example"`
	OTelExporterEndpoint string `env:"OTEL_EXPORTER_OTLP_ENDPOINT" envDefault:"localhost:55680"` // use "localhost:16685"

	Logger zerolog.Logger `json:"-"`
}

func Parse(serviceName string) Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Panic().Err(err).Msg("failed to load environment variables")
	}

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	cfg.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
		With().
		Str("service", serviceName).
		Logger()

	cfg.ServiceName = serviceName

	bs, _ := json.MarshalIndent(cfg, "", "  ")
	cfg.Logger.Debug().Msgf("config: %s", bs)

	return cfg
}
