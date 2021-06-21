package plugin

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Subcommand is an interface that defines the common base for subcommands returned by plugins
type Subcommand interface {
	// Name returns the subcommand's name
	Name() string
	// BindFlags binds the subcommand's flags to the CLI. This allows each subcommand to define its own
	// command line flags.
	BindFlags(fs *pflag.FlagSet)
	// Run runs the subcommand.
	Run(args []string) error
}

// RequiresValidation is a subcommand that requires pre run
type RequiresPreRun interface {
	// PreRun runs before command's Run().
	// It can be used to verify that the command can be run
	PreRun(args []string) error
}

// RequiresValidation is a subcommand that requires post run
type RequiresPostRun interface {
	// PostRun runs after command's Run().
	PostRun(args []string) error
}

type InitHook func(*cobra.Command, Subcommand) error
