package database

import (
	"authentication/custom_models"
	"authentication/graph/model"
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client
}

func Connect() *DB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoUri := os.Getenv("MONGO_DB_URI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	if err != nil {
		panic(err)
	}
	db := DB{client: client}
	return &db
}

func (d *DB) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return d.client.Disconnect(ctx)
}

func (d *DB) Save(doc *model.UserInput) (bool, error) {
	coll := d.client.Database("auth").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := coll.InsertOne(ctx, doc)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (d *DB) FindOneByUsername(username string) (*custom_models.User, error) {
	coll := d.client.Database("auth").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var res *custom_models.User
	filter := bson.M{"username": username}
	err := coll.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (d *DB) FindOneByPhonenumber(phonenumber float64) (*custom_models.User, error) {
	coll := d.client.Database("auth").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var res *custom_models.User
	filter := bson.M{"phonenumber": phonenumber}
	err := coll.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (d *DB) FindOneByUserId(id string) (*custom_models.User, error) {
	coll := d.client.Database("auth").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var res *custom_models.User
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectId}
	err = coll.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (d *DB) FindOneAndUpdateUsername(id string, newusername string) (bool, error) {
	coll := d.client.Database("auth").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}
	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{"username": newusername}}
	var res *custom_models.Token
	err = coll.FindOneAndUpdate(ctx, filter, update).Decode(&res)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
