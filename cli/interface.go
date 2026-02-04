package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Command represents a CLI command that can be executed.
//
// A command implements at least Name() and Run() methods. Commands that also implement
// Options or ComplexOptions interfaces get additional functionality like flag binding
// and validation.
//
// Example:
//
//	var _ cli.Command = &MyCommand{}
//
//	type MyCommand struct{}
//
//	func (c *MyCommand) Name() string { return "mycommand" }
//	func (c *MyCommand) Run(cmd *cobra.Command, args []string) error { ... }
type Command interface {
	// Name returns the command's name used for CLI invocation.
	Name() string

	// Run executes the command logic with the parsed cobra.Command and arguments.
	// Called after flags are bound and validated (if implemented).
	Run(cmd *cobra.Command, args []string) error
}

// Options provides flag binding capability for commands.
//
// Implementing this interface allows NewCobraCommand to automatically call BindFlags
// during command construction.
type Options interface {
	// BindFlags binds the command's flags to the pflag.FlagSet.
	// This allows each subcommand to define its own command line flags.
	BindFlags(fs *pflag.FlagSet)
}

// ComplexOptions extends Options with lifecycle methods for initialization and validation.
//
// This interface is for commands that require additional setup after flag parsing
// and validation before execution. The lifecycle is:
//  1. BindFlags() - from Options (if implemented)
//  2. Complete() - initialize resources, parse args, set up logger, etc.
//  3. Validate() - validate the configuration
type ComplexOptions interface {
	Options

	// Complete initializes the command's options after flags have been parsed.
	// This is the right place to:
	//   - Initialize resources (loggers, clients, etc.)
	//   - Load configuration from files
	//   - Parse and validate arguments
	//   - Set up derived state
	Complete(cmd *cobra.Command, args []string) error

	// Validate validates the command's options after completion.
	// Return an error if the configuration is invalid.
	// Called automatically by NewCobraCommand before Run().
	Validate() error
}
