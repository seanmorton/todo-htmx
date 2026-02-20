package serializers

import (
	"net/http"

	"github.com/seanmorton/todo-htmx/internal/domain"
)

func ParseProject(project *domain.Project, r *http.Request) error {
	var errs []string
	project.Name = parseString(r, "name", &errs)
	return validationErr(errs)
}
