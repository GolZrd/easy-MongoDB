package main

import (
	"context"
	"log"
	"time"
	"web_mongo/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Client struct {
	ID    string `json:"id" bson:"_id,omitempty"`
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
	Age   string `json:"age" bson:"age"`
}

var mongoClient *mongo.Client

func main() {
	mongoClient, err := db.ConnectToMongo()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	New(mongoClient)

	viewsEngine := html.New("./templates", ".html")
	app := fiber.New(fiber.Config{Views: viewsEngine})

	app.Get("/", IndexHandler)
	app.Get("/create", showCreateHandler)
	app.Post("/create", createHandler)
	app.Get("/edit/:id", EditPage)
	app.Post("/edit/:id", EditHandler)
	app.Post("/delete/:id", DeleteHandler)

	app.Listen(":8000")
}

func New(mongo *mongo.Client) {
	mongoClient = mongo
}

func returnCollectionPointer(collection string) *mongo.Collection {
	return mongoClient.Database("test3").Collection(collection)
}

func IndexHandler(c *fiber.Ctx) error {
	collection := returnCollectionPointer("clients")

	filter := bson.D{}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Println("ошибка в курсоре", err)
		c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
		return err
	}
	defer cursor.Close(context.Background())

	var clients []Client

	for cursor.Next(context.Background()) {
		var client Client
		if err := cursor.Decode(&client); err != nil {
			log.Fatal(err)
		}

		clients = append(clients, client)
	}

	return c.Render("index", clients)

}

func showCreateHandler(c *fiber.Ctx) error {
	return c.Render("create", nil)
}

func createHandler(c *fiber.Ctx) error {
	collection := returnCollectionPointer("clients")

	var newclient Client

	if err := c.BodyParser(&newclient); err != nil {
		log.Println(err)
		c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
		return err
	}

	res, err := collection.InsertOne(context.Background(), newclient)
	if err != nil {
		log.Println(err)
		c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
		return err
	} else {
		log.Println("Inserted document with ID:", res.InsertedID)
		c.Redirect("/", 301)
	}
	return nil
}

func EditPage(c *fiber.Ctx) error {
	collection := returnCollectionPointer("clients")

	var client Client

	id := c.Params("id")
	mongoId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
		return err
	}

	err = collection.FindOne(context.Background(), bson.M{"_id": mongoId}).Decode(&client)
	if err != nil {
		log.Println(err)
		c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
		return err
	}

	c.Render("edit", client)
	return nil
}

func EditHandler(c *fiber.Ctx) error {
	collection := returnCollectionPointer("clients")

	id := c.Params("id")

	mongoId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
		return err
	}

	var client Client
	if err := c.BodyParser(&client); err != nil {
		log.Println(err)
		c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
		return err
	}

	update := bson.D{
		{"$set", bson.D{
			{"name", client.Name},
			{"email", client.Email},
			{"age", client.Age},
		}},
	}

	res, err := collection.UpdateByID(context.Background(), mongoId, update)
	if err != nil {
		log.Println(err)
		c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
		return err
	} else {
		log.Println("Updated document count:", res.ModifiedCount)
		c.Redirect("/", 301)
	}

	return nil
}

func DeleteHandler(c *fiber.Ctx) error {
	collection := returnCollectionPointer("clients")

	id := c.Params("id")

	mongoId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
		return err
	}

	filter := bson.D{{"_id", mongoId}}

	deleted, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Println(err)
		c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
		return err
	} else {
		log.Println("Deleted document count:", deleted.DeletedCount)
		c.Redirect("/", 301)
	}

	return err
}
