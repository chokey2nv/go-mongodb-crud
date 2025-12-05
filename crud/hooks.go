package crud

import (
	"time"

	"github.com/beevik/guid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ApplyBeforeInsertHooks(data bson.M) {
	data["id"] = guid.NewString()
	data["createdAt"] = primitive.NewDateTimeFromTime(time.Now())
	data["updatedAt"] = primitive.NewDateTimeFromTime(time.Now())
}

func ApplyBeforeUpdateHooks(update bson.M) {
	if set, ok := update["$set"].(bson.M); ok {
		set["updatedAt"] = primitive.NewDateTimeFromTime(time.Now())
	} else {
		update["$set"] = bson.M{"updatedAt": primitive.NewDateTimeFromTime(time.Now())}
	}
}

func ApplyAfterReadHooks(target interface{}) {
	// optional: convert ObjectID → string etc
}

func ApplyBeforeReadHooks(target interface{}) {
	// optional: convert ObjectID → string etc
}
