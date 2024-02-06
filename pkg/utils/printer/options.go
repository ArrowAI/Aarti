package printer

import (
	"reflect"
	"text/tabwriter"
)

type Encoder func(v any) ([]byte, error)

type Option func(*printer)

type printer struct {
	format         Format
	max            int
	noHeaders      bool
	UpperHeaders   bool
	LowerHeaders   bool
	UpperValues    bool
	LowerValues    bool
	writer         *tabwriter.Writer
	json           Encoder
	yaml           Encoder
	formatters     map[string]func(v any) string
	typeFormatters map[reflect.Type]func(v any) string
	fields         []string
}

func WithJSON() Option {
	return func(p *printer) {
		p.format = JSON
	}
}

func WithYAML() Option {
	return func(p *printer) {
		p.format = YAML
	}
}

func WithFormat(format Format) Option {
	return func(p *printer) {
		p.format = format
	}
}

func WithMax(max int) Option {
	return func(p *printer) {
		if max > 0 {
			p.max = max
		}
	}
}

func WithWriter(writer *tabwriter.Writer) Option {
	return func(p *printer) {
		if writer != nil {
			p.writer = writer
		}
	}
}

func WithJSONMarshaler(fn Encoder) Option {
	return func(p *printer) {
		if fn != nil {
			p.json = fn
		}
	}
}

func WithYAMLMarshaler(fn Encoder) Option {
	return func(p *printer) {
		if fn != nil {
			p.yaml = fn
		}
	}
}

func WithNoHeaders() Option {
	return func(p *printer) {
		p.noHeaders = true
	}
}

func WithUpperHeaders() Option {
	return func(p *printer) {
		p.UpperHeaders = true
	}
}

func WithLowerHeaders() Option {
	return func(p *printer) {
		p.LowerHeaders = true
	}
}

func WithUpperValues() Option {
	return func(p *printer) {
		p.UpperValues = true
	}
}

func WithLowerValues() Option {
	return func(p *printer) {
		p.LowerValues = true
	}
}

func WithFormatter(fieldName string, fn func(v any) string) Option {
	return func(p *printer) {
		if fn != nil {
			p.formatters[fieldName] = fn
		}
	}
}

func WithTypeFormatter(t any, fn func(v any) string) Option {
	return func(p *printer) {
		if fn != nil {
			p.typeFormatters[reflect.TypeOf(t)] = fn
		}
	}
}

func WithFields(fields ...string) Option {
	return func(p *printer) {
		p.fields = fields
	}
}
