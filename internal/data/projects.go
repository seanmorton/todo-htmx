package data

import (
	"github.com/seanmorton/todo-htmx/internal/domain"
)

func (d *DB) CreateProject(project domain.Project) (domain.Project, error) {
	result, err := d.dbConn.Exec("INSERT INTO projects(name) VALUES(?)", project.Name)
	if err != nil {
		return domain.Project{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return domain.Project{}, err
	}
	project.Id = id
	return project, nil
}

func (d *DB) DeleteProject(id int64) (bool, error) {
	result, err := d.dbConn.Exec("DELETE FROM projects WHERE id = ?", id)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected == 0 {
		return false, nil

	}
	return true, nil
}

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
