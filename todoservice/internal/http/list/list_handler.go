package list

import (
	"encoding/json"
	"fmt"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/db"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/lists"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/log"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type Handler struct {
	service  lists.ListService
	database *sqlx.DB
}

func NewHandler(service lists.ListService, database *sqlx.DB) *Handler {
	return &Handler{service: service, database: database}
}

func (h *Handler) CreateList(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("create list")
	var list models.List
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		log.C(r.Context()).Errorf("error while creating list handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("error while creating list handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	createdList, err := h.service.CreateList(ctx, list)
	log.C(r.Context()).Debugf("create list handler with id: %v", createdList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("error while committing transaction in create list handler: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdList); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetList(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("get list")
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("error while getting list handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	list, err := h.service.GetList(ctx, id)
	log.C(r.Context()).Debugf("get list handler with id: %v", list.ID)
	if err != nil {
		log.C(r.Context()).Errorf("error while getting list handler: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("erorr while get list handler transaction is committing: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(list); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateList(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("update list")
	vars := mux.Vars(r)
	id := vars["id"]

	var list models.List
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		log.C(r.Context()).Errorf("error while updating list handler: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	list.ID = id
	fmt.Println(list)

	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("error while updating list handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	_, err = h.service.GetList(ctx, list.ID)
	log.C(r.Context()).Debugf("update list handler with id: %v", list.ID)
	if err != nil {
		log.C(r.Context()).Errorf("update list handler: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = h.service.UpdateList(ctx, list)
	log.C(r.Context()).Debugf("update list handler with id: %v", list.ID)
	if err != nil {
		log.C(r.Context()).Errorf("update list handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("error while committing transaction in update list handler: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteList(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("delete list")
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("error while deleting list handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	_, err = h.service.GetList(ctx, id)
	log.C(r.Context()).Debugf("delete list handler with id: %v", id)
	if err != nil {
		log.C(r.Context()).Errorf("error while deleting list handler: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := h.service.DeleteList(ctx, id); err != nil {
		log.C(r.Context()).Errorf("error while deleting list handle: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("error while committing transaction in delete list handler: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListAllByUser(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("list all by user id")
	vars := mux.Vars(r)
	id := vars["user_id"]

	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("error while listing all list by user id handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	listsByUserID, err := h.service.ListAllByUserID(ctx, id)
	log.C(r.Context()).Debugf("list all by user id handler: %v", listsByUserID)
	if err != nil {
		log.C(r.Context()).Errorf("error while list all by user id handler: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("error while committing transaction in list all by user id: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(listsByUserID); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) GetAllLists(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("get all lists")
	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("error while getting all lists handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	result, err := h.service.GetAllLists(ctx)
	if err != nil {
		log.C(r.Context()).Errorf("error while getting all lists handler: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("error while committing transaction in get all lists: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetUsersByListID(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("get users by list id")
	vars := mux.Vars(r)
	id := vars["id"]
	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("error while getting users by list id handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	result, err := h.service.GetUsersByListID(ctx, id)
	log.C(r.Context()).Debugf("get users by list id handler: %v", result)
	if err != nil {
		log.C(r.Context()).Errorf("error while getting users by list id handler: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("error while getting users by list id transaction does not commit: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetListOwnerID(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("get list owner id")
	vars := mux.Vars(r)
	id := vars["id"]
	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("error while getting list owner id handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	result, err := h.service.GetListOwnerID(ctx, id)
	log.C(r.Context()).Debugf("get list owner id handler: %v", result)
	if err != nil {
		log.C(r.Context()).Errorf("error while getting list owner id handler: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("error while getting list owner id transaction does not commit: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreateAccess(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("create access")
	var access models.Access
	if err := json.NewDecoder(r.Body).Decode(&access); err != nil {
		log.C(r.Context()).Errorf("error while creating access handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("error while creating access handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	createdAccess, err := h.service.CreateAccess(ctx, access)
	log.C(r.Context()).Debugf("create access handler: %v", createdAccess)
	if err != nil {
		log.C(r.Context()).Errorf("error while creating access handler: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("error while creating access transaction does not commit: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdAccess); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteAccess(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("delete access")
	vars := mux.Vars(r)
	listID := vars["list_id"]
	userID := vars["user_id"]

	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("error while deleting access handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	if err := h.service.DeleteAccess(ctx, listID, userID); err != nil {
		log.C(r.Context()).Errorf("error while deleting access handler: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("error while deleting access transaction does not commit: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) AcceptList(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("accept list")
	vars := mux.Vars(r)
	listID := vars["list_id"]
	userID := vars["user_id"]

	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("error while getting access handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	err = h.service.AcceptList(ctx, listID, userID)
	if err != nil {
		log.C(r.Context()).Errorf("error while getting accepting list handler: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	log.C(r.Context()).Debugf("get access handler - userID: %s, listID: %s", userID, listID)

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("error while getting access transaction does not commit: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetAccess(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("get access")
	vars := mux.Vars(r)
	listID := vars["list_id"]
	userID := vars["user_id"]

	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("error while getting access handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	access, err := h.service.GetAccess(ctx, listID, userID)
	if err != nil {
		log.C(r.Context()).Errorf("error while getting access handler: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	log.C(r.Context()).Debugf("get access handler - userID: %s, listID: %s", access.UserID, access.ListID)

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("error while getting access transaction does not commit: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(access); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateListDescription(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("update list description handler")
	listID := mux.Vars(r)["id"]

	var updateData struct {
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		log.C(r.Context()).Errorf("error while updating list description handler: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if updateData.Description == "" {
		log.C(r.Context()).Error("error while updating list description handler: empty description")
		http.Error(w, "Description cannot be empty", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("error while updating list description handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	_, err = h.service.GetList(ctx, listID)
	if err != nil {
		log.C(r.Context()).Errorf("error while updating list description handler: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	updatedList, err := h.service.UpdateListDescription(ctx, listID, updateData.Description)
	log.C(r.Context()).Debugf("update list description handler with list: %v", updatedList)
	if err != nil {
		log.C(r.Context()).Errorf("error while updating list description handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("error while updating list description handler transaction does not commit: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedList); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateListName(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("update list name handler")
	listID := mux.Vars(r)["id"]

	var updateData struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		log.C(r.Context()).Errorf("error while updating list name handler: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if updateData.Name == "" {
		log.C(r.Context()).Error("error while updating list name handler: empty name")
		http.Error(w, "Description cannot be empty", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("error while updating list name handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	_, err = h.service.GetList(ctx, listID)
	if err != nil {
		log.C(r.Context()).Errorf("error while updating list name handler: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	updatedList, err := h.service.UpdateListName(ctx, listID, updateData.Name)
	log.C(r.Context()).Debugf("update list name handler with list: %v", updatedList)
	if err != nil {
		log.C(r.Context()).Errorf("error while updating list name handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("error while updating list name handler transaction does not commit: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedList); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetListsByUser(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("get lists by user handler")
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		log.C(r.Context()).Errorf("error while getting lists by user handler missing user id in the context")
		http.Error(w, "there is no user id in the context:"+http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("error while getting lists by user id handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	access, err := h.service.ListAllByUserID(ctx, userID)
	var result []models.List
	for _, el := range access {
		l, err := h.service.GetList(ctx, el.ListID)
		if err != nil {
			log.C(r.Context()).Errorf("error while getting the list with id: %s; %v", el.ListID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result = append(result, l)
	}
	log.C(r.Context()).Debugf("get lists by user id handler: %v", result)
	if err != nil {
		log.C(r.Context()).Errorf("error while getting list by user id handler: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("error while getting list by user id transaction does not commit: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetAcceptedLists(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("get accepted lists by user handler")
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		log.C(r.Context()).Errorf("error while getting lists by user handler missing user id in the context")
		http.Error(w, "there is no user id in the context:"+http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("error while getting lists by user id handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	access, err := h.service.GetAcceptedLists(ctx, userID)
	var result []models.List
	for _, el := range access {
		l, err := h.service.GetList(ctx, el.ListID)
		if err != nil {
			log.C(r.Context()).Errorf("error while getting the list with id: %s; %v", el.ListID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result = append(result, l)
	}
	log.C(r.Context()).Debugf("get lists by user id handler: %v", result)
	if err != nil {
		log.C(r.Context()).Errorf("error while getting list by user id handler: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("error while getting list by user id transaction does not commit: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetPendingLists(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("get pending lists by user handler")
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		log.C(r.Context()).Errorf("error while getting pending lists by user handler missing user id in the context")
		http.Error(w, "there is no user id in the context:"+http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("error while getting pending lists by user id handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	access, err := h.service.GetPendingLists(ctx, userID)
	var result []models.List
	for _, el := range access {
		l, err := h.service.GetList(ctx, el.ListID)
		if err != nil {
			log.C(r.Context()).Errorf("error while getting the pending list with id: %s; %v", el.ListID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result = append(result, l)
	}
	log.C(r.Context()).Debugf("get pending lists by user id handler: %v", result)
	if err != nil {
		log.C(r.Context()).Errorf("error while getting pending list by user id handler: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("error while getting pending list by user id transaction does not commit: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetAllUserTodos(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("get all todos by user handler")
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		log.C(r.Context()).Errorf("error while getting all todos by user handler missing user id in the context")
		http.Error(w, "there is no user id in the context:"+http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("error while getting all todos by user id handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	access, err := h.service.ListAllByUserID(ctx, userID)
	accepted, err := h.service.GetAcceptedLists(ctx, userID)
	combined := append(access, accepted...)
	var result []models.Todo
	for _, el := range combined {
		todosForList, err := h.service.GetAllTodosForList(ctx, el.ListID)
		if err != nil {
			log.C(r.Context()).Errorf("error while getting the list with id: %s; %v", el.ListID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result = append(result, todosForList...)
	}
	log.C(r.Context()).Debugf("get lists by user id handler: %v", result)
	if err != nil {
		log.C(r.Context()).Errorf("error while getting list by user id handler: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("error while getting list by user id transaction does not commit: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetAccessesByListID(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("get accesses by list id handler")
	vars := mux.Vars(r)
	listID := vars["list_id"]

	ctx := r.Context()

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(r.Context()).Errorf("error while getting access handler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	accesses, err := h.service.GetAccessesByListID(ctx, listID)
	if err != nil {
		log.C(r.Context()).Errorf("error while getting access handler: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	log.C(r.Context()).Debugf("get accesses handler %v", accesses)

	if err = tx.Commit(); err != nil {
		log.C(r.Context()).Errorf("error while getting access transaction does not commit: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(accesses); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
