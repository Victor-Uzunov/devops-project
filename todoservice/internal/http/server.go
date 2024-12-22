package http

import (
	httplist "github.com/Victor-Uzunov/devops-project/todoservice/internal/http/list"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/http/todo"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/http/user"
	listsdomain "github.com/Victor-Uzunov/devops-project/todoservice/internal/lists"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/oauth2"
	tododomain "github.com/Victor-Uzunov/devops-project/todoservice/internal/todos"
	userdomain "github.com/Victor-Uzunov/devops-project/todoservice/internal/users"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	token "github.com/Victor-Uzunov/devops-project/todoservice/pkg/jwt"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/time"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/uid"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
	"log"
	"net/http"
)

type Server struct {
	ListHandler   *httplist.Handler
	TodoHandler   *todo.Handler
	UserHandler   *user.Handler
	Oauth2Handler *oauth2.Handler
	Middleware    Middlewares
}

func NewServer(db *sqlx.DB, config token.ConfigOAuth2) *Server {
	listRepo := listsdomain.NewSQLXListRepository()
	todoRepo := tododomain.NewSQLXTodoRepository()
	userRepo := userdomain.NewSQLXUserRepository()

	uuidServer := uid.NewService()
	timeServer := time.Time{}

	listService := listsdomain.NewService(listRepo, uuidServer, timeServer)
	todoService := tododomain.NewService(todoRepo, uuidServer, timeServer)
	userService := userdomain.NewService(userRepo, uuidServer, timeServer)

	listHandler := httplist.NewHandler(listService, db)
	todoHandler := todo.NewHandler(todoService, db)
	userHandler := user.NewHandler(userService, db)

	oauth2Handler := oauth2.NewOAuth2(config, userService, db)
	tokenParser := token.NewTokenParser(config)
	middleware := NewMiddleware(userService, listService, todoService, tokenParser, db)

	return &Server{
		ListHandler:   listHandler,
		TodoHandler:   todoHandler,
		UserHandler:   userHandler,
		Oauth2Handler: oauth2Handler,
		Middleware:    middleware,
	}
}

func NewServerWithServices(db *sqlx.DB, listService listsdomain.ListService, todoService tododomain.TodoService, userService userdomain.UserService, middleware Middlewares) *Server {
	listHandler := httplist.NewHandler(listService, db)
	todoHandler := todo.NewHandler(todoService, db)
	userHandler := user.NewHandler(userService, db)

	return &Server{
		ListHandler: listHandler,
		TodoHandler: todoHandler,
		UserHandler: userHandler,
		Middleware:  middleware,
	}
}

func (s *Server) RegisterRoutes(router *mux.Router) {
	loginRouter := router.PathPrefix("/login").Subrouter()
	loginRouter.HandleFunc("/", s.Oauth2Handler.RootHandler).Methods(http.MethodGet)
	loginRouter.HandleFunc("/github", s.Oauth2Handler.GithubLoginHandler).Methods(http.MethodGet)
	loginRouter.HandleFunc("/github/callback", s.Oauth2Handler.GithubCallbackHandler).Methods(http.MethodGet)
	loginRouter.HandleFunc("/refresh-token", s.Oauth2Handler.RefreshTokenHandler).Methods(http.MethodGet)
	loginRouter.HandleFunc("/logout", s.UserHandler.Logout).Methods(http.MethodPost)

	protectedRouter := router.PathPrefix("").Subrouter()
	protectedRouter.Use(s.Middleware.JWTMiddleware)
	protectedRouter.Handle("/lists/create", s.Middleware.Protected(http.HandlerFunc(s.ListHandler.CreateList), constants.Writer, constants.NoRestriction)).Methods(http.MethodPost)
	protectedRouter.Handle("/lists_access/create/{list_id:[a-zA-Z0-9-]+}/{user_id:[a-zA-Z0-9-]+}", s.Middleware.Protected(http.HandlerFunc(s.ListHandler.CreateAccess), constants.Reader, constants.HasAccessList)).Methods(http.MethodPost)
	protectedRouter.Handle("/lists_access/list/{list_id:[a-zA-Z0-9-]+}", s.Middleware.Protected(http.HandlerFunc(s.ListHandler.GetAccessesByListID), constants.Reader, constants.HasAccessList)).Methods(http.MethodGet)
	protectedRouter.Handle("/lists_access/{list_id:[a-zA-Z0-9-]+}/{user_id:[a-zA-Z0-9-]+}", s.Middleware.Protected(http.HandlerFunc(s.ListHandler.AcceptList), constants.Reader, constants.NoRestriction)).Methods(http.MethodPost)
	protectedRouter.Handle("/lists_access/{list_id:[a-zA-Z0-9-]+}/{user_id:[a-zA-Z0-9-]+}", s.Middleware.Protected(http.HandlerFunc(s.ListHandler.GetAccess), constants.Reader, constants.NoRestriction)).Methods(http.MethodGet)
	protectedRouter.Handle("/lists_access/{list_id:[a-zA-Z0-9-]+}/{user_id:[a-zA-Z0-9-]+}", s.Middleware.Protected(http.HandlerFunc(s.ListHandler.DeleteAccess), constants.Reader, constants.NoRestriction)).Methods(http.MethodDelete)
	protectedRouter.Handle("/lists/all", s.Middleware.Protected(http.HandlerFunc(s.ListHandler.GetAllLists), constants.Admin, constants.NoRestriction)).Methods(http.MethodGet)
	protectedRouter.Handle("/lists/user/all", s.Middleware.Protected(http.HandlerFunc(s.ListHandler.GetListsByUser), constants.Reader, constants.NoRestriction)).Methods(http.MethodGet)
	protectedRouter.Handle("/lists/user/accepted", s.Middleware.Protected(http.HandlerFunc(s.ListHandler.GetAcceptedLists), constants.Reader, constants.NoRestriction)).Methods(http.MethodGet)
	protectedRouter.Handle("/lists/pending/all", s.Middleware.Protected(http.HandlerFunc(s.ListHandler.GetPendingLists), constants.Reader, constants.NoRestriction)).Methods(http.MethodGet)
	protectedRouter.Handle("/lists/{list_id:[a-zA-Z0-9-]+}/todos", s.Middleware.Protected(http.HandlerFunc(s.TodoHandler.ListTodosByListID), constants.Reader, constants.HasAccessList)).Methods(http.MethodGet)
	protectedRouter.Handle("/lists/{id:[a-zA-Z0-9-]+}/owner", s.Middleware.Protected(http.HandlerFunc(s.ListHandler.GetListOwnerID), constants.Reader, constants.NoRestriction)).Methods(http.MethodGet)
	protectedRouter.Handle("/lists/{id:[a-zA-Z0-9-]+}/users", s.Middleware.Protected(http.HandlerFunc(s.ListHandler.GetUsersByListID), constants.Reader, constants.NoRestriction)).Methods(http.MethodGet)
	protectedRouter.Handle("/lists/{id:[a-zA-Z0-9-]+}/description", s.Middleware.Protected(http.HandlerFunc(s.ListHandler.UpdateListDescription), constants.Writer, constants.HasAccessList)).Methods(http.MethodPatch)
	protectedRouter.Handle("/lists/{id:[a-zA-Z0-9-]+}/name", s.Middleware.Protected(http.HandlerFunc(s.ListHandler.UpdateListName), constants.Writer, constants.HasAccessList)).Methods(http.MethodPatch)
	protectedRouter.Handle("/lists/{id:[a-zA-Z0-9-]+}", s.Middleware.Protected(http.HandlerFunc(s.ListHandler.GetList), constants.Reader, constants.HasAccessList)).Methods(http.MethodGet)
	protectedRouter.Handle("/lists/{id:[a-zA-Z0-9-]+}", s.Middleware.Protected(http.HandlerFunc(s.ListHandler.UpdateList), constants.Writer, constants.HasAccessList)).Methods(http.MethodPut)
	protectedRouter.Handle("/lists/{id:[a-zA-Z0-9-]+}", s.Middleware.Protected(http.HandlerFunc(s.ListHandler.DeleteList), constants.Writer, constants.HasAccessList)).Methods(http.MethodDelete)

	protectedRouter.Handle("/todos/create", s.Middleware.Protected(http.HandlerFunc(s.TodoHandler.CreateTodo), constants.Writer, constants.HasAccessTodo)).Methods(http.MethodPost)
	protectedRouter.Handle("/todos/all", s.Middleware.Protected(http.HandlerFunc(s.TodoHandler.GetAllTodos), constants.Admin, constants.NoRestriction)).Methods(http.MethodGet)
	protectedRouter.Handle("/todos/user/all", s.Middleware.Protected(http.HandlerFunc(s.ListHandler.GetAllUserTodos), constants.Reader, constants.NoRestriction)).Methods(http.MethodGet)
	protectedRouter.Handle("/todos/{id:[a-zA-Z0-9-]+}/complete", s.Middleware.Protected(http.HandlerFunc(s.TodoHandler.CompleteTodo), constants.Writer, constants.HasAccessTodo)).Methods(http.MethodPatch)
	protectedRouter.Handle("/todos/{id:[a-zA-Z0-9-]+}/title", s.Middleware.Protected(http.HandlerFunc(s.TodoHandler.UpdateTodoTitle), constants.Writer, constants.HasAccessTodo)).Methods(http.MethodPatch)
	protectedRouter.Handle("/todos/{id:[a-zA-Z0-9-]+}/assign_to", s.Middleware.Protected(http.HandlerFunc(s.TodoHandler.UpdateAssignedTo), constants.Writer, constants.HasAccessTodo)).Methods(http.MethodPatch)
	protectedRouter.Handle("/todos/{id:[a-zA-Z0-9-]+}/priority", s.Middleware.Protected(http.HandlerFunc(s.TodoHandler.UpdateTodoPriority), constants.Writer, constants.HasAccessTodo)).Methods(http.MethodPatch)
	protectedRouter.Handle("/todos/{id:[a-zA-Z0-9-]+}/description", s.Middleware.Protected(http.HandlerFunc(s.TodoHandler.UpdateTodoDescription), constants.Writer, constants.HasAccessTodo)).Methods(http.MethodPatch)
	protectedRouter.Handle("/todos/{id:[a-zA-Z0-9-]+}", s.Middleware.Protected(http.HandlerFunc(s.TodoHandler.GetTodo), constants.Reader, constants.HasAccessTodo)).Methods(http.MethodGet)
	protectedRouter.Handle("/todos/{id:[a-zA-Z0-9-]+}", s.Middleware.Protected(http.HandlerFunc(s.TodoHandler.UpdateTodo), constants.Writer, constants.HasAccessTodo)).Methods(http.MethodPut)
	protectedRouter.Handle("/todos/{id:[a-zA-Z0-9-]+}", s.Middleware.Protected(http.HandlerFunc(s.TodoHandler.DeleteTodo), constants.Writer, constants.HasAccessTodo)).Methods(http.MethodDelete)

	protectedRouter.Handle("/users/create", s.Middleware.Protected(http.HandlerFunc(s.UserHandler.CreateUser), constants.Admin, constants.NoRestriction)).Methods(http.MethodPost)
	protectedRouter.Handle("/users/all", s.Middleware.Protected(http.HandlerFunc(s.UserHandler.GetAllUsers), constants.Reader, constants.NoRestriction)).Methods(http.MethodGet)
	protectedRouter.Handle("/users/{id:[a-zA-Z0-9-]+}", s.Middleware.Protected(http.HandlerFunc(s.UserHandler.UpdateUser), constants.Admin, constants.NoRestriction)).Methods(http.MethodPut)
	protectedRouter.Handle("/users/{id:[a-zA-Z0-9-]+}", s.Middleware.Protected(http.HandlerFunc(s.UserHandler.DeleteUser), constants.Admin, constants.NoRestriction)).Methods(http.MethodDelete)
	protectedRouter.Handle("/users/{id:[a-zA-Z0-9-]+}", s.Middleware.Protected(http.HandlerFunc(s.UserHandler.GetUser), constants.Reader, constants.NoRestriction)).Methods(http.MethodGet)
	protectedRouter.Handle("/users/email/{email:.+}", s.Middleware.Protected(http.HandlerFunc(s.UserHandler.GetUserByEmail), constants.Reader, constants.NoRestriction)).Methods(http.MethodGet)

}

func (s *Server) Start() {
	router := mux.NewRouter()
	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	router.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ready"))

	})
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})
	router.Use(HandlePreflight)
	handler := c.Handler(router)
	s.RegisterRoutes(router)
	log.Fatal(http.ListenAndServe(":5000", handler))
}
