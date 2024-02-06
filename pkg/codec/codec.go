package codec

type Codec[T any] interface {
	Encode(v T) ([]byte, error)
	Decode(b []byte) (T, error)
	Name() string
}

type Funcs[T any] struct {
	Format     string
	EncodeFunc func(v T) ([]byte, error)
	DecodeFunc func(b []byte) (T, error)
}

func (c Funcs[T]) Encode(a T) ([]byte, error) {
	return c.EncodeFunc(a)
}

func (c Funcs[T]) Decode(b []byte) (T, error) {
	return c.DecodeFunc(b)
}

func (c Funcs[T]) Name() string {
	return c.Format
}
