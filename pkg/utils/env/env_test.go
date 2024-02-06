package env

import (
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

type TestCase[T Value] struct {
	key   string
	value string
	want  T
	def   T
}

func run[T Value](t *testing.T, tests []TestCase[T]) {
	t.Run("", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.key, func(t *testing.T) {
				var zero T
				if err := Unset(tt.key); err != nil {
					t.Fatal(err)
				}
				if got := Get[T](tt.key); fmt.Sprintf("%v", got) != fmt.Sprintf("%v", zero) {
					t.Errorf("Get() = %v, want %v", got, zero)
				}
				if got := GetDefault(tt.key, tt.def); fmt.Sprintf("%v", got) != fmt.Sprintf("%v", tt.def) {
					t.Errorf("%s: GetDefault() = %v, want %v", tt.value, got, tt.def)
				}
				if err := Set(tt.key, tt.value); err != nil {
					t.Fatalf("Setenv() = %v", err)
				}
				if fmt.Sprintf("%v", tt.want) != fmt.Sprintf("%v", tt.def) {
					if got := Get[T](tt.key); fmt.Sprintf("%v", got) != fmt.Sprintf("%v", tt.want) {
						t.Errorf("Get() = %v, want %v", got, tt.want)
					}
				}
				if got := GetDefault(tt.key, tt.def); fmt.Sprintf("%v", got) != fmt.Sprintf("%v", tt.want) {
					t.Errorf("%s: GetDefault() = %v, want %v", tt.value, got, tt.want)
				}
			})
		}
	})
}

type TestSliceCase[T Value] struct {
	key   string
	value string
	want  []T
	def   []T
}

func join[T any](v []T, sep string) string {
	var vals []string
	for _, v := range v {
		vals = append(vals, fmt.Sprintf("%v", v))
	}
	return strings.Join(vals, sep)
}

func runSlice[T Value](t *testing.T, tests []TestSliceCase[T]) {
	t.Run("", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.key, func(t *testing.T) {
				if err := Unset(tt.key); err != nil {
					t.Fatal(err)
				}
				if got := GetSliceDefault(tt.key, tt.def); join(got, ",") != join(tt.def, ",") {
					t.Errorf("%s: GetDefault() = %v, want %v", tt.value, got, tt.def)
				}
				if err := Set(tt.key, tt.value); err != nil {
					t.Fatalf("Setenv() = %v", err)
				}
				if got := GetSliceDefault(tt.key, tt.def); join(got, ",") != join(tt.want, ",") {
					t.Errorf("%s: GetDefault() = %v, want %v", tt.value, got, tt.want)
				}
				if err := SetSlice(tt.key, tt.want); err != nil {
					t.Fatalf("Setenv() = %v", err)
				}
				if got := GetSliceDefault(tt.key, tt.def); join(got, ",") != join(tt.want, ",") {
					t.Errorf("%s: GetDefault() = %v, want %v", tt.value, got, tt.want)
				}
			})
		}
	})
}

func TestEnvFloatSlice(t *testing.T) {
	tests := []TestSliceCase[float32]{
		{"TEST", "4.2", []float32{4.2}, []float32{1}},
		{"TEST", "0", []float32{0}, []float32{0}},
	}
	runSlice(t, tests)
}

func TestEnvBoolSlice(t *testing.T) {
	tests := []TestSliceCase[bool]{
		{"TEST", "true,false, on,    of", []bool{true, false, true, false}, []bool{false}},
		{"TEST", "false", []bool{false}, []bool{true}},
		{"TEST", "1", []bool{true}, []bool{false}},
		{"TEST", "0", []bool{false}, []bool{true}},
		{"TEST", "yes", []bool{true}, []bool{false}},
		{"TEST", "no", []bool{false}, []bool{true}},
		{"TEST", "", []bool{true, false}, []bool{true, false}},
		{"TEST", ",,", []bool{true, false}, []bool{true, false}},
		{"TEST", "false,,true,false", []bool{false, true, false}, []bool{true, false}},
		{"TEST", ",false,", []bool{false}, []bool{true, false}},
	}
	runSlice(t, tests)
}

func TestEnvFloat(t *testing.T) {
	tests := []TestCase[float32]{
		{"TEST", "4.2", 4.2, 1},
		{"TEST", "0", 0, 0},
	}
	run(t, tests)
}

func TestEnvInt(t *testing.T) {
	tests := []TestCase[int]{
		{"TEST", "42", 42, 1},
		{"TEST", "0", 0, 0},
	}
	run(t, tests)
}

func TestEnvBool(t *testing.T) {
	tests := []TestCase[bool]{
		{"TEST", "true", true, false},
		{"TEST", "false", false, true},
		{"TEST", "1", true, false},
		{"TEST", "0", false, true},
		{"TEST", "yes", true, false},
		{"TEST", "no", false, true},
		{"TEST", "", true, true},
	}
	run(t, tests)
}

func parseIPNet(v string) net.IPNet {
	ip, ipnet, err := net.ParseCIDR(v)
	if err != nil {
		panic(err)
	}
	ipnet.IP = ip
	return *ipnet
}

func TestEnvIP(t *testing.T) {
	tests := []TestCase[net.IPNet]{
		{"TEST", "192.168.0.0/24  ", parseIPNet("192.168.0.0/24"), parseIPNet("10.0.0.0/24")},
		{"TEST", "  192.168.0.10/24", parseIPNet("192.168.0.10/24"), parseIPNet("10.0.0.0/24")},
	}
	run(t, tests)
}

func TestEnvDuration(t *testing.T) {
	tests := []TestCase[time.Duration]{
		{"TEST", "1h", time.Hour, time.Minute},
		{"TEST", "3600", 3600 * time.Millisecond, time.Minute},
	}
	run(t, tests)
}

func TestEnvDate(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	tests := []TestCase[time.Time]{
		{"TEST", now.Format(time.RFC3339), now, time.Now().Add(-time.Hour)},
	}
	run(t, tests)
}
