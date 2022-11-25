package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	//	"github.com/gorilla/mux"
	"github.com/travas-io/travas/model"
	"github.com/travas-io/travas/pkg/config"
	"github.com/travas-io/travas/pkg/hash"
	"github.com/travas-io/travas/pkg/token"
	"github.com/travas-io/travas/query"
	"github.com/travas-io/travas/query/repo"
	"go.mongodb.org/mongo-driver/mongo"
	//	"gopkg.in/mgo.v2/bson"
)

type Travas struct {
	App *config.Tools
	DB  query.TravasDBRepo
}

func NewTravas(app *config.Tools, db *mongo.Client) *Travas {
	return &Travas{
		App: app,
		DB:  repo.NewTravasDB(app, db),
	}
}

func (tr *Travas) Welcome() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Todo : render the home page of the application
		ctx.JSON(http.StatusOK, gin.H{})
	}
}

func (tr *Travas) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	}
}

func (tr *Travas) ProcessRegister() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user model.Tourist

		if err := ctx.Request.ParseForm(); err != nil {
			_ = ctx.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
		}
		user.Email = ctx.Request.Form.Get("email")
		user.FirstName = ctx.Request.Form.Get("first_name")
		user.LastName = ctx.Request.Form.Get("last_name")
		user.Phone = ctx.Request.Form.Get("phone")
		user.Password = ctx.Request.Form.Get("password")
		user.CheckPassword = ctx.Request.Form.Get("check_password")
		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		if user.Password != user.CheckPassword {
			_ = ctx.AbortWithError(http.StatusInternalServerError, errors.New("passwords did not match"))
		}

		user.Password, _ = hash.Encrypt(user.Password)
		user.CheckPassword, _ = hash.Encrypt(user.CheckPassword)

		if err := tr.App.Validator.Struct(&user); err != nil {
			if _, ok := err.(*validator.InvalidValidationError); !ok {
				_ = ctx.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
				log.Println(err)
				return
			}
		}
		tours := []model.Tour{}
		track, userID, err := tr.DB.InsertUser(user, tours)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("error while adding new user"))
			return
		}
		cookieData := sessions.Default(ctx)

		data := model.IntraData{
			ID:       userID,
			Email:    user.Email,
			Password: user.Password,
		}
		cookieData.Set("data", data)

		if err := cookieData.Save(); err != nil {
			log.Println("error from the session storage")
			_ = ctx.AbortWithError(http.StatusNotFound, gin.Error{Err: err})
			return
		}
		switch track {
		case 1:
			// add the user id to session
			// redirect to the home page of the application
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Existing Account, Go to the Login page",
			})
		case 0:
			//	after inserting new user to the database
			//  notify the user to verify their  details via mail
			//  OR
			//  Send notification message on the page for them to login
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Registered Successfully",
			})
		}
	}
}

func (tr *Travas) LoginPage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	}
}

func (tr *Travas) ProcessLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := ctx.Request.ParseForm(); err != nil {
			_ = ctx.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
		}
		email := ctx.Request.Form.Get("email")
		password := ctx.Request.Form.Get("password")

		cookieData := sessions.Default(ctx)
		data := cookieData.Get("data").(model.IntraData)

		verified, err := hash.Verify(password, data.Password)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, errors.New("cannot verify user input password"))
		}
		if verified {
			switch {
			case email == data.Email:
				_, checkErr := tr.DB.CheckForUser(data.ID)

				if checkErr != nil {
					_ = ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("unregistered user %v", checkErr))
				}
				// generate the jwt token
				t1, t2, err := token.Generate(data.Email, data.ID)
				if err != nil {
					_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("token no generated : %v ", err))
				}

				var tk map[string]string
				tk = map[string]string{"t1": t1, "t2": t2}

				// update the database adding the token to user database
				_, updateErr := tr.DB.UpdateInfo(data.ID, tk)
				if updateErr != nil {
					_ = ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("unregistered user %v", updateErr))
				}

				ctx.SetCookie("authorization", t1, 60*60*24*7, "/", "localhost", false, true)
				ctx.JSON(http.StatusOK, gin.H{"message": "Welcome to user homepage"})

				// ctx.Redirect(http.StatusSeeOther, "api/auth/user/home")
			}
		}
	}
}

func (tr *Travas) Main() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Welcome to user homepage"})
	}
}

func (tr *Travas) CreateTour() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var tour model.Tour
		if err := ctx.Request.ParseForm(); err != nil {
			_ = ctx.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
		}
		tour.OperatorID = ctx.Request.Form.Get("operator_id")
		tour.TourTitle = ctx.Request.Form.Get("tour_title")
		tour.MeetingPoint = ctx.Request.Form.Get("meeting_point")
		tour.StartTime = ctx.Request.Form.Get("start_time")
		tour.LanguageOffered = ctx.Request.Form.Get("language_offered")
		tour.NumberOfTourist = ctx.Request.Form.Get("number_of_tourist")
		tour.Description = ctx.Request.Form.Get("description")
		tour.TourGuide = ctx.Request.Form.Get("tour_guide")
		tour.TourOperator = ctx.Request.Form.Get("tour_operator")
		tour.OperatorContact = ctx.Request.Form.Get("operator_contact")
		tour.Date = ctx.Request.Form.Get("date")

		tours := []model.Tour{}
		_, tourID, err := tr.DB.InsertTour(tour, tours)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("error while adding new user"))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"CreatedTour_ID": tourID,
		})

	}
}

func (tr *Travas) DeleteTour() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		_, err := tr.DB.DeleteTour(id)
		if err != nil {
			ctx.JSON(406, gin.H{"message": "Tour could not be deleted", "error": err.Error()})
			ctx.Abort()
			return
		}
		ctx.JSON(200, gin.H{"message": "Tour deleted"})
	}

}

func (tr *Travas) UpdateTour() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		tour := model.Tour{}

		if ctx.BindJSON(&tour) != nil {
			ctx.JSON(406, gin.H{"message": "Invalid Parameters"})
			ctx.Abort()
			return
		}
		_, err := tr.DB.UpdateTour(id, tour)
		if err != nil {
			ctx.JSON(406, gin.H{"message": "tour count not be updated", "error": err.Error()})
			ctx.Abort()
			return
		}
		ctx.JSON(200, gin.H{"message": "tour updated"})

	}

}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
