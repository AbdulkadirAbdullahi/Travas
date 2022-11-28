package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/travas-io/travas/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// database queries is done in this file

func (td *TravasDB) InsertUser(user model.Tourist) (int, primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.D{{Key: "email", Value: user.Email}}

	var res bson.M
	err := TouristData(td.DB, "tourist").FindOne(ctx, filter).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			user.ID = primitive.NewObjectID()
			_, insertErr := TouristData(td.DB, "tourist").InsertOne(ctx, user)
			if insertErr != nil {
				td.App.ErrorLogger.Fatalf("cannot add user to the database : %v ", insertErr)
			}
			return 0, user.ID, nil
		}
		td.App.ErrorLogger.Fatal(err)
	}

	id := (res["_id"]).(primitive.ObjectID)

	return 1, id, err
}

func (td *TravasDB) CheckForUser(userID primitive.ObjectID) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var result bson.M

	filter := bson.D{{Key: "_id", Value: userID}}
	err := TouristData(td.DB, "tourist").FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			td.App.ErrorLogger.Println("no document found for this query")
			return false, err
		}
		td.App.ErrorLogger.Fatalf("cannot execute the database query perfectly : %v ", err)
	}

	return true, nil
}

func (td *TravasDB) UpdateInfo(userID primitive.ObjectID, tk map[string]string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: userID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "token", Value: tk["t1"]}, {Key: "new_token", Value: tk["t2"]}}}}

	_, err := TouristData(td.DB, "tourist").UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (td *TravasDB) LoadPackage(res []model.Operator) ([]model.Operator, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	cursor, err := OperatorData(td.DB, "operator").Find(ctx, bson.D{{}})
	if err != nil {
		return res, fmt.Errorf("cannot find document in the database %v ", err)
	}

	if err = cursor.All(ctx, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (td *TravasDB) AddTourPackage(userID primitive.ObjectID, tour model.Tour) (bool, error){
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: userID}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "tour_list", Value: tour}}}}

	_, err:= TouristData(td.DB, "tourist").UpdateOne(ctx, filter, update)
	if err != nil{
		return false, fmt.Errorf("cannot update document : %v ", err)
	}
	return  true, nil
}
