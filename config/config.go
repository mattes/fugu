package config

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrInvalidStringSliceValue = fmt.Errorf("Invalid StringSlice")
)

type Value interface {
	// yaml.Setter
	Names() []string
	Set(value interface{}) error

	// Args returns a slice of strings to be used in exec args.
	Arg() string
}

type StringSliceValue struct {
	Name    []string
	Value   []string
	Defined bool
}

func (v *StringSliceValue) Arg() (out string) {
	if v.Defined {
		args := make([]string, 0)
		for _, v2 := range v.Value {
			args = append(args, fmt.Sprintf(`--%s="%s"`, v.Name[0], v2))
		}
		out = strings.Join(args, " ")
	}
	return
}

func (v *StringSliceValue) Set(value interface{}) error {
	v.Defined = true
	switch value.(type) {
	case []string:
		for _, v2 := range value.([]string) {
			v.Value = append(v.Value, v2)
		}
	case []interface{}:
		for _, v2 := range value.([]interface{}) {
			v.Value = append(v.Value, fmt.Sprintf("%v", v2))
		}
	default:
		return ErrInvalidStringSliceValue
	}
	return nil
}

func (v *StringSliceValue) Names() []string {
	return v.Name
}

type BoolValue struct {
	Name    []string
	Value   bool
	Defined bool
}

func (v *BoolValue) Arg() (out string) {
	if v.Defined {
		out = fmt.Sprintf(`--%s=%v`, v.Name[0], v.Value)
	}
	return
}

func (v *BoolValue) Set(value interface{}) error {
	v.Defined = true
	v2, err := strconv.ParseBool(value.(string))
	if err != nil {
		return err
	}
	v.Value = v2
	return nil
}

func (v *BoolValue) Names() []string {
	return v.Name
}

type Int64Value struct {
	Name    []string
	Value   int64
	Defined bool
}

func (v *Int64Value) Arg() (out string) {
	if v.Defined {
		out = fmt.Sprintf(`--%s=%v`, v.Name[0], v.Value)
	}
	return
}

func (v *Int64Value) Set(value interface{}) error {
	v.Defined = true
	v2, err := strconv.ParseInt(value.(string), 10, 0)
	if err != nil {
		return err
	}
	v.Value = v2
	return nil
}

func (v *Int64Value) Names() []string {
	return v.Name
}

type StringValue struct {
	Name    []string
	Value   string
	Defined bool
}

func (v *StringValue) Arg() (out string) {
	if v.Defined {
		out = fmt.Sprintf(`--%s="%v"`, v.Name[0], v.Value)
	}
	return
}

func (v *StringValue) Set(value interface{}) error {
	v.Defined = true
	v.Value = fmt.Sprintf("%v", value)
	return nil
}

func (v *StringValue) Names() []string {
	return v.Name
}
