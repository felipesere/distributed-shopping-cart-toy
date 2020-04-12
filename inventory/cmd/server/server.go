package main

import (
	"context"
	"fmt"
	"github.com/felipesere/inventory/v0/pkg/products"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/http"
	"os"
	"time"
)

type server struct {
	client     *mongo.Client
	database   string
	collection string
}

func with(c *gin.Context) func(error) {
	return func(err error) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (s *server) connected(operation func(ctx context.Context) error) error {
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := operation(c)
	if err != nil {
		return fmt.Errorf("failure when running operation: %w", err)
	}

	return nil
}

func (s *server) WithCollection(f func(ctx context.Context, collection *mongo.Collection) error) error {
	return s.connected(func(ctx context.Context) error {
		return f(ctx, s.client.Database(s.database).Collection(s.collection))
	})
}

func (s *server) check(c *gin.Context) {
	fail := with(c)
	err := s.connected(func(ctx context.Context) error {
		return s.client.Ping(ctx, readpref.Primary())
	})
	if err != nil {
		fail(err)
		return
	}
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (s *server) addToInventory(c *gin.Context) {
	fail := with(c)
	var newAddition products.Product
	err := c.BindJSON(&newAddition)

	if err != nil {
		fail(err)
		return
	}
	newAddition.Added = time.Now()

	err = s.WithCollection(func(ctx context.Context, collection *mongo.Collection) error {
		_, err := collection.InsertOne(ctx, newAddition)
		return err
	})

	if err != nil {
		fail(err)
		return
	}
}

func (s *server) listInventory(c *gin.Context) {
	fail := with(c)

	category := c.Query("category")

	result := make([]products.Product, 0)
	err := s.WithCollection(func(ctx context.Context, collection *mongo.Collection) error {
		query := bson.M{}
		if category != "" {
			query = bson.M{"categories": bson.M{"$all": bson.A{category}}}
		}

		cursor, err := collection.Find(ctx, query)
		if err != nil {
			return fmt.Errorf("unable to read items from collection: %w", err)
		}

		for cursor.Next(ctx) {
			var p products.Product
			err := bson.Unmarshal(cursor.Current, &p)
			if err != nil {
				return fmt.Errorf("unable to read product from collection: %w", err)
			}
			result = append(result, p)
		}

		return nil
	})

	if err != nil {
		fail(err)
	}

	c.JSON(http.StatusOK, result)
}

func (s *server) clearInventory(c *gin.Context) {
	fail := with(c)

	err := s.WithCollection(func(ctx context.Context, collection *mongo.Collection) error {
		_, err := collection.DeleteMany(ctx, bson.D{})
		if err != nil {
			return fmt.Errorf("failure when deleting: %w", err)
		}
		return nil
	})

	if err != nil {
		fail(err)
		return
	}

	c.Status(http.StatusOK)
}

func main() {
	mongoDBcluster := os.Getenv("MONGO_DB_CLUSTER")
	auth := options.Credential{
		Username:    "root",
		Password:    "rootpassword",
		PasswordSet: true,
	}
	client, err := mongo.NewClient(options.Client().SetAuth(auth).ApplyURI(mongoDBcluster))
	if err != nil {
		panic(err.Error())
	}
	err = client.Connect(context.Background())
	if err != nil {
		panic(err.Error())
	}

	s := server{
		client:     client,
		database:   "inventory",
		collection: "available",
	}
	r := gin.Default()
	r.GET("/check", s.check)
	r.POST("/inventory/available", s.addToInventory)
	r.GET("/inventory/available", s.listInventory)
	r.DELETE("/inventory/all", s.clearInventory)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
