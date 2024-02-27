package brute

import (
	"sync"

	"github.com/streadway/amqp"
	"github.com/vicsobdev/LittleBruteClient/queue"
	"go.uber.org/zap"
)

type Brute struct {
	WordListPath string
	ProxyPath    string
	OutPath      string
	Verbose      bool
	Debug        bool
	logger       zap.Logger
	rabbit       Rabbit
	stats        Stats
}

type Stats struct {
	hits   queue.Queue
	bad    queue.Queue
	total  int32
	errors []string
	mx     sync.Mutex
}

type RabbitConfig struct {
	Username string
	Password string
	Host     string
	Port     string
}

type Rabbit struct {
	Config        RabbitConfig
	Conn          *amqp.Connection
	Channel       *amqp.Channel
	RetrieveQueue amqp.Queue
	PublishQueue  amqp.Queue
}

type response struct {
	Item    string   `json:"item,omitempty"`
	Capture string   `json:"capture,omitempty"`
	Status  int      `json:"status,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

func NewBrute(listPath, proxyPath, outPath string, verbose bool) (*Brute, error) {

	return &Brute{
		WordListPath: listPath,
		ProxyPath:    proxyPath,
		OutPath:      outPath,
		Verbose:      verbose,
	}, nil
}
