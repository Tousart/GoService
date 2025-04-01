package main

import (
	"httpServer/server/API/http"
	"httpServer/server/config"
	"httpServer/server/repository/postgres"
	"httpServer/server/repository/rabbitMQ"
	"httpServer/server/repository/redis"
	"httpServer/server/usecases/service"
	"log"

	pkgHttp "httpServer/server/pkg/http"

	_ "httpServer/server/docs"

	_ "github.com/lib/pq"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/chi/v5"
)

// @title My API
// @version 1.0
// @description http server
// @host localhost:8080
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath /

func main() {
	// addr := flag.String("addr", ":8080", "address")
	httpFlags := config.ParseFlags()
	var cfg config.ServerConfig
	config.MustLoad(httpFlags.HTTPConfigPath, &cfg)

	// Сессии
	// sessionsRepo, err := ramrepository.NewSessionsRepository()
	// sessionsRepo, err := redis.NewSessionsRepository("redis:6379", "password", 0, 24*time.Hour)
	sessionsRepo, err := redis.NewSessionsRepository(cfg.Redis)
	if err != nil {
		log.Fatalf("failed to create sessions repository: %v", err)
	}
	sessionsService := service.NewSessions(sessionsRepo)

	// Пользователи
	// usersRepo := ramrepository.NewUsersRepository()
	// usersRepo, err := postgres.NewUsersRepository("postgres://user:password@data_base:5432/postgres_db?sslmode=disable")
	usersRepo, err := postgres.NewUsersRepository(cfg.Postgres)
	if err != nil {
		log.Fatalf("failed to create users repository: %v", err)
	}
	usersService := service.NewUsers(usersRepo)
	usersNewHandler := http.NewUsersHandler(usersService, sessionsService)

	// Таски
	tasksSender, err := rabbitMQ.NewRabbitMQSender(cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("failed to create tasks sender: %v", err)
	}

	// tasksRepo, err := ramrepository.NewTasks()
	// tasksRepo, err := postgres.NewTasksRepository("postgres://user:password@data_base:5432/postgres_db?sslmode=disable")
	tasksRepo, err := postgres.NewTasksRepository(cfg.Postgres)
	if err != nil {
		log.Fatalf("failed to create tasks repository: %v", err)
	}
	tasksService := service.NewTasks(tasksRepo, tasksSender)
	tasksNewHandler := http.NewTasksHandler(tasksService, sessionsService)

	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	http.Health(r)
	tasksNewHandler.WithTasksHandlers(r)
	usersNewHandler.WithUsersHandlers(r)

	log.Printf("Starting server on %s", cfg.HTTPConfig.Address)
	if err := pkgHttp.CreateAndRunServer(r, cfg.HTTPConfig); err != nil {
		log.Fatalf("Failed to start server %v", err)
	}
}
