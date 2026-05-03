package crud

import (
	"time"

	"github.com/beevik/guid"
	"go.mongodb.org/mongo-driver/bson"
)

// primitive.NewDateTimeFromTime(time.Now())
func ApplyBeforeInsertHooks(data bson.M) {
	data["id"] = guid.NewString()
	data["createdAt"] = time.Now()
	data["updatedAt"] = time.Now()
}

func ApplyBeforeUpdateHooks(update bson.M) {
	if set, ok := update["$set"].(bson.M); ok {
		set["updatedAt"] = time.Now()
	} else {
		update["$set"] = bson.M{"updatedAt": time.Now()}
	}
}

func ApplyAfterReadHooks(target interface{}) {
	// optional: convert ObjectID → string etc
}

func ApplyBeforeReadHooks(target interface{}) {
	// optional: convert ObjectID → string etc
}
