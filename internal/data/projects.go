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

func (d *DB) GetProject(id int64) (*domain.Project, error) {
	var project domain.Project
	err := d.dbConn.QueryRow(
		"SELECT id, name, created_at, deleted_at FROM projects WHERE id = ?", id,
	).Scan(&project.Id, &project.Name, &project.CreatedAt, &project.DeletedAt)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (d *DB) UpdateProject(project domain.Project) (*domain.Project, error) {
	res, err := d.dbConn.Exec(
		"UPDATE projects SET name = ?, deleted_at = ? WHERE id = ?",
		project.Name, project.DeletedAt, project.Id,
	)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, nil
	}
	return &project, nil
}

func (d *DB) ListProjects(filters domain.ProjectFilters) ([]domain.Project, error) {
	query := "SELECT id, name, created_at, deleted_at FROM projects"
	if !filters.ShowDeleted {
		query += " WHERE deleted_at IS NULL"
	}
	query += " ORDER BY name"

	rows, err := d.dbConn.Query(query)
	if err != nil {
		return nil, err
	}

	var projects []domain.Project
	defer rows.Close()
	for rows.Next() {
		var project domain.Project
		if err := rows.Scan(&project.Id, &project.Name, &project.CreatedAt, &project.DeletedAt); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}
