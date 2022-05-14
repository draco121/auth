package database

import (
	"auth/custom_models"
	"auth/graph/model"
	"auth/startup"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client
}

func Connect() *DB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(startup.Config.Mongodburi))
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

func (d *DB) InsertToken(token *custom_models.Token) (bool, error) {
	coll := d.client.Database("auth").Collection("session")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := coll.InsertOne(ctx, token)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (d *DB) FindOneAndDeleteToken(token string) (bool, error) {
	coll := d.client.Database("auth").Collection("session")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"token": token}
	err := coll.FindOneAndDelete(ctx, filter).Decode(nil)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (d *DB) IsTokenExists(token string) (bool, error) {
	coll := d.client.Database("auth").Collection("session")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"token": token}
	var res *custom_models.Token
	err := coll.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
