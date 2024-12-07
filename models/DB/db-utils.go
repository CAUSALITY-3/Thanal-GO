package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MongoFindAll[T any](mongoCollecion *mongo.Collection, doc T) ([]T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var results []T

	cursor, err := mongoCollecion.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		if err := cursor.Decode(&doc); err != nil {
			log.Fatal(err)
			return nil, err
		}
		results = append(results, doc)
	}
	return results, nil
}

func MongoCreate[T any](mongoCollecion *mongo.Collection, doc T) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := mongoCollecion.InsertOne(ctx, &doc)
	if err != nil {
		return doc, err
	}
	insertedUser := mongoCollecion.FindOne(ctx, bson.M{"_id": result.InsertedID})
	if err := insertedUser.Decode(&doc); err != nil {
		return doc, err
	}
	return doc, nil
}

func MongoFindOne[T any](mongoCollecion *mongo.Collection, doc T, filter primitive.M) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := mongoCollecion.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return doc, err
	}
	return doc, nil
}

func MongoFindOneAndUpdate[T any](mongoCollecion *mongo.Collection, doc T, filter primitive.M, update primitive.M) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	err := mongoCollecion.FindOneAndUpdate(ctx, &filter, &update, opts).Decode(&doc)
	if err != nil {
		return doc, err

	}
	return doc, nil
}
