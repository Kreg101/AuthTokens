package storage

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

// User represents fields in database
type User struct {
	Token  string    `json:"token"`
	Expire time.Time `json:"expire"`
}

// Mongo is an implementation of Repository interface
type Mongo struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewMongo creates new structure for db connection
func NewMongo(dbName, collectionName string, client *mongo.Client) (*Mongo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	collection := client.Database(dbName).Collection(collectionName)

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"expire": 1},
		Options: options.Index().SetExpireAfterSeconds(0),
	}
	_, err = collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return nil, err
	}

	return &Mongo{
		client:     client,
		collection: collection,
	}, nil
}

// InsertRefresh inserts refresh token in db
func (s *Mongo) InsertRefresh(refresh string) error {
	ctx := context.Background()
	session, err := s.startSession(ctx, refresh)
	if err != nil {
		return err
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		user := User{
			Token:  refresh,
			Expire: time.Now().Add(time.Hour * 72),
		}
		_, err = s.collection.InsertOne(context.Background(), user)
		if err != nil {
			return err
		}

		err = session.CommitTransaction(ctx)
		if err != nil {
			return err
		}

		return nil
	}

	if err != nil {
		return err
	}

	err = s.update(ctx, session, refresh, refresh)
	if err != nil {
		return err
	}

	return nil
}

// CheckRefresh refreshes token in db
func (s *Mongo) CheckRefresh(oldR, newR string) (bool, error) {
	ctx := context.Background()
	session, err := s.startSession(ctx, oldR)
	if err != nil {
		return false, err
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	err = s.update(ctx, session, oldR, newR)
	if err != nil {
		return false, nil
	}

	return true, nil
}

// startSession starts transaction
func (s *Mongo) startSession(ctx context.Context, refresh string) (mongo.Session, error) {
	session, err := s.client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	err = session.StartTransaction()
	if err != nil {
		return nil, err
	}
	defer session.AbortTransaction(ctx)

	filter := bson.M{"token": refresh}
	var user User
	err = s.collection.FindOne(context.Background(), filter).Decode(&user)
	return nil, err
}

// update updates token in db
func (s *Mongo) update(ctx context.Context, session mongo.Session, oldR, newR string) error {
	update := bson.M{"$set": bson.M{"token": newR, "expire": time.Now().Add(time.Hour * 72)}}
	_, err := s.collection.UpdateOne(context.Background(), bson.M{"token": oldR}, update)
	if err != nil {
		return err
	}

	err = session.CommitTransaction(ctx)
	if err != nil {
		return err
	}
	return nil
}
