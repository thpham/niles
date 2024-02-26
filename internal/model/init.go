package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/thpham/niles/internal/model/tabs/workflowtab"
	"github.com/thpham/niles/internal/model/tabs/workflowtemplate"
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		workflowtab.QuickGetArgo(m.Log),
		workflowtemplate.GetTemplateList(m.Globals.ConfigContainer.TemplateDirs, m.Log),
	)
}
