package workflowtab

import (
	"encoding/json"
	"log"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thpham/niles/internal/command"
	"github.com/thpham/niles/internal/config"
)

var (
	cc                config.ConfigContainer
	ArgoCmdSwitches = []string{"-a", "--json"}
)

func NewCmdCC(config config.ConfigContainer) {
	cc = config
}

// Calls `argo` to get job information for Workflow Tab
func GetArgo(t time.Time) tea.Msg {

	var argoJson ArgoJSON

	cmd := cc.Binpaths["argo"]
	out, err := exec.Command(cmd, ArgoCmdSwitches...).Output()
	if err != nil {
		return command.ErrorMsg{
			From:    "GetArgo",
			ErrHelp: "Failed to run argo command, check your niles.conf and set the correct paths there.",
			OrigErr: err,
		}
	}

	err = json.Unmarshal(out, &argoJson)
	if err != nil {
		return command.ErrorMsg{
			From:    "GetArgo",
			ErrHelp: "argo JSON failed to parse, note your argo version and open an issue with us here: https://github.com/thpham/niles/issues/new/choose",
			OrigErr: err,
		}
	}

	return argoJson
}

func TimedGetArgo(l *log.Logger) tea.Cmd {
	l.Printf("TimedGetArgo() start, tick: %d\n", cc.GetTick())
	return tea.Tick(cc.GetTick()*time.Second, GetArgo)
}

func QuickGetArgo(l *log.Logger) tea.Cmd {
	l.Printf("QuickGetArgo() start\n")
	return tea.Tick(0*time.Second, GetArgo)
}
