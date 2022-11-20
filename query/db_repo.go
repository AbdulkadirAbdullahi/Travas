package query

import (
	"github.com/travas-io/travas/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// todo -> all our queries method are to implement the interface

type TravasDBRepo interface {
	InsertUser(user model.Tourist) (int, primitive.ObjectID, error)
}
