package storage

import "database/sql"

func (s *Storage) DeleteUser(uid string) error {
	_, err := s.db.Exec("DELETE FROM user WHERE id=?", uid)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) CheckToken(token string) (bool, error) {
	// Токен == id пользователя в системе
	// мне просто было лень делать нормальные токены

	// проверяем есть ли пользователь с таким токен-id в системе
	var userInfo struct {
		UserType string `db:"user_type"`
	}
	err := s.db.Get(&userInfo, `
		SELECT user_type FROM user WHERE id=?`,
		token)
	if err == sql.ErrNoRows {
		// такого пользователя в принципе НЕТУ
		return false, nil
	}
	if err != nil {
		// случилась серверная ошибка, выходим
		return false, err
	}
	// проверяем его user_type
	if userInfo.UserType == "аdmin" {
		return true, nil
	} else {
		return false, nil
	}
}
