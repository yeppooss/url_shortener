package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"url-shortener-api/internal/lib/hasher"
	"url-shortener-api/internal/lib/user"
)

type UrlShortener struct {
	Url   string
	Alias string
}

type Storage struct {
	db *mongo.Client
}

func New(connectionString string) (*Storage, error) {
	const op = "storage.mongodb.New"

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: client}, nil
}

func (s *Storage) Disconnect() error {
	const op = "storage.mongodb.New"

	err := s.db.Disconnect(context.TODO())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) SaveUrl(urlToSave, alias string) (int64, error) {
	const op = "storage.mongodb.SaveUrl"

	coll := s.db.Database("urls").Collection("urls")
	newUrl := UrlShortener{Url: urlToSave, Alias: alias}

	_, err := coll.InsertOne(context.TODO(), newUrl)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return 1, nil
}

func (s *Storage) GetUrl(alias string) (string, error) {
	const op = "storage.mongodb.GetUrl"
	coll := s.db.Database("urls").Collection("urls")
	filter := bson.D{{"alias", alias}}

	var result UrlShortener

	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return result.Url, nil
}

func (s *Storage) DeleteUrl(alias string) error {
	const op = "storage.mongodb.DeleteUrl"

	coll := s.db.Database("urls").Collection("urls")
	filter := bson.D{{"alias", alias}}

	_, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) CreateUser(email, password string) (string, error) {
	const op = "storage.mongodb.CreateUser"

	coll := s.db.Database("urls").Collection("users")
	passwordToSave := hasher.GetMD5Hash(password)
	newUrl := user.User{Email: email, Password: passwordToSave}

	res, err := coll.InsertOne(context.TODO(), newUrl)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s *Storage) GetUser(email string) (user.User, error) {
	const op = "storage.mongodb.GetUser"

	coll := s.db.Database("urls").Collection("users")
	filter := bson.D{{"email", email}}

	var userRes user.User

	err := coll.FindOne(context.TODO(), filter).Decode(&userRes)
	if err != nil {
		return user.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return userRes, nil
}
