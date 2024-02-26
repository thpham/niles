package workflowtab

import (
	"log"
	"time"

	"github.com/thpham/niles/internal/command"
	"github.com/thpham/niles/internal/table"
)

var ArgoTabCols = []table.Column{
	{
		Title: "Job ID",
		Width: 10,
	},
	{
		Title: "Job Name",
		Width: 60,
	},
	{
		Title: "Account",
		Width: 10,
	},
	{
		Title: "User Name",
		Width: 20,
	},
	{
		Title: "Job State",
		Width: 10,
	},
	{
		Title: "Priority",
		Width: 10,
	},
}

type ArgoJSON string
type TableRows []table.Row

func (argoJson *ArgoJSON) FilterSqueueTable(f string, l *log.Logger) (*TableRows, *command.ErrorMsg) {
	var (
		argoTabRows      = TableRows{}
		errMsg         *command.ErrorMsg
	)

	t := time.Now()
	
	l.Printf("Filter QUEUE end in %.3f seconds\n", time.Since(t).Seconds())

	return &argoTabRows, errMsg
}
