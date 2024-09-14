package data

import (
	"github.com/seanmorton/todo-htmx/internal/domain"
)

func (d *DB) ListUsers() ([]domain.User, error) {
	var users []domain.User
	rows, err := d.dbConn.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(
			&user.Id, &user.Name, &user.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
