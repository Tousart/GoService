package http

import (
	"httpServer/server/API/http/types"
	"httpServer/server/usecases"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Users struct {
	serviseUsers    usecases.Users
	serviceSessions usecases.Sessions
}

func NewUsersHandler(users usecases.Users, sessions usecases.Sessions) *Users {
	return &Users{
		serviseUsers:    users,
		serviceSessions: sessions,
	}
}

// @Summary Post a Register
// @Description Register user
// @Tags user
// @Accept json
// @Produce plain
// @Param user body types.User true "User"
// @Success 201 {string} string "Пользователь {login} зарегистрирован."
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /register [post]
func (s *Users) postRegisterHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateUserRequestHandler(r)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = s.serviseUsers.PostRegister(req.Login, req.Password)
	types.ProcessErrorRegister(w, err, req.Login)
}

// @Summary Post a Login
// @Description Login user
// @Tags user
// @Accept json
// @Produce json
// @Param user body types.User true "User"
// @Success 200 {object} types.GetSessionIdHandler
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /login [post]
func (s *Users) postLoginHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateUserRequestHandler(r)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	userId, err := s.serviseUsers.PostLogin(req.Login, req.Password)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	sessionId, err := s.serviceSessions.PostSessionId(userId)
	types.ProcessErrorLogin(w, err, &types.GetSessionIdHandler{SessionId: sessionId})
}

func (s *Users) WithUsersHandlers(r chi.Router) {
	r.Post("/register", s.postRegisterHandler)
	r.Post("/login", s.postLoginHandler)
}
