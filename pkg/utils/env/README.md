# Env

[![Language: Go](https://img.shields.io/badge/lang-Go-6ad7e5.svg?style=flat-square&logo=go)](https://golang.org/)
[![Go Reference](https://pkg.go.dev/badge/go.linka.cloud/env.svg)](https://pkg.go.dev/go.linka.cloud/env)

*A generic typed implementation of environment variables.*

```go
package env // import "go.linka.cloud/env"


// FUNCTIONS

func Get[T Value](name string) T
func GetDefault[T Value](key string, defaultVal T) T
func GetSlice[T Value](name string) []T
func GetSliceDefault[T Value](name string, def []T) []T
func Set[T Value](name string, v T) error
func SetSlice[T Value](name string, v []T) error
func Unset(name string) error

// TYPES

type Value interface {
float32 | float64 |
uint | uint8 | uint16 | uint32 | uint64 |
int | int8 | int16 | int32 | int64 |
bool |
string |
time.Time | time.Duration |
net.IP | net.IPNet | netip.Addr | netip.Prefix | netip.AddrPort
}
```
