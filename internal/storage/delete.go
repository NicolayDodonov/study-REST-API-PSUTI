package storage

import "database/sql"

func (s *Storage) DeleteUser(uid string) error {
	_, err := s.db.Exec(`
	DELETE FROM users WHERE id = $1`, uid)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) CheckUserRoot(uid string) (bool, error) {
	var userInfo struct {
		UserType string `db:"user_type"`
	}
	err := s.db.Get(&userInfo, `
		SELECT user_type FROM users WHERE id=$1;`,
		uid)
	if err == sql.ErrNoRows {
		// такого пользователя в принципе НЕТУ
		return false, nil
	}
	if err != nil {
		// случилась серверная ошибка, выходим
		return false, err
	}
	// проверяем его user_type
	if userInfo.UserType == "admin" {
		return true, nil
	} else {
		return false, nil
	}
}
