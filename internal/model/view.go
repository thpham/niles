package model

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/thpham/niles/internal/keybindings"
	"github.com/thpham/niles/internal/styles"
	"github.com/thpham/niles/internal/version"
)

// genTabs() generates top tabs
func (m Model) genTabs() string {

	var doc strings.Builder

	tlist := make([]string, len(tabs))
	for i, v := range tabs {
		if i == int(m.ActiveTab) {
			tlist = append(tlist, styles.TabActiveTab.Render(v))
		} else {
			tlist = append(tlist, styles.Tab.Render(v))
		}
	}
	row := lipgloss.JoinHorizontal(lipgloss.Top, tlist...)

	//gap := tabGap.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(row)-2)))
	gap := styles.TabGap.Render(strings.Repeat(" ", max(0, m.winW-lipgloss.Width(row)-2)))
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
	doc.WriteString(row + "\n")

	return doc.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m Model) tabAbout() string {

	s := "Version: " + version.BuildVersion + "\n"
	s += "Commit : " + version.BuildCommit + "\n"

	s += `

	Acknowledgments:
		- John Doe
`

	return s
}

func (m *Model) genTabHelp() string {
	var th string
	switch m.ActiveTab {
	case tabWorkflows:
		th = "List of workflows in the queue"
	case tabWorkflowTemplate:
		th = "Edit and submit one of the workflow templates"
	default:
		th = "Niles"
	}
	return th + "\n"
}

func (m Model) View() string {

	var (
		header     strings.Builder
		MainWindow strings.Builder
	)

	// HEADER / TABS
	header.WriteString(m.genTabs())
	header.WriteString(m.genTabHelp())

	if m.Debug {
		// One debug line
		header.WriteString(fmt.Sprintf("%s Width: %d Height: %d ErrorMsg: %s\n", styles.TextRed.Render("DEBUG ON:"), m.Globals.winW, m.Globals.winH, m.Globals.ErrorMsg))
	}

	if m.Globals.ErrorHelp != "" {
		m.Log.Println("Got error")
		header.WriteString(styles.ErrorHelp.Render(fmt.Sprintf("ERROR: %s", m.Globals.ErrorHelp)))
	} else {
		m.Log.Println("Got NO error")
	}

	// PICK and RENDER ACTIVE TAB
	switch m.ActiveTab {
	case tabWorkflows:
		m.Log.Printf("CALL WorkflowTab.View()\n")
		MainWindow.WriteString(m.WorkflowTab.View(m.Log))

	case tabWorkflowTemplate:
		m.Log.Printf("CALL WorkflowTemplateTab.View()\n")
		MainWindow.WriteString(m.WorkflowTemplateTab.View(m.Log))

	case tabAbout:
		MainWindow.WriteString(m.tabAbout())
		// TODO: default
	}

	return lipgloss.JoinVertical(lipgloss.Left, header.String(), styles.MainWindow.Render(MainWindow.String()), styles.HelpWindow.Render(m.Help.View(keybindings.DefaultKeyMap)))
}
