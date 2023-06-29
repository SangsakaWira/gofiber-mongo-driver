package controller

import(
	"log"
    "context"
	"github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    model "rest-api/model"
    "rest-api/config"
)

var ctx = context.Background()

func CreateBook (c *fiber.Ctx) error {

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Define a struct to match the JSON structure
	var book model.Book

	// Parse the JSON request body into the struct
	if err := c.BodyParser(&book); err != nil {
		log.Fatal(err)
	}

	_, err = db.Collection("books").InsertOne(ctx,model.Book{book.Judul,book.Harga,book.Deskripsi})
	if err != nil {
		log.Fatal(err.Error())
	}

	return c.Send(c.Body()) // []byte("user=john")
}

func FindBook (c *fiber.Ctx) error {
	db, err := database.Connect()
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
}

func FindBookById (c *fiber.Ctx) error {
	db, err := database.Connect()
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
}

func UpdateBook (c *fiber.Ctx) error {

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Define a struct to match the JSON structure
	var book model.BookUpdate

	// Parse the JSON request body into the struct
	if err := c.BodyParser(&book); err != nil {
		log.Fatal(err)
	}

	objID, _ := primitive.ObjectIDFromHex(book.ID)

	var selector = bson.M{"_id": objID}

	_, err = db.Collection("books").UpdateOne(ctx,selector,bson.M{"$set":model.Book{book.Judul,book.Harga,book.Deskripsi}})
	
	if err != nil {
		log.Fatal(err.Error())
	}

	return c.SendString("Successfully Update Data!") //
}

func DeleteBook (c *fiber.Ctx) error {
	db, err := database.Connect()
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
}