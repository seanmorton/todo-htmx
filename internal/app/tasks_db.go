package app

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/seanmorton/todo-htmx/internal/domain"
)

type TasksDB struct {
	db *sql.DB
}

func NewTasksDB(db *sql.DB) TasksDB {
	return TasksDB{db: db}
}

func (t *TasksDB) Create(task domain.Task) (domain.Task, error) {
	result, err := t.db.Exec(
		"INSERT INTO tasks(title, description, due_date, completed_at, recur_policy) VALUES(?, ?, ?, ?, ?)",
		task.Title, task.Description, task.DueDate, task.CompletedAt, task.RecurPolicy,
	)
	if err != nil {
		return domain.Task{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return domain.Task{}, err
	}
	task.ID = id
	return task, nil
}

func (t *TasksDB) Get(id int64) (domain.Task, error) {
	var task domain.Task
	row := t.db.QueryRow("SELECT * FROM tasks WHERE id = ?", id)
	err := row.Scan(
		&task.ID, &task.Title, &task.Description,
		&task.Assignee, &task.DueDate, &task.CompletedAt,
		&task.RecurPolicy, &task.CreatedAt,
	)
	return task, err
}

func (t *TasksDB) Update(task domain.Task) (domain.Task, error) {
	_, err := t.db.Exec(
		`UPDATE tasks
     SET title = ?, description = ?, due_date = ?, completed_at = ?, recur_policy = ?
     WHERE id = ?`,
		task.Title, task.Description, task.DueDate, task.CompletedAt, task.RecurPolicy, task.ID,
	)
	return task, err
}

func (t *TasksDB) Delete(id int64) error {
	res, err := t.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == int64(1) {
		return errors.New("task not found")
	}
	return err
}

func (t *TasksDB) Query(filters map[string]any) ([]domain.Task, error) {
	query := "SELECT * FROM tasks"
	args := make([]any, 0, len(filters))
	if len(filters) > 0 {
		query += " WHERE"
		count := 0
		for col, val := range filters {
			var clause string
			if val == nil {
				clause = fmt.Sprintf("%s IS NULL", col)
			} else {
				clause = fmt.Sprintf("%s = ?", col)
			}
			args = append(args, val)
			if count < 1 {
				query += fmt.Sprintf(" %s", clause)
			} else {
				query += fmt.Sprintf(" AND %s", clause)
			}
			count++
		}
	}
	query += " ORDER BY COALESCE(due_date, '9999-9-9') ASC, created_at DESC"
	fmt.Println(query)

	var tasks []domain.Task
	rows, err := t.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var task domain.Task
		if err := rows.Scan(
			&task.ID, &task.Title, &task.Description,
			&task.Assignee, &task.DueDate, &task.CompletedAt,
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
