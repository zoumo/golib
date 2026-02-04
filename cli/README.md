# CLI

cli is a simple package for building Cobra commands with common patterns.

## Example

### Simple Command (No Flags)

```go
package main

import (
    "fmt"

    "github.com/spf13/cobra"
    "github.com/zoumo/golib/cli"
)

var _ cli.Command = &GreetCommand{}

// GreetCommand is the simplest command - no flags, just positional args.
type GreetCommand struct{}

func (c *GreetCommand) Name() string {
    return "greet"
}

func (c *GreetCommand) Run(_ *cobra.Command, args []string) error {
    if len(args) == 0 {
        fmt.Println("Hello, World!")
        return nil
    }
    fmt.Printf("Hello, %s!\n", args[0])
    return nil
}

func main() {
    root := &cobra.Command{Use: "myapp"}
    root.AddCommand(cli.NewCobraCommand(&GreetCommand{}))
    root.Execute()
}
```

### Command with Options (BindFlags)

```go
package main

import (
    "fmt"

    "github.com/spf13/cobra"
    "github.com/spf13/pflag"
    "github.com/zoumo/golib/cli"
)

var (
    _ cli.Command = &EchoCommand{}
    _ cli.Options = &EchoCommand{}
)

// EchoCommand demonstrates a command with flags.
type EchoCommand struct {
    count int
}

func (c *EchoCommand) Name() string {
    return "echo"
}

func (c *EchoCommand) BindFlags(fs *pflag.FlagSet) {
    fs.IntVar(&c.count, "count", 1, "Number of times to echo")
}

func (c *EchoCommand) Run(_ *cobra.Command, args []string) error {
    if len(args) == 0 {
        return fmt.Errorf("usage: echo <text>")
    }
    for i := 0; i < c.count; i++ {
        fmt.Println(args[0])
    }
    return nil
}

func main() {
    root := &cobra.Command{Use: "myapp"}
    root.AddCommand(cli.NewCobraCommand(&EchoCommand{}))
    root.Execute()
}
```

### Complex Command (With CommonOptions, Complete, Validate)

```go
package main

import (
    "errors"
    "fmt"

    "github.com/spf13/cobra"
    "github.com/spf13/pflag"
    "github.com/zoumo/golib/cli"
)

var _ cli.ComplexOptions = &QueryOptions{}

// QueryOptions implements Options and ComplexOptions interfaces.
// By embedding CommonOptions, it inherits Workspace and Logger fields.
type QueryOptions struct {
    cli.CommonOptions

    Resource string
    Limit    int
}

// BindFlags implements Options interface.
func (o *QueryOptions) BindFlags(fs *pflag.FlagSet) {
    fs.StringVar(&o.Resource, "resource", "", "Resource to query (required)")
    fs.IntVar(&o.Limit, "limit", 10, "Maximum results")
}

// Complete implements ComplexOptions interface.
// MUST call embedded CommonOptions.Complete to initialize Logger.
func (o *QueryOptions) Complete(cmd *cobra.Command, args []string) error {
    return o.CommonOptions.Complete(cmd, args)
}

// Validate implements ComplexOptions interface.
// MUST call embedded CommonOptions.Validate first, then add custom validation.
func (o *QueryOptions) Validate() error {
    // Call parent first
    if err := o.CommonOptions.Validate(); err != nil {
        return err
    }
    // Custom validation
    if o.Resource == "" {
        return errors.New("--resource is required")
    }
    if o.Limit < 0 {
        return errors.New("--limit must be non-negative")
    }
    return nil
}

var (
    _ cli.Command        = &QueryCommand{}
    _ cli.ComplexOptions = &QueryCommand{}
)

// QueryCommand implements Command interface.
// Through pointer embedding of *QueryOptions, it also satisfies Options and ComplexOptions.
type QueryCommand struct {
    *QueryOptions
}

// Name implements Command interface.
func (c *QueryCommand) Name() string {
    return "query"
}

// Run implements Command interface.
func (c *QueryCommand) Run(_ *cobra.Command, args []string) error {
    // Access fields directly through pointer embedding
    c.Logger.Info("Querying", "resource", c.Resource, "limit", c.Limit)
    fmt.Printf("Querying: %s (limit: %d)\n", c.Resource, c.Limit)
    return nil
}

func main() {
    root := &cobra.Command{Use: "myapp"}
    root.AddCommand(cli.NewCobraCommand(&QueryCommand{
        QueryOptions: &QueryOptions{},
    }))
    root.Execute()
}
```

**Key Points:**
- `QueryOptions` implements both `Options` (via `BindFlags`) and `ComplexOptions` (via `Complete`, `Validate`)
- `QueryOptions` embeds `CommonOptions` to inherit `Workspace` and `Logger` fields
- `QueryOptions.Complete()` and `QueryOptions.Validate()` MUST call the embedded `CommonOptions` methods to avoid shadowing
- `QueryCommand` implements `Command` interface and embeds `*QueryOptions` pointer
- Through pointer embedding, `QueryCommand` gets promoted methods and fields from `QueryOptions`, satisfying all type assertions

**Compile-time interface assertions** (the `var _` declarations):
- Go's idiomatic way to verify a type implements an interface at compile time
- Provides immediate feedback if interface requirements change
- No runtime overhead - declarations are optimized away
- Example: `var _ cli.Command = &MyCommand{}`

## API

### Interfaces

```go
type Command interface {
    Name() string
    Run(cmd *cobra.Command, args []string) error
}

type Options interface {
    BindFlags(fs *pflag.FlagSet)
}

type ComplexOptions interface {
    Options
    Complete(cmd *cobra.Command, args []string) error
    Validate() error
}
```

### Functions

```go
func NewCobraCommand(c Command) *cobra.Command
```

Creates a `*cobra.Command` from a `Command`. Automatically:
- Binds flags if the command implements `Options`
- Calls `Complete()` before `Run()` if command implements `ComplexOptions`
- Calls `Validate()` before `Run()` if command implements `ComplexOptions`

### CommonOptions

```go
type CommonOptions struct {
    Workspace string
    Logger    log.Logger
}

func (c *CommonOptions) BindFlags(fs *pflag.FlagSet)
func (c *CommonOptions) Complete(cmd *cobra.Command, args []string) error
func (c *CommonOptions) Validate() error
```

Provides workspace directory and logger for your command.