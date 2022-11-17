package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/travas-io/travas/pkg/config"
	"github.com/travas-io/travas/pkg/controller"
	"github.com/travas-io/travas/pkg/db"

	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
)

var app config.TravasConfig

// var session
func main() {

	err := godotenv.Load(".env")
	if err != nil {
		app.ErrorLogger.Fatalf("cannot load up the env file : %v", err)
	}

	ErrorLogger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	InfoLogger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	app.ErrorLogger = ErrorLogger
	app.InfoLogger = InfoLogger

	port := os.Getenv("PORT")
	uri := os.Getenv("TRAVAS_DB_URI")

	app.InfoLogger.Println("*---------- Connecting to the travas cloud database --------")

	client := db.DatabaseConnection(uri)

	// close database connection
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {

		}
	}(client, context.TODO())

	app.InfoLogger.Println("*---------- Starting Travas Web Server -----------*")

	router := gin.Default()
	err = router.SetTrustedProxies([]string{"127.0.0.0"})

	if err != nil {
		app.ErrorLogger.Fatalf("untrusted proxy address : %v", err)
	}

	handler := controller.NewTravasHandler(&app, client)
	Routes(router, *handler)

	app.InfoLogger.Println("*---------- Starting Travas Web Server -----------*")
	err = router.Run(port)
	if err != nil {
		app.ErrorLogger.Fatalf("cannot start the server : %v", err)
	}
}