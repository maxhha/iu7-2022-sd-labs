package models

import (
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"github.com/hashicorp/go-multierror"
)

func MarshalDict(d map[string]string) graphql.Marshaler {
	if d == nil {
		return graphql.Null
	}

	return graphql.WriterFunc(func(writer io.Writer) {
		writer.Write([]byte("{"))
		first := true
		for k, v := range d {
			if first {
				first = false
			} else {
				writer.Write([]byte(","))
			}

			writer.Write([]byte(fmt.Sprintf(`"%s":"%s"`, k, v)))
		}
		writer.Write([]byte("}"))
	})
}

func UnmarshalDict(v interface{}) (map[string]string, error) {
	d, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%T is not map[string]interface{}", v)
	}

	res := make(map[string]string, len(d))
	var err error

	for k, kv := range d {
		switch kv := kv.(type) {
		case string:
			res[k] = kv
		default:
			err = multierror.Append(err, fmt.Errorf(`field '%s' of type %T is not string`, k, kv))
		}
	}

	return res, err
}
