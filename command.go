package clc

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

// Program is a interface specifying a sub functionality of a given cli program.
// It should have a Exec method.
type Program interface {
	Exec() error
}

// Command struct contains all the commands and the associated functions to run
type Command struct {
	Args Program
	help string
}

func (c *Command) Exec() error {
	return c.Args.Exec()
}

// ParseArgs is used to parse the given arguments into appropriate command argument structure.
func (c *Command) ParseArgs(args []string) error {
	currentOption := ""
	rs := reflect.ValueOf(c.Args).Elem()
	for i := 0; i < len(args); i++ {
		val := args[i]
		if len(val) < 2 {
			return errors.New("Invalid Option: '" + val + "'")
		}
		if val[0] != '-' {
			currentField := rs.FieldByName("RestArgs")
			if !currentField.IsValid() {
				continue
			}
			if currentField.Kind() != reflect.Array && currentField.Kind() != reflect.Slice {
				return errors.New("the RestArgs should be an array or a slice element of strings")
			}
			rags, ok := currentField.Interface().([]string)
			if !ok {
				return errors.New("the RestArgs should be an array or a slice element of strings")
			}
			rags = append(rags, val)
			currentField.Set(reflect.ValueOf(rags))
			continue
		}
		currentOption = val[1:]
		currentField := rs.FieldByName(capitalize(currentOption))
		if !currentField.IsValid() {
			i++
			if i >= len(args) {
				return nil
			}
			val = args[i]
			if len(val) < 2 {
				return errors.New("Invalid Option: '" + val + "'")
			}
			if val[0] == '-' {
				i--
				continue
			}
			continue
		}
		if currentField.Type().Kind() == reflect.Bool {
			currentField.SetBool(true)
			continue
		}
		i++
		if i >= len(args) {
			return nil
		}
		val = args[i]
		if len(val) < 2 {
			return errors.New("Invalid Option: '" + val + "'")
		}
		switch currentField.Type().Kind() {
		case reflect.Int:
			value, err := strconv.Atoi(val)
			if err != nil {
				return errors.New("Expected int<4>, got '" + val + "'")
			}
			currentField.Set(reflect.ValueOf(value))
		case reflect.Int8:
			value, err := strconv.ParseInt(val, 10, 8)
			if err != nil {
				return errors.New("Expected int<8>, got '" + val + "'")
			}
			currentField.Set(reflect.ValueOf(int8(value)))
		case reflect.Int16:
			value, err := strconv.ParseInt(val, 10, 16)
			if err != nil {
				return errors.New("Expected int<16>, got '" + val + "'")
			}
			currentField.Set(reflect.ValueOf(int16(value)))
		case reflect.Int32:
			value, err := strconv.ParseInt(val, 10, 32)
			if err != nil {
				return errors.New("Expected int<32>, got '" + val + "'")
			}
			currentField.Set(reflect.ValueOf(int32(value)))
		case reflect.Int64:
			value, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return errors.New("Expected int<64>, got '" + val + "'")
			}
			currentField.Set(reflect.ValueOf(int64(value)))
		// Unsigned....
		case reflect.Uint:
			value, err := strconv.ParseUint(val, 10, 4)
			if err != nil {
				return errors.New("Expected int<4>, got '" + val + "'")
			}
			currentField.Set(reflect.ValueOf(value))
		case reflect.Uint8:
			value, err := strconv.ParseUint(val, 10, 8)
			if err != nil {
				return errors.New("Expected uint<8>, got '" + val + "'")
			}
			currentField.Set(reflect.ValueOf(uint8(value)))
		case reflect.Uint16:
			value, err := strconv.ParseUint(val, 10, 16)
			if err != nil {
				return errors.New("Expected uint<16>, got '" + val + "'")
			}
			currentField.Set(reflect.ValueOf(uint16(value)))
		case reflect.Uint32:
			value, err := strconv.ParseUint(val, 10, 32)
			if err != nil {
				return errors.New("Expected uint<32>, got '" + val + "'")
			}
			currentField.Set(reflect.ValueOf(uint32(value)))
		case reflect.Uint64:
			value, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				return errors.New("Expected uint<64>, got '" + val + "'")
			}
			currentField.Set(reflect.ValueOf(uint64(value)))
		case reflect.String:
			currentField.Set(reflect.ValueOf(val))
		case reflect.Float32:
			value, err := strconv.ParseFloat(val, 32)
			if err != nil {
				return errors.New("Expected float<32>, got '" + val + "'")
			}
			currentField.Set(reflect.ValueOf(float32(value)))
		case reflect.Float64:
			value, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return errors.New("Expected float<64>, got '" + val + "'")
			}
			currentField.Set(reflect.ValueOf(value))
		}
	}
	return nil
}

// Capitalize is used to just capitalize the first letter of the
// passed string.
func capitalize(a string) string {
	if len(a) <= 1 {
		return strings.ToUpper(a)
	}
	return strings.ToUpper(string(a[0])) + string(a[1:])
}
