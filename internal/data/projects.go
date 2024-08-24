package data

import (
	"github.com/seanmorton/todo-htmx/internal/domain"
)

func (d *DB) ListProjects() ([]domain.Project, error) {
	var projects []domain.Project
	rows, err := d.dbConn.Query("SELECT * FROM projects")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var project domain.Project
		if err := rows.Scan(
			&project.Id, &project.Name, &project.CreatedAt,
		); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}
