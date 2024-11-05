package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Name  string `bson:"name"`
	Email string `bson:"email"`
	Age   string `bson:"age"`
}

func main() {
	// Строка подключения MongoDB
	connectStr := "mongodb://localhost:27017"

	// Создаем контекс с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Создаем клиента MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectStr))
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		log.Fatal(err)
	}

	// Проверяем подключение
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")

	// Закрываем соединение c MongoDB
	defer client.Disconnect(ctx)

	// Получаем ссылку на коллекцию "users". Если ее нет, то создаем.
	userCollection := client.Database("mydb").Collection("users")

	// Создаем новый документ user
	newUser := User{Name: "John", Email: "test@example.com", Age: "38"}

	// Добавляем новый документ user в коллекцию
	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		log.Fatal(err)
	}

	// Выводим ID добавленного документа
	fmt.Println("Inserted document with ID:", result.InsertedID)

	// Указываем пустой фильтр для получения всех документов
	filter := bson.D{}

	// Выполняем поиск всех документов user в коллекции
	cursor, err := userCollection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	// Перебираем все документы и выводим их
	var users []User
	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	// Выводим полученные данные
	fmt.Println(users)

	// Определяем фильтр для поиска по имени
	filter = bson.D{{"name", "John"}}

	// Определим операцию для обновления документа
	update := bson.D{{"$set", bson.D{{"name", "Johnny"}}}}

	// Выполняем обновление документа
	res, err := userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}

	// Выводим количество обновленных документов
	fmt.Println("Updated document count:", res.ModifiedCount)

	// Определяем фильтр для удаления документа
	filter = bson.D{{"name", "Johnny"}}

	// Выполняем удаление документа
	deleted, err := userCollection.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	// Выводим количество удаленных документов
	fmt.Println("Deleted document count:", deleted.DeletedCount)

	// Определяем фильтр для поиска с использованием регулярного выражения
	filter = bson.D{{"name", bson.D{{"$regex", "^J"}}}}

	// Выполняем поиск всех документов user в коллекции
	cursor, err = userCollection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var usersJ []User
	// Перебираем все документы и выводим их
	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			log.Fatal(err)
		}
		usersJ = append(usersJ, user)
	}

	fmt.Println(usersJ)
}
