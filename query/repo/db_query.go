package repo

import (
	"context"
	"time"

	"github.com/travas-io/travas/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// database queries is done in this file

func (td *TravasDB) InsertTour(tour model.Tour, tours []model.Tour) (int, primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var id primitive.ObjectID

	filter := bson.D{{Key: "tour_title", Value: tour.TourTitle}}
	var res bson.M
	err := TouristData(td.DB, "tours").FindOne(ctx, filter).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			tour.ID = primitive.NewObjectID()
			doc := bson.D{
				{Key: "_id", Value: tour.ID},
				{Key: "operator_id", Value: tour.OperatorID},
				{Key: "tour_title", Value: tour.TourTitle},
				{Key: "meeeting_point", Value: tour.MeetingPoint},
				{Key: "start_time", Value: tour.StartTime},
				{Key: "language_offered", Value: tour.LanguageOffered},
				{Key: "number_of_tourist", Value: tour.NumberOfTourist},
				{Key: "description", Value: tour.Description},
				{Key: "tour_guide", Value: tour.TourGuide},
				{Key: "tour_operator", Value: tour.TourOperator},
				{Key: "operator_contact", Value: tour.OperatorContact},
				{Key: "date", Value: tour.Date},
			}
			_, insertErr := TouristData(td.DB, "tours").InsertOne(ctx, doc)
			if insertErr != nil {
				td.App.ErrorLogger.Fatalf("cannot add user to the database : %v ", insertErr)
			}
			return 0, tour.ID, nil
		}
		td.App.ErrorLogger.Fatal(err)
	}

	for k, v := range res {
		if k == "_id" {
			switch tourID := v.(type) {
			case primitive.ObjectID:
				id = tourID
			}
		}
	}
	return 1, id, err
}

func (td *TravasDB) DeleteTour(tourID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var result bson.M

	filter := bson.D{{Key: "_id", Value: tourID}}
	err := TouristData(td.DB, "tours").FindOneAndDelete(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			td.App.ErrorLogger.Println("no document found for this query")
			return false, err
		}
		td.App.ErrorLogger.Fatalf("cannot execute the database query perfectly : %v ", err)
	}
	// td.App.InfoLogger.Printf("found document %v", result)
	return true, nil
}

func (td *TravasDB) UpdateTour(tourID string, tour model.Tour) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: tourID}}
	_, err := TouristData(td.DB, "tours").UpdateByID(ctx, filter, tour)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (td *TravasDB) InsertUser(user model.Tourist, tours []model.Tour) (int, primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var id primitive.ObjectID

	filter := bson.D{{Key: "email", Value: user.Email}}

	var res bson.M
	err := TouristData(td.DB, "tourist").FindOne(ctx, filter).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			user.ID = primitive.NewObjectID()
			doc := bson.D{
				{Key: "_id", Value: user.ID},
				{Key: "email", Value: user.Email},
				{Key: "first_name", Value: user.FirstName},
				{Key: "last_name", Value: user.LastName},
				{Key: "phone", Value: user.Phone},
				{Key: "password", Value: user.Password},
				{Key: "check_password", Value: user.CheckPassword},
				{Key: "tours_list", Value: tours},
				{Key: "created_at", Value: user.CreatedAt},
				{Key: "updated_at", Value: user.UpdatedAt},
				{Key: "geo_location", Value: ""},
				{Key: "token", Value: ""},
				{Key: "new_token", Value: ""},
			}
			_, insertErr := TouristData(td.DB, "tourist").InsertOne(ctx, doc)
			if insertErr != nil {
				td.App.ErrorLogger.Fatalf("cannot add user to the database : %v ", insertErr)
			}
			return 0, user.ID, nil
		}
		td.App.ErrorLogger.Fatal(err)
	}

	for k, v := range res {
		if k == "_id" {
			switch userID := v.(type) {
			case primitive.ObjectID:
				id = userID
			}
		}
	}
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
	// td.App.InfoLogger.Printf("found document %v", result)
	return true, nil
}

func (td *TravasDB) UpdateInfo(userID primitive.ObjectID, tk map[string]string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// opt := options.Update().SetUpsert(true)
	filter := bson.D{{Key: "_id", Value: userID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "token", Value: tk["t1"]}, {Key: "new_token", Value: tk["t2"]}}}}

	_, err := TouristData(td.DB, "tourist").UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}
	return true, nil
}
