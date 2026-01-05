# TOML Parser HOWTO

This guide provides best practices and usage patterns for the Tideland TOML parser.

## Table of Contents

1. [Choosing Access Patterns](#choosing-access-patterns)
2. [Working with Nested Structures](#working-with-nested-structures)
3. [Handling Errors](#handling-errors)
4. [Table Arrays](#table-arrays)
5. [Type Conversions](#type-conversions)
6. [Configuration Management](#configuration-management)
7. [Common Patterns](#common-patterns)

## Choosing Access Patterns

The parser provides two access patterns. Choose based on your use case:

### Flat Access - When to Use

Use flat dot notation when:
- You know the exact path to the value
- You're accessing a few specific values
- You want concise, readable code

```go
// Direct and simple
host := must(doc.GetString("server.host"))
port := must(doc.GetInt("server.port"))
timeout := must(doc.GetInt("server.timeout"))
```

### Hierarchical Navigation - When to Use

Use hierarchical navigation when:
- You're exploring unknown structure
- You need to iterate over sections or keys
- You're building dynamic configuration systems
- You want to validate configuration structure

```go
// Explore all servers
for _, name := range doc.Sections() {
    section, _ := doc.Section(name)
    if section.Has("host") && section.Has("port") {
        // It's a server configuration
        host, _ := section.GetString("host")
        port, _ := section.GetInt("port")
        fmt.Printf("Server %s: %s:%d\n", name, host, port)
    }
}
```

### Combining Patterns

Often the best approach combines both:

```go
// Use flat access for known paths
dbSection, _ := doc.Section("database")

// Then navigate hierarchically
for _, connName := range dbSection.Sections() {
    conn, _ := dbSection.Section(connName)
    
    // Back to flat for specific values
    host, _ := conn.GetString("host")
    port, _ := conn.GetInt("port")
    
    registerDatabase(connName, host, port)
}
```

## Working with Nested Structures

### Dotted Keys Create Structure

```toml
# These are equivalent:
[physical]
color = "orange"
shape = "round"

# vs

physical.color = "orange"
physical.shape = "round"
```

Both create the same structure:

```go
physical, _ := doc.Section("physical")
color, _ := physical.GetString("color")
shape, _ := physical.GetString("shape")
```

### Deep Nesting

For deeply nested structures, dot notation is cleaner:

```toml
app.server.http.port = 8080
app.server.http.host = "localhost"
app.server.https.port = 8443
app.server.https.enabled = true
```

```go
// Simple access
httpPort, _ := doc.GetInt("app.server.http.port")
httpsEnabled, _ := doc.GetBool("app.server.https.enabled")
```

### Checking Existence

Always check if paths exist before accessing:

```go
if doc.Has("optional.feature.enabled") {
    enabled, _ := doc.GetBool("optional.feature.enabled")
    if enabled {
        // Configure optional feature
    }
}
```

## Handling Errors

### Error Checking Pattern

Always check errors for file I/O and parsing:

```go
doc, err := toml.Parse("config.toml")
if err != nil {
    return fmt.Errorf("failed to parse config: %w", err)
}
```

### Detailed Error Information

Extract details from parse errors:

```go
doc, err := toml.Parse("config.toml")
if err != nil {
    if parseErr, ok := err.(*toml.ParseError); ok {
        log.Printf("Parse error in %s at line %d, column %d: %v",
            parseErr.Path,
            parseErr.Line,
            parseErr.Col,
            parseErr.Err)
        
        // Handle specific error types
        switch parseErr.Code {
        case toml.ErrSyntax:
            // Syntax error - check TOML file
        case toml.ErrType:
            // Type mismatch - wrong getter used
        case toml.ErrNotFound:
            // Missing key - check path or make optional
        }
    }
    return err
}
```

### Must Pattern for Required Config

For required configuration values:

```go
func must[T any](val T, err error) T {
    if err != nil {
        panic(fmt.Sprintf("required config missing: %v", err))
    }
    return val
}

// Use for critical configuration
host := must(doc.GetString("server.host"))
port := must(doc.GetInt("server.port"))
```

### Defaults Pattern for Optional Config

For optional values with defaults:

```go
func getStringOr(doc *toml.Document, path, defaultVal string) string {
    if !doc.Has(path) {
        return defaultVal
    }
    val, err := doc.GetString(path)
    if err != nil {
        return defaultVal
    }
    return val
}

// Use for optional configuration
logLevel := getStringOr(doc, "logging.level", "info")
maxRetries := getIntOr(doc, "client.max_retries", 3)
```

## Table Arrays

### Iterating Table Arrays

```toml
[[servers]]
name = "alpha"
ip = "10.0.0.1"

[[servers]]
name = "beta"
ip = "10.0.0.2"
```

```go
servers, err := doc.TableArray("servers")
if err != nil {
    log.Fatal(err)
}

for i, server := range servers {
    name, _ := server.GetString("name")
    ip, _ := server.GetString("ip")
    fmt.Printf("Server %d: %s (%s)\n", i+1, name, ip)
}
```

### Nested Table Arrays

```toml
[[products]]
name = "Hammer"

  [[products.variants]]
  sku = "HAM-001"
  color = "red"

  [[products.variants]]
  sku = "HAM-002"
  color = "blue"

[[products]]
name = "Screwdriver"

  [[products.variants]]
  sku = "SCR-001"
  size = "small"
```

```go
products, _ := doc.TableArray("products")
for _, product := range products {
    name, _ := product.GetString("name")
    fmt.Printf("Product: %s\n", name)
    
    variants, _ := product.TableArray("variants")
    for _, variant := range variants {
        sku, _ := variant.GetString("sku")
        fmt.Printf("  SKU: %s\n", sku)
    }
}
```

### Building Structures from Table Arrays

```go
type Server struct {
    Name string
    IP   string
    Port int
}

func loadServers(doc *toml.Document) ([]Server, error) {
    serverTables, err := doc.TableArray("servers")
    if err != nil {
        return nil, err
    }
    
    servers := make([]Server, len(serverTables))
    for i, st := range serverTables {
        servers[i] = Server{
            Name: must(st.GetString("name")),
            IP:   must(st.GetString("ip")),
            Port: must(st.GetInt("port")),
        }
    }
    return servers, nil
}
```

## Type Conversions

### Arrays

Arrays in TOML can be homogeneous or mixed:

```toml
numbers = [1, 2, 3, 4, 5]
mixed = [1, "two", 3.0, true]
nested = [[1, 2], [3, 4], [5, 6]]
```

```go
// Homogeneous arrays
numbers, _ := doc.GetArray("numbers")
for _, n := range numbers {
    fmt.Printf("%d ", n.(int64))
}

// Mixed arrays - type assert carefully
mixed, _ := doc.GetArray("mixed")
for _, val := range mixed {
    switch v := val.(type) {
    case int64:
        fmt.Printf("int: %d\n", v)
    case string:
        fmt.Printf("string: %s\n", v)
    case float64:
        fmt.Printf("float: %f\n", v)
    case bool:
        fmt.Printf("bool: %t\n", v)
    }
}
```

### Datetime Values

```toml
offset = 1979-05-27T07:32:00Z
local = 1979-05-27T07:32:00
date = 1979-05-27
time = 07:32:00
```

```go
offset, _ := doc.GetTime("offset")
fmt.Printf("Offset: %v (UTC offset: %d)\n", offset, offset.Unix())

local, _ := doc.GetTime("local")
fmt.Printf("Local: %v\n", local.Format("2006-01-02 15:04:05"))
```

### Inline Tables

```toml
point = {x = 1, y = 2, z = 3}
```

```go
// Access as a value
pointVal := doc.root.keys["point"]
table, _ := pointVal.Table()

x, _ := table["x"].Int()
y, _ := table["y"].Int()
z, _ := table["z"].Int()

fmt.Printf("Point: (%d, %d, %d)\n", x, y, z)
```

## Configuration Management

### Wrapper Pattern

Create a configuration struct that wraps the document:

```go
type AppConfig struct {
    doc *toml.Document
}

func LoadConfig(path string) (*AppConfig, error) {
    doc, err := toml.Parse(path)
    if err != nil {
        return nil, err
    }
    return &AppConfig{doc: doc}, nil
}

func (c *AppConfig) ServerHost() string {
    return must(c.doc.GetString("server.host"))
}

func (c *AppConfig) ServerPort() int {
    port, _ := c.doc.GetInt("server.port")
    if port == 0 {
        return 8080 // default
    }
    return int(port)
}

func (c *AppConfig) DatabaseURL() string {
    return must(c.doc.GetString("database.url"))
}
```

Usage:

```go
config, err := LoadConfig("app.toml")
if err != nil {
    log.Fatal(err)
}

server := &http.Server{
    Addr: fmt.Sprintf("%s:%d", 
        config.ServerHost(),
        config.ServerPort()),
}
```

### Environment Override Pattern

Combine TOML with environment variables:

```go
func (c *AppConfig) ServerPort() int {
    // Check environment first
    if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
        if port, err := strconv.Atoi(portStr); err == nil {
            return port
        }
    }
    
    // Fall back to TOML
    port, _ := c.doc.GetInt("server.port")
    if port == 0 {
        return 8080 // default
    }
    return int(port)
}
```

## Common Patterns

### Validation

Validate configuration after loading:

```go
func (c *AppConfig) Validate() error {
    required := []string{
        "server.host",
        "server.port",
        "database.url",
    }
    
    for _, path := range required {
        if !c.doc.Has(path) {
            return fmt.Errorf("missing required config: %s", path)
        }
    }
    
    // Value validation
    port, _ := c.doc.GetInt("server.port")
    if port < 1 || port > 65535 {
        return fmt.Errorf("invalid port: %d", port)
    }
    
    return nil
}
```

### Hot Reload

Watch for config file changes:

```go
type ConfigWatcher struct {
    path   string
    config *AppConfig
    mu     sync.RWMutex
}

func (w *ConfigWatcher) Reload() error {
    newConfig, err := LoadConfig(w.path)
    if err != nil {
        return err
    }
    
    if err := newConfig.Validate(); err != nil {
        return err
    }
    
    w.mu.Lock()
    w.config = newConfig
    w.mu.Unlock()
    
    return nil
}

func (w *ConfigWatcher) Get() *AppConfig {
    w.mu.RLock()
    defer w.mu.RUnlock()
    return w.config
}
```

### Multi-Environment Configs

Load different configs per environment:

```go
func LoadEnvironmentConfig(env string) (*AppConfig, error) {
    paths := []string{
        "config.toml",              // Base config
        fmt.Sprintf("config.%s.toml", env), // Environment-specific
    }
    
    var config *AppConfig
    for _, path := range paths {
        if _, err := os.Stat(path); err == nil {
            cfg, err := LoadConfig(path)
            if err != nil {
                return nil, err
            }
            if config == nil {
                config = cfg
            } else {
                // Merge configs (simplified)
                config = mergeConfigs(config, cfg)
            }
        }
    }
    
    return config, nil
}
```

## Best Practices

1. **Always check errors** when parsing files and accessing values
2. **Use Has() before Get()** for optional configuration
3. **Provide defaults** for non-critical settings
4. **Validate after loading** to fail fast on bad configuration
5. **Use wrapper types** to encapsulate configuration access
6. **Document your TOML schema** with comments in example files
7. **Test with invalid configs** to ensure error handling works
8. **Keep TOML files simple** - avoid excessive nesting
9. **Use table arrays** for lists of similar items
10. **Prefer flat access** for known paths, hierarchical for exploration

## See Also

- Package documentation: `go doc tideland.dev/go/toml`
- TOML specification: https://toml.io/en/v1.0.0
- Test files for more examples: `*_test.go`
