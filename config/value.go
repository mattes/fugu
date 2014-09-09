package config

import (
	"fmt"
	"strconv"
)

var (
	ErrInvalidStringSliceValue = fmt.Errorf("Invalid StringSlice")
)

type Value interface {
	Names() []string
	Set(value interface{}) error
	Get() interface{}

	// Args returns a slice of strings to be used in exec args.
	Arg() []string
}

type StringSliceValue struct {
	Name    []string
	Value   []string
	Defined bool
}

func (v *StringSliceValue) Arg() (out []string) {
	if v.Defined {
		out = make([]string, 0)
		for _, v2 := range v.Value {
			out = append(out, fmt.Sprintf(`--%s="%s"`, v.Name[0], v2))
		}
	}
	return
}

func (v *StringSliceValue) Set(value interface{}) error {
	// replace existing value if already set
	if v.Defined == true {
		v.Value = make([]string, 0)
	}

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

func (v *StringSliceValue) Get() interface{} {
	return v.Value
}

func (v *StringSliceValue) Names() []string {
	return v.Name
}

type BoolValue struct {
	Name    []string
	Value   bool
	Defined bool
}

func (v *BoolValue) Arg() (out []string) {
	if v.Defined {
		if v.Value == false {
			out = []string{fmt.Sprintf(`--%s=%v`, v.Name[0], v.Value)}
		} else {
			out = []string{fmt.Sprintf(`--%s`, v.Name[0])}
		}
	}
	return
}

func (v *BoolValue) Set(value interface{}) error {
	v.Defined = true
	switch value.(type) {
	case string:
		v2, err := strconv.ParseBool(value.(string))
		if err != nil {
			return err
		}
		v.Value = v2
	case bool:
		v.Value = value.(bool)
	default:
		return ErrInvalidStringSliceValue
	}
	return nil
}

func (v *BoolValue) Get() interface{} {
	return v.Value
}

func (v *BoolValue) Names() []string {
	return v.Name
}

type Int64Value struct {
	Name    []string
	Value   int64
	Defined bool
}

func (v *Int64Value) Arg() (out []string) {
	if v.Defined {
		out = []string{fmt.Sprintf(`--%s=%v`, v.Name[0], v.Value)}
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

func (v *Int64Value) Get() interface{} {
	return v.Value
}

func (v *Int64Value) Names() []string {
	return v.Name
}

type StringValue struct {
	Name    []string
	Value   string
	Defined bool
}

func (v *StringValue) Arg() (out []string) {
	if v.Defined {
		out = []string{fmt.Sprintf(`--%s="%v"`, v.Name[0], v.Value)}
	}
	return
}

func (v *StringValue) Set(value interface{}) error {
	v.Defined = true
	v.Value = fmt.Sprintf("%v", value)
	return nil
}

func (v *StringValue) Get() interface{} {
	return v.Value
}

func (v *StringValue) Names() []string {
	return v.Name
}
