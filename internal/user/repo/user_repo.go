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
	Find(email string) (*model.User, error)
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepo {
	return &userRepo{db: db}
}

// Create inserts a new user record into the database and returns the generated ID.
func (r *userRepo) Create(user *model.User) error {
	query := `INSERT INTO users (user_name, email, password) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRow(query, user.UserName, user.Email, user.Password).Scan(&user.ID)
}

// GetByID fetches a user by primary key from the users table.
func (r *userRepo) GetByID(id int64) (*model.User, error) {
	user := &model.User{}
	query := "SELECT id, user_name, email, password FROM users WHERE id=$1"
	err := r.db.Get(user, query, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Update modifies an existing user's profile fields.
func (r *userRepo) Update(user *model.User) error {
	query := `UPDATE users SET user_name=$1, email=$2, password=$3 WHERE id=$4`
	_, err := r.db.Exec(query, user.UserName, user.Email, user.Password, user.ID)
	return err
}

// Delete removes a user row from the users table.
func (r *userRepo) Delete(id int64) error {
	query := `DELETE FROM users WHERE id=$1`
	_, err := r.db.Exec(query, id)
	return err
}

// Find retrieves a user by email from the users table.
func (r *userRepo) Find(email string) (*model.User, error) {
	user := &model.User{}
	query := "SELECT id, user_name, email, password FROM users WHERE email=$1"
	err := r.db.Get(user, query, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
