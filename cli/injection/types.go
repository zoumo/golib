package injection

import (
	"github.com/go-logr/logr"
)

// RequiresValidation indicate the subcommand requires a logger
type RequiresLogger interface {
	InjectLogger(obj logr.Logger)
}

// RequiresValidation indicate the subcommand requires workspace
type RequiresWorkspace interface {
	InjectWorkspace(ws string)
}
