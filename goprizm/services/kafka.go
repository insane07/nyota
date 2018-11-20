package services

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"goprizm/fileutils"
	"goprizm/log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Reference - https://gowalker.org/github.com/Shopify/sarama#Config
type kafkaOptions struct {
	connTimeout   time.Duration // Connection timeout in seconds
	bufferSize    int           // Buffer size of the message store
	maxRetries    int           // Maximum retries to connect to Kafka broker
	retryInterval time.Duration // Retry interval in seconds on a failure
	retryBackOff  time.Duration // Retry backoff in seconds
	serverCert    string        // Kafka server certificate
	caCerts       string        // CA Certs bundle
	tlsEnable     bool          // Flag to control TLS
}

var (
	kafkaOpts = kafkaOptions{
		connTimeout:   60 * time.Second,
		bufferSize:    4096,
		maxRetries:    10,
		retryInterval: 30 * time.Second,
		retryBackOff:  1 * time.Second,
		serverCert:    filepath.Join(Etc, "/certs/server.pem"),
		caCerts:       filepath.Join(Etc, "/certs/cabundle.pem"),
		tlsEnable:     false,
	}
)

// Return Kafka configuration. Instantiate the instance from the
// list of brokers. If brokers is empty then the localhost will be
// used
func Kafka() ([]string, *sarama.Config) {
	config := sarama.NewConfig()
	config.Net.DialTimeout = kafkaOpts.connTimeout

	// Producer settings
	// The level of ack reliability needed from broker. The default value from
	// the library is WaitForLocal. Change this to as needed. Changing this from
	// NoResponse will also need to incorporate Producer.Timeout setting.
	config.Producer.RequiredAcks = sarama.NoResponse
	config.Producer.Retry.Max = kafkaOpts.maxRetries
	// Delivered messages will be returned on the default channel. Default is false
	// and setting it to true
	config.Producer.Return.Successes = false
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	// Metadata settings
	config.Metadata.Retry.Max = kafkaOpts.maxRetries
	config.Metadata.Retry.Backoff = kafkaOpts.retryBackOff
	config.Metadata.Full = true

	// Consumer settings
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// Miscellaneous
	// This should be set to actual TenantId for every request sent to the
	// broker. This facilitates logging, debugging and troubleshooting aspects.
	//config.ClientID = "tenant_id"
	config.ChannelBufferSize = kafkaOpts.bufferSize

	// Setup TLS
	if kafkaOpts.tlsEnable {
		err := kfAddTLSConfig(config)
		if err != nil {
			log.Errorf("Failed to setup TLS with Kafka")
		} else {
			log.Debugf("Successfully enabled TLS with Kafka")
		}
	} else {
		config.Net.TLS.Config = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	brokers := os.Getenv("KAFKA")
	if brokers == "" {
		brokers = "localhost:9092"
	}

	return strings.Split(brokers, ","), config
}

// Initialize the TLS configuration
func kfAddTLSConfig(config *sarama.Config) error {
	cert, err := fileutils.ReadBytes(kafkaOpts.serverCert)
	if err != nil {
		return err
	}

	certs, err := initSystemCertPool()
	if err != nil {
		return err
	}

	config.Net.TLS.Config = &tls.Config{
		RootCAs: certs,
	}

	config.Net.TLS.Enable = true

	if !config.Net.TLS.Config.RootCAs.AppendCertsFromPEM(cert) {
		return fmt.Errorf("failed to append server cert")
	}

	return nil
}

// Initialize the System Certs pool
func initSystemCertPool() (*x509.CertPool, error) {
	certs, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}

	if certs == nil {
		return x509.NewCertPool(), nil
	}

	return certs, nil
}

// Connect to Kafka brokers and return the AsyncProducer
func ConnectAsync(brokers []string, config *sarama.Config) (producer sarama.AsyncProducer, err error) {
	err = WithRetry(kafkaOpts.maxRetries, kafkaOpts.retryInterval, func() (e error) {
		producer, e = sarama.NewAsyncProducer(brokers, config)
		return
	})

	if err != nil {
		log.Errorf("Failed to connect to Kafka. Error :%v", err)
	} else {
		log.Debugf("Connected to Kafka")
	}

	return
}

// Connect to Kafka brokers and return the SyncProducer
func ConnectSync(brokers []string, config *sarama.Config) (producer sarama.SyncProducer, err error) {
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	err = WithRetry(kafkaOpts.maxRetries, kafkaOpts.retryInterval, func() (e error) {
		producer, e = sarama.NewSyncProducer(brokers, config)
		return
	})

	if err != nil {
		log.Errorf("Failed to connect to Kafka. Error :%v", err)
	} else {
		log.Debugf("Connected to Kafka")
	}

	return
}

// Generic function to retry connection
func WithRetry(maxRetries int, interval time.Duration, do func() error) (err error) {
	for count := 0; count < maxRetries; count++ {
		err = do()
		if err != nil {
			log.Warnf("Retry attempt#%d | error=%v", count, err)
			time.Sleep(interval)
		} else {
			break
		}
	}
	return err
}

var (
	// KafkaMsgChanSize - default size of consumer msg chan
	KafkaMsgChanSize = 1000
)

func Brokers() string {
	brokers := os.Getenv("KAFKA")
	if brokers == "" {
		brokers = "localhost:9092"
	}
	return brokers
}

type KafkaProducer struct {
	*kafka.Producer
}

func NewKafkaProducer() (kf KafkaProducer, err error) {
	kf.Producer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": Brokers(),
	})
	if err != nil {
		return kf, err
	}

	// Consume errors and other events.
	go func() {
		for ev := range kf.Events() {
			switch e := ev.(type) {
			case *kafka.Message:
				if err := e.TopicPartition.Error; err != nil {
					log.Errorf("kafka - send msg:%+v err:%v", e, err)
				}

			default:
				log.Errorf("kafka - producer unknown event:%+v", e)
			}
		}
	}()

	log.Printf("kafka - producer connected")
	return kf, nil
}

type KafkaConsumer struct {
	*kafka.Consumer
	Messages chan *kafka.Message
}

func NewKafkaConsumer(group string, topics []string) (kf KafkaConsumer, err error) {
	kf.Consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               Brokers(),
		"group.id":                        group,
		"session.timeout.ms":              10000,
		"go.events.channel.enable":        true,
		"log_level":                       4,
		"enable.partition.eof":            false,
		"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"},
		"go.application.rebalance.enable": true,
		//Default values
		//enable.auto.commit:              true,
		//auto.commit.interval.ms:         5000,
	})

	if err != nil {
		return kf, err
	}
	if err = kf.SubscribeTopics(topics, nil); err != nil {
		return kf, nil
	}

	kf.Messages = make(chan *kafka.Message, 1000)

	go func() {
		for ev := range kf.Events() {
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				log.Printf("kafka - assign partition event:%+v", e.Partitions)
				kf.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				log.Printf("kafka - unassign partition event:%+v", ev)
				kf.Unassign()
			case *kafka.Message:
				kf.Messages <- e
			case kafka.PartitionEOF:
				log.Printf("kafka - partition eof %+v", e)
			case kafka.Error:
				log.Printf("kafka - error event %+v", e)
			}
		}
	}()

	log.Printf("kafka - consumer connected")
	return kf, nil
}
