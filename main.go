package main

import "github.com/gofiber/fiber/v2"
import (
    "context"
    "log"
    "fmt"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
)

var ctx = context.Background()

type book struct {
    judul  string `bson:"judul"`
    harga string    `bson:"harga"`
    deskripsi string    `bson:"deskripsi"`
}

func connect() (*mongo.Database, error) {
    clientOptions := options.Client()
    clientOptions.ApplyURI("mongodb://localhost:27017")
    client, err := mongo.NewClient(clientOptions)
    if err != nil {
        return nil, err
    }

    err = client.Connect(ctx)
    if err != nil {
        return nil, err
    }

    return client.Database("libnation-back-end"), nil
}

func find() {
    db, err := connect()
    if err != nil {
        log.Fatal(err.Error())
    }

    csr, err := db.Collection("books").Find(ctx, bson.M{"judul": "Harry Potter: The Chamber of Secret"})
    if err != nil {
        log.Fatal(err.Error())
    }
    defer csr.Close(ctx)

    result := make([]book, 0)
    for csr.Next(ctx) {
        var row book
        err := csr.Decode(&row)
        if err != nil {
            log.Fatal(err.Error())
        }

        result = append(result, row)
    }

    if len(result) > 0 {
        fmt.Println("judul  :", result[0].judul)
        fmt.Println("harga :", result[0].harga)
    }
}


func main() {
    app := fiber.New()

    app.Get("/", func(c *fiber.Ctx) error {
        // Set up the MongoDB connection
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ping the MongoDB server to verify the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Access the MongoDB database and collection
	database := client.Database("libnation-back-end")
	collection := database.Collection("books")

    // Find all documents in the collection
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.Background())

	// Create a slice to hold the retrieved documents
	var documents []bson.M

	// Iterate over the cursor to retrieve documents
	for cursor.Next(context.Background()) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			log.Fatal(err)
		}
		// Append each document to the slice
		documents = append(documents, doc)
	}

	// Check if any errors occurred during iteration
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	// Process the retrieved documents as needed
	// For example, you can return them as JSON response
	return c.JSON(documents)
    })

    app.Listen(":3000")
}
