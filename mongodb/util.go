package mongodb

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getConnection(conn string) (*mongo.Client, error) {
	conOpts := options.Client().ApplyURI(conn)
	if err := conOpts.Validate(); err != nil {
		os.Exit(1)
	}
	cl, err := mongo.Connect(context.TODO(), conOpts)
	if err != nil {
		return nil, err
	}
	return cl, nil

}
