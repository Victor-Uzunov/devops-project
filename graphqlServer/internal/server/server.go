package server

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	graph "github.com/Victor-Uzunov/devops-project/graphqlServer/generated"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/internal/client"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/internal/resolvers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
)

type Server struct {
	Port   string
	Router http.Handler
}

func NewServer(config client.APIConfig) *Server {
	todoServiceClient := client.NewTodoServiceClient(&http.Client{}, config)
	directives := resolvers.NewDirective(todoServiceClient)

	rootResolver := resolvers.NewRootResolver(
		todoServiceClient,
	)
	gqlCfg := graph.Config{
		Resolvers: rootResolver,
		Directives: graph.DirectiveRoot{
			Validate: directives.ValidateDirective,
		},
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(gqlCfg))
	router := mux.NewRouter()
	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	router.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ready"))

	})
	corsRouter := router.PathPrefix("").Subrouter()
	corsRouter.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"POST", "GET", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		Debug:            true,
	}).Handler)

	corsRouter.HandleFunc("/", playground.Handler("GraphQL playground", "/query"))
	corsRouter.Handle("/query", JWTMiddleware(srv))

	return &Server{
		Port:   config.Port,
		Router: router,
	}
}

func (s *Server) Start() {
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", s.Port)
	log.Fatal(http.ListenAndServe(":"+s.Port, s.Router))
}
