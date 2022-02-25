package db_service

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/thinkerajay/tracker-service/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var collection *mongo.Collection

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	collection = client.Database("tracker-service").Collection("views_data")


}

func Consume(ctx context.Context) {
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	writeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	partitionConsumer, err := consumer.ConsumePartition("tracker-service-views-data", 1, 0)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()


	messageCounter := 0
	var documents []interface{}
	var event types.Event
	defer cancel()
ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Consumed message %s with offset %d\n", string(msg.Value), msg.Offset)
			if messageCounter == 1e3{
				_, err := collection.InsertMany(writeCtx, documents, nil)
				if err != nil {
					log.Println(err)
				}
				messageCounter = 0
			}
			err = json.Unmarshal(msg.Value, &event)
			documents = append(documents, event)
			messageCounter++

		case <-ctx.Done():
			break ConsumerLoop
		}
	}

}
