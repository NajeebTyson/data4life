package repository

import (
	"context"
	"data4life/pkg/token"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	hostMongodb    = "localhost"
	portMongodb    = 27017
	dbnameMongodb  = "data4life"
	collection     = "tokens"
	tokenAttribute = "token"
)

type TokenStoreMongodb struct {
	conn *mongo.Client
	db   *mongo.Database
}

func (s *TokenStoreMongodb) Close() error {
	if s.conn != nil {
		if err := s.conn.Disconnect(context.TODO()); err != nil {
			return err
		}
	}
	return nil
}

func NewTokenStoreMongodb(dbName string) (*TokenStoreMongodb, error) {
	url := fmt.Sprintf("mongodb://%s:%d", hostMongodb, portMongodb)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}

	return &TokenStoreMongodb{
		conn: client,
		db:   client.Database(dbName),
	}, nil

}

func (s *TokenStoreMongodb) AddToken(t *token.Token) error {
	tt := bson.D{{tokenAttribute, t.Data}}
	if _, err := s.db.Collection(collection).InsertOne(context.TODO(), tt); err != nil {
		return err
	}
	return nil
}

func (s *TokenStoreMongodb) AddTokenBatch(tokens []string) error {
	data := make([]interface{}, len(tokens))
	for i, t := range tokens {
		data[i] = bson.D{{tokenAttribute, t}}
	}
	ordered := false
	if _, err := s.db.Collection(collection).InsertMany(context.TODO(), data, &options.InsertManyOptions{Ordered: &ordered}); err != nil {
		return err
	}
	return nil
}

func (s *TokenStoreMongodb) GetToken(t string) (*token.Token, error) {
	filter := bson.D{{tokenAttribute, t}}
	var result struct {
		Token string
	}

	err := s.db.Collection(collection).FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &token.Token{Data: result.Token}, nil
}

func (s *TokenStoreMongodb) DeleteToken(t string) error {
	filter := bson.D{{tokenAttribute, t}}
	_, err := s.db.Collection(collection).DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}
