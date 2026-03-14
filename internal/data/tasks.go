package data

import (
	"database/sql"
	"time"

	"github.com/seanmorton/todo-htmx/internal/domain"
	"github.com/seanmorton/todo-htmx/pkg"
)

func (d *DB) CreateTask(task domain.Task) (domain.Task, error) {
	result, err := d.dbConn.Exec(
		`INSERT INTO tasks(title, project_id, assignee_id, description, due_date, completed_at, recur_policy)
		VALUES(?, ?, ?, ?, ?, ?, ?)`,
		task.Title, task.ProjectId, task.AssigneeId, task.Description, task.DueDate, task.CompletedAt, task.RecurPolicy,
	)
	if err != nil {
		return domain.Task{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return domain.Task{}, err
	}
	task.Id = id
	return task, nil
}

func (d *DB) GetTask(id int64) (*domain.Task, error) {
	var task domain.Task
	row := d.dbConn.QueryRow(
		`SELECT id, project_id, assignee_id, title, description, due_date, completed_at, recur_policy, created_at
		FROM tasks
		WHERE id = ?`,
		id,
	)
	err := row.Scan(
		&task.Id, &task.ProjectId, &task.AssigneeId,
		&task.Title, &task.Description, &task.DueDate, &task.CompletedAt,
		&task.RecurPolicy, &task.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &task, err
}

func (d *DB) UpdateTask(task domain.Task) (*domain.Task, error) {
	res, err := d.dbConn.Exec(
		`UPDATE tasks
     SET title = ?, project_id = ?, assignee_id = ?, description = ?, due_date = ?, recur_policy = ?, completed_at = ?
     WHERE id = ?`,
		task.Title, task.ProjectId, task.AssigneeId, task.Description, task.DueDate, task.RecurPolicy, task.CompletedAt, task.Id,
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
	return &task, err
}

func (d *DB) DeleteTask(id int64) (bool, error) {
	res, err := d.dbConn.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return false, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected == 0 {
		return false, nil

	}
	return true, nil
}

func (d *DB) QueryTasks(filter domain.TaskFilters) ([]domain.Task, error) {
	query := `SELECT tasks.id, tasks.project_id, tasks.assignee_id, tasks.title, tasks.description,
		tasks.due_date, tasks.completed_at, tasks.recur_policy, tasks.created_at
		FROM tasks
		JOIN projects ON tasks.project_id = projects.id
		WHERE projects.deleted_at IS NULL`
	var args []any

	if filter.ProjectID != nil {
		query += " AND tasks.project_id = ?"
		args = append(args, *filter.ProjectID)
	}
	if filter.AssigneeID != nil {
		query += " AND tasks.assignee_id = ?"
		args = append(args, *filter.AssigneeID)
	}
	if filter.Completed {
		query += " AND tasks.completed_at IS NOT NULL"
	} else {
		query += " AND tasks.completed_at IS NULL"
	}
	if filter.Search != nil {
		query += " AND (tasks.title LIKE ? OR tasks.description LIKE ?)"
		wildcard := "%" + *filter.Search + "%"
		args = append(args, wildcard, wildcard)
	}
	if filter.NextMonthOnly {
		nextMonth := time.Now().AddDate(0, 1, 0)
		query += " AND (tasks.due_date < ? OR tasks.due_date IS NULL)"
		args = append(args, pkg.DateStr(&nextMonth))
	}

	if filter.Completed {
		query += " ORDER BY tasks.completed_at DESC"
	} else {
		query += " ORDER BY COALESCE(tasks.due_date, '9999-9-9') ASC, tasks.created_at DESC"
	}

	var tasks []domain.Task
	rows, err := d.dbConn.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var task domain.Task
		if err := rows.Scan(
			&task.Id, &task.ProjectId, &task.AssigneeId,
			&task.Title, &task.Description, &task.DueDate, &task.CompletedAt,
			&task.RecurPolicy, &task.CreatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
