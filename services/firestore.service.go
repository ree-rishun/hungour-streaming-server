package services

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func BuildApp(ctx context.Context) (*firestore.Client, error) {
	sa := option.WithoutAuthentication()
	conf := &firebase.Config{ProjectID: os.Getenv("PROJECT_ID")}
	app, err := firebase.NewApp(ctx, conf, sa)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	return client, err
}
