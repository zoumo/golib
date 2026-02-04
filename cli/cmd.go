package cli

import "github.com/spf13/cobra"

// NewCobraCommand creates a cobra.Command from a Command implementation.
//
// This function automatically handles the full command lifecycle:
//   1. Sets the command name from c.Name()
//   2. Binds flags if the command implements Options interface
//   3. Calls Complete() before Run() if command implements ComplexOptions
//   4. Calls Validate() before Run() if command implements ComplexOptions
//   5. Calls Run() to execute the command
//
// Note: The type assertions work on the concrete type, not the Command interface.
// This means if you return a concrete type that embeds *QueryOptions pointer,
// the promoted methods will be found through method promotion.
func NewCobraCommand(c Command) *cobra.Command {
	cmd := &cobra.Command{
		Use: c.Name(),
		RunE: func(cc *cobra.Command, args []string) error {
			// If command implements ComplexOptions, run the full lifecycle
			if o, ok := c.(ComplexOptions); ok {
				// Complete: initialize resources (logger, workspace, etc.)
				if err := o.Complete(cc, args); err != nil {
					return err
				}
				// Validate: validate the configuration
				if err := o.Validate(); err != nil {
					return err
				}
			}
			// Run: execute the command logic
			return c.Run(cc, args)
		},
	}
	// Bind flags if the command implements Options interface
	if o, ok := c.(Options); ok {
		o.BindFlags(cmd.Flags())
	}
	return cmd
}
