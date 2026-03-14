package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/seanmorton/todo-htmx/internal/domain"
	"github.com/seanmorton/todo-htmx/internal/serializers"
	"github.com/seanmorton/todo-htmx/internal/templates"
)

func (s *Server) createProject(w http.ResponseWriter, r *http.Request) *httpErr {
	project := domain.Project{}
	validationErr := serializers.ParseProject(&project, r)
	if validationErr != nil {
		return &httpErr{validationErr.Error(), 400, validationErr}
	}

	project, err := s.db.CreateProject(project)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return &httpErr{"project name already exists", 409, err}
		} else {
			return &httpErr{"failed creating project", 500, err}
		}
	}

	s.hxEvent(w, "projectChange")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (s *Server) deleteProject(w http.ResponseWriter, r *http.Request) *httpErr {
	project, fetchErr := s.fetchProject(r)
	if fetchErr != nil {
		return fetchErr
	}

	now := time.Now()
	project.DeletedAt = &now
	res, err := s.db.UpdateProject(project)
	if err != nil {
		return &httpErr{"failed deleting project", 500, err}
	}
	if res == nil {
		return &httpErr{"project not found", 404, nil}
	}

	s.hxEvent(w, "projectChange")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (s *Server) restoreProject(w http.ResponseWriter, r *http.Request) *httpErr {
	project, fetchErr := s.fetchProject(r)
	if fetchErr != nil {
		return fetchErr
	}

	project.DeletedAt = nil
	res, err := s.db.UpdateProject(project)
	if err != nil {
		return &httpErr{"failed restoring project", 500, err}
	}
	if res == nil {
		return &httpErr{"project not found", 404, nil}
	}

	s.hxEvent(w, "projectChange")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (s *Server) projects(w http.ResponseWriter, r *http.Request) *httpErr {
	filters := domain.ProjectFilters{ShowDeleted: r.FormValue("show_deleted") == "on"}
	projects, err := s.db.ListProjects(filters)
	if err != nil {
		return &httpErr{"failed getting projects", 500, err}
	}

	s.hxRender(w, r, templates.Projects(projects, filters))
	return nil
}

func (s *Server) newProject(w http.ResponseWriter, r *http.Request) *httpErr {
	templates.ProjectForm(domain.Project{}).Render(r.Context(), w)
	return nil
}

func (s *Server) getProject(w http.ResponseWriter, r *http.Request) *httpErr {
	project, fetchErr := s.fetchProject(r)
	if fetchErr != nil {
		return fetchErr
	}

	templates.ProjectForm(project).Render(r.Context(), w)
	return nil
}

func (s *Server) updateProject(w http.ResponseWriter, r *http.Request) *httpErr {
	project, fetchErr := s.fetchProject(r)
	if fetchErr != nil {
		return fetchErr
	}
	validationErr := serializers.ParseProject(&project, r)
	if validationErr != nil {
		return &httpErr{validationErr.Error(), 400, validationErr}
	}

	res, err := s.db.UpdateProject(project)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return &httpErr{"project name already exists", 409, err}
		}
		return &httpErr{"failed updating project", 500, err}
	}
	if res == nil {
		return &httpErr{"project not found", 404, nil}
	}

	s.hxEvent(w, "projectChange")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (s *Server) projectRows(w http.ResponseWriter, r *http.Request) *httpErr {
	filters := domain.ProjectFilters{ShowDeleted: r.FormValue("show_deleted") == "on"}
	projects, err := s.db.ListProjects(filters)
	if err != nil {
		return &httpErr{"failed getting projects", 500, err}
	}

	s.hxRender(w, r, templates.ProjectRows(projects))
	return nil
}

func (s *Server) fetchProject(r *http.Request) (domain.Project, *httpErr) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return domain.Project{}, &httpErr{"invalid id", 400, nil}
	}
	project, err := s.db.GetProject(id)
	if err != nil {
		return domain.Project{}, &httpErr{"failed getting project", 500, err}
	}
	if project == nil {
		return domain.Project{}, &httpErr{"project not found", 404, nil}
	}
	return *project, nil
}
