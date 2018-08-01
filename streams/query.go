package streams

import (
	"reflect"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

const (
	min = "min"
	max = "max"
)

// Query struct wraps query parameters and provides "query builder"
//  as a convenient way to generate DB query.
type Query struct {
	Name       string      `alias:"name"`
	StreamType string      `alias:"type"`
	Coords     [][]float64 `alias:"coords"`
	Page       uint64      `alias:"page"`
	Limit      uint64      `alias:"limit"`
	MinPrice   *uint64     `alias:"minPrice"`
	MaxPrice   *uint64     `alias:"maxPrice"`
}

func generateRange(query *bson.M, field string, val uint64) *bson.M {
	prefix := field[:3]
	q := *query
	if prefix != min && prefix != max {
		return query
	}
	name := strings.ToLower(field[3:])
	kind := "$gte"
	if prefix == max {
		kind = "$lt"
	}
	if entry, ok := q[name]; ok {
		entry := entry.(bson.M)
		entry[kind] = val
		return query
	}

	q[name] = bson.M{
		kind: val,
	}

	return query
}

// GenerateQuery extracts a database query from
// query  parameters.
func GenerateQuery(q *Query) *bson.M {
	val := reflect.ValueOf(q).Elem()
	reqType := reflect.TypeOf(*q)
	query := bson.M{}
	for i := 0; i < reqType.NumField(); i++ {
		structField := reqType.Field(i)
		field := val.FieldByName(structField.Name)
		name := structField.Tag.Get("alias")
		fi := field.Interface()
		switch fi.(type) {
		case string:
			// Create LIKE regex.
			v := field.String()
			if v != "" {
				query[name] = bson.RegEx{v, ""}
			}
		case *uint64:
			if fi.(*uint64) != nil {
				generateRange(&query, name, *fi.(*uint64))
			}
		}
	}

	if q.Coords != nil {
		query["location"] = bson.M{
			"$within": bson.M{
				"$polygon": q.Coords,
			},
		}
	}

	return &query
}
