package helper

import (
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func ApplyDateConv(pipeline mongo.Pipeline) mongo.Pipeline {
	return append(pipeline, bson.D{
		{Key: "$addFields", Value: bson.D{
			{Key: "createdAt", Value: bson.D{
				{Key: "$dateToString", Value: bson.D{
					{Key: "format", Value: "%Y-%m-%dT%H:%M:%S.%LZ"},
					{Key: "date", Value: "$createdAt"},
				}},
			}},
		}}},
	)
}
func ApplyArrayDateConv(pipeline mongo.Pipeline, input string) mongo.Pipeline {
	return append(pipeline, bson.D{
		{Key: "$addFields", Value: bson.D{
			{Key: input, Value: bson.D{
				{Key: "$map", Value: bson.D{
					{Key: "input", Value: "$data"},
					{Key: "as", Value: "item"},
					{Key: "in", Value: bson.D{
						{Key: "$mergeObjects", Value: bson.A{
							"$$item",
							bson.D{
								{Key: "createdAt", Value: bson.D{
									{Key: "$dateToString", Value: bson.D{
										{Key: "format", Value: "%Y-%m-%dT%H:%M:%S.%LZ"},
										{Key: "date", Value: "$$item.createdAt"},
									}},
								}},
								{Key: "updatedAt", Value: bson.D{
									{Key: "$dateToString", Value: bson.D{
										{Key: "format", Value: "%Y-%m-%dT%H:%M:%S.%LZ"},
										{Key: "date", Value: "$$item.updatedAt"},
									}},
								}},
								{Key: "deletedAt", Value: bson.D{
									{Key: "$dateToString", Value: bson.D{
										{Key: "format", Value: "%Y-%m-%dT%H:%M:%S.%LZ"},
										{Key: "date", Value: "$$item.deletedAt"},
									}},
								}},
							},
						}},
					}},
				}},
			}},
		}},
	})
}
func ApplyNamedDateConvs(pipeline mongo.Pipeline, fields []string) mongo.Pipeline {
	for _, field := range fields {
		pipeline = ApplyNamedDateConv(pipeline, field)
	}
	return pipeline
}
func ApplyNamedDateConv(pipeline mongo.Pipeline, field string) mongo.Pipeline {
	return append(pipeline, bson.D{
		{Key: "$addFields", Value: bson.D{
			{Key: field, Value: bson.D{
				{Key: "$dateToString", Value: bson.D{
					{Key: "format", Value: "%Y-%m-%dT%H:%M:%S.%LZ"},
					{Key: "date", Value: fmt.Sprintf("$%s", field)},
				}},
			}},
		}},
	})
}

func CleanUpBSON(bsonResult *primitive.M, stepKey string) {
	for key, value := range *bsonResult {
		if value == nil || isZero(value) {
			delete(*bsonResult, key)
			continue
		}
		switch v := value.(type) {
		case primitive.A:
			if len(v) == 0 {
				delete(*bsonResult, key)
			} else if !allObjects(v) {
				(*bsonResult)[key] = v
			} else {
				for key1, value1 := range v {
					if m, ok := value1.(primitive.M); ok {
						nextStepKey := fmt.Sprintf("%s.%d.", key, key1)
						CleanUpBSON(&m, nextStepKey)
					}
				}
			}
		case primitive.M:
			delete(*bsonResult, key)
			for key1, value1 := range v {
				nextStepKey := key + "." + key1
				if stepKey != "" {
					nextStepKey = stepKey + "." + nextStepKey
				}
				if m, ok := value1.(primitive.M); ok {
					CleanUpBSON(&m, nextStepKey)
				}
				(*bsonResult)[nextStepKey] = value1
			}
		}
	}
}

func StructToBSON(structObject interface{}, bsonResult *bson.M) error {
	byteValue, err := bson.Marshal(structObject)
	if err != nil {
		return err
	}
	err = bson.Unmarshal(byteValue, bsonResult)
	if err != nil {
		return err
	}
	CleanUpBSON(bsonResult, "")
	return nil
}
func isZero(value interface{}) bool {
	return reflect.DeepEqual(value, reflect.Zero(reflect.TypeOf(value)).Interface())
}
func allObjects(arr primitive.A) bool {
	for _, v := range arr {
		if _, ok := v.(primitive.M); !ok {
			return false
		}
	}
	return true
}
