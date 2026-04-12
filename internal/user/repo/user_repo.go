package repo

import (
	"grpc-crud/internal/user/model"

	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	Create(user *model.User) error
	GetByID(id int64) (*model.User, error)
	Update(user *model.User) error
	Delete(id int64) error
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepo {
	return &userRepo{db: db}
}

// Create inserts a new user record into the database and returns the generated ID.
func (r *userRepo) Create(user *model.User) error {
	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRow(query, user.Name, user.Email).Scan(&user.ID)
}

// GetByID fetches a user by primary key from the users table.
func (r *userRepo) GetByID(id int64) (*model.User, error) {
	user := &model.User{}
	err := r.db.Get(user, "SELECT id, name, email FROM users WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Update modifies an existing user's name and email.
func (r *userRepo) Update(user *model.User) error {
	query := `UPDATE users SET name=$1, email=$2 WHERE id=$3`
	_, err := r.db.Exec(query, user.Name, user.Email, user.ID)
	return err
}

// Delete removes a user row from the users table.
func (r *userRepo) Delete(id int64) error {
	query := `DELETE FROM users WHERE id=$1`
	_, err := r.db.Exec(query, id)
	return err
}
