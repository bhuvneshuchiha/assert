# 🛡️ Assert

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/bhuvneshuchiha/assert)](https://goreportcard.com/report/github.com/bhuvneshuchiha/assert)

A powerful assertion library for Go that provides runtime assertions with rich debugging information, context data collection, and graceful error handling.

## ✨ Features

- 🔍 **Rich Debugging Output** - Detailed stack traces and context information
- 📊 **Context Data Collection** - Attach and collect debugging data from multiple sources
- 💧 **Flush Support** - Automatic flushing of buffers and resources before assertion failures  
- 🎯 **Multiple Assertion Types** - Comprehensive set of assertion functions
- 🔧 **Configurable Output** - Customizable output writers and formatting
- 🚀 **Zero Dependencies** - Uses only Go standard library

## 📦 Installation

```bash
go get github.com/bhuvneshuchiha/assert
```

## 🚀 Quick Start

### Basic Assertions

```go
package main

import "github.com/bhuvneshuchiha/assert"

func main() {
    x := 42
    y := 42
    
    // Basic assertion
    assert.Assert(x == y, "x should equal y", "x", x, "y", y)
    
    // Error checking
    err := someFunction()
    assert.NoError(err, "someFunction should not return error")
    
    // Nil checking
    var ptr *int
    assert.Nil(ptr, "pointer should be nil")
    
    ptr = &x
    assert.NotNil(ptr, "pointer should not be nil")
}
```

### Advanced Usage with Context Data

```go
// Implement AssertData interface for custom debugging info
type DatabaseState struct {
    ConnectionCount int
    ActiveQueries   []string
}

func (d *DatabaseState) Dump() string {
    return fmt.Sprintf("Connections: %d, Active: %v", 
        d.ConnectionCount, d.ActiveQueries)
}

// Add context data that will be included in assertion failures
dbState := &DatabaseState{ConnectionCount: 5, ActiveQueries: []string{"SELECT * FROM users"}}
assert.AddAssertData("database", dbState)

// This assertion failure will include database state
assert.Assert(false, "Something went wrong")
```

### Custom Flushers

```go
// Implement AssertFlush interface for cleanup before assertions
type BufferFlusher struct {
    buffer *bytes.Buffer
}

func (b *BufferFlusher) Flush() {
    // Flush buffer contents before assertion failure
    fmt.Print(b.buffer.String())
    b.buffer.Reset()
}

// Add flusher that will be called before assertion failures
flusher := &BufferFlusher{buffer: myBuffer}
assert.AddAssertFlush(flusher)
```

## 📋 Available Assertions

### `Assert(condition bool, msg string, data ...any)`
Basic assertion that checks if a condition is true.

```go
assert.Assert(len(slice) > 0, "slice should not be empty", "length", len(slice))
```

### `NoError(err error, msg string, data ...any)`
Asserts that an error is nil.

```go
file, err := os.Open("config.json")
assert.NoError(err, "failed to open config file", "filename", "config.json")
```

### `Nil(item any, msg string, data ...any)`
Asserts that an item is nil.

```go
var result *Result
assert.Nil(result, "result should be nil before initialization")
```

### `NotNil(item any, msg string, data ...any)`
Asserts that an item is not nil (handles both nil values and nil pointers).

```go
config := loadConfig()
assert.NotNil(config, "config should not be nil after loading")
```

### `Never(msg string, data ...any)`
Always triggers an assertion failure. Useful for code paths that should never be reached.

```go
switch status {
case StatusActive:
    // handle active
case StatusInactive:
    // handle inactive
default:
    assert.Never("unknown status encountered", "status", status)
}
```

## 🔧 Configuration

### Custom Output Writer

```go
// Redirect assertion output to a custom writer
var buf bytes.Buffer
assert.ToWriter(&buf)

// Or redirect to a file
logFile, _ := os.Create("assertions.log")
assert.ToWriter(logFile)
```

### Managing Context Data

```go
// Add debugging context
assert.AddAssertData("user_session", sessionData)
assert.AddAssertData("request_info", requestData)

// Remove when no longer needed
assert.RemoveAssertData("user_session")
```

## 📊 Output Format

When an assertion fails, you'll see detailed output including:

```
ARGS: [user_id 12345 operation login]
ASSERT
   msg=user authentication failed
   area=Assert
   user_id=12345
   operation=login
   database=Connections: 5, Active: [SELECT * FROM users]
   session=SessionID: abc123, UserID: 12345

goroutine 1 [running]:
runtime/debug.Stack()
    /usr/local/go/src/runtime/debug/stack.go:24 +0x5e
github.com/bhuvneshuchiha/assert.runAssert(...)
    /path/to/assert/assert.go:45
github.com/bhuvneshuchiha/assert.Assert(...)
    /path/to/assert/assert.go:58
main.authenticateUser(...)
    /path/to/main.go:123
```

## 🏗️ Interfaces

### AssertData Interface
Implement this interface to provide custom debugging information:

```go
type AssertData interface {
    Dump() string
}

type MyDebugInfo struct {
    State string
    Count int
}

func (m *MyDebugInfo) Dump() string {
    return fmt.Sprintf("State: %s, Count: %d", m.State, m.Count)
}
```

### AssertFlush Interface
Implement this interface to perform cleanup before assertion failures:

```go
type AssertFlush interface {
    Flush()
}

type LogFlusher struct {
    logger *log.Logger
}

func (l *LogFlusher) Flush() {
    // Ensure all logs are written before assertion
    l.logger.Writer().(*os.File).Sync()
}
```

## ⚠️ Important Notes

- **Program Termination**: All assertion failures call `os.Exit(1)` to terminate the program
- **Reentrant Safety**: The library has protection against reentrant assertion calls during flush operations
- **Production Use**: Consider the performance impact of context data collection in production environments
- **Stack Traces**: Full stack traces are included in assertion output for debugging

## 🔄 Best Practices

1. **Use Descriptive Messages**: Always provide clear, actionable assertion messages
   ```go
   // ✅ Good
   assert.Assert(user.IsActive(), "user must be active to perform this operation", 
       "user_id", user.ID, "status", user.Status)
   
   // ❌ Bad
   assert.Assert(user.IsActive(), "assertion failed")
   ```

2. **Clean Up Context Data**: Remove assertion data when it's no longer relevant
   ```go
   assert.AddAssertData("request", requestData)
   defer assert.RemoveAssertData("request")
   ```

3. **Use Appropriate Assertion Types**: Choose the most specific assertion for your use case
   ```go
   // ✅ Use specific assertions
   assert.NoError(err, "database connection failed")
   assert.NotNil(user, "user lookup returned nil")
   
   // ❌ Less clear
   assert.Assert(err == nil, "error occurred")
   assert.Assert(user != nil, "user is nil")
   ```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 🙏 Acknowledgments

- Inspired by Python's/Lua's assert statement and other assertion libraries
- Built with Go's excellent standard library
- Designed for maximum debugging utility

---

**Note**: This library is designed for development and testing environments. Consider the performance implications of context data collection in production systems.
