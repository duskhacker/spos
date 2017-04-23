package cafe

import (
	"log"
	"sync"

	"github.com/bitly/go-nsq"
	"github.com/duskhacker/cqrsnu/serializer"
	"github.com/p4tin/goaws/app"
)

var (
	nsqConfig     = nsq.NewConfig()
	connectToNSQD bool
	mutex         sync.RWMutex

	lookupdHTTPAddrs app.StringArray
	nsqdTCPAddr      string

	s serializer.Serializer

	consumers []*nsq.Consumer
	producer  *nsq.Producer
)

func init() {
	s = serializer.NewSerializer()
}

func Init() {
	Tabs = NewTabs()
	InitConsumers()
}

func InitConsumers() {
	NewConsumer(OpenTabTopic, OpenTabTopic+"Consumer", OpenTabHandler)
	NewConsumer(PlaceOrderTopic, PlaceOrderTopic+"Consumer", PlaceOrderHandler)
	NewConsumer(MarkDrinksServedTopic, MarkDrinksServedTopic+"Consumer", MarkDrinksServedHandler)
	NewConsumer(MarkFoodPreparedTopic, MarkFoodPreparedTopic+"Consumer", MarkFoodPreparedHandler)
	NewConsumer(MarkFoodServedTopic, MarkFoodServedTopic+"Consumer", MarkFoodServedHandler)
	NewConsumer(FoodServedTopic, FoodServedTopic+"Consumer", FoodServedHandler)
	NewConsumer(DrinksServedTopic, DrinksServedTopic+"Consumer", DrinksServedHandler)
	NewConsumer(CloseTabTopic, CloseTabTopic+"Consumer", CloseTabHandler)
}

func SetNsqdTCPAddr(address string) {
	nsqdTCPAddr = address
}

func SetLookupdHTTPAddrs(addresses app.StringArray) {
	lookupdHTTPAddrs = addresses
}

func SetConnectToNSQD(v bool) {
	connectToNSQD = v
}

func NewConsumer(topic, channel string, handler func(*nsq.Message) error) *nsq.Consumer {
	nsqConfig.UserAgent = channel

	consumer, err := nsq.NewConsumer(topic, channel, nsqConfig)
	if err != nil {
		log.Fatalf("%s:%s; NewConsumer: %s", topic, channel, err)
	}
	consumer.SetLogger(nil, 0)

	consumer.AddHandler(nsq.HandlerFunc(nsq.HandlerFunc(handler)))

	if connectToNSQD {
		if err = consumer.ConnectToNSQD(nsqdTCPAddr); err != nil {
			log.Fatalf("%s:%s; ConnectToNSQLookupds: %s", topic, channel, err)
		}

	} else {
		if err = consumer.ConnectToNSQLookupds(lookupdHTTPAddrs); err != nil {
			log.Fatalf("%s:%s; ConnectToNSQLookupds: %s", topic, channel, err)
		}
	}

	consumers = append(consumers, consumer)
	return consumer
}

func StopAllConsumers() {
	for _, consumer := range consumers {
		consumer.Stop()
	}
}

func newProducer() *nsq.Producer {
	nsqConfig.UserAgent = "cqrsnuProducer"

	producer, err := nsq.NewProducer(nsqdTCPAddr, nsqConfig)
	if err != nil {
		log.Fatalf("error creating nsq.Producer: %s", err)
	}
	producer.SetLogger(nil, 0)

	if err = producer.Ping(); err != nil {
		log.Fatalf("error pinging nsqd: %s\n", err)
	}

	return producer
}

func Send(topic string, message interface{}) {
	if producer == nil {
		producer = newProducer()
	}

	if err := producer.Publish(topic, s.Serialize(message)); err != nil {
		log.Fatalf("Send %s error: %s\n", topic, err)
	}
}
