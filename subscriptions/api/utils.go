package api

import (
	"net/url"
	"reflect"
	"strconv"

	"github.com/datapace/datapace/streams"
)

func parseInt(val string) (uint64, error) {
	v, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, err
	}

	return v, nil
}

func searchFields(req *searchSubsReq, q url.Values) error {
	val := reflect.ValueOf(req).Elem()
	reqType := reflect.TypeOf(*req)
	for i := 0; i < reqType.NumField(); i++ {
		structField := reqType.Field(i)
		field := val.FieldByName(structField.Name)
		name := structField.Tag.Get("alias")

		val := q[name]
		if len(val) > 1 {
			return streams.ErrMalformedData
		}
		if len(val) == 0 {
			continue
		}

		fi := field.Interface()
		switch fi.(type) {
		case string:
			field.SetString(val[0])
		case uint64:
			v, err := parseInt(val[0])
			if err != nil {
				return streams.ErrMalformedData
			}
			field.SetUint(v)
		}
	}
	// Silently prevent too large limit value.
	if req.Limit > maxLimit {
		req.Limit = maxLimit
	}

	return nil
}
