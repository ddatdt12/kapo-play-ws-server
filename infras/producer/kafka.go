package producer

import (
	"os"

	systemlog "log"

	"github.com/IBM/sarama"
	"github.com/ddatdt12/kapo-play-ws-server/configs"
	"github.com/rs/zerolog/log"
)

type DataCollector struct {
	Collector sarama.SyncProducer
}

func NewDataCollector() *DataCollector {
	// brokerListStr := configs.EnvConfigs.BROKER_LIST
	// version, err := sarama.ParseKafkaVersion(configs.EnvConfigs.KAFKA_VERSION)

	// if err != nil {
	// 	log.Fatal().Err(err).Msg("Error parsing Kafka version")
	// }
	sarama.Logger = systemlog.New(os.Stdout, "[sarama] ", systemlog.LstdFlags)
	brokerList := []string{"pkc-312o0.ap-southeast-1.aws.confluent.cloud:9092"}

	// For the data collector, we are looking for strong consistency semantics.
	// Because we don't change the flush settings, sarama will try to produce messages
	// as fast as possible to keep latency low.

	log.Info().Msgf("Creating data collector with brokers %v", configs.EnvConfigs)
	config := sarama.NewConfig()
	config.Version = sarama.V1_0_1_0
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 3                    // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true
	config.ClientID = "kapo-play-ws-server"
	// config.Metadata.Full = true
	config.Net.SASL.Enable = true
	config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	config.Net.SASL.User = "MVMZADWS2NZNNB34"
	config.Net.SASL.Password = "I749pubsuHG9/kkHjvXR1nb2Ut8K9P5loD/3xaK50o/Qwv1g+CWpOVKh0zLzRWqb"
	config.Net.SASL.Handshake = true

	// On the broker side, you may want to change the following settings to get
	// stronger consistency guarantees:
	// - For your broker, set `unclean.leader.election.enable` to false
	// - For the topic, you could increase `min.insync.replicas`.

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating sync producer")
	}

	log.Info().Msg("Data collector created")

	return &DataCollector{
		Collector: producer,
	}
}
