// Package cli is adapted from https://github.com/rancher/wrangler-cli
package cli

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unsafe"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"aarti/pkg/utils/signals"
)

var (
	caseRegexp = regexp.MustCompile("([a-z])([A-Z])")
)

type PersistentPreRunnable interface {
	PersistentPre(cmd *cobra.Command, args []string) error
}

type PreRunnable interface {
	Pre(cmd *cobra.Command, args []string) error
}

type Runnable interface {
	Run(cmd *cobra.Command, args []string) error
}

type customizer interface {
	Customize(cmd *cobra.Command)
}

type fieldInfo struct {
	FieldType  reflect.StructField
	FieldValue reflect.Value
}

func fields(obj interface{}) []fieldInfo {
	ptrValue := reflect.ValueOf(obj)
	objValue := ptrValue.Elem()

	var result []fieldInfo

	for i := 0; i < objValue.NumField(); i++ {
		fieldType := objValue.Type().Field(i)
		if !fieldType.IsExported() {
			continue
		}
		if fieldType.Anonymous && fieldType.Type.Kind() == reflect.Struct {
			result = append(result, fields(objValue.Field(i).Addr().Interface())...)
		} else if !fieldType.Anonymous {
			result = append(result, fieldInfo{
				FieldValue: objValue.Field(i),
				FieldType:  objValue.Type().Field(i),
			})
		}
	}

	return result
}

func Name(obj interface{}) string {
	ptrValue := reflect.ValueOf(obj)
	objValue := ptrValue.Elem()
	commandName := strings.Replace(objValue.Type().Name(), "Command", "", 1)
	commandName, _ = name(commandName, "", "")
	return commandName
}

func Main(cmd *cobra.Command) {
	ctx := signals.SetupSignalHandler()
	if err := cmd.ExecuteContext(ctx); err != nil {
		logrus.Fatal(err)
	}
}

func makeEnvVar[T comparable](to []func(), name string, vars []string, defValue T, flags *pflag.FlagSet, fn func(flag string) (T, error)) []func() {
	for _, v := range vars {
		to = append(to, func() {
			v := os.Getenv(v)
			if v == "" {
				return
			}
			fv, err := fn(name)
			if err == nil && fv == defValue {
				flags.Set(name, v)
			}
		})
	}
	return to
}

// Command populates a obj.Command() object by extracting args from struct tags of the
// Runnable obj passed.  Also the Run method is assigned to the RunE of the cli.
// name = Override the struct field with
func Command(obj Runnable, c *cobra.Command) *cobra.Command {
	var (
		envs     []func()
		arrays   = map[string]reflect.Value{}
		slices   = map[string]reflect.Value{}
		maps     = map[string]reflect.Value{}
		ptrValue = reflect.ValueOf(obj)
		objValue = ptrValue.Elem()
	)

	if c.Use == "" {
		c.Use = Name(obj)
	}

	for _, info := range fields(obj) {
		fieldType := info.FieldType
		v := info.FieldValue

		name, alias := name(fieldType.Name, fieldType.Tag.Get("name"), fieldType.Tag.Get("short"))
		usage := fieldType.Tag.Get("usage")
		envVars := strings.Split(fieldType.Tag.Get("env"), ",")
		defValue := fieldType.Tag.Get("default")
		if len(envVars) == 1 && envVars[0] == "" {
			envVars = nil
		}
		for _, v := range envVars {
			if v == "" {
				continue
			}
			usage += fmt.Sprintf("%s [$%s]", usage, v)
		}
		defInt, err := strconv.Atoi(defValue)
		if err != nil {
			defInt = 0
		}
		defValueLower := strings.ToLower(defValue)
		defBool := defValueLower == "true" || defValueLower == "1" || defValueLower == "yes" || defValueLower == "y"

		flags := c.PersistentFlags()
		switch fieldType.Type.Kind() {
		case reflect.Int:
			flags.IntVarP((*int)(unsafe.Pointer(v.Addr().Pointer())), name, alias, defInt, usage)
			envs = append(envs, makeEnvVar(envs, name, envVars, defInt, flags, flags.GetInt)...)
		case reflect.String:
			flags.StringVarP((*string)(unsafe.Pointer(v.Addr().Pointer())), name, alias, defValue, usage)
			envs = append(envs, makeEnvVar(envs, name, envVars, defValue, flags, flags.GetString)...)
		case reflect.Slice:
			// env is not supported for slices
			switch fieldType.Tag.Get("split") {
			case "false":
				arrays[name] = v
				flags.StringArrayP(name, alias, nil, usage)
			default:
				slices[name] = v
				flags.StringSliceP(name, alias, nil, usage)
			}
		case reflect.Map:
			maps[name] = v
			flags.StringSliceP(name, alias, nil, usage)
		case reflect.Bool:
			flags.BoolVarP((*bool)(unsafe.Pointer(v.Addr().Pointer())), name, alias, defBool, usage)
			envs = append(envs, makeEnvVar(envs, name, envVars, defBool, flags, flags.GetBool)...)
		default:
			panic("Unknown kind on field " + fieldType.Name + " on " + objValue.Type().Name())
		}

		if len(envVars) == 0 {
			continue
		}
	}

	if p, ok := obj.(PersistentPreRunnable); ok {
		c.PersistentPreRunE = p.PersistentPre
	}

	if p, ok := obj.(PreRunnable); ok {
		c.PreRunE = p.Pre
	}

	c.RunE = obj.Run
	c.PersistentPreRunE = bind(c.PersistentPreRunE, arrays, slices, maps, envs)
	c.PreRunE = bind(c.PreRunE, arrays, slices, maps, envs)
	c.RunE = bind(c.RunE, arrays, slices, maps, envs)

	cust, ok := obj.(customizer)
	if ok {
		cust.Customize(c)
	}

	return c
}

func assignMaps(app *cobra.Command, maps map[string]reflect.Value) error {
	for k, v := range maps {
		k = contextKey(k)
		s, err := app.Flags().GetStringSlice(k)
		if err != nil {
			return err
		}
		if s != nil {
			values := map[string]string{}
			for _, part := range s {
				parts := strings.SplitN(part, "=", 2)
				if len(parts) == 1 {
					values[parts[0]] = ""
				} else {
					values[parts[0]] = parts[1]
				}
			}
			v.Set(reflect.ValueOf(values))
		}
	}
	return nil
}

func assignSlices(app *cobra.Command, slices map[string]reflect.Value) error {
	for k, v := range slices {
		k = contextKey(k)
		s, err := app.Flags().GetStringSlice(k)
		if err != nil {
			return err
		}
		if s != nil {
			v.Set(reflect.ValueOf(s[:]))
		}
	}
	return nil
}

func assignArrays(app *cobra.Command, arrays map[string]reflect.Value) error {
	for k, v := range arrays {
		k = contextKey(k)
		s, err := app.Flags().GetStringArray(k)
		if err != nil {
			return err
		}
		if s != nil {
			v.Set(reflect.ValueOf(s[:]))
		}
	}
	return nil
}

func contextKey(name string) string {
	parts := strings.Split(name, ",")
	return parts[len(parts)-1]
}

func name(name, setName, short string) (string, string) {
	if setName != "" {
		return setName, short
	}
	parts := strings.Split(name, "_")
	i := len(parts) - 1
	name = caseRegexp.ReplaceAllString(parts[i], "$1-$2")
	name = strings.ToLower(name)
	result := append([]string{name}, parts[0:i]...)
	for i := 0; i < len(result); i++ {
		result[i] = strings.ToLower(result[i])
	}
	if short == "" && len(result) > 1 {
		short = result[1]
	}
	return result[0], short
}

func bind(next func(*cobra.Command, []string) error,
	arrays map[string]reflect.Value,
	slices map[string]reflect.Value,
	maps map[string]reflect.Value,
	envs []func()) func(*cobra.Command, []string) error {
	if next == nil {
		return nil
	}
	return func(cmd *cobra.Command, args []string) error {
		for _, envCallback := range envs {
			envCallback()
		}
		if err := assignArrays(cmd, arrays); err != nil {
			return err
		}
		if err := assignSlices(cmd, slices); err != nil {
			return err
		}
		if err := assignMaps(cmd, maps); err != nil {
			return err
		}

		if next != nil {
			return next(cmd, args)
		}

		return nil
	}
}
