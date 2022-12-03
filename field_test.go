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
	"errors"
	"testing"

	"github.com/mediaexchange-io/assert"
)

func TestF_WithBoolTrue(t *testing.T) {
	Assert := assert.With(t)

	field := F("name", true)

	Assert.That(field.Name).IsEqualTo("name")
	Assert.That(field.StringValue).IsEqualTo("true")
}

func TestF_WithBoolFalse(t *testing.T) {
	Assert := assert.With(t)

	field := F("name", false)

	Assert.That(field.Name).IsEqualTo("name")
	Assert.That(field.StringValue).IsEqualTo("false")
}

func TestF_WithString(t *testing.T) {
	Assert := assert.With(t)

	field := F("name", "foobarbaz")

	Assert.That(field.Name).IsEqualTo("name")
	Assert.That(field.StringValue).IsEqualTo("foobarbaz")
}

func TestF_WithInt8(t *testing.T) {
	Assert := assert.With(t)

	field := F("name", -128)

	Assert.That(field.Name).IsEqualTo("name")
	Assert.That(field.StringValue).IsEqualTo("-128")
}

func TestF_WithInt16(t *testing.T) {
	Assert := assert.With(t)

	field := F("name", 65535)

	Assert.That(field.Name).IsEqualTo("name")
	Assert.That(field.StringValue).IsEqualTo("65535")
}

func TestF_WithInt32(t *testing.T) {
	Assert := assert.With(t)

	field := F("name", -2147483647)

	Assert.That(field.Name).IsEqualTo("name")
	Assert.That(field.StringValue).IsEqualTo("-2147483647")
}

func TestErr(t *testing.T) {
	Assert := assert.With(t)

	field := Err(errors.New("sample error"))

	Assert.That(field.Name).IsEqualTo("error")
	Assert.That(field.StringValue).IsEqualTo("sample error")
}

func TestField_Json(t *testing.T) {
	Assert := assert.With(t)
	var field Field

	field = F("name", true)
	Assert.That(field.Json()).IsEqualTo("\"name\":true")

	field = F("name", -127)
	Assert.That(field.Json()).IsEqualTo("\"name\":-127")

	field = F("name", -32767)
	Assert.That(field.Json()).IsEqualTo("\"name\":-32767")

	field = F("name", -2147483647)
	Assert.That(field.Json()).IsEqualTo("\"name\":-2147483647")

	field = F("name", -9223372036854775807)
	Assert.That(field.Json()).IsEqualTo("\"name\":-9223372036854775807")

	field = F("name", "value")
	Assert.That(field.Json()).IsEqualTo("\"name\":\"value\"")

	err := errors.New("error")
	field = Err(err)
	Assert.That(field.Json()).IsEqualTo("\"error\":\"error\"")
}

func TestField_String(t *testing.T) {
	Assert := assert.With(t)
	var field Field

	field = F("name", true)
	Assert.That(field.String()).IsEqualTo("name=true")

	field = F("name", -127)
	Assert.That(field.String()).IsEqualTo("name=-127")

	field = F("name", -32767)
	Assert.That(field.String()).IsEqualTo("name=-32767")

	field = F("name", -2147483647)
	Assert.That(field.String()).IsEqualTo("name=-2147483647")

	field = F("name", -9223372036854775807)
	Assert.That(field.String()).IsEqualTo("name=-9223372036854775807")

	field = F("name", "value")
	Assert.That(field.String()).IsEqualTo("name=value")

	err := errors.New("error")
	field = Err(err)
	Assert.That(field.String()).IsEqualTo("error=error")
}

func BenchmarkField_Json(b *testing.B) {
	var field Field

	err := errors.New("error message")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		field = F("name", true)
		field.Json()

		field = F("name", -127)
		field.Json()

		field = F("name", -32767)
		field.Json()

		field = F("name", -2147483647)
		field.Json()

		field = F("name", -9223372036854775807)
		field.Json()

		field = F("name", "value")
		field.Json()

		field = Err(err)
		field.Json()
	}
}

func BenchmarkField_String(b *testing.B) {
	var field Field

	err := errors.New("error message")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		field = F("name", true)
		field.String()

		field = F("name", -127)
		field.String()

		field = F("name", -32767)
		field.String()

		field = F("name", -2147483647)
		field.String()

		field = F("name", -9223372036854775807)
		field.String()

		field = F("name", "value")
		field.String()

		field = Err(err)
		field.String()
	}
}
