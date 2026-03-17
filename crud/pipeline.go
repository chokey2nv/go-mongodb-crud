package crud

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Paginate(limit, skip int64) mongo.Pipeline {
	return mongo.Pipeline{
		bson.D{{Key: "$skip", Value: skip}},
		bson.D{{Key: "$limit", Value: limit}},
	}
}

// AggregatePageResult is the shape returned by our facet pipeline.
type AggregatePageResult[T any] struct {
	Data  []T `bson:"data"`
	Total []struct {
		Count int64 `bson:"count"`
	} `bson:"total"`
}

// ParseAggregateResult extracts data slice and total from the query result.
func ParseAggregateResult[T any](res []AggregatePageResult[T]) ([]T, int64) {
	if len(res) == 0 {
		return nil, 0
	}

	r := res[0]

	var total int64
	if len(r.Total) > 0 {
		total = r.Total[0].Count
	}

	return r.Data, total
}

// FacetDataTotal constructs a $facet stage with data and total count.
// Limit/skip are applied inside the facet so you can still get total.
func FacetDataTotal(limit, skip int64) bson.D {
	return bson.D{{Key: "$facet", Value: bson.D{
		{Key: "data", Value: bson.A{
			bson.D{{Key: "$sort", Value: bson.D{{Key: "createdAt", Value: -1}}}},
			bson.D{{Key: "$skip", Value: skip}},
			bson.D{{Key: "$limit", Value: limit}},
		}},
		{Key: "total", Value: bson.A{
			bson.D{{Key: "$count", Value: "count"}},
		}},
	}}}
}
