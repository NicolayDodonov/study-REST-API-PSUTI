package handler

import (
	"encoding/json"
	"net/http"
	"study-REST-API-PSUTI/internal/logger"
	"study-REST-API-PSUTI/internal/model"
	"study-REST-API-PSUTI/internal/storage"
)

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
	_ = json.NewEncoder(w).Encode(resp)
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
	_ = json.NewEncoder(w).Encode(uid)
	h.logger.Info("register success")
}

// GetUser - получение списка пользователей в системе
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Читаем сообщение
	var Parameters model.GetUserParams
	if err := json.NewDecoder(r.Body).Decode(&Parameters); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error(err.Error())
		return
	}

	// Проверяем

}

// UpdateUser Обновить какие то данные в конкретном агенте по UID
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {

}

// DeleteUser Удалить конкретного пользователя
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {

}

func TODO(w http.ResponseWriter, r *http.Request) {}
