package env

import (
	"fmt"
	"net"
	"net/netip"
	"os"
	"strconv"
	"strings"
	"time"
)

type Value interface {
	float32 | float64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		int | int8 | int16 | int32 | int64 |
		bool |
		string |
		time.Time | time.Duration |
		net.IP | net.IPNet | netip.Addr | netip.Prefix | netip.AddrPort
}

func Set[T Value](name string, v T) error {
	return os.Setenv(name, fmt.Sprintf("%v", v))
}

func SetSlice[T Value](name string, v []T) error {
	var s []string
	for _, v := range v {
		s = append(s, fmt.Sprintf("%v", v))
	}
	return os.Setenv(name, strings.Join(s, ","))
}

func Unset(name string) error {
	return os.Unsetenv(name)
}

func GetSlice[T Value](name string) []T {
	var v []T
	for _, s := range strings.Split(os.Getenv(name), ",") {
		var t T
		setValue(s, &t)
		v = append(v, t)
	}
	return v
}

func GetSliceDefault[T Value](name string, def []T) []T {
	v, ok := os.LookupEnv(name)
	if !ok {
		return def
	}
	vs := strings.Split(v, ",")
	var s []string
	for _, v := range vs {
		if v = strings.TrimSpace(v); v != "" {
			s = append(s, v)
		}
	}
	if len(s) == 0 {
		return def
	}
	var out []T
	for i, v := range s {
		var val T
		if i < len(def) {
			val = def[i]
		}
		setValue(v, any(&val))
		out = append(out, val)
	}
	return out
}

func Get[T Value](name string) T {
	var v T
	setValue(os.Getenv(name), any(&v))
	return v
}

func GetDefault[T Value](key string, defaultVal T) T {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}
	setValue(value, any(&defaultVal))
	return defaultVal
}

func setValue(s string, v any) {
	s = strings.TrimSpace(s)
	switch v.(type) {
	case *float32:
		if f, err := strconv.ParseFloat(s, 32); err == nil {
			*v.(*float32) = float32(f)
		}
	case *float64:
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			*v.(*float64) = f
		}
	case *uint:
		if u, err := strconv.ParseUint(s, 10, 64); err == nil {
			*v.(*uint) = uint(u)
		}
	case *uint8:
		if u, err := strconv.ParseUint(s, 10, 8); err == nil {
			*v.(*uint8) = uint8(u)
		}
	case *uint16:
		if u, err := strconv.ParseUint(s, 10, 16); err == nil {
			*v.(*uint16) = uint16(u)
		}
	case *uint32:
		if u, err := strconv.ParseUint(s, 10, 32); err == nil {
			*v.(*uint32) = uint32(u)
		}
	case *uint64:
		if u, err := strconv.ParseUint(s, 10, 64); err == nil {
			*v.(*uint64) = u
		}
	case *int:
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			*v.(*int) = int(i)
		}
	case *int8:
		if i, err := strconv.ParseInt(s, 10, 8); err == nil {
			*v.(*int8) = int8(i)
		}
	case *int16:
		if i, err := strconv.ParseInt(s, 10, 16); err == nil {
			*v.(*int16) = int16(i)
		}
	case *int32:
		if i, err := strconv.ParseInt(s, 10, 32); err == nil {
			*v.(*int32) = int32(i)
		}
	case *int64:
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			*v.(*int64) = i
		}
	case *bool:
		switch strings.ToLower(s) {
		case "true", "yes", "on", "1":
			*v.(*bool) = true
		case "false", "no", "off", "0":
			*v.(*bool) = false
		}
	case *string:
		*v.(*string) = s
	case *net.IP:
		if ip := net.ParseIP(s); ip != nil {
			*v.(*net.IP) = ip
		}
	case *net.IPNet:
		if ip, ipnet, err := net.ParseCIDR(s); err == nil {
			ipnet.IP = ip
			*v.(*net.IPNet) = *ipnet
		}
	case *netip.Addr:
		if addr, err := netip.ParseAddr(s); err == nil {
			*v.(*netip.Addr) = addr
		}
	case *netip.Prefix:
		if prefix, err := netip.ParsePrefix(s); err == nil {
			*v.(*netip.Prefix) = prefix
		}
	case *netip.AddrPort:
		if addrPort, err := netip.ParseAddrPort(s); err == nil {
			*v.(*netip.AddrPort) = addrPort
		}
	case *time.Time:
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			*v.(*time.Time) = t
		}
	case *time.Duration:
		if d, err := time.ParseDuration(s); err == nil {
			*v.(*time.Duration) = d
			break
		}
		if n, err := strconv.ParseInt(s, 10, 64); err == nil {
			*v.(*time.Duration) = time.Duration(n) * time.Millisecond
		}
	}
}
