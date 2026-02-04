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

// Package complex demonstrates advanced usage with CommonOptions.
package complex

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

// NewQueryCommand returns a command with default options.
func NewQueryCommand() cli.Command {
	return &QueryCommand{
		QueryOptions: &QueryOptions{},
	}
}
