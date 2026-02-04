// Copyright 2025 jim.zoumo@gmail.com
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

// Package simple demonstrates basic usage of the cli package.
package simple

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zoumo/golib/cli"
)

var _ cli.Command = &GreetCommand{}

// GreetCommand is the simplest command - no flags, just positional args.
type GreetCommand struct{}

// Name returns the command name.
func (c *GreetCommand) Name() string {
	return "greet"
}

// Run executes the command.
func (c *GreetCommand) Run(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		fmt.Println("Hello, World!")
		return nil
	}
	fmt.Printf("Hello, %s!\n", args[0])
	return nil
}

// NewGreetCommand creates a new greet command.
func NewGreetCommand() cli.Command {
	return &GreetCommand{}
}

var (
	_ cli.Command = &EchoCommand{}
	_ cli.Options = &EchoCommand{}
)

// EchoCommand demonstrates a command with flags.
type EchoCommand struct {
	count int
}

// Name returns the command name.
func (c *EchoCommand) Name() string {
	return "echo"
}

// BindFlags binds the command flags.
func (c *EchoCommand) BindFlags(fs *pflag.FlagSet) {
	fs.IntVar(&c.count, "count", 1, "Number of times to echo")
}

// Run executes the command.
func (c *EchoCommand) Run(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: echo <text>")
	}
	for i := 0; i < c.count; i++ {
		fmt.Println(args[0])
	}
	return nil
}

// NewEchoCommand creates a new echo command.
func NewEchoCommand() cli.Command {
	return &EchoCommand{}
}
