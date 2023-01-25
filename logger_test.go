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
	"bytes"
	"net"
	"strings"
	"sync"
	"testing"

	"github.com/mediaexchange-io/assert"
)

func TestSetLevel_WithWarning(t *testing.T) {
	// Since `go test` doesn't reset the test context for each unit test,
	// all subsequent tests should emit only warning and error messages.
	SetLevel(WARN)
}

func TestEmit_WithInfo_ShouldBeEmpty(t *testing.T) {
	Assert := assert.With(t)

	var buf bytes.Buffer
	SetWriter(&buf)

	Info("test")

	Assert.That(buf.String()).IsEmpty()
}

func TestEmit_WithWarn_ShouldNotBeEmpty(t *testing.T) {
	Assert := assert.With(t)

	var buf bytes.Buffer
	SetWriter(&buf)

	Warn("test")

	Assert.That(buf.String()).IsNotEmpty()
}

func TestEmit_WithWarnAndFields_ShouldNotBeEmpty(t *testing.T) {
	Assert := assert.With(t)

	var buf bytes.Buffer
	SetWriter(&buf)

	Warn("test", F("port", 9000), F("hostname", "localhost"))

	s := buf.String()
	Assert.That(s).IsNotEmpty()

	if !strings.Contains(s, "port=9000") {
		t.Fatal("Port field was not in the log message")
	}

	if !strings.Contains(s, "hostname=\"localhost\"") {
		t.Fatal("Hostname field was not in the log message")
	}
}

func TestSetServer_WithServer_ShouldReceive(t *testing.T) {
	Assert := assert.With(t)

	// Create a simple UDP server that exits after one message is received.
	addr, err := net.ResolveUDPAddr("udp", "localhost:9999")
	Assert.That(err).IsOk()

	sock, err := net.ListenUDP("udp", addr)
	Assert.That(err).IsOk()

	defer sock.Close()

	// Waits for the server goroutine to finish.
	var wait sync.WaitGroup

	go func() {
		buf := make([]byte, 4096)
		len, _, err := sock.ReadFromUDP(buf)
		Assert.That(err).IsOk()
		Assert.That(len).IsGreaterThan(0)
		Assert.That(buf).IsNotNil()

		s := string(buf)
		Assert.That(s).IsNotEmpty()

		if !strings.Contains(s, "\"port\":9000") {
			t.Fatal("Port field was not in the log message")
		}

		if !strings.Contains(s, "\"hostname\":\"localhost\"") {
			t.Fatal("Hostname field was not in the log message")
		}

		// Signal that the test is complete.
		wait.Done()
	}()

	// Prevent output on the console.
	var buf bytes.Buffer
	SetWriter(&buf)

	// Send a message to the server.
	SetServer("localhost:9999")
	wait.Add(1)
	Warn("test", F("port", 9000), F("hostname", "localhost"))

	// Wait for the test to complete.
	wait.Wait()
}
