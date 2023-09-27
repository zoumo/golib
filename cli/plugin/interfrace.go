// Copyright 2023 jim.zoumo@gmail.com
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
