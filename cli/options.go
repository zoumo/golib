package cli

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/zoumo/golib/log"
)

var _ ComplexOptions = &CommonOptions{}

// CommonOptions provides common functionality for CLI commands.
//
// It implements the ComplexOptions interface and provides:
//   - Workspace: The working directory for the command
//   - Logger: A structured logger with the command name
//
// Commands can embed this struct to reuse these common features.
//
// Example:
//
//	type MyCommand struct {
//	    cli.CommonOptions
//	    MyField string
//	}
//
//	func (c *MyCommand) Complete(cmd *cobra.Command, args []string) error {
//	    // Always call parent first to initialize Logger and Workspace
//	    return c.CommonOptions.Complete(cmd, args)
//	}
type CommonOptions struct {
	// Workspace is the working directory for the command.
	// Defaults to current working directory if not set.
	Workspace string

	// Logger is a structured logger for the command.
	// Initialized in Complete() with the command name.
	Logger log.Logger
}

// BindFlags implements ComplexOptions interface.
//
// This is a no-op for CommonOptions, allowing derived types
// to add their own flags.
func (c *CommonOptions) BindFlags(fs *pflag.FlagSet) {
}

// Complete implements ComplexOptions interface.
//
// Initializes the Logger with the command name and sets Workspace
// to the current working directory if not already set.
//
// IMPORTANT: Derived types MUST call this method first in their
// own Complete() implementation to ensure proper initialization.
func (c *CommonOptions) Complete(cmd *cobra.Command, args []string) error {
	// Initialize logger with command name for better log filtering
	c.Logger = log.Log.WithName(cmd.Name())

	// Set workspace to current directory if not provided
	if c.Workspace == "" {
		ws, err := os.Getwd()
		if err != nil {
			return err
		}
		c.Workspace = ws
	}
	return nil
}

// Validate implements ComplexOptions interface.
//
// Performs base validation. Derived types should call this method
// first in their own Validate() implementation before adding
// custom validation logic.
func (c *CommonOptions) Validate() error {
	return nil
}
