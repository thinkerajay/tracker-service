package queue_service

import (
	"encoding/json"
	. "github.com/Shopify/sarama"
	"github.com/thinkerajay/tracker-service/types"
	"log"
)

var enqueuer SyncProducer
var err error

func init() {
	enqueuer, err = NewSyncProducer([]string{"localhost:9092"}, nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func Enqueue(event *types.Event) error {
	defer func(enqueuer SyncProducer) {
		err := enqueuer.Close()
		if err != nil {
			log.Println(err)
		}
	}(enqueuer)
	eventData, err := json.Marshal(event)
	if err != nil {
		log.Println(err)
		return err
	}
	msg := &ProducerMessage{Topic: "tracker-service-views-data", Value: StringEncoder(eventData), Key: StringEncoder(event.PageUrl), Partition: 1}
	partition, offset, err := enqueuer.SendMessage(msg)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Enqueued %+v to %d with %d", msg, partition, offset)
	return nil
}
