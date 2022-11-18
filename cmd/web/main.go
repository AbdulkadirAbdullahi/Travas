package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/travas-io/travas/cmd/web/routes"
	"github.com/travas-io/travas/pkg/config"
	"github.com/travas-io/travas/pkg/controller"
	"github.com/travas-io/travas/pkg/db"
	"go.mongodb.org/mongo-driver/mongo"
)

var app config.TravasConfig

var session *scs.SessionManager

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		app.ErrorLogger.Fatalf("cannot load up the env file : %v", err)
	}

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.HttpOnly = true
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteStrictMode
	session.Cookie.Secure = true

	ErrorLogger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	InfoLogger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	app.ErrorLogger = ErrorLogger
	app.InfoLogger = InfoLogger
	app.Session = session

	//hardcoding the credentials
	//port := os.Getenv("PORT")
	//uri := os.Getenv("TRAVAS_DB_URI")

	Mongo_URL := "mongodb+srv://travas:211920@Travas.com@cluster0.yowd966.mongodb.net/test"

	app.InfoLogger.Println("*---------- Connecting to the travas cloud database --------")

	client := db.DatabaseConnection(Mongo_URL)

	// close database connection
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {

		}
	}(client, context.TODO())

	app.InfoLogger.Println("*---------- Starting Travas Web Server -----------*")

	router := gin.Default()
	err = router.SetTrustedProxies([]string{"127.0.0.1"})

	if err != nil {
		app.ErrorLogger.Fatalf("untrusted proxy address : %v", err)
	}

	handler := controller.NewTravasHandler(&app, client)
	routes.Routes(router, *handler)

	app.InfoLogger.Println("*---------- Starting Travas Web Server -----------*")
	//err = router.Run(port)
	err = router.Run("localhost:8080")
	if err != nil {
		app.ErrorLogger.Fatalf("cannot start the server : %v", err)
	}
}
