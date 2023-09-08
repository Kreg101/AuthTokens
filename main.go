package main

import (
	"context"
	"fmt"
	"github.com/Kreg101/AuthJwt/internal/api"
	"github.com/Kreg101/AuthJwt/internal/storage"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func main() {

	config, err := LoadConfig(".")
	if err != nil {
		log.Fatalf("can't read config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.DBHost))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	fmt.Println(config.DBName, config.CollectionName)
	rep, err := storage.NewMongo(config.DBName, config.CollectionName, client)
	if err != nil {
		panic(err)
	}

	fmt.Println(config.ServerHost)
	s := api.NewServer(config.ServerHost, rep)
	err = s.Run()
	if err != nil {
		panic(err)
	}
}

type Config struct {
	AccessKey      string `mapstructure:"ACCESS_KEY"`
	RefreshKey     string `mapstructure:"REFRESH_KEY"`
	DBHost         string `mapstructure:"DB_HOST"`
	DBName         string `mapstructure:"DB_NAME"`
	CollectionName string `mapstructure:"COLLECTION_NAME"`
	ServerHost     string `mapstructure:"SERVER_HOST"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	api.AccessKey = config.AccessKey

	return
}
