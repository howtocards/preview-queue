package main

import (
	"facade/src/internal/api"
	"facade/src/internal/config"
	"facade/src/internal/queue"
	"fmt"

	"github.com/powerman/structlog"
	flag "github.com/spf13/pflag"
	"github.com/streadway/amqp"
)

var (
	log = structlog.New()

	cfg struct {
		api   api.Config
		queue queue.Config
	}
)

func init() {
	flag.StringVar(&cfg.api.Host, "host", config.APIHost, "server host")
	flag.IntVar(&cfg.api.Port, "port", config.APIPort, "server port")
	flag.StringVar(&cfg.api.BasePath, "base-path", config.APIBasePath, "base path swagger")

	flag.StringVar(&cfg.queue.Name, "queue.name", "event", "queue name")
	flag.BoolVar(&cfg.queue.Durable, "durable", true, "the queue will survive a broker restart")
	flag.BoolVar(&cfg.queue.AutoDelete, "auto-delete", false, "queue that has had at least one consumer is deleted when last consumer unsubscribes")
	flag.BoolVar(&cfg.queue.Exclusive, "exclusive", false, "used by only one connection and the queue will be deleted when that connection closes")
	flag.BoolVar(&cfg.queue.NoWait, "no-wait", false, "")

	flag.StringVar(&cfg.queue.Pass, "rabbit.pass", "rabbitmq", "RabbitMQ password")
	flag.StringVar(&cfg.queue.User, "rabbit.user", "rabbitmq", "RabbitMQ user")
	flag.StringVar(&cfg.queue.Host, "rabbit.host", "localhost", "RabbitMQ host")
	flag.IntVar(&cfg.queue.Port, "rabbit.port", 5672, "RabbitMQ port")

	flag.Parse()
}

func main() {
	log.Fatal(run())
}

func run() error {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d",
		cfg.queue.User,
		cfg.queue.Pass,
		cfg.queue.Host,
		cfg.queue.Port,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	defer log.WarnIfFail(conn.Close)

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer log.WarnIfFail(ch.Close)

	_, err = ch.QueueDeclare(cfg.queue.Name, cfg.queue.Durable, cfg.queue.AutoDelete,
		cfg.queue.Exclusive, cfg.queue.NoWait, nil)
	if err != nil {
		return err
	}

	server, err := api.NewServer(log, queue.New(ch, cfg.queue), cfg.api)
	if err != nil {
		return err
	}

	return server.Serve()
}
