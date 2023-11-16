package main

import (
	"context"
	"fmt"
	"github.com/startcodextech/goauth/users"
	"github.com/startcodextech/goevents/config"
	"github.com/startcodextech/goevents/system"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("users exitted abnormally: %s\n", err)
		os.Exit(1)
	}
}

func run() (err error) {
	ctx := context.Background()
	var cfg config.AppConfig
	cfg, err = config.InitConfig()
	if err != nil {
		return err
	}

	s, err := system.NewSystem(cfg)
	if err != nil {
		return err
	}
	defer func(db *mongo.Client, ctx context.Context) {
		if err = db.Disconnect(ctx); err != nil {
			return
		}
	}(s.MongoDB(), ctx)

	if users.Ro
}
