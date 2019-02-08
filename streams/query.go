package streams

import (
	"reflect"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

const (
	id        = "id"
	plain     = "plain"
	min       = "min"
	max       = "max"
	kindTag   = "kind"
	dbTag     = "db"
	gte       = "$gte"
	lt        = "$lt"
	ne        = "$ne"
	not       = "$not"
	prefixNot = "-"
)

// Query struct wraps query parameters and provides
// "query builder" as a convenient way to generate DB query.
// ID fields are treated as EQUALS, not LIKE and they are of kind "id".
// String values are treated as LIKE and they are of kind "plain".
// It's possible to search string and ID values using logical NOT operator
// by adding character "-" as the value prefix (e.g. name=-myName).
// However, if the value itself starts with "-", just double negation and
// the first "-" will be ignored (i.e. name=--myName wil perform name=-myName search).
// Range fields ar treated as a RANGE, not EQUALS so they are of kind "min" or "max".
type Query struct {
	Page       uint64
	Limit      uint64
	Coords     [][]float64
	Partners   []string
	Owner      string  `kind:"id" db:"owner"`
	Name       string  `kind:"plain" db:"name"`
	StreamType string  `kind:"plain" db:"type"`
	MinPrice   *uint64 `kind:"min" db:"price"`
	MaxPrice   *uint64 `kind:"max" db:"price"`
}

func setRange(query *bson.M, kind, dbName string, value interface{}) {
	if val, ok := value.(*uint64); ok && val != nil {
		q := *query
		v := *val

		if entry, ok := q[dbName].(bson.M); ok {
			entry[kind] = v
			return
		}

		q[dbName] = bson.M{
			kind: v,
		}
	}
}

func setString(query *bson.M, dbName string, value interface{}, isID bool) {
	if v, ok := value.(string); ok && v != "" {
		q := *query
		if strings.HasPrefix(v, prefixNot) {
			v = v[1:]
			if !strings.HasPrefix(v, prefixNot) {
				if isID {
					q[dbName] = bson.M{
						ne: v,
					}
					return
				}
				q[dbName] = bson.M{
					not: bson.RegEx{v, ""},
				}
				return
			}
		}

		if isID {
			q[dbName] = v
			return
		}

		q[dbName] = bson.RegEx{v, ""}
	}
}

func genQuery(qType reflect.Type, qVal reflect.Value) bson.M {
	query := bson.M{}
	for i := 0; i < qType.NumField(); i++ {
		structField := qType.Field(i)
		field := qVal.FieldByName(structField.Name)
		kind := structField.Tag.Get(kindTag)
		dbName := structField.Tag.Get(dbTag)
		switch kind {
		case id:
			setString(&query, dbName, field.Interface(), true)
		case plain:
			setString(&query, dbName, field.Interface(), false)
		case min:
			setRange(&query, gte, dbName, field.Interface())
		case max:
			setRange(&query, lt, dbName, field.Interface())
		}
	}

	return query
}

// GenQuery extracts a database query from query parameters.
func GenQuery(q *Query) *bson.M {
	v := reflect.ValueOf(q).Elem()
	t := reflect.TypeOf(*q)
	query := genQuery(t, v)
	if q.Coords != nil {
		query["location"] = bson.M{
			"$within": bson.M{
				"$polygon": q.Coords,
			},
		}
	}
	query["owner"] = bson.M{
		"$in": q.Partners,
	}

	return &query
}
