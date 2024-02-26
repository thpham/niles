package model

import (
	"log"

	"github.com/charmbracelet/bubbles/help"
	"github.com/thpham/niles/internal/config"
	"github.com/thpham/niles/internal/model/tabs/abouttab"
	"github.com/thpham/niles/internal/model/tabs/workflowtab"
	"github.com/thpham/niles/internal/model/tabs/workflowtemplate"
)

const (
	tabWorkflows = iota
	tabWorkflowTemplate
	tabAbout
)

var tabs = []string{
	"Workflows",
	"Workflow Template",
	"About",
}

type ActiveTabKeys interface {
	SetupKeys()
	DisableKeys()
}

var tabKeys = []ActiveTabKeys{
	&workflowtab.KeyMap,
	&workflowtemplate.KeyMap,
	&abouttab.KeyMap,
}

// TODO: in structures below:
// - make embedding and accessing leafs uniform (shorthand notation vs Full path)
type Model struct {
	Globals
	workflowtab.WorkflowTab
	workflowtemplate.WorkflowTemplateTab
}

type Globals struct {
	ActiveTab uint
	UpdateCnt uint64
	Debug     bool
	DebugMsg  string
	lastKey   string
	winW      int
	winH      int
	Log       *log.Logger
	Help      help.Model
	UserName  string
	UAccounts []string
	config.ConfigContainer
	ErrorMsg  error
	ErrorHelp string
	SizeErr   string
}
