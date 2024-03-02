/*
Using the bottom up approach with

			DB -> repository -> service -> handler
*/

package user

import (
	"context"
	"database/sql"
)

// DBTX is an interface that abstracts database operations, providing flexibility to use either a database connection or a transaction.
type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

// repository is an abstraction over the database, facilitating decoupling of databases.
type repository struct {
	db DBTX
}

// NewRepository creates and returns a new repository with the provided database interface.
func NewRepository(db DBTX) Repository {
	return &repository{db: db}
}

// CreateUser adds a user's information into the database and returns the ID given to the user.
func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	var lastReturnedId int

	// Insert the user into the database
	query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id"
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).Scan(&lastReturnedId)
	if err != nil {
		return &User{}, err
	}

	user.ID = int64(lastReturnedId)
	return user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	userFromDb := User{}

	query := "SELECT id, username, password FROM users WHERE email=$1"
	// error is here
	err := r.db.QueryRowContext(ctx, query, email).Scan(&userFromDb.ID, &userFromDb.Username, &userFromDb.Password)
	if err != nil {
		return &User{}, err
	}

	return &userFromDb, nil
}
