package crud

import "go.mongodb.org/mongo-driver/bson"

type QueryBuilder struct {
	ands []bson.M
	ors  []bson.M
}

func NewQuery() *QueryBuilder {
	return &QueryBuilder{
		ands: make([]bson.M, 0),
		ors:  make([]bson.M, 0),
	}
}

// AND
func (q *QueryBuilder) And(cond bson.M) {
	if cond == nil || len(cond) == 0 {
		return
	}
	q.ands = append(q.ands, cond)
}

// OR
func (q *QueryBuilder) Or(cond bson.M) {
	if cond == nil || len(cond) == 0 {
		return
	}
	q.ors = append(q.ors, cond)
}

func (q *QueryBuilder) Build() bson.M {
	ands := q.ands
	ors := q.ors

	// No conditions
	if len(ands) == 0 && len(ors) == 0 {
		return bson.M{}
	}

	// Only ANDs
	if len(ors) == 0 {
		if len(ands) == 1 {
			return ands[0]
		}
		return bson.M{"$and": ands}
	}

	// Only ORs
	if len(ands) == 0 {
		if len(ors) == 1 {
			return ors[0]
		}
		return bson.M{"$or": ors}
	}

	// Mixed
	filter := bson.M{}
	filter["$and"] = ands
	filter["$or"] = ors
	return filter
}

// internal helper
func (q *QueryBuilder) addAnd(cond bson.M) *QueryBuilder {
	if len(cond) > 0 {
		q.ands = append(q.ands, cond)
	}
	return q
}

func (q *QueryBuilder) addOr(cond bson.M) *QueryBuilder {
	if len(cond) > 0 {
		q.ors = append(q.ors, cond)
	}
	return q
}

// Equity
func (q *QueryBuilder) Eq(field string, v any) *QueryBuilder {
	return q.addAnd(bson.M{field: v})
}
func (q *QueryBuilder) Ne(field string, v any) *QueryBuilder {
	return q.addAnd(bson.M{field: bson.M{"$ne": v}})
}

// Comparison
func (q *QueryBuilder) Gt(field string, v any) *QueryBuilder {
	return q.addAnd(bson.M{field: bson.M{"$gt": v}})
}
func (q *QueryBuilder) Gte(field string, v any) *QueryBuilder {
	return q.addAnd(bson.M{field: bson.M{"$gte": v}})
}

func (q *QueryBuilder) Lt(field string, v any) *QueryBuilder {
	return q.addAnd(bson.M{field: bson.M{"$lt": v}})
}
func (q *QueryBuilder) Lte(field string, v any) *QueryBuilder {
	return q.addAnd(bson.M{field: bson.M{"$lte": v}})
}

// Array
func (q *QueryBuilder) In(field string, values []any) *QueryBuilder {
	return q.addAnd(bson.M{field: bson.M{"$in": values}})
}
func (q *QueryBuilder) Nin(field string, values []any) *QueryBuilder {
	return q.addAnd(bson.M{field: bson.M{"$nin": values}})
}

// exists
func (q *QueryBuilder) Exists(field string, exists bool) *QueryBuilder {
	return q.addAnd(bson.M{field: bson.M{"$exists": exists}})
}

// Regex
func (q *QueryBuilder) Regex(field, pattern string) *QueryBuilder {
	return q.addAnd(bson.M{
		field: bson.M{"$regex": pattern, "$options": "i"},
	})
}

// not
func (q *QueryBuilder) Not(field string, v any) *QueryBuilder {
	return q.addAnd(bson.M{
		field: bson.M{"$not": bson.M{"$eq": v}},
	})
}

func (q *QueryBuilder) AddSearch(fields []string, keyword string) *QueryBuilder {
	if keyword == "" || len(fields) == 0 {
		return q
	}
	or := make([]bson.M, 0, len(fields))
	for _, f := range fields {
		or = append(or, bson.M{f: bson.M{"$regex": keyword, "$options": "i"}})
	}
	return q.addOr(bson.M{"$or": or})
}

func (q *QueryBuilder) AddIDs(field string, ids []string) *QueryBuilder {
	if len(ids) == 0 {
		return q
	}
	return q.addAnd(bson.M{field: bson.M{"$in": ids}})
}

func (q *QueryBuilder) Add(cond bson.M) *QueryBuilder {
	return q.addAnd(cond)
}

// Or methods
func (q *QueryBuilder) OrEq(field string, v any) *QueryBuilder {
	return q.addOr(bson.M{field: v})
}

func (q *QueryBuilder) OrNe(field string, v any) *QueryBuilder {
	return q.addOr(bson.M{field: bson.M{"$ne": v}})
}

func (q *QueryBuilder) OrGt(field string, v any) *QueryBuilder {
	return q.addOr(bson.M{field: bson.M{"$gt": v}})
}

func (q *QueryBuilder) OrRegex(field, pattern string) *QueryBuilder {
	return q.addOr(bson.M{
		field: bson.M{"$regex": pattern, "$options": "i"},
	})
}
