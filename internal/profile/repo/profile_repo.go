package repo

import (
	"grpc-crud/internal/profile/model"

	"github.com/jmoiron/sqlx"
)

type ProfileRepo interface {
	Create(p *model.Profile) error
	Get(userID int64) (*model.Profile, error)
	Update(p *model.Profile) error
	Delete(userID int64) error
}

type profileRepo struct {
	db *sqlx.DB
}

func NewProfileRepo(db *sqlx.DB) ProfileRepo {
	return &profileRepo{db: db}
}

func (r *profileRepo) Create(p *model.Profile) error {
	query := `
	INSERT INTO profiles (user_id, name, full_name, email, bio)
	VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(query, p.UserID, p.Name, p.FullName, p.Email, p.Bio)
	return err
}

func (r *profileRepo) Get(userID int64) (*model.Profile, error) {
	query := `
	SELECT user_id, name, full_name, email, bio
	FROM profiles
	WHERE user_id = $1
	`

	var p model.Profile
	err := r.db.Get(&p, query, userID)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *profileRepo) Update(p *model.Profile) error {
	query := `
	UPDATE profiles
	SET name = $1,
	    full_name = $2,
	    bio = $3,
	    updated_at = CURRENT_TIMESTAMP
	WHERE user_id = $4
	`

	_, err := r.db.Exec(query, p.Name, p.FullName, p.Bio, p.UserID)
	return err
}

func (r *profileRepo) Delete(userID int64) error {
	query := `DELETE FROM profiles WHERE user_id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}
