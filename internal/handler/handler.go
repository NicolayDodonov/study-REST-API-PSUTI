package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"study-REST-API-PSUTI/internal/logger"
	"study-REST-API-PSUTI/internal/model"
	"study-REST-API-PSUTI/internal/storage"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var key = []byte("I shouldn't store the token here.")

type Handler struct {
	s      *storage.Storage
	logger *logger.Logger
}

func New(s *storage.Storage, logger *logger.Logger) *Handler {
	return &Handler{s: s, logger: logger}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	// Читаем сообщение
	var msg model.Login
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error(err.Error())
		return
	}

	// Проверяем ключевые поля на пустоту
	if msg.Username == "" || msg.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error("login or password is empty")
		return
	}

	// Проводим логирование пользователя
	resp, err := h.s.Login(msg.Username, msg.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(err.Error())
		return
	}

	// Создаём ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(struct {
		Token string `json:"token"`
	}{
		Token: h.generateToken(resp),
	})
	h.logger.Info("login success")
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	// Читаем сообщение
	var UserInfo model.UserInfo
	if err := json.NewDecoder(r.Body).Decode(&UserInfo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error(err.Error())
		return
	}

	// Проверяем поля на пустоту
	if UserInfo.Login == "" || UserInfo.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error("login or password is empty")
		return
	}

	// Регистрируем пользователя в БД
	uid, err := h.s.Registration(&UserInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(err.Error())
		return
	}

	// Создаём ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(struct {
		Token string `json:"token"`
	}{
		Token: h.generateToken(uid),
	})
	h.logger.Info("register success")
}

// GetUser - получение списка пользователей в системе
// с учётом агрегирования
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Читаем сообщение
	var Parameters model.GetUserParams
	if err := json.NewDecoder(r.Body).Decode(&Parameters); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error(err.Error())
		return
	}

	data, err := h.s.GetUsers(Parameters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
	h.logger.Info("get user success")
}

// UpdateUser Обновить какие то данные
// в конкретном агенте по UID
// Обновляются только числовые параметры
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// считываем входящее сообщение
	var UserInfo model.UserInfo
	if err := json.NewDecoder(r.Body).Decode(&UserInfo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error(err.Error())
		return
	}

	token := r.URL.Query().Get("token")
	f, _ := h.checkToken(token)
	if !f {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error("token is invalid. token:" + token)
		return
	}

	// проверяем данные
	flag, err := h.s.CheckUserByID(UserInfo.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(err.Error())
		return
	}
	if !flag {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error("user not exist")
		return
	}
	// обновляем пользователя
	err = h.s.UpdateUserData(&UserInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(UserInfo)
	h.logger.Info("update user success")
}

// DeleteUser Удалить конкретного пользователя
// удаляем конкретного пользователя по UID
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// получаем uid пользователя на удаление и токен пользователя
	uid := r.URL.Query().Get("uid")
	token := r.URL.Query().Get("token")
	// проверяем id на пустоту
	if uid == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error("uid is empty")
		return
	}
	// проверяем токен на пустоту
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error("token is empty")
		return
	}
	// проверяем права запрашивающего
	flag, err := h.s.CheckToken(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(err.Error())
	}
	if !flag {
		w.WriteHeader(http.StatusUnauthorized)
		h.logger.Info("token is invalid")
	}

	// удаляем пользователя
	err = h.s.DeleteUser(uid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(err.Error())
	}

	// Возвращаем информацию пользователю
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
		Uid     string `json:"uid"`
	}{
		Message: "User deleted",
		Uid:     uid,
	})
	h.logger.Info("delete user success")
}

// генерация токен строки по uid на 1 час
func (h *Handler) generateToken(uid string) string {
	claims := jwt.MapClaims{
		"uid":  uid,
		"time": time.Now().Add(time.Hour * time.Duration(1)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, _ := token.SignedString(key)
	return t
}

// проверяет токен на валидность
func (h *Handler) checkToken(tokenString string) (bool, error) {
	token, err := h.parseToken(tokenString)
	if err != nil {
		return false, err
	}
	if !token.Valid {
		return false, nil
	}
	return true, nil
}

// получение токен объекта из токена строки
func (h *Handler) parseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TODO(w http.ResponseWriter, r *http.Request) {}
