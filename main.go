package main

import (
	"flag"
	"httpServer/API/http"
	ramrepository "httpServer/repository/ram_repository"
	"httpServer/usecases/service"
	"log"

	pkgHttp "httpServer/pkg/http"

	_ "httpServer/docs"

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
	addr := flag.String("addr", ":8080", "address")

	usersRepo := ramrepository.NewUsersRepository()
	usersService := service.NewUsers(usersRepo)
	usersNewHandler := http.NewUsersHandler(usersService)

	tasksRepo := ramrepository.NewTasks()
	tasksService := service.NewTasks(tasksRepo)
	tasksNewHandler := http.NewTasksHandler(tasksService, usersService)

	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	tasksNewHandler.WithTasksHandlers(r)
	usersNewHandler.WithUsersHandlers(r)

	log.Printf("Starting server on %s", *addr)
	if err := pkgHttp.CreateAndRunServer(r, *addr); err != nil {
		log.Fatalf("Failed to start server %v", err)
	}
}
