package storage

import "database/sql"

func (s *Storage) Login(username, password string) (any, error) {
	var Uid struct {
		Uid string `db:"uid"`
	}
	err := s.db.Get(Uid, `
		SELECT id FROM users 
		WHERE login = &1 and password = &2
	`,
		username,
		password)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return Uid, nil
}
