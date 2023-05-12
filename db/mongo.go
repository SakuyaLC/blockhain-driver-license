package db

import (
	"blockchain/main/internal/model"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertBlock(block model.Block) {

	// Подключение к MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Определение коллекции
	collection := client.Database("blockchainDB").Collection("blocks")

	// Проверка наличия блока в коллекции
	filter := bson.M{"index": block.Index}
	existingBlock := &model.Block{}
	err = collection.FindOne(ctx, filter).Decode(existingBlock)

	// Если блок уже есть в коллекции, то не сохраняем его
	if err == nil {
		fmt.Println("Block already in the DB: ")
		fmt.Println(block)
		return
	} else {

		// Сохранение блока в коллекции
		_, err = collection.InsertOne(ctx, block)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Block inserted into the DB: ")
		fmt.Println(block)

	}

	client.Disconnect(ctx)
}

func UpdateBlockchain() []model.Block {

	// Подключение к MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Определение коллекции
	collection := client.Database("blockchainDB").Collection("blocks")

	// Получение всех документов из коллекции
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	// Закрытие курсора при выходе из функции
	defer cursor.Close(ctx)

	// Инициализация массива блоков
	var Blockchain []model.Block

	// Проход по всем документам в коллекции и добавление блоков в массив
	for cursor.Next(ctx) {
		var block model.Block
		err := cursor.Decode(&block)
		if err != nil {
			log.Fatal(err)
		}
		Blockchain = append(Blockchain, block)
	}

	client.Disconnect(ctx)

	return Blockchain
}

func CheckIfCollectionHasRecords(ctx context.Context, collectionName string) (bool, error) {
	// Подключение к MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return false, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return false, err
	}
	defer client.Disconnect(ctx)

	// Определение коллекции
	collection := client.Database("blockchainDB").Collection(collectionName)

	// Получаем количество записей в коллекции
	count, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return false, err
	}

	// Если записей нет, возвращаем false, иначе true
	return count > 0, nil
}
