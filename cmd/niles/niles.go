package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/thpham/niles/internal/cmdline"
	"github.com/thpham/niles/internal/command"
	"github.com/thpham/niles/internal/config"
	"github.com/thpham/niles/internal/logger"
	"github.com/thpham/niles/internal/model"
	"github.com/thpham/niles/internal/model/tabs/workflowtab"
	"github.com/thpham/niles/internal/model/tabs/workflowtemplate"
	"github.com/thpham/niles/internal/styles"
	"github.com/thpham/niles/internal/table"
	"github.com/thpham/niles/internal/version"
)

func main() {

	var (
		debugSet bool = false
		args *cmdline.CmdArgs
	)

	fmt.Printf("Welcome to Niles!\n\n")

	cc := config.NewConfigContainer()
	err := cc.GetConfig()
	if err != nil {
		log.Printf("ERROR: parsing config files: %s\n", err)
	}

	args, err = cmdline.NewCmdArgs()
	if err != nil {
		log.Fatalf("ERROR: parsing cmdline args: %s\n", err)
	}

	if *args.Version {
		version.DumpVersion()
		os.Exit(0)
	}

	// TODO: JFT We have the CMDline switches and config, now overwrite/append what's changed
	//log.Println(cc.DumpConfig())

	log.Printf("INFO: %s\n", cc.DumpConfig())
	// TODO: this is ugly, but quick. Rework, use model...
	command.NewCmdCC(*cc)
	workflowtab.NewCmdCC(*cc)

	// TODO: move all this away to view/styles somewhere...
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(styles.BlueGrey).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Background(styles.Blue).
		Foreground(styles.Yellow).
		Bold(false)

	// Filter TextInput
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 30
	ti.Width = 30

	// logging
	debugSet, l := logger.SetupLogger()

	// setup help
	hlp := help.New()
	hlp.Styles.ShortKey = styles.TextYellow
	hlp.Styles.ShortDesc = styles.TextCyan
	hlp.Styles.ShortSeparator = styles.TextBlueGrey

	/// JD viewport
	vp := viewport.New(10, 10)
	vp.Style = styles.JDviewportBox

	m := model.Model{
		Globals: model.Globals{
			Help:            hlp,
			ActiveTab:       0,
			Log:             l,
			Debug:           debugSet,
			ConfigContainer: *cc,
		},
		WorkflowTab: workflowtab.WorkflowTab{
			QueueTable: table.New(table.WithColumns(workflowtab.ArgoTabCols), table.WithRows(workflowtab.TableRows{}), table.WithStyles(s)),
			Filter: ti,
		},
		WorkflowTemplateTab: workflowtemplate.WorkflowTemplateTab{
			EditTemplate: false,
			NewWorkflowScript: "",
			TemplatesTable: table.New(
				table.WithColumns(workflowtemplate.TemplatesListCols),
				table.WithRows(workflowtemplate.TemplatesListRows{}),
				table.WithStyles(s),
			),
		},
	}

	p := tea.NewProgram(tea.Model(m), tea.WithAltScreen())
	ret, err := p.Run()
	if err != nil {
		fmt.Printf("Error starting program: %s", err)
		os.Exit(1)
	}

	if retMod, ok := ret.(model.Model); ok && retMod.Globals.SizeErr != "" {
		fmt.Printf("%s\n", retMod.Globals.SizeErr)
	}
	fmt.Printf("Goodbye!\n")
}