package printer

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func parseTag(index int, fieldName, tag string) (column, error) {
	parts := strings.Split(tag, ",")
	if len(parts) > 2 {
		return column{}, fmt.Errorf("invalid tag for field %q: %q, tag should be `print:\"(name string),(order int)\"` or `print:\"-\"", fieldName, tag)
	}
	c := column{name: fieldName, header: fieldName, order: index}
	for _, v := range parts {
		v = strings.TrimSpace(v)
		if v == "-" {
			c.disabled = true
			return c, nil
		}
		if n, err := strconv.Atoi(v); err == nil {
			c.order = n
			continue
		}
		if v != "" {
			c.header = v
		}
	}
	return c, nil
}

func arr(v any) []any {
	r := reflect.ValueOf(v)
	if r.Kind() == reflect.Ptr && r.Elem().Kind() == reflect.Slice {
		r = r.Elem()
	}
	if r.Kind() != reflect.Slice {
		return []any{v}
	}
	var out []any
	for i := 0; i < r.Len(); i++ {
		out = append(out, r.Index(i).Interface())
	}
	return out
}

func derefValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

func derefType(v reflect.Type) reflect.Type {
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}
