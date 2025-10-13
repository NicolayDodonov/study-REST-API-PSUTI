package storage

import "database/sql"

func (s *Storage) Login(username, password string) (string, error) {
	var Uid struct {
		Uid string `db:"id"`
	}
	err := s.db.Get(&Uid, `
		SELECT id FROM users 
		WHERE login = $1 and password = $2
	`,
		username,
		password)

	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return Uid.Uid, nil
}
