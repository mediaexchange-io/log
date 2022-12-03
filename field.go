/*
 * Copyright 2019 MediaExchange.io
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package log

import (
	"reflect"
	"strconv"
)

// Field stores a name/value pair to be formatted by the emitter.
type Field struct {
	Name        string
	StringValue string
	Quoted      bool
}

// F returns a field with any value converted to a string, ready for logging.
func F(name string, value any) Field {
	f := Field{Name: name, Quoted: false}
	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)

	switch t.Kind() {
	case reflect.Bool:
		if v.Bool() {
			f.StringValue = "true"
		} else {
			f.StringValue = "false"
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		f.StringValue = strconv.FormatInt(v.Int(), 10)
	case reflect.String:
		f.StringValue = v.String()
		f.Quoted = true
	default:
		f.StringValue = "<[Unknown Type " + t.Kind().String() + "]>"
	}

	return f
}

// Err returns a Field that contains the message from an error.
func Err(value error) Field {
	return Field{Name: "error", Quoted: true, StringValue: value.Error()}
}

// String returns the contents of the Field as `key=value`.
func (field Field) String() string {
	return field.Name + "=" + field.StringValue
}

// Json returns the contents of the Field as `"key":value`.
func (field Field) Json() string {
	return "\"" + field.Name + "\":" + quotedValue(field)
}

// quotedValue converts the field's value into a string.
func quotedValue(field Field) string {
	if field.Quoted {
		return "\"" + field.StringValue + "\""
	}
	return field.StringValue
}
