package workflowtab

import (
	"log"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/thpham/niles/internal/generic"
	"github.com/thpham/niles/internal/stats"
	"github.com/thpham/niles/internal/table"
)

type WorkflowTab struct {
	InfoOn                bool
	CountsOn              bool
	StatsOn               bool
	FilterOn              bool
	QueueTable            table.Model
	Queue                 ArgoJSON
	Filter                textinput.Model
	SelectedWorkflow      string
	SelectedWorkflowState string
	MenuOn                bool
	MenuChoice            MenuItem
	Menu                  list.Model
	Stats
	Breakdowns
}

type Stats struct {
	// TODO: also perhaps: count by user? account?
	StateCnt map[string]uint
	AvgWait  time.Duration
	MinWait  time.Duration
	MaxWait  time.Duration
	MedWait  time.Duration
	AvgRun   time.Duration
	MinRun   time.Duration
	MaxRun   time.Duration
	MedRun   time.Duration
}

type Breakdowns struct {
	Top5user   generic.CountItemSlice
	Top5acc    generic.CountItemSlice
	JobPerQos  generic.CountItemSlice
	JobPerPart generic.CountItemSlice
}

func (t *WorkflowTab) AdjTableHeight(h int, l *log.Logger) {
	l.Printf("FixTableHeight(%d) from %d\n", h, t.QueueTable.Height())
	if t.InfoOn || t.CountsOn || t.FilterOn {
		t.QueueTable.SetHeight(h - 30)
	} else {
		t.QueueTable.SetHeight(h - 15)
	}
	l.Printf("FixTableHeight to %d\n", t.QueueTable.Height())
}

func (t *WorkflowTab) GetStatsFiltered(l *log.Logger) {

	top5user := generic.CountItemMap{}
	top5acc := generic.CountItemMap{}
	jpq := generic.CountItemMap{}
	jpp := generic.CountItemMap{}

	t.Stats.StateCnt = map[string]uint{}
	tmp := []time.Duration{}
	tmpRun := []time.Duration{}
	t.AvgWait = 0
	t.MedWait = 0

	// sort & filter breakdowns
	t.Breakdowns.Top5user = generic.Top5(generic.SortItemMapBySel("Count", &top5user))
	t.Breakdowns.Top5acc = generic.Top5(generic.SortItemMapBySel("Count", &top5acc))
	t.Breakdowns.JobPerPart = generic.SortItemMapBySel("Count", &jpp)
	t.Breakdowns.JobPerQos = generic.SortItemMapBySel("Count", &jpq)

	//l.Printf("TOP5USER: %#v\n", t.Breakdowns.Top5user)
	//l.Printf("TOP5ACC: %#v\n", t.Breakdowns.Top5acc)
	//l.Printf("JobPerQos: %#v\n", t.Breakdowns.JobPerQos)
	//l.Printf("JobPerPart: %#v\n", t.Breakdowns.JobPerPart)

	t.MedWait, t.MinWait, t.MaxWait = stats.Median(tmp)
	t.MedRun, t.MinRun, t.MaxRun = stats.Median(tmpRun)
	t.AvgWait = stats.Avg(tmp)
	t.AvgRun = stats.Avg(tmpRun)

	l.Printf("GetStatsFiltered end\n")
}
