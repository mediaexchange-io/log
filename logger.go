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
	"io"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	ISO8601Micro = "2006-01-02T15:04:05.000000Z0700"
)

var (
	conn        *net.UDPConn
	emitJson    bool
	level       Level
	programName string
	writer      io.Writer
)

// init creates a default console logger that can be used immediately with no
// further configuration necessary.
func init() {
	// By default the logger emits formatted text, not JSON
	emitJson = false

	// Set the minimum logging level emitted to INFO
	level = INFO

	// Set the name of the program.
	programName = path.Base(os.Args[0])

	// Default to stderr which is the Posix standard as it's unbuffered.
	writer = os.Stderr
}

// SetEmitJson changes the type of output sent to the aggregator.
func SetEmitJson(b bool) {
	emitJson = b
}

// SetLevel changes the minimum level that is emitted. This is used to prevent
// logging from becoming too noisy. By default the minimum level is set to
// INFO so that DEBUG messages don't overwhelm the aggregator.
func SetLevel(l Level) {
	level = l
}

// SetWriter sets the output for all future log messages not bound for an
// aggregator.
func SetWriter(w io.Writer) {
	if w == nil {
		panic("SetWriter: Parameter was nil")
	}
	writer = w
}

// SetServer starts a UDP client that sends messages to a log aggregation
// server. This occurs in parallel with console logging.
func SetServer(address string) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		panic(err)
	}

	conn, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		panic(err)
	}

	emitJson = true
}

// Debug emits a message with the DEBUG level.
func Debug(msg string, fields ...Field) {
	if level <= DEBUG {
		emit(DEBUG, msg, fields)
	}
}

// Info emits a message with the INFO level.
func Info(msg string, fields ...Field) {
	if level <= INFO {
		emit(INFO, msg, fields)
	}
}

// Warn emits a message with the WARN level.
func Warn(msg string, fields ...Field) {
	if level <= WARN {
		emit(WARN, msg, fields)
	}
}

// Error emits a message with the ERROR level.
func Error(msg string, fields ...Field) {
	if level <= ERROR {
		emit(ERROR, msg, fields)
	}
}

// emit generates the log message and sends it.
func emit(level Level, message string, fields []Field) {
	// Pull the timestamp now so the UDP aggregator and the console are consistent.
	t := time.Now().UTC()

	// Send the message to the log aggregator. Note that UDP is fast, but
	// unreliable. Messages may be received out-of-order or not at all.
	if conn != nil {
		go conn.Write(json(t, message, fields))
	}

	// Send the message to the console logger.
	if emitJson {
		writer.Write(json(t, message, fields))
	} else {
		writer.Write(text(t, message, fields))
	}
}

func json(t time.Time, message string, fields []Field) []byte {
	return []byte("{\"time\":" + strconv.FormatInt(t.UnixNano(), 10) + ",\"name\":\"" + programName + "\",\"level\":\"" + level.String() + "\",\"message\":\"" + message + "\",\"fields\":" + fieldJson(fields) + "}")
}

func text(t time.Time, message string, fields []Field) []byte {
	return []byte(t.Format(ISO8601Micro) + " [" + programName + "] " + level.String() + " " + message + fieldString(fields) + "\n")
}

func fieldJson(fields []Field) string {
	var builder strings.Builder
	builder.WriteRune('{')
	addComma := false
	for _, field := range fields {
		if addComma {
			builder.WriteRune(',')
		} else {
			addComma = true
		}

		builder.WriteString(field.Json())
	}
	builder.WriteRune('}')
	return builder.String()
}

func fieldString(fields []Field) string {
	var builder strings.Builder
	for _, field := range fields {
		builder.WriteRune(' ')
		builder.WriteString(field.String())
	}
	return builder.String()
}
