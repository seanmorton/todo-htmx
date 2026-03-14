package handlers

import (
	"net/http"
	"strconv"
	"strings"

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
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return &httpErr{"invalid id", 400, nil}
	}

	deleted, err := s.db.DeleteProject(id)
	if err != nil {
		return &httpErr{"failed deleting project", 500, err}
	}
	if !deleted {
		return &httpErr{"project not found", 404, nil}
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (s *Server) projects(w http.ResponseWriter, r *http.Request) *httpErr {
	projects, err := s.db.ListProjects()
	if err != nil {
		return &httpErr{"failed getting projects", 500, err}
	}

	s.hxRender(w, r, templates.Projects(projects))
	return nil
}

func (s *Server) newProject(w http.ResponseWriter, r *http.Request) *httpErr {
	templates.ProjectForm(domain.Project{}).Render(r.Context(), w)
	return nil
}

func (s *Server) getProject(w http.ResponseWriter, r *http.Request) *httpErr {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return &httpErr{"invalid id", 400, nil}
	}

	project, err := s.db.GetProject(id)
	if err != nil {
		return &httpErr{"project not found", 404, err}
	}

	templates.ProjectForm(*project).Render(r.Context(), w)
	return nil
}

func (s *Server) updateProject(w http.ResponseWriter, r *http.Request) *httpErr {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return &httpErr{"invalid id", 400, nil}
	}

	project := domain.Project{}
	validationErr := serializers.ParseProject(&project, r)
	if validationErr != nil {
		return &httpErr{validationErr.Error(), 400, validationErr}
	}

	updated, err := s.db.UpdateProject(id, project.Name)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return &httpErr{"project name already exists", 409, err}
		}
		return &httpErr{"failed updating project", 500, err}
	}
	if !updated {
		return &httpErr{"project not found", 404, nil}
	}

	s.hxEvent(w, "projectChange")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (s *Server) projectRows(w http.ResponseWriter, r *http.Request) *httpErr {
	projects, err := s.db.ListProjects()
	if err != nil {
		return &httpErr{"failed getting projects", 500, err}
	}

	s.hxRender(w, r, templates.ProjectRows(projects))
	return nil
}
