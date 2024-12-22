package http

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/db"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/lists"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/todos"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/users"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/jwt"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/log"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

//go:generate mockery --name=Middlewares --output=automock --with-expecter=true --outpkg=automock --case=underscore --disable-version-string --with-expecter=true
type Middlewares interface {
	Protected(next http.Handler, neededRole constants.Role, accessibility constants.Accessibility) http.Handler
	JWTMiddleware(next http.Handler) http.Handler
}

type Middleware struct {
	userService users.UserService
	listService lists.ListService
	todoService todos.TodoService
	tokenParser *jwt.TokenParser
	database    *sqlx.DB
}

func NewMiddleware(userService users.UserService, listService lists.ListService, todoService todos.TodoService, tokenParser *jwt.TokenParser, database *sqlx.DB) Middlewares {
	return &Middleware{
		userService: userService,
		listService: listService,
		todoService: todoService,
		tokenParser: tokenParser,
		database:    database,
	}
}

func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func HandlePreflight(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w)
		log.C(r.Context()).Debugf("setCORSHeaders successfully")
		if r.Method == http.MethodOptions {
			logrus.Debug("Setting CORS Headers")
			w.WriteHeader(http.StatusOK)
			return
		}

		nextHandler.ServeHTTP(w, r)
	})
}

func (m *Middleware) Protected(next http.Handler, neededRole constants.Role, accessibility constants.Accessibility) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.C(r.Context()).Info("Protected middleware")
		claim, ok := r.Context().Value("user").(*jwt.Claims)
		if !ok {
			http.Error(w, "there is no user claim in the context", http.StatusUnauthorized)
			return
		}
		log.C(r.Context()).Debugf("the email and the role of the user are: %s, %s", claim.Email, claim.Role)
		if constants.RolePower(neededRole) > constants.RolePower(pkg.StringToRole(claim.Role)) {
			log.C(r.Context()).Errorf("user`s role is %v, but the needed role is %v", claim.Role, neededRole)
			http.Error(w, http.StatusText(http.StatusForbidden)+"you do not have the needed role for this action", http.StatusForbidden)
			return
		}

		switch accessibility {
		case constants.IsOwner:
			m.isOwner(w, r, next)
		case constants.HasAccessTodo:
			m.hasTodoAccess(w, r, next)
		case constants.HasAccessList:
			m.hasListAccess(w, r, next)
		case constants.NoRestriction:
			next.ServeHTTP(w, r)
		default:
			log.C(r.Context()).Errorf("accessibility `%s` is not valid one", accessibility)
			http.Error(w, http.StatusText(http.StatusInternalServerError)+"not a valid accessibility type", http.StatusInternalServerError)
			return
		}
	})
}

func (m *Middleware) isOwner(w http.ResponseWriter, r *http.Request, next http.Handler) {
	log.C(r.Context()).Info("isOwner middleware")
	ctx := r.Context()
	vars := mux.Vars(r)
	var getList = true
	var ownerID string
	id := vars["id"]
	if id == "" {
		id = vars["list_id"]
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.C(ctx).Errorf("cannot read request body: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if id == "" {
		var list models.List
		if err := json.Unmarshal(bodyBytes, &list); err != nil {
			log.C(ctx).Errorf("cannot get list from the body in hasAccess middleware: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id = list.ID
		getList = false
		ownerID = list.OwnerID
	}

	var claim *jwt.Claims
	if user, isAdmin := authorizeAdmin(r); isAdmin {
		log.C(ctx).Debugf("user is admin: %v", claim)
		next.ServeHTTP(w, r)
		return
	} else if user == nil {
		log.C(ctx).Error("error: there is no user in the context")
		http.Error(w, "there is no user in the context", http.StatusUnauthorized)
		return
	} else {
		claim = user
	}
	log.C(ctx).Debugf("the email and role are: %s, %s", claim.Email, claim.Role)

	tx, err := m.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(ctx).Errorf("isOwner middleware transaction failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	user, err := m.userService.GetUserByEmail(ctx, claim.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if getList {
		list, err := m.listService.GetList(ctx, id)
		if err != nil {
			log.C(r.Context()).Errorf("middleware cannot get list: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ownerID = list.OwnerID
	}

	if user.ID != ownerID {
		log.C(ctx).Errorf("userID: %s is not owner of the listID: %s, ownnerID is: %s", user.ID, id, ownerID)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	err = tx.Commit()
	if err != nil {
		log.C(ctx).Errorf("JWTMiddleware transaction failed to commit: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError)+" error committing the transaction", http.StatusInternalServerError)
		return
	}
	next.ServeHTTP(w, r)
}

func (m *Middleware) hasListAccess(w http.ResponseWriter, r *http.Request, next http.Handler) {
	ctx := r.Context()
	log.C(ctx).Info("has access in list middleware")
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		id = vars["list_id"]
	}

	var claim *jwt.Claims
	if user, isAdmin := authorizeAdmin(r); isAdmin {
		log.C(ctx).Debugf("user is admin: %v", claim)
		next.ServeHTTP(w, r)
		return
	} else if user == nil {
		log.C(ctx).Debug("there is no user in the context")
		http.Error(w, "there is no user in the context", http.StatusUnauthorized)
		return
	} else {
		claim = user
	}
	log.C(ctx).Debugf("the email and role are: %s, %s", claim.Email, claim.Role)

	tx, err := m.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(ctx).Errorf("isOwner middleware transaction failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	user, err := m.userService.GetUserByEmail(ctx, claim.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	listsAccess, err := m.listService.ListAllByUserID(ctx, user.ID)
	if err != nil {
		log.C(r.Context()).Errorf("middleware cannot get all lists for a user: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	listsPending, err := m.listService.GetPendingLists(ctx, user.ID)
	if err != nil {
		log.C(r.Context()).Errorf("middleware cannot get pending lists for a user: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	listsAccepted, err := m.listService.GetAcceptedLists(ctx, user.ID)
	if err != nil {
		log.C(r.Context()).Errorf("middleware cannot get accepted lists for a user: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hasAccess := false
	for _, access := range listsAccess {
		if access.ListID == id {
			hasAccess = true
			break
		}
	}
	for _, access := range listsPending {
		if access.ListID == id {
			hasAccess = true
			break
		}
	}
	for _, access := range listsAccepted {
		if access.ListID == id {
			hasAccess = true
			break
		}
	}

	if !hasAccess {
		log.C(ctx).Errorf("user do not have access for list with ID: %s", id)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	err = tx.Commit()
	if err != nil {
		log.C(ctx).Errorf("JWTMiddleware transaction failed to commit: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError)+" error committing the transaction", http.StatusInternalServerError)
		return
	}
	next.ServeHTTP(w, r)

}

func (m *Middleware) hasTodoAccess(w http.ResponseWriter, r *http.Request, next http.Handler) {
	ctx := r.Context()
	log.C(ctx).Info("has access in todo middleware")
	vars := mux.Vars(r)
	id := vars["id"]
	getTodo := true
	var listID string
	var todo models.Todo

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.C(ctx).Errorf("cannot read request body: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if id == "" {
		if err := json.Unmarshal(bodyBytes, &todo); err != nil {
			log.C(ctx).Errorf("cannot get todo from the body in hasAccess middleware: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id = todo.ID
		listID = todo.ListID
		getTodo = false
	}
	var claim *jwt.Claims
	if user, isAdmin := authorizeAdmin(r); isAdmin {
		log.C(ctx).Debugf("user is admin: %v", claim)
		next.ServeHTTP(w, r)
		return
	} else if user == nil {
		log.C(ctx).Debug("there is no user in the context")
		http.Error(w, "there is no user in the context", http.StatusUnauthorized)
		return
	} else {
		claim = user
	}
	log.C(ctx).Debugf("the email and role are: %s, %s", claim.Email, claim.Role)

	tx, err := m.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(ctx).Errorf("isOwner middleware transaction failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)
	log.C(ctx).Debugf("the email and role are: %s", claim.Email)
	user, err := m.userService.GetUserByEmail(ctx, claim.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	listsAccess, err := m.listService.ListAllByUserID(ctx, user.ID)
	if err != nil {
		log.C(ctx).Errorf("middleware cannot get all lists for a user: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	listsAccepted, err := m.listService.GetAcceptedLists(ctx, user.ID)
	if err != nil {
		log.C(r.Context()).Errorf("middleware cannot get accepted lists for a user: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if getTodo {
		todo, err = m.todoService.GetTodo(ctx, id)
		if err != nil {
			log.C(ctx).Errorf("middleware cannot get todo for a user: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		listID = todo.ListID
	}

	hasAccess := false
	for _, access := range listsAccess {
		if access.ListID == listID {
			hasAccess = true
			break
		}
	}
	for _, access := range listsAccepted {
		if access.ListID == listID {
			hasAccess = true
			break
		}
	}
	if !hasAccess {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	err = tx.Commit()
	if err != nil {
		log.C(ctx).Errorf("JWTMiddleware transaction failed to commit: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError)+" error committing the transaction", http.StatusInternalServerError)
		return
	}
	next.ServeHTTP(w, r.WithContext(ctx))
}

func (m *Middleware) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.C(ctx).Info("JWTMiddleware rest server")
		token := r.Header.Get(constants.AuthorizationHeader)
		if token == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "user is not authorized",
			})
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")
		claim, err := m.tokenParser.ParseJWT(token)
		if err != nil {
			log.C(ctx).Errorf("error parsing token: %v", err)
			http.Error(w, "error while parsing the token: "+http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, "user", claim)
		ctx = context.WithValue(ctx, "user_id", claim.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func authorizeAdmin(r *http.Request) (*jwt.Claims, bool) {
	claim, ok := r.Context().Value("user").(*jwt.Claims)
	if !ok {
		return nil, false
	}

	if claim.Role == string(constants.Admin) {
		return claim, true
	}

	return claim, false
}
