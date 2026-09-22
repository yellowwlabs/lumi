package auth

import "database/sql"

type SessionRepo struct {
	db *sql.DB
}

func NewSessionRepo(db *sql.DB) *SessionRepo {
	return &SessionRepo{db: db}
}

func (r *SessionRepo) CreateSession(s *Session) error {
	query := `INSERT INTO sessions (id, user_id, token, token_type, ip_address, expires_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(query, s.ID, s.UserID, s.Token, s.TokenType, s.IpAddress, s.ExpiresAt, s.CreatedAt, s.UpdatedAt)
	return err
}

func (r *SessionRepo) GetSessionByToken(token string, tokenType TokenType) (*Session, error) {
	query := `SELECT id, user_id, token, token_type, ip_address, expires_at, created_at, updated_at FROM sessions WHERE token = $1 AND token_type = $2`
	row := r.db.QueryRow(query, token, tokenType)

	var s Session
	err := row.Scan(&s.ID, &s.UserID, &s.Token, &s.TokenType, &s.IpAddress, &s.ExpiresAt, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	return &s, nil
}

func (r *SessionRepo) DeleteSession(id string) error {
	_, err := r.db.Exec(`DELETE FROM sessions WHERE id = $1`, id)
	return err
}

func (r *SessionRepo) DeleteSessionsByUser(userID string, tokenType TokenType) error {
	_, err := r.db.Exec(`DELETE FROM sessions WHERE user_id = $1 AND token_type = $2`, userID, tokenType)
	return err
}
