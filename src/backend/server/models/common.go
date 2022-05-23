package models

import (
	"database/sql"
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/hashicorp/go-multierror"
	"github.com/shopspring/decimal"
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

const dateTimeLayout = `2006-01-02T15:04:05.000Z07:00`

func MarshalTime(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(`"`))
		w.Write([]byte(t.UTC().Format(dateTimeLayout)))
		w.Write([]byte(`"`))
	})
}

func UnmarshalTime(v interface{}) (time.Time, error) {
	str, ok := v.(string)

	if !ok {
		return time.Time{}, fmt.Errorf("convert to string")
	}

	t, err := time.Parse(dateTimeLayout, str)

	if err != nil {
		return time.Time{}, fmt.Errorf("time parse: %w", err)
	}

	return t, nil
}

func MarshalDecimal(d decimal.Decimal) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(d.String()))
	})
}

func UnmarshalDecimal(v interface{}) (decimal.Decimal, error) {
	valf, okf := v.(float64)
	if okf {
		return decimal.NewFromFloatWithExponent(valf, -2), nil
	}

	vali, oki := v.(int64)
	if oki {
		return decimal.NewFromInt(vali), nil
	}

	return decimal.Decimal{}, fmt.Errorf("fail convert to float or int %#v", v)
}

func MarshalNullDecimal(d decimal.NullDecimal) graphql.Marshaler {
	if d.Valid {
		return MarshalDecimal(d.Decimal)
	}
	return graphql.Null
}

func UnmarshalNullDecimal(v interface{}) (decimal.NullDecimal, error) {
	if v == nil {
		return decimal.NullDecimal{}, nil
	}

	d, err := UnmarshalDecimal(v)

	if err != nil {
		return decimal.NullDecimal{}, err
	}

	return decimal.NullDecimal{
		Decimal: d,
		Valid:   true,
	}, nil
}

func MarshalNullTime(t sql.NullTime) graphql.Marshaler {
	if t.Valid {
		return MarshalTime(t.Time)
	} else {
		return graphql.Null
	}
}

func UnmarshalNullTime(v interface{}) (sql.NullTime, error) {
	if v == nil {
		return sql.NullTime{}, nil
	}

	time, err := UnmarshalTime(v)

	if err != nil {
		return sql.NullTime{}, err
	}

	return sql.NullTime{
		Time:  time,
		Valid: true,
	}, nil
}
