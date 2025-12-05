package crud

import (
	"context"
	"reflect"

	"github.com/chokey2nv/go-mongodb-crud/helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Identifiable interface {
	GetId() string
	// SetId(string)
}
type BaseModel[T Identifiable] struct {
	root *RootModel[T]
}

func NewBaseModel[T Identifiable](db *mongo.Database, collection string) *BaseModel[T] {
	return &BaseModel[T]{
		root: NewRootModel[T](db, collection),
	}
}

func (m *BaseModel[T]) Insert(ctx context.Context, item T) (T, error) {
	b := bson.M{}
	helper.StructToBSON(item, &b)
	return m.root.Insert(ctx, b)
}
func (m *BaseModel[T]) Update(ctx context.Context, id string, item T) (T, error) {
	b := bson.M{}
	helper.StructToBSON(item, &b)
	var out T
	err := m.root.Update(ctx, bson.M{"id": id}, bson.M{"$set": b})
	if err != nil {
		return out, err
	}
	return m.root.FindOne(ctx, bson.M{"id": id})
}

func (m *BaseModel[T]) Get(ctx context.Context, where T) (T, error) {
	q := NewQuery()

	if where.GetId() != "" {
		q.Eq("id", where.GetId())
	} else {
		bsonData := bson.M{}
		helper.StructToBSON(where, &bsonData)
		q.Add(bsonData)
	}
	return m.root.FindOne(ctx, q.Build())
}

func (m *BaseModel[T]) Exists(ctx context.Context, item T) (bool, error) {
	bsonData := bson.M{}
	helper.StructToBSON(item, &bsonData)
	out, err := m.root.FindOne(ctx, bsonData)
	if err != nil {
		return false, err
	}
	return !reflect.ValueOf(out).IsZero(), nil
}

// count
func (m *BaseModel[T]) Count(ctx context.Context, item T) (int64, error) {
	bsonData := bson.M{}
	helper.StructToBSON(item, &bsonData)
	return m.root.Count(ctx, bsonData)
}
func (m *BaseModel[T]) Delete(ctx context.Context, id string) error {
	return m.root.HardDelete(ctx, id)
}
func (m *BaseModel[T]) Archive(ctx context.Context, id string) error {
	return m.root.SoftDelete(ctx, id)
}

type ListOptions[T Identifiable] struct {
	Limit          int64
	Skip           int64
	SortBy         string // "createdAt", "name", etc
	SortDesc       bool
	Ids            []string
	Search         string
	SearchIn       []string // fields to search across
	Filter         T
	CustomQuery    func(q *QueryBuilder) // custom filters
	CustomPipeline func(pipeline mongo.Pipeline) mongo.Pipeline
}

// list
func (m *BaseModel[T]) FindMany(
	ctx context.Context,
	opt *ListOptions[T],
) ([]T, int64, error) {
	q := NewQuery()
	if opt.Search != "" {
		q.AddSearch(opt.SearchIn, opt.Search)
	}

	if len(opt.Ids) > 0 {
		q.AddIDs("id", opt.Ids)
	}

	if !reflect.ValueOf(opt.Filter).IsZero() {
		bsonData := bson.M{}
		helper.StructToBSON(opt.Filter, &bsonData)
		q.Add(bsonData)
	}

	opt.CustomQuery(q)

	return m.root.List(ctx, q.Build(), opt.Limit, opt.Skip, bson.D{{Key: opt.SortBy, Value: -1}})
}
func (m *BaseModel[T]) List(
	ctx context.Context,
	opt *ListOptions[T],
) ([]T, int64, error) {
	q := NewQuery()
	if opt.Search != "" {
		q.AddSearch(opt.SearchIn, opt.Search)
	}

	if len(opt.Ids) > 0 {
		q.AddIDs("id", opt.Ids)
	}

	if !reflect.ValueOf(opt.Filter).IsZero() {
		b := bson.M{}
		if err := helper.StructToBSON(opt.Filter, &b); err == nil {
			q.Add(b)
		}
		q.Add(b)
	}

	if opt.CustomQuery != nil {
		opt.CustomQuery(q)
	}

	pipeline := mongo.Pipeline{}

	if match := q.Build(); len(match) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: match}})
	}

	if opt.CustomPipeline != nil {
		pipeline = opt.CustomPipeline(pipeline)
	}

	pipeline = append(pipeline, FacetDataTotal(opt.Limit, opt.Skip))
	pipeline = helper.ApplyArrayDateConv(pipeline, "data")

	var res []AggregatePageResult[T]
	if err := m.root.Aggregate(ctx, pipeline, &res); err != nil {
		return nil, 0, err
	}
	data, total := ParseAggregateResult(res)
	return data, total, nil
}
