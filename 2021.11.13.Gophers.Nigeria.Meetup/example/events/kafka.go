package events

import (
	"time"

	"github.com/AndersonQ/talks/2021.11.13/example/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/propagation"
)

type Headers map[string]string
type Event struct {
	Headers Headers
	Key     string
	Payload string
}

// Ensure Headers implements OTel propagation.TextMapCarrier
var _ = propagation.TextMapCarrier(Headers{})

func (h Headers) KafkaHeader() []kafka.Header {
	var kh []kafka.Header

	for k, v := range h {
		kh = append(kh, kafka.Header{
			Key:   k,
			Value: []byte(v),
		})
	}

	return kh
}

func (h Headers) Get(key string) string {
	return h[key]
}

func (h Headers) Set(key, value string) {
	h[key] = value
}

func (h Headers) Keys() []string {
	keys := make([]string, len(h))
	for k := range h {
		keys = append(keys, k)
	}

	return keys
}

type Consumer struct {
	kc   *kafka.Consumer
	l    zerolog.Logger
	Done chan struct{}
}

func NewProducer(cfg config.Config) *kafka.Producer {
	conf := &kafka.ConfigMap{
		"bootstrap.servers":  cfg.KafkaBootstrapServer,
		"message.timeout.ms": 6000,
	}

	p, err := kafka.NewProducer(conf)
	if err != nil {
		cfg.Logger.Panic().Err(err).Msg("failed to create kafkaProducer")
	}

	return p
}
func NewConsumer(cfg config.Config) Consumer {
	conf := &kafka.ConfigMap{
		"group.id":           cfg.KafkaGroupID,
		"bootstrap.servers":  cfg.KafkaBootstrapServer,
		"session.timeout.ms": 6000,
		"auto.offset.reset":  "earliest",
	}
	consumer, err := kafka.NewConsumer(conf)
	if err != nil {
		cfg.Logger.Panic().Err(err).Msg("could not create kafka consumer")
	}

	if err := consumer.SubscribeTopics([]string{cfg.Topic}, nil); err != nil {
		cfg.Logger.Panic().Err(err).Msg("\"could not subscribe to topics")
	}

	return Consumer{
		kc:   consumer,
		l:    cfg.Logger,
		Done: make(chan struct{}),
	}
}

func (c Consumer) Run() chan Event {
	ch := make(chan Event, 1)
	go func() {
		timeoutMs := int(time.Second)
		for {
			select {
			case <-c.Done:
				return
			default:
				kev := c.kc.Poll(timeoutMs)
				switch kmt := kev.(type) {
				case *kafka.Message:
					ch <- newEvent(kmt)
				case kafka.Error:
					if kmt.Code() != kafka.ErrTimedOut {
						c.l.Error().Msgf("failed to read kafka message: %v", kmt)
					}
				}
			}
		}
	}()

	return ch
}

func newEvent(m *kafka.Message) Event {
	h := Headers{}
	for _, kh := range m.Headers {
		h[kh.Key] = string(kh.Value)
	}

	return Event{
		Headers: h,
		Key:     string(m.Key),
		Payload: string(m.Value),
	}
}
