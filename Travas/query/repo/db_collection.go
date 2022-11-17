package repo

import "go.mongodb.org/mongo-driver/mongo"

// todo -> this file is where the application will be connected to collection in database

func TouristData(db *mongo.Client, collection string) *mongo.Collection {
	var touristData = db.Database("travas").Collection("tourist")
	return touristData
}
