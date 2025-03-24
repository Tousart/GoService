package main

import (
	"httpServer/server/API/http"
	"httpServer/server/config"
	"httpServer/server/repository/rabbitMQ"
	ramrepository "httpServer/server/repository/ram_repository"
	"httpServer/server/usecases/service"
	"log"

	pkgHttp "httpServer/server/pkg/http"

	_ "httpServer/server/docs"

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
	sessionsRepo := ramrepository.NewSessionsRepository()
	sessionsService := service.NewSessions(sessionsRepo)

	// Пользователи
	usersRepo := ramrepository.NewUsersRepository()
	usersService := service.NewUsers(usersRepo)
	usersNewHandler := http.NewUsersHandler(usersService, sessionsService)

	// Таски
	tasksSender, err := rabbitMQ.NewRabbitMQSender(cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("failed %v", err)
	}
	tasksRepo := ramrepository.NewTasks()
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
