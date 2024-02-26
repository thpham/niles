package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	Blue = lipgloss.Color("#268bd2")
	Yellow = lipgloss.Color("#b58900")
	Red = lipgloss.Color("#dc322f")
	Green = lipgloss.Color("#859900")
	BlueGrey = lipgloss.Color("#778899")
	Cyan = lipgloss.Color("#2aa198")

	// Generic text color styles
	TextRed          = lipgloss.NewStyle().Foreground(Red)
	TextYellow       = lipgloss.NewStyle().Foreground(Yellow)
	TextGreen        = lipgloss.NewStyle().Foreground(Green)
	TextBlue         = lipgloss.NewStyle().Foreground(Blue)
	TextBlueGrey     = lipgloss.NewStyle().Foreground(BlueGrey)
	TextCyan         = lipgloss.NewStyle().Foreground(Cyan)
	TextYellowOnBlue = lipgloss.NewStyle().Foreground(Yellow).Background(Blue).Underline(true)

	// ErrorHelp Box
	//ErrorHelp = lipgloss.NewStyle().Foreground(red).Border(lipgloss.RoundedBorder()).BorderForeground(red)
	ErrorHelp = lipgloss.NewStyle().Foreground(Red)

	// TABS
	Tab = lipgloss.NewStyle().
		Border(TabTabBorder, true).
		BorderForeground(TabColor).
		Padding(0, 1)
	TabColor           = Cyan
	TabActiveTab       = Tab.Copy().Border(TabActiveTabBorder, true).Foreground(Yellow)
	TabActiveTabBorder = lipgloss.ThickBorder()
	TabTabBorder       = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}
	TabGap = Tab.Copy().
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)

	// (S)tats Box Style
	StatsBoxStyle       = lipgloss.NewStyle().Padding(0, 1).BorderStyle(lipgloss.DoubleBorder()).BorderForeground(Cyan)
	StatsSeparatorTitle = lipgloss.NewStyle().Foreground(Yellow).Background(Cyan)

	// JobDetails viewport box
	//JDviewportBox = lipgloss.NewStyle().Border(lipgloss.DoubleBorder(), true, false).BorderForeground(Yellow).Padding(1, 1)
	JDviewportBox = lipgloss.NewStyle()

	MenuBoxStyle      = lipgloss.NewStyle().Padding(1, 1).BorderStyle(lipgloss.DoubleBorder()).BorderForeground(Cyan)
	MenuTitleStyle    = lipgloss.NewStyle().Foreground(Yellow)
	MenuNormalTitle   = lipgloss.NewStyle().Foreground(Cyan)
	MenuSelectedTitle = lipgloss.NewStyle().Foreground(Yellow).Background(Cyan)
	MenuNormalDesc    = lipgloss.NewStyle().Foreground(Yellow).Background(Cyan)
	MenuSelectedDesc  = lipgloss.NewStyle().Foreground(Yellow)

	CountsBox = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 0).BorderForeground(Cyan)

	// Main Window area
	MainWindow = lipgloss.NewStyle().MaxHeight(80)
	HelpWindow = lipgloss.NewStyle().Padding(0, 0).Border(lipgloss.RoundedBorder(), true, false, false).Height(2).MaxHeight(3).BorderForeground(Cyan)

	// Workflow Templates, template not found
	NotFound = lipgloss.NewStyle().Foreground(Red)

	// Workflow tab, infobox
	WorkflowInfoBox         = lipgloss.NewStyle()
	WorkflowInfoInBox       = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(Cyan).MaxHeight(7)
	WorkflowInfoInBottomBox = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(Cyan).MaxHeight(7)

	// Workflow steps
	WorkflowStepBoxStyle        = lipgloss.NewStyle().Padding(1, 2).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(Cyan)
	WorkflowStepExitStatusRed   = lipgloss.NewStyle().Foreground(Red)
	WorkflowStepExitStatusGreen = lipgloss.NewStyle().Foreground(Green)

	//TresBox = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(Cyan).Width(40)
	TresBox = lipgloss.NewStyle()
)
