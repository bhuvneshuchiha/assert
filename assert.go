// // Package myassert provides a simple assertion function for standard Go code,
// // similar to Python's assert statement.
// package assert


// // Enabled controls whether assertions are active. Set to false to disable.
// var Enabled = true

// // Assert panics if the condition is false, with an optional formatted message.
// // If Enabled is false, the assertion is skipped.
// // Usage: myassert.Assert(x == y, "x should equal y, got %d, want %d", x, y)
// func Assert(condition bool, format string, args ...any) {
//     if !Enabled {
//         return
//     }
//     if !condition {
//         if format == "" {
//             panic("assertion failed")
//         }
//         panic(fmt.Errorf(format, args...))
//     }
// }


package assert

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"reflect"
	"runtime/debug"
)

// TODO using slog for logging
type AssertData interface {
    Dump() string
}
type AssertFlush interface {
    Flush()
}

var flushes []AssertFlush = []AssertFlush{}
var assertData map[string]AssertData = map[string]AssertData{}
var writer io.Writer

func AddAssertData(key string, value AssertData) {
	assertData[key] = value
}

func RemoveAssertData(key string) {
	delete(assertData, key)
}

func AddAssertFlush(flusher AssertFlush) {
    flushes = append(flushes, flusher)
}

func ToWriter(w io.Writer) {
	writer = w
}

func runAssert(msg string, args ...interface{}) {
    // There is a bit of a issue here.  if you flush you cannot assert
    // cannot be reentrant
    // TODO I am positive i could create some sort of latching that prevents the
    // reentrant problem
    for _, f := range flushes {
        f.Flush()
    }

    slogValues := []interface{}{
        "msg",
        msg,
        "area",
        "Assert",
    }
    slogValues = append(slogValues, args...)
    fmt.Fprintf(os.Stderr, "ARGS: %+v\n", args)

	for k, v := range assertData {
        slogValues = append(slogValues, k, v.Dump())
	}

    fmt.Fprintf(os.Stderr, "ASSERT\n")
    for i := 0; i < len(slogValues); i += 2 {
        fmt.Fprintf(os.Stderr, "   %s=%v\n", slogValues[i], slogValues[i + 1])
    }
    fmt.Fprintln(os.Stderr, string(debug.Stack()))
    os.Exit(1)
}

// TODO Think about passing around a context for debugging purposes
func Assert(truth bool, msg string, data ...any) {
	if !truth {
		runAssert(msg, data...)
	}
}

func Nil(item any, msg string, data ...any) {
    slog.Info("Nil Check", "item", item)
	if item == nil {
        return
    }

    slog.Error("Nil#not nil encountered")
    runAssert(msg, data...)
}

func NotNil(item any, msg string, data ...any) {
	if item == nil || reflect.ValueOf(item).Kind() == reflect.Ptr && reflect.ValueOf(item).IsNil() {
		slog.Error("NotNil#nil encountered")
		runAssert(msg, data...)
	}
}

func Never(msg string, data ...any) {
    runAssert(msg, data...)
}

func NoError(err error, msg string, data ...any) {
	if err != nil {
        data = append(data, "error", err)
		runAssert(msg, data...)
	}
}

