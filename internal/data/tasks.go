package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/seanmorton/todo-htmx/internal/domain"
	"github.com/seanmorton/todo-htmx/pkg"
)

func (d *DB) CreateTask(task domain.Task) (domain.Task, error) {
	result, err := d.dbConn.Exec(
		"INSERT INTO tasks(title, project_id, assignee_id, description, due_date, completed_at, recur_policy) VALUES(?, ?, ?, ?, ?, ?, ?)",
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
	row := d.dbConn.QueryRow("SELECT * FROM tasks WHERE id = ?", id)
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

func (d *DB) QueryTasks(params map[string]any, nextMonthOnly bool) ([]domain.Task, error) {
	query := "SELECT * FROM tasks"
	args := make([]any, 0, len(params))
	if len(params) > 0 {
		query += " WHERE"
		count := 0

		for col, val := range params {
			col := pkg.CamelToSnake(col)
			var clause string

			switch val {
			case nil:
				clause = fmt.Sprintf("%s IS NULL", col)
			case "NOT NULL":
				clause = fmt.Sprintf("%s IS NOT NULL", col)
			default:
				clause = fmt.Sprintf("%s = ?", col)
				args = append(args, val)
			}
			if count < 1 {
				query += fmt.Sprintf(" %s", clause)
			} else {
				query += fmt.Sprintf(" AND %s", clause)
			}
			count++
		}
	}

	if nextMonthOnly {
		nextMonth := time.Now().AddDate(0, 1, 0)
		query += fmt.Sprintf(" AND (due_date < '%s' OR due_date IS NULL)", pkg.DateStr(&nextMonth))
	}

	if params["completed_at"] == "NOT NULL" {
		query += " ORDER BY completed_at DESC"
	} else {
		query += " ORDER BY COALESCE(due_date, '9999-9-9') ASC, created_at DESC"
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
