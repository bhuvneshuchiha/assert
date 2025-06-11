// Package myassert provides a simple assertion function for standard Go code,
// similar to Python's assert statement.
package assert
import "fmt"

// Enabled controls whether assertions are active. Set to false to disable.
var Enabled = true

// Assert panics if the condition is false, with an optional formatted message.
// If Enabled is false, the assertion is skipped.
// Usage: myassert.Assert(x == y, "x should equal y, got %d, want %d", x, y)
func Assert(condition bool, format string, args ...any) {
    if !Enabled {
        return
    }
    if !condition {
        if format == "" {
            panic("assertion failed")
        }
        panic(fmt.Errorf(format, args...))
    }
}
