package command

import (
	"github.com/thpham/niles/internal/config"
)

var cc config.ConfigContainer

// NewCmdCC sets the package variable ConfigContainer to the values from config files/defaults.
// To be used in the package locally without being passed around.
// Not the smartest choice, consider refactoring later.
func NewCmdCC(config config.ConfigContainer) {
	cc = config
}

// usage: cmd := cc.Binpaths["oc"]

type ErrorMsg struct {
	From    string
	ErrHelp string
	OrigErr error
}
