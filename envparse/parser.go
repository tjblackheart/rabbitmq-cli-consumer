package envparse

import (
	"os"

	"github.com/ricbra/rabbitmq-cli-consumer/config"
)

func OverrideConfig(cfg *config.Config) {

	if os.Getenv("RMQ_HOST") != "" {
		cfg.RabbitMq.Host = os.Getenv("RMQ_HOST")
	}

	if os.Getenv("RMQ_USER") != "" {
		cfg.RabbitMq.Username = os.Getenv("RMQ_USER")
	}

	if os.Getenv("RMQ_PASSWORD") != "" {
		cfg.RabbitMq.Password = os.Getenv("RMQ_PASSWORD")
	}

	if os.Getenv("RMQ_PORT") != "" {
		cfg.RabbitMq.Port = os.Getenv("RMQ_PORT")
	}

	if os.Getenv("RMQ_VHOST") != "" {
		cfg.RabbitMq.Vhost = os.Getenv("RMQ_VHOST")
	}

	if os.Getenv("RMQ_QUEUE") != "" {
		cfg.RabbitMq.Queue = os.Getenv("RMQ_QUEUE")
	}

	if os.Getenv("RMQ_EXCHANGE") != "" {
		cfg.Exchange.Name = os.Getenv("RMQ_EXCHANGE")
	}

	if os.Getenv("RMQ_ROUTING_KEY") != "" {
		cfg.QueueSettings.Routingkey = os.Getenv("RMQ_ROUTING_KEY")
	}

	if os.Getenv("RMQ_LOG_INFO") != "" {
		cfg.Logs.Info = os.Getenv("RMQ_LOG_INFO")
	}

	if os.Getenv("RMQ_LOG_ERR") != "" {
		cfg.Logs.Error = os.Getenv("RMQ_LOG_ERR")
	}

}
