package main

import "github.com/gofiber/fiber/v2"
import (
    "context"
    "log"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

var ctx = context.Background()

type Book struct {
    Judul  string `bson:"judul"`
    Harga string    `bson:"harga"`
    Deskripsi string    `bson:"deskripsi"`
}

type BookUpdate struct {
    ID string `bson:"id"`
    Judul  string `bson:"judul"`
    Harga string    `bson:"harga"`
    Deskripsi string    `bson:"deskripsi"`
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

func main() {
    app := fiber.New()

    app.Post("/", func(c *fiber.Ctx) error {

        db, err := connect()
        if err != nil {
            log.Fatal(err.Error())
        }

        // Define a struct to match the JSON structure
		var book Book

		// Parse the JSON request body into the struct
		if err := c.BodyParser(&book); err != nil {
			log.Fatal(err)
		}

        _, err = db.Collection("books").InsertOne(ctx,Book{book.Judul,book.Harga,book.Deskripsi})
        if err != nil {
            log.Fatal(err.Error())
        }
    
        return c.Send(c.Body()) // []byte("user=john")
    })

    app.Get("/", func(c *fiber.Ctx) error {
        db, err := connect()
        if err != nil {
            log.Fatal(err.Error())
        }

        cursor, err := db.Collection("books").Find(context.Background(), bson.D{})
        if err != nil {
            log.Fatal(err)
        }

        defer cursor.Close(context.Background())

        var documents []bson.M

        for cursor.Next(context.Background()) {
            var doc bson.M
            if err := cursor.Decode(&doc); err != nil {
                log.Fatal(err)
            }
            documents = append(documents, doc)
        }

        if err := cursor.Err(); err != nil {
            log.Fatal(err)
        }

        return c.JSON(documents)
    })

    app.Get("/:id", func(c *fiber.Ctx) error {
        db, err := connect()
        if err != nil {
            log.Fatal(err.Error())
        }

        collection := db.Collection("books")

        id := c.Params("id")

        objID, _ := primitive.ObjectIDFromHex(id)

        cursor, err := collection.Find(context.Background(), bson.M{"_id":objID})

        if err != nil {
            log.Fatal(err)
        }

        defer cursor.Close(context.Background())

        var documents []bson.M

        for cursor.Next(context.Background()) {
            var doc bson.M
            if err := cursor.Decode(&doc); err != nil {
                log.Fatal(err)
            }
            documents = append(documents, doc)
        }

        if err := cursor.Err(); err != nil {
            log.Fatal(err)
        }

        return c.JSON(documents)
    })

    app.Patch("/", func(c *fiber.Ctx) error {

        db, err := connect()
        if err != nil {
            log.Fatal(err.Error())
        }

        // Define a struct to match the JSON structure
		var book BookUpdate

		// Parse the JSON request body into the struct
		if err := c.BodyParser(&book); err != nil {
			log.Fatal(err)
		}

        objID, _ := primitive.ObjectIDFromHex(book.ID)

        var selector = bson.M{"_id": objID}

        _, err = db.Collection("books").UpdateOne(ctx,selector,bson.M{"$set":Book{book.Judul,book.Harga,book.Deskripsi}})
        
        if err != nil {
            log.Fatal(err.Error())
        }
    
        return c.SendString("Successfully Update Data!") //
    })

    app.Delete("/:id", func(c *fiber.Ctx) error {
        db, err := connect()
        if err != nil {
            log.Fatal(err.Error())
        }

        collection := db.Collection("books")

        id := c.Params("id")

        objID, _ := primitive.ObjectIDFromHex(id)

        _, err = collection.DeleteOne(context.Background(), bson.M{"_id":objID})

        if err != nil {
            log.Fatal(err)
        }

        return c.SendString("Remove success!")
    })

    app.Listen(":3000")
}
