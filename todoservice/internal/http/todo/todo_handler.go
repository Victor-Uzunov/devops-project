package todo

import (
	"bytes"
	"encoding/json"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/db"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/todos"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/converters"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/log"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"io"
	"net/http"
)

type Handler struct {
	service  todos.TodoService
	database *sqlx.DB
}

func NewHandler(service todos.TodoService, database *sqlx.DB) *Handler {
	return &Handler{service: service, database: database}
}

func (h *Handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("todo handler create request")
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.C(r.Context()).Errorf("error reading request body: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(bodyBytes) == 0 {
		log.C(r.Context()).Error("request body is empty")
		http.Error(w, "Request body cannot be empty", http.StatusBadRequest)
		return
	}

	log.C(r.Context()).Infof("Request body: %s", bodyBytes)

	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		log.C(r.Context()).Errorf("error while todo handler create request body err: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("error while todo handler create tx err: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)
	createdTodo, err := h.service.CreateTodo(ctx, todo)
	log.C(r.Context()).Debugf("todo handler create success, createdTodo: %v", createdTodo)
	if err != nil {
		log.C(r.Context()).Errorf("error while todo handler create err: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("erorr while todo handler create tx err: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdTodo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetTodo(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("todo handler get request")
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("error while todo handler get tx err: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	todo, err := h.service.GetTodo(ctx, id)
	log.C(r.Context()).Debugf("todo handler get success, todo: %v", todo)
	if err != nil {
		log.C(r.Context()).Errorf("error while todo handler get err: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("error while todo handler get tx err: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("todo handler update request")
	var todo models.Todo
	vars := mux.Vars(r)
	todoID := vars["id"]

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		log.C(r.Context()).Errorf("error while todo handler update req body err: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	todo.ID = todoID

	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("error while todo handler update tx err: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	_, err = h.service.GetTodo(ctx, todoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = h.service.UpdateTodo(ctx, todo)
	log.C(r.Context()).Debugf("todo handler update success, todo: %v", todo)
	if err != nil {
		log.C(r.Context()).Errorf("error while todo handler update err: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("error while todo handler update tx err: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("todo handler delete request")
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("erorr while todo handler delete tx err: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	_, err = h.service.GetTodo(ctx, id)
	log.C(r.Context()).Debugf("todo handler delete success, todo: %v", id)
	if err != nil {
		log.C(r.Context()).Errorf("erorr while todo handler delete err: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := h.service.DeleteTodo(ctx, id); err != nil {
		log.C(r.Context()).Errorf("erorr while todo handler delete err: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("erorr while todo handler delete tx err: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListTodosByListID(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("todo handler list request")
	vars := mux.Vars(r)
	listID := vars["list_id"]

	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("todo handler list tx err: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	todosByUser, err := h.service.ListTodosByListID(ctx, listID)
	if err != nil {
		log.C(r.Context()).Errorf("erorr while todo handler list tx err: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("erorr while todo handler list tx err: %v", err)
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

func (h *Handler) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("todo handler get all request")
	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("erorr while todo handler get tx err: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	todosByUser, err := h.service.GetAllTodos(ctx)
	log.C(r.Context()).Debugf("todo handler get all success, todos: %v", todosByUser)
	if err != nil {
		log.C(r.Context()).Errorf("erorr while todo handler get all err: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("erorr while todo handler get all tx err: %v", err)
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

func (h *Handler) CompleteTodo(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("complete todo handler")
	todoID := mux.Vars(r)["id"]

	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("erorr while completing todo handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	_, err = h.service.GetTodo(ctx, todoID)
	if err != nil {
		log.C(r.Context()).Errorf("erorr while completing todo handler, there is no such todo: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	updatedTodo, err := h.service.CompleteTodo(ctx, todoID)
	log.C(r.Context()).Debugf("complete todo handler for todo: %v", updatedTodo)
	if err != nil {
		log.C(r.Context()).Errorf("erorr while completing todo handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("erorr while completing todo handler transaction does not commit: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedTodo); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateTodoDescription(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("update todo description handler")
	todoID := mux.Vars(r)["id"]

	var updateData struct {
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		log.C(r.Context()).Errorf("erorr while updating todo description handler: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if updateData.Description == "" {
		log.C(r.Context()).Error("error while updating todo description handler, no description provided")
		http.Error(w, "Description cannot be empty", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("erorr while updating todo description handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	_, err = h.service.GetTodo(ctx, todoID)
	if err != nil {
		log.C(r.Context()).Errorf("erorr while updating todo description handler: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	updatedTodo, err := h.service.UpdateTodoDescription(ctx, todoID, updateData.Description)
	log.C(r.Context()).Debugf("update todo description handler with list: %v", updatedTodo)
	if err != nil {
		log.C(r.Context()).Errorf("erorr while updating toso description handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("erorr while updating todo description handler transaction does not commit: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedTodo); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateTodoTitle(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("update todo title handler")
	todoID := mux.Vars(r)["id"]

	var updateData struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		log.C(r.Context()).Errorf("erorr while updating todo title handler: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if updateData.Title == "" {
		log.C(r.Context()).Error("error while updating todo title handler, no title provided")
		http.Error(w, "Title cannot be empty", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("erorr while updating todo title handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	_, err = h.service.GetTodo(ctx, todoID)
	if err != nil {
		log.C(r.Context()).Errorf("error while updating todo title handler: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	updatedTodo, err := h.service.UpdateTodoTitle(ctx, todoID, updateData.Title)
	log.C(r.Context()).Debugf("update todo title handler with list: %v", updatedTodo)
	if err != nil {
		log.C(r.Context()).Errorf("erorr while updating todo title handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("erorr while updating todo title handler transaction does not commit: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedTodo); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateTodoPriority(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("update todo priority handler")
	todoID := mux.Vars(r)["id"]

	var updateData struct {
		Priority string `json:"priority"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		log.C(r.Context()).Errorf("erorr while updating todo priority handler: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	priority, err := converters.ToPriorityLevel(updateData.Priority)
	if err != nil {
		log.C(r.Context()).Errorf("invalid priority level: %v", err)
		http.Error(w, "Invalid priority level", http.StatusBadRequest)
		return
	}
	log.C(r.Context()).Debugf("update todo priority handler with priority: %v", priority)

	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("erorr while updating todo priority handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	_, err = h.service.GetTodo(ctx, todoID)
	if err != nil {
		log.C(r.Context()).Errorf("error while updating todo priority handler: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	updatedTodo, err := h.service.UpdateTodoPriority(ctx, todoID, priority)
	log.C(r.Context()).Debugf("update todo priority handler with list: %v", updatedTodo)
	if err != nil {
		log.C(r.Context()).Errorf("error updating todo priority handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("error while updating todo priority handler transaction does not commit: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedTodo); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateAssignedTo(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("update todo assigned_to handler")
	todoID := mux.Vars(r)["id"]

	var updateData struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		log.C(r.Context()).Errorf("erorr while updating todo assigned_to handler: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := pkg.ValidateUUID(updateData.UserID); err != nil {
		log.C(r.Context()).Errorf("invalid user id: %v", err)
		http.Error(w, "Invalid uuid for the user", http.StatusBadRequest)
		return
	}

	log.C(r.Context()).Debugf("update todo assigned_to handler with id: %s", updateData.UserID)

	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("erorr while updating todo assigned_to handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	_, err = h.service.GetTodo(ctx, todoID)
	if err != nil {
		log.C(r.Context()).Errorf("error while updating todo priority handler, there is no such todo: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	updatedTodo, err := h.service.UpdateAssignedTo(ctx, todoID, updateData.UserID)
	log.C(r.Context()).Debugf("update todo assigned_to handler with list: %v", updatedTodo)
	if err != nil {
		log.C(r.Context()).Errorf("error while updating todo assigned_to handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("erorr while updating todo assigned_to handler transaction does not commit: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedTodo); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
