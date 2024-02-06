package printer

import (
	"fmt"
	"reflect"
	"sort"
	"time"

	"aarti/pkg/utils/printer/internal/slices"
)

type column struct {
	index    int
	name     string
	header   string
	order    int
	disabled bool
}

type columns []column

func (c columns) Len() int {
	return len(c)
}

func (c columns) Less(i, j int) bool {
	return c[i].order < c[j].order
}

func (c columns) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c columns) Exported() columns {
	return slices.Filter(c, func(v column) bool {
		return !v.disabled
	})
}

func (c columns) Sort() columns {
	sort.Sort(c)
	return c[:]
}

func (c columns) Filter(fields ...string) columns {
	if len(fields) == 0 {
		return c
	}
	return slices.Filter(c, func(v column) bool {
		return slices.Contains(fields, v.name)
	})
}

func (c columns) Headers() []string {
	return slices.Map(c.Exported().Sort(), func(v column) string {
		return v.header
	})
}

func (c columns) Values(v any, f map[string]func(v any) string, tf map[reflect.Type]func(v any) string) []string {
	val := derefValue(reflect.ValueOf(v))
	return slices.Map(c.Exported().Sort(), func(c column) string {
		for k, v := range tf {
			if k.AssignableTo(val.Field(c.index).Type()) {
				return v(val.Field(c.index).Interface())
			}
		}
		v := derefValue(val.Field(c.index))
		if !v.IsValid() {
			// TODO(adphi): add options to print nil as "-" or "nil"
			return ""
		}
		i := v.Interface()
		if fn, ok := f[c.name]; ok {
			return fn(i)
		}
		switch i := i.(type) {
		case string:
			return i
		// TODO(adphi): add options to print []byte as hex or base64 or string
		case []byte:
			return fmt.Sprintf("%v", i)
		// timestamppb.Timestamp
		case interface{ AsTime() time.Time }:
			return i.AsTime().String()
		// durationpb.Duration
		case interface{ AsDuration() time.Duration }:
			return i.AsDuration().String()
		case fmt.Stringer:
			return i.String()
		default:
			return fmt.Sprintf("%v", i)
		}
	})
}

func makeColumns(v any) (columns, error) {
	var columns columns
	t := derefType(reflect.TypeOf(v))
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("print")
		c, err := parseTag(i, field.Name, tag)
		if err != nil {
			return nil, err
		}
		c.index = i
		if !field.IsExported() {
			c.disabled = true
		}
		columns = append(columns, c)
	}
	return columns, nil
}
