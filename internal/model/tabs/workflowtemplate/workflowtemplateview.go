package workflowtemplate

import (
	"log"
	"strings"

	"github.com/thpham/niles/internal/styles"
)

func (wft *WorkflowTemplateTab) tabWorkflowTemplate() string {

	if wft.EditTemplate {
		return wft.TemplateEditor.View()
	} else {
		if len(wft.TemplatesList) == 0 {
			return styles.NotFound.Render("\nNo templates found!\n")
		} else {
			return wft.TemplatesTable.View()
		}
	}
}

func (wft *WorkflowTemplateTab) View(l *log.Logger) string {
	var (
		MainWindow strings.Builder
	)

	MainWindow.WriteString(wft.tabWorkflowTemplate())

	return MainWindow.String()
}
