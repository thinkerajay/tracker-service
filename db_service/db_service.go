package db_service

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/thinkerajay/tracker-service/constants"
	"github.com/thinkerajay/tracker-service/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)


var collection *mongo.Collection

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(constants.DB_URL))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	collection = client.Database(constants.DB_NAME).Collection(constants.COLLECTION_NAME)

}

func FetchViewsCount(startDate string, endDate string)(int64, error){
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	records, err := collection.CountDocuments(ctx, bson.M{"created_at":bson.M{"$gt": startDate, "$lt": endDate}})
	if err != nil{
		log.Println(err)
		return 0,err
	}

	return records, nil
}

func Consume(ctx context.Context) {
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	//writeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	partitionConsumer, err := consumer.ConsumePartition(constants.KAFKA_TOPIC, 1, 0)
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
	//defer cancel()
ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Consumed message %s with offset %d\n", string(msg.Value), msg.Offset)
			if messageCounter > 0 {
				_, err := collection.InsertMany(context.TODO(), documents, nil)
				if err != nil {
					log.Println(err)
				}
				messageCounter = 0
				documents = []interface{}{}
			}
			err = json.Unmarshal(msg.Value, &event)
			documents = append(documents, bson.M{
				"created_at": event.CreatedAt,
				"type":       event.Type,
				"page_url":   event.PageUrl,
				"user_id":    event.UserId,
			})
			messageCounter++

		case <-ctx.Done():
			break ConsumerLoop
		}
	}

}
