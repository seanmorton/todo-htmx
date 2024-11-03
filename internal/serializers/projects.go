package serializers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/seanmorton/todo-htmx/internal/domain"
)

func ParseProjectForm(project *domain.Project, r *http.Request) error {
	var errMessages []string
	name := r.FormValue("name")
	if name == "" {
		errMessages = append(errMessages, "name is required")
	} else {
		project.Name = name
	}

	if errMessages != nil {
		return errors.New(strings.Join(errMessages, "; "))
	}
	return nil
}
