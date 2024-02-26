package model

import (
	"fmt"
	"log"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/thpham/niles/internal/command"
	"github.com/thpham/niles/internal/keybindings"
	"github.com/thpham/niles/internal/model/tabs/workflowtab"
	"github.com/thpham/niles/internal/model/tabs/workflowtemplate"
	"github.com/thpham/niles/internal/styles"
	"github.com/thpham/niles/internal/table"
)

type errMsg error

type activeTabType interface {
	AdjTableHeight(int, *log.Logger)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		brk                 bool = false
		activeTab           activeTabType
		activeTable         *table.Model
		activeFilter        *textinput.Model
		activeFilterOn      *bool
		activeUserInputsOn  *bool
		activeJDViewport    bool
	)

	// This shortens the testing for table movement keys
	switch m.ActiveTab {
	case tabWorkflows:
		activeTab = &m.WorkflowTab
		activeTable = &m.WorkflowTab.QueueTable
		activeFilter = &m.WorkflowTab.Filter
		activeFilterOn = &m.WorkflowTab.FilterOn
	case tabWorkflowTemplate:
		activeTable = &m.WorkflowTemplateTab.TemplatesTable
	}

	// Filter is turned on, take care of this first
	// TODO: revisit this for filtering on multiple tabs
	switch {
	case activeJDViewport:
		// catch only up/down keys, leave the rest to fallthrough
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keybindings.DefaultKeyMap.Up),
				key.Matches(msg, keybindings.DefaultKeyMap.Down),
				key.Matches(msg, keybindings.DefaultKeyMap.PageUp),
				key.Matches(msg, keybindings.DefaultKeyMap.PageDown):
				m.Log.Printf("VIEWPORT: up/down msg\n")
				var cmd tea.Cmd
				return m, cmd
			}
		}

	case activeFilterOn != nil && *activeFilterOn:
		m.Log.Printf("Filter is ON")
		switch msg := msg.(type) {

		case tea.KeyMsg:
			switch msg.Type {
			// TODO: when filter is set/cleared, trigger refresh with new filtered data
			case tea.KeyEnter:
				// finish & apply entering filter
				*activeFilterOn = false
				brk = true

			case tea.KeyEsc:
				// abort entering filter
				*activeFilterOn = false
				activeFilter.SetValue("")
				brk = true
			}

			if brk {
				activeTable.SetCursor(0)
				activeTab.AdjTableHeight(m.winH, m.Log)
				//m.Log.Printf("ActiveTable = %v\n", activeTable)
				m.Log.Printf("Update: Filter set, setcursor(0), activetable.Cursor==%d\n", activeTable.Cursor())
				switch m.ActiveTab {
					case tabWorkflows:
						return m, nil

					default:
						return m, nil
				}
			}
		}

		tmp, cmd := activeFilter.Update(msg)
		*activeFilter = tmp
		return m, cmd

	case activeUserInputsOn != nil && *activeUserInputsOn:
		m.Log.Printf("UserInputs is ON")
		switch msg := msg.(type) {

		case tea.KeyMsg:
			switch msg.Type {
				// TODO: when filter is set/cleared, trigger refresh with new filtered data
				case tea.KeyEnter:
					// finish & apply entering filter
					*activeUserInputsOn = false
					brk = true

				case tea.KeyEsc:
					// abort entering filter
					*activeUserInputsOn = false
					brk = true

				case tea.KeyUp, tea.KeyDown, tea.KeyTab:
				
			}

			if brk {
				activeTable.SetCursor(0)
				activeTab.AdjTableHeight(m.winH, m.Log)
				m.Log.Printf("Update: Param set, setcursor(0), activetable.Cursor==%d\n", activeTable.Cursor())
				switch m.ActiveTab {

				default:
					return m, nil
				}
			}
		}

		var cmd tea.Cmd
		return m, cmd

	case m.WorkflowTab.MenuOn:
		m.Log.Printf("Update: In Menu\n")
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.WorkflowTab.Menu.SetWidth(msg.Width)
			return m, nil

		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "esc":
				m.WorkflowTab.MenuOn = false
				return m, nil
			case "ctrl+c":
				//m.quitting = true
				m.WorkflowTab.MenuOn = false
				//return m, tea.Quit
				return m, nil

			case "enter":
				m.WorkflowTab.MenuOn = false
				i, ok := m.WorkflowTab.Menu.SelectedItem().(workflowtab.MenuItem)
				if ok {
					m.WorkflowTab.MenuChoice = workflowtab.MenuItem(i)
					if m.WorkflowTab.MenuChoice.GetAction() == "INFO" {
						// TODO: IF Stats==ON AND NxM, turn it of, can't have both on below NxM
						m.WorkflowTab.InfoOn = true
						if m.WorkflowTab.StatsOn && m.Globals.winH < 60 {
							m.Log.Printf("Toggle InfoBox: Height %d too low (<60). Turn OFF Stats\n", m.Globals.winH)
							// We have to turn off stats otherwise screen will break at this Height!
							m.WorkflowTab.StatsOn = false
							// TODO: send a message via ErrMsg
						}
					}
					// host is needed for ssh command
					activeTab.AdjTableHeight(m.winH, m.Log)
					retCmd := m.WorkflowTab.MenuChoice.ExecMenuItem(m.WorkflowTab.SelectedWorkflow, m.Log)
					return m, retCmd
				}
				//return m, tea.Quit
				return m, nil
			}
		}

		var cmd tea.Cmd
		m.WorkflowTab.Menu, cmd = m.WorkflowTab.Menu.Update(msg)
		return m, cmd

	case m.EditTemplate:
		// TODO: move this code to a function/method
		var cmds []tea.Cmd
		var cmd tea.Cmd

		m.Log.Printf("Update: In EditTemplate: %#v\n", msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			m.Log.Printf("Update: m.EditTemplate case tea.KeyMsg\n")
			switch msg.Type {
			case tea.KeyEsc:
				m.EditTemplate = false
				workflowtemplate.EditorKeyMap.DisableKeys()
				tabKeys[m.ActiveTab].SetupKeys()
				//if m.TemplateEditor.Focused() {
				//	m.TemplateEditor.Blur()
				//} else {
				//	m.EditTemplate = false
				//}

			case tea.KeyCtrlS:
				// TODO:
				// 1. Exit editor
				// 2. Save content to file
				// 3. Notify user about generated filename from 2.
				// 4. Submit job
				m.Log.Printf("EditTemplate: Ctrl+s pressed\n")
				m.EditTemplate = false
				workflowtemplate.EditorKeyMap.DisableKeys()
				tabKeys[m.ActiveTab].SetupKeys()
				return m, nil

			case tea.KeyCtrlC:
				return m, tea.Quit

			default:
				if !m.TemplateEditor.Focused() {
					cmd = m.TemplateEditor.Focus()
					cmds = append(cmds, cmd)
				}
			}

		// We handle errors just like any other message
		case errMsg:
			//m.err = msg
			return m, nil
		}

		m.TemplateEditor, cmd = m.TemplateEditor.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)

	}

	switch msg := msg.(type) {

	// TODO: https://pkg.go.dev/github.com/charmbracelet/bubbletea#WindowSizeMsg
	// ToDo:
	// prevent updates for non-selected tabs

	// ERROR msg
	case command.ErrorMsg:
		m.Log.Printf("ERROR msg, from: %s\n", msg.From)
		m.Log.Printf("ERROR msg, original error: %q\n", msg.OrigErr)
		m.Globals.ErrorMsg = msg.OrigErr
		m.Globals.ErrorHelp = msg.ErrHelp
		// cases when this is BAD and we can't continue
		switch msg.From {
		case "GetUserName", "GetUserAssoc":
			return m, tea.Quit
		}
		return m, nil

	// Get initial job template list
	case workflowtemplate.TemplatesListRows:
		m.Log.Printf("Update: Got TemplatesListRows msg: %#v\n", msg)
		if msg != nil {
			// if it's not empty, append to table
			m.WorkflowTemplateTab.TemplatesTable.SetRows(msg)
			m.WorkflowTemplateTab.TemplatesList = msg
		}
		return m, nil

	// getting initial template text
	case workflowtemplate.TemplateText:
		m.Log.Printf("Update: Got TemplateText msg: %#v\n", msg)
		// HERE: we initialize the new textarea editor and flip the EditTemplate switch to ON
		tabKeys[m.ActiveTab].DisableKeys()
		workflowtemplate.EditorKeyMap.SetupKeys()
		m.EditTemplate = true
		m.TemplateEditor = textarea.New()
		m.TemplateEditor.SetWidth(m.winW - 15)
		m.TemplateEditor.SetHeight(m.winH - 15)
		m.TemplateEditor.SetValue(string(msg))
		m.TemplateEditor.Focus()
		m.TemplateEditor.CharLimit = 0
		return m, workflowtemplate.EditorOn()

	// Windows resize
	case tea.WindowSizeMsg:
		m.Log.Printf("Update: got WindowSizeMsg: %d %d\n", msg.Width, msg.Height)
		// TODO: if W<195 || H<60 we can't really run without breaking view, so quit and inform user
		// 187x44 == 13" MacBook Font 14 iTerm (HUGE letters!)
		if msg.Height < 43 || msg.Width < 185 {
			m.Log.Printf("FATAL: Window too small to run without breaking view. Have %dx%d. Need at least 185x43.\n", msg.Width, msg.Height)
			m.Globals.SizeErr = fmt.Sprintf("FATAL: Window too small to run without breaking view. Have %dx%d. Need at least 185x43.\nIncrease your terminal window and/or decrease font size.", msg.Width, msg.Height)
			return m, tea.Quit
		}
		m.winW = msg.Width
		m.winH = msg.Height
		// TODO: set also maxheight/width here on change?
		styles.MainWindow = styles.MainWindow.Height(m.winH - 10)
		styles.MainWindow = styles.MainWindow.Width(m.winW - 15)
		styles.HelpWindow = styles.HelpWindow.Width(m.winW)
		styles.WorkflowStepBoxStyle = styles.WorkflowStepBoxStyle.Width(m.winW - 20)
		// InfoBox
		w := ((m.Globals.winW - 25) / 3) * 3
		styles.WorkflowInfoInBox = styles.WorkflowInfoInBox.Width(w / 3).Height(5)
		styles.WorkflowInfoInBottomBox = styles.WorkflowInfoInBottomBox.Width(w + 4).Height(5)

		// Adjust ALL tables
		m.WorkflowTab.AdjTableHeight(m.winH, m.Log)

		// Adjust StatBoxes
		m.Log.Printf("CTB Width = %d\n", styles.StatsBoxStyle.GetWidth())

	// WorkflowTab update
	case workflowtab.ArgoJSON:
		m.Log.Printf("U(): got ArgoJSON\n")
		
		m.UpdateCnt++
		// if active window != this, don't trigger new refresh
		if m.ActiveTab == tabWorkflows {
			return m, workflowtab.TimedGetArgo(m.Log)
		} else {
			return m, nil
		}

	// Keys pressed
	case tea.KeyMsg:
		switch {

		// Counters
		case key.Matches(msg, keybindings.DefaultKeyMap.Count):
			// Depends at which tab we're at
			m.Log.Printf("Toggle Counters pressed at %d\n", m.ActiveTab)
			switch m.ActiveTab {
			case tabWorkflows:
				m.WorkflowTab.InfoOn = false
				toggleSwitch(&m.WorkflowTab.CountsOn)
			}
			activeTab.AdjTableHeight(m.winH, m.Log)
			return m, nil

		// UP
		// TODO: what if it's a list?
		case key.Matches(msg, keybindings.DefaultKeyMap.Up):
			activeTable.MoveUp(1)
			m.lastKey = "up"

		// DOWN
		case key.Matches(msg, keybindings.DefaultKeyMap.Down):
			activeTable.MoveDown(1)
			m.lastKey = "down"

		// PAGE DOWN
		case key.Matches(msg, keybindings.DefaultKeyMap.PageDown):
			activeTable.MoveDown(activeTable.Height())
			m.lastKey = "pgdown"

		// PAGE UP
		case key.Matches(msg, keybindings.DefaultKeyMap.PageUp):
			activeTable.MoveUp(activeTable.Height())
			m.lastKey = "pgup"

		// 1..6 Tab Selection keys
		case key.Matches(msg, keybindings.DefaultKeyMap.TtabSel):
			k, _ := strconv.Atoi(msg.String())
			tabKeys[m.ActiveTab].DisableKeys()
			m.ActiveTab = uint(k) - 1
			tabKeys[m.ActiveTab].SetupKeys()
			m.lastKey = msg.String()

			// clear error states
			m.Globals.ErrorHelp = ""
			m.Globals.ErrorMsg = nil

			switch m.ActiveTab {
			case tabWorkflows:
				return m, workflowtab.TimedGetArgo(m.Log)
			default:
				return m, nil
			}

		// TAB
		case key.Matches(msg, keybindings.DefaultKeyMap.Tab):
			tabKeys[m.ActiveTab].DisableKeys()
			// switch tab
			m.ActiveTab = (m.ActiveTab + 1) % uint(len(tabs))
			// setup keys
			tabKeys[m.ActiveTab].SetupKeys()
			m.lastKey = "tab"

			// clear error states
			m.Globals.ErrorHelp = ""
			m.Globals.ErrorMsg = nil

			switch m.ActiveTab {
			case tabWorkflows:
				return m, workflowtab.TimedGetArgo(m.Log)
			default:
				return m, nil
			}

		// Shift+TAB
		case key.Matches(msg, keybindings.DefaultKeyMap.ShiftTab):
			tabKeys[m.ActiveTab].DisableKeys()
			// switch tab
			if m.ActiveTab == 0 {
				m.ActiveTab = uint(len(tabs) - 1)
			} else {
				m.ActiveTab -= 1
			}
			// setup keys
			tabKeys[m.ActiveTab].SetupKeys()
			m.lastKey = "tab"

			// clear error states
			m.Globals.ErrorHelp = ""
			m.Globals.ErrorMsg = nil

			switch m.ActiveTab {
			case tabWorkflows:
				return m, workflowtab.TimedGetArgo(m.Log)
			default:
				return m, nil
			}

		// SLASH
		case key.Matches(msg, keybindings.DefaultKeyMap.Slash):
			m.Log.Printf("Filter key pressed\n")
			switch m.ActiveTab {
			case tabWorkflows:
				m.WorkflowTab.FilterOn = true
			}
			activeTab.AdjTableHeight(m.winH, m.Log)
			return m, nil

		// t
		case key.Matches(msg, keybindings.DefaultKeyMap.TimeRange):
			m.Log.Printf("time-range key pressed\n")
			switch m.ActiveTab {
			}
			activeTab.AdjTableHeight(m.winH, m.Log)
			return m, nil

		// ENTER
		case key.Matches(msg, keybindings.DefaultKeyMap.Enter):
			switch m.ActiveTab {

			// Workflow tab: Open Workflow menu
			case tabWorkflows:
				// Check if there is anything in the filtered table and if cursor is on a valid item
				n := m.WorkflowTab.QueueTable.Cursor()
				m.Log.Printf("Update ENTER key @ jobqueue table\n")
				if n == -1 {
					m.Log.Printf("Update ENTER key @ jobqueue table, no jobs selected/empty table\n")
					return m, nil
				}
				// IF Info==ON AND NxM, turn it of, can't have both on below NxM
				if m.WorkflowTab.InfoOn && m.Globals.winH < 60 {
					m.Log.Printf("Toggle MenuBox: Height %d too low (<60). Turn OFF Info\n", m.Globals.winH)
					m.WorkflowTab.InfoOn = false
				}
				// If yes, turn on menu
				m.WorkflowTab.MenuOn = true
				m.WorkflowTab.SelectedWorkflow = m.WorkflowTab.QueueTable.SelectedRow()[0]
				m.WorkflowTab.SelectedWorkflowState = m.WorkflowTab.QueueTable.SelectedRow()[4]
				// Create new menu
				m.WorkflowTab.Menu = workflowtab.NewMenu(m.WorkflowTab.SelectedWorkflowState, m.Log)
				return m, nil

			// Job from Template tab: Open template for editing
			case tabWorkflowTemplate:
				m.Log.Printf("Update ENTER key @ workflowtemplate table\n")
				// return & handle editing there
				if len(m.WorkflowTemplateTab.TemplatesList) != 0 {
					return m, workflowtemplate.GetTemplate(m.WorkflowTemplateTab.TemplatesTable.SelectedRow()[2], m.Log)
				} else {
					return m, nil
				}
			}

		// Refresh the View
		case key.Matches(msg, keybindings.DefaultKeyMap.Refresh):
			switch m.ActiveTab {
			}
			return m, nil

		// Info - toggle on/off
		case key.Matches(msg, keybindings.DefaultKeyMap.Info):
			m.Log.Println("Toggle InfoBox")

			// TODO: IF Stats==ON AND NxM, turn it of, can't have both on below NxM
			if m.WorkflowTab.StatsOn && m.Globals.winH < 60 {
				m.Log.Printf("Toggle InfoBox: Height %d too low (<60). Turn OFF Stats\n", m.Globals.winH)
				// We have to turn off stats otherwise screen will break at this Height!
				m.WorkflowTab.StatsOn = false
				// TODO: send a message via ErrMsg
			}

			m.WorkflowTab.CountsOn = false
			toggleSwitch(&m.WorkflowTab.InfoOn)
			m.WorkflowTab.AdjTableHeight(m.Globals.winH, m.Log)
			return m, nil

		// Stats - toggle on/off
		case key.Matches(msg, keybindings.DefaultKeyMap.Stats):
			switch m.ActiveTab {
			case tabWorkflows:
				m.Log.Printf("WorkflowTab toggle from: %v\n", m.WorkflowTab.StatsOn)
				toggleSwitch(&m.WorkflowTab.StatsOn)
				// IF Info==ON AND NxM, turn it of, can't have both on below NxM
				if m.WorkflowTab.InfoOn && m.Globals.winH < 60 {
					m.Log.Printf("Toggle StatsBox: Height %d too low (<60). Turn OFF Info\n", m.Globals.winH)
					m.WorkflowTab.InfoOn = false
				}
				m.Log.Printf("WorkflowTab toggle to: %v\n", m.WorkflowTab.StatsOn)
			}
			return m, nil

		// QUIT
		case key.Matches(msg, keybindings.DefaultKeyMap.Quit):
			m.Log.Printf("Quit key pressed\n")
			return m, tea.Quit
		}
	}

	return m, nil
}

func toggleSwitch(b *bool) {
	if *b {
		*b = false
	} else {
		*b = true
	}
}
