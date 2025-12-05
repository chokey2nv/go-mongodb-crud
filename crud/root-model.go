package crud

import (
	"context"
	"time"

	"github.com/chokey2nv/go-mongodb-crud/helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RootModel[T any] struct {
	Collection *mongo.Collection
}

func NewRootModel[T any](db *mongo.Database, col string) *RootModel[T] {
	return &RootModel[T]{Collection: db.Collection(col)}
}

func (r *RootModel[T]) NewID() primitive.ObjectID {
	return primitive.NewObjectID()
}

// count
func (r *RootModel[T]) Count(ctx context.Context, filter bson.M) (int64, error) {
	return r.Collection.CountDocuments(ctx, filter)
}
func (r *RootModel[T]) Insert(ctx context.Context, data bson.M) (T, error) {
	ApplyBeforeInsertHooks(data)

	insertRes, err := r.Collection.InsertOne(ctx, data)

	out, err := r.FindOne(ctx, bson.M{"_id": insertRes.InsertedID})
	return out, err
}

func (r *RootModel[T]) Update(ctx context.Context, filter, update bson.M) error {
	ApplyBeforeUpdateHooks(update)

	_, err := r.Collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *RootModel[T]) FindOne(ctx context.Context, filter bson.M) (T, error) {
	var outs []T
	var out T
	pipeline := helper.ApplyNamedDateConvs(
		mongo.Pipeline{},
		[]string{"createdAt", "updatedAt"},
	)

	pipeline = append(pipeline, bson.D{
		{Key: "$match", Value: filter},
	})

	err := r.Aggregate(ctx, pipeline, &outs)

	if err != nil {
		return out, err
	}
	if len(outs) > 0 {
		out = outs[0]
	}
	ApplyAfterReadHooks(out)
	return out, nil
}

func (r *RootModel[T]) List(
	ctx context.Context,
	filter bson.M,
	skip, limit int64,
	sort bson.D,
) ([]T, int64, error) {

	opts := options.Find().SetSkip(skip).SetLimit(limit).SetSort(sort)
	cur, err := r.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	var out []T
	for cur.Next(ctx) {
		var v T
		if err := cur.Decode(&v); err != nil {
			return nil, 0, err
		}
		ApplyAfterReadHooks(&v)
		out = append(out, v)
	}

	count, _ := r.Collection.CountDocuments(ctx, filter)
	return out, count, nil
}

func (r *RootModel[T]) Aggregate(
	ctx context.Context,
	pipeline mongo.Pipeline,
	out interface{},
) error {
	cur, err := r.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	return cur.All(ctx, out)
}

// Soft delete: keeps doc but marks deleted
func (r *RootModel[T]) SoftDelete(ctx context.Context, id string) error {
	_, err := r.Collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{
		"$set": bson.M{"isDeleted": true, "deletedAt": primitive.NewDateTimeFromTime(time.Now())},
	})
	return err
}

func (r *RootModel[T]) HardDelete(ctx context.Context, id string) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

// ------------------------------------------------------------------
// TRANSACTION SUPPORT
// ------------------------------------------------------------------

// func (r *RootModel[T]) WithTransaction(
// 	ctx context.Context,
// 	fn func(sessCtx mongo.SessionContext) error,
// ) error {
// 	client := r.db.Client()
// 	session, err := client.StartSession()
// 	if err != nil {
// 		return err
// 	}
// 	defer session.EndSession(ctx)

// 	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
// 		return nil, fn(sessCtx)
// 	})
// 	return err
// }
