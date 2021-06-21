package plugin

import (
	"github.com/spf13/cobra"
)

func RunSubcommand(cmd Subcommand, args []string) error {
	if preRun, ok := cmd.(RequiresPreRun); ok {
		if err := preRun.PreRun(args); err != nil {
			return err
		}
	}
	if err := cmd.Run(args); err != nil {
		return err
	}
	if postRun, ok := cmd.(RequiresPostRun); ok {
		if err := postRun.PostRun(args); err != nil {
			return err
		}
	}
	return nil
}

func NewCobraSubcommand(subcmd Subcommand, hooks ...InitHook) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:          subcmd.Name(),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunSubcommand(subcmd, args)
		},
	}

	for _, hook := range hooks {
		if err := hook(cmd, subcmd); err != nil {
			return nil, err
		}
	}

	subcmd.BindFlags(cmd.Flags())

	return cmd, nil
}

func NewCobraSubcommandOrDie(subcmd Subcommand, hooks ...InitHook) *cobra.Command {
	cmd, err := NewCobraSubcommand(subcmd, hooks...)
	if err != nil {
		panic(err)
	}
	return cmd
}
