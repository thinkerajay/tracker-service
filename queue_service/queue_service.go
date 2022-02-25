package queue_service

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/thinkerajay/tracker-service/constants"
	"github.com/thinkerajay/tracker-service/types"
	"log"
)

var enqueuer sarama.SyncProducer
var err error

func init() {
	config := sarama.NewConfig()
	admin, err := sarama.NewClusterAdmin([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatal("Error while creating cluster admin: ", err.Error())
	}
	defer func() { _ = admin.Close() }()
	err = admin.CreateTopic(constants.KAFKA_TOPIC, &sarama.TopicDetail{
		NumPartitions:     2,
		ReplicationFactor: 1,
	}, false)
	enqueuer, err = sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func Enqueue(event *types.Event) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		log.Println(err)
		return err
	}
	msg := &sarama.ProducerMessage{Topic: constants.KAFKA_TOPIC, Value: sarama.StringEncoder(eventData), Key: sarama.StringEncoder(event.PageUrl), Partition: 1}
	partition, offset, err := enqueuer.SendMessage(msg)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Enqueued %+v to %d with %d", msg, partition, offset)
	return nil
}
