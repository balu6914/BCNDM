package http

import (
	"monetasa/streams"
	"net/url"
	"reflect"
	"strconv"
)

func parseInt(val string) (uint64, error) {
	v, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func stringInSlice(a string, list []string) bool {
	for _, s := range list {
		if s == a {
			return true
		}
	}

	return false
}

func containsPoints(q url.Values) bool {
	for _, v := range locationPoints {
		if len(q[v[0]]) != 0 || len(q[v[1]]) != 0 {
			return true
		}
	}

	return false
}

func searchFields(req *searchStreamsReq, q url.Values) error {
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
		case uint64, *uint64:
			v, err := parseInt(val[0])
			if err != nil {
				return streams.ErrMalformedData
			}

			if field.Kind() == reflect.Uint64 {
				field.SetUint(v)
				break
			}

			field.Set(reflect.ValueOf(&v))
		}
	}
	// Silently prevent too big limit value.
	if req.Limit > maxLimit {
		req.Limit = maxLimit
	}

	return nil
}

func locationFields(req *searchStreamsReq, q url.Values) error {
	if !containsPoints(q) {
		return nil
	}

	for i, v := range locationPoints {
		// X and Y coordinates are q[v[0]] and q[v[1]] respectively.
		if len(q[v[0]]) != 1 || len(q[v[1]]) != 1 {
			return streams.ErrMalformedData
		}
		req.Coords = append(req.Coords, []float64{0, 0})
		for j := range v {
			var err error
			req.Coords[i][j], err = strconv.ParseFloat(q[v[j]][0], 64)
			if err != nil {
				return streams.ErrMalformedData
			}
		}
	}

	return nil
}
