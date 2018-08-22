package subscriptions

import (
	"reflect"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

const (
	id        = "id"
	kindTag   = "kind"
	dbTag     = "db"
	ne        = "$ne"
	prefixNot = "-"
)

// Query struct wraps query parameters and provides
// "query builder" as a convenient way to generate DB query.
// ID fields are treated as EQUALS, not LIKE and they are of kind "id".
// String values are treated as LIKE and they are of kind "plain".
// Range fields ar treated as a RANGE, not EQUALS so they are of kind "min" or "max".
type Query struct {
	Page        uint64
	Limit       uint64
	StreamID    string `kind:"id" db:"stream_id"`
	StreamOwner string `kind:"id" db:"stream_owner"`
	UserID      string `kind:"id" db:"user_id"`
}

func setString(query *bson.M, dbName string, value interface{}) {
	if v, ok := value.(string); ok && v != "" {
		q := *query

		if strings.HasPrefix(v, prefixNot) {
			v = v[1:]
			if !strings.HasPrefix(v, prefixNot) {
				q[dbName] = bson.M{
					ne: v,
				}
				return
			}
		}

		q[dbName] = v
		return
	}
}

// GenQuery extracts a database query
// from query parameters.
func GenQuery(q *Query) *bson.M {
	qVal := reflect.ValueOf(q).Elem()
	qType := reflect.TypeOf(*q)
	query := bson.M{}
	for i := 0; i < qType.NumField(); i++ {
		structField := qType.Field(i)
		field := qVal.FieldByName(structField.Name)
		kind := structField.Tag.Get(kindTag)
		dbName := structField.Tag.Get(dbTag)
		if kind == id {
			setString(&query, dbName, field.Interface())
		}
	}

	return &query
}
