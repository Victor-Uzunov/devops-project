package user

import (
	"context"
	"encoding/json"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/db"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/users"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/log"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"time"
)

type Handler struct {
	service  users.UserService
	database *sqlx.DB
}

func NewHandler(service users.UserService, database *sqlx.DB) *Handler {
	return &Handler{service: service, database: database}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("create user handler called")
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.C(r.Context()).Errorf("error while creating user json decoding failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("error while creating user transaction failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	createdUser, err := h.service.CreateUser(ctx, user)
	log.C(r.Context()).Debugf("create user success: %v", createdUser)
	if err != nil {
		log.C(r.Context()).Errorf("error while creating user failed: %v", err)
		http.Error(w, "Failed to create user. Please ensure all required fields are correctly filled and there is no such a user already. Error details: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("error while creating user transaction failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("get user handler called")
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("error while getting user transaction failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	user, err := h.service.GetUser(ctx, id)
	if err != nil {
		log.C(r.Context()).Errorf("erorr while getting user failed: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("error while getting user transaction failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("update user handler called")
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.C(r.Context()).Errorf("error while updating user json decoding failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("erorr while updating user transaction failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	_, err = h.service.GetUser(ctx, user.ID)
	log.C(r.Context()).Debugf("update user success: %v", user)
	if err != nil {
		log.C(r.Context()).Errorf("erorr while updating user failed: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = h.service.UpdateUser(ctx, user)
	log.C(r.Context()).Debugf("update user success: %v", user)
	if err != nil {
		log.C(r.Context()).Errorf("erorr while updating user failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("erorr while updating user transaction failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("delete user handler called")
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("delete user transaction failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	_, err = h.service.GetUser(ctx, id)
	log.C(r.Context()).Debugf("get user success: %v", id)
	if err != nil {
		log.C(r.Context()).Errorf("erorr while deleting user failed: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := h.service.DeleteUser(ctx, id); err != nil {
		log.C(r.Context()).Errorf("erorr while deleting user failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.C(r.Context()).Debugf("delete user success: %v", id)

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("erorr while deleting user transaction failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("get all users handler called")
	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("erorr while getting all users transaction failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	todosByUser, err := h.service.GetAllUsers(ctx)
	log.C(r.Context()).Debugf("get all users success: %v", todosByUser)
	if err != nil {
		log.C(r.Context()).Errorf("erorr while getting all users failed: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("erorr while getting all users transaction failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(todosByUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("get user by email handler called")
	vars := mux.Vars(r)
	email := vars["email"]
	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("erorr while getting all users transaction failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	user, err := h.service.GetUserByEmail(ctx, email)
	log.C(r.Context()).Debugf("get user`s email success: %v", user)
	if err != nil {
		log.C(r.Context()).Errorf("erorr while getting user`s email failed: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("erorr while getting user`s email transaction failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type LogoutRequest struct {
	Email string `json:"email"`
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("logout handler called")
	ctx := r.Context()

	var logoutRequest LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&logoutRequest); err != nil {
		log.C(ctx).Errorf("error unmarshalling request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("error while getting user transaction failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	err = h.service.Logout(ctx, logoutRequest.Email)
	if err != nil {
		log.C(r.Context()).Errorf("erorr while logging out user failed: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	log.C(ctx).Info("Delete user cookies for access token, refresh token and user role")
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "user_role",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: false,
		Path:     "/",
	})

	//h.deleteUserCookies(ctx, &w)

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("error while getting user transaction failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) deleteUserCookies(ctx context.Context, w *http.ResponseWriter) {
	log.C(ctx).Info("Delete user cookies for access token, refresh token and user role")
	http.SetCookie(*w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
	})
	http.SetCookie(*w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
	})
	http.SetCookie(*w, &http.Cookie{
		Name:     "user_role",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: false,
		Path:     "/",
	})

	log.C(ctx).Debug("Cookies for access token, refresh token and user role are deleted successfully")
}
