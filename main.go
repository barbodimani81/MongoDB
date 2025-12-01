package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"mongo/generator"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	generatedCount := flag.Int("generatedCount", 10000, "Generated count of documents")
	batchLimit := flag.Int("batchLimit", 1000, "Batch limit")
	flag.Parse()

	// 1. Connect to Mongo
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		log.Fatalf("connect error: %v", err)
	}
	defer func() {
		_ = client.Disconnect(context.Background())
	}()

	collection := client.Database("testdb").Collection("items")

	// 2. Use generator to produce N random items as JSON
	ch := generator.Generate(*generatedCount)

	// 3. For each generated JSON payload, unmarshal and insert into Mongo
	start := time.Now()
	var batch []any
	for payload := range ch {
		var doc bson.M
		if err := json.Unmarshal(payload, &doc); err != nil {
			log.Fatalf("unmarshal error: %v", err)
		}

		batch = append(batch, doc)
		if len(batch) >= *batchLimit {
			_, err := collection.InsertMany(ctx, batch)
			if err != nil {
				log.Fatalf("insert error: %v", err)
			}
			fmt.Printf("Inserted: %d items\n", *batchLimit)
			batch = batch[:0]
		}

		if len(batch) > 0 {
			if _, err := collection.InsertMany(ctx, batch); err != nil {
				log.Fatalf("final insert error: %v", err)
			}
		}
	}
	fmt.Printf("Inserted final batch of %d items\n", len(batch))

	elapsed := time.Since(start)
	fmt.Println("Done inserting in: ", elapsed)
}
