package auth

import "database/sql"

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(user *User) error {
	query := `INSERT INTO users (id, display_name, username, password_hash, email, email_verified, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(query, user.ID, user.DisplayName, user.Username, user.PasswordHash, user.Email, user.EmailVerified, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *UserRepo) GetUserByUsername(username string) (*User, error) {
	query := `SELECT id, display_name, username, password_hash, email, email_verified, created_at, updated_at FROM users WHERE username = $1`
	row := r.db.QueryRow(query, username)

	var user User
	err := row.Scan(&user.ID, &user.DisplayName, &user.Username, &user.PasswordHash, &user.Email, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) GetUserByID(id string) (*User, error) {
	query := `SELECT id, display_name, username, password_hash, email, email_verified, created_at, updated_at FROM users WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var user User
	err := row.Scan(&user.ID, &user.DisplayName, &user.Username, &user.PasswordHash, &user.Email, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, display_name, username, password_hash, email, email_verified, created_at, updated_at FROM users WHERE email = $1`
	row := r.db.QueryRow(query, email)

	var user User
	err := row.Scan(&user.ID, &user.DisplayName, &user.Username, &user.PasswordHash, &user.Email, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) UpdateUser(user *User) error {
	query := `UPDATE users SET display_name = $1, username = $2, password_hash = $3, email = $4, email_verified = $5, updated_at = $6 WHERE id = $7`
	_, err := r.db.Exec(query, user.DisplayName, user.Username, user.PasswordHash, user.Email, user.EmailVerified, user.UpdatedAt, user.ID)
	return err
}

func (r *UserRepo) DeleteUser(userID string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}
