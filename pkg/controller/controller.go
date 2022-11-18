package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/travas-io/travas/pkg/config"
	"github.com/travas-io/travas/query"
	"github.com/travas-io/travas/query/repo"

	"go.mongodb.org/mongo-driver/mongo"
)

type TravasHandler struct {
	App *config.TravasConfig
	DB  query.TravasDBRepo
}

func NewTravasHandler(app *config.TravasConfig, db *mongo.Client) *TravasHandler {
	return &TravasHandler{
		App: app,
		DB:  repo.NewTravasDB(app, db),
	}
}

// todo -> this is where all our handler/ controller logic will be done

func (tv *TravasHandler) HomePage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println("Travas HomePage")
		if err := ctx.Request.ParseForm(); err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
		}

	}
}
func (tv *TravasHandler) GetUsers() gin.HandlerFunc {
	// we will get the list of users from db_collection using the db_query
	return func(ctx *gin.Context) {
		fmt.Println("Users list")
	}
}
