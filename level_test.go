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
	"github.com/MediaExchange/assert"
	"testing"
)

func TestLevel_String(t *testing.T) {
	var actual string
	var assert = assert.With(t)

	actual = DEBUG.String()
	assert.That(actual).IsEqualTo("DEBUG")

	actual = INFO.String()
	assert.That(actual).IsEqualTo("INFO")

	actual = WARN.String()
	assert.That(actual).IsEqualTo("WARN")

	actual = ERROR.String()
	assert.That(actual).IsEqualTo("ERROR")

	assert.ThatPanics(func() { Level(-1).String() })
}

func TestParse(t *testing.T) {
	var actual Level
	var assert = assert.With(t)

	actual = Parse("DEBUG")
	assert.That(actual).IsEqualTo(DEBUG)

	actual = Parse("INFO")
	assert.That(actual).IsEqualTo(INFO)

	actual = Parse("WARN")
	assert.That(actual).IsEqualTo(WARN)

	actual = Parse("ERROR")
	assert.That(actual).IsEqualTo(ERROR)

	assert.ThatPanics(func() { Parse("foobar") })
}
