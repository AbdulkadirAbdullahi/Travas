package query

import (
	"github.com/travas-io/travas/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// todo -> all our queries method are to implement the interface

type TravasDBRepo interface {
	InsertUser(user model.Tourist, tours []model.Tour) (int, primitive.ObjectID, error)
	CheckForUser(userID primitive.ObjectID) (bool, error)
	UpdateInfo(userID primitive.ObjectID, tk map[string]string) (bool, error)
	InsertTour(tour model.Tour, tours []model.Tour) (int, primitive.ObjectID, error)
	DeleteTour(tourID string) (bool, error)
	UpdateTour(tourID string, tour model.Tour) (bool, error)
	GetTour(tourID string) (tour model.Tour, err error)
	FindAllTours() (tours []model.Tour, err error)
}
