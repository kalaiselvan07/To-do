package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kalaiselvan07/todo/authentication"
	"github.com/kalaiselvan07/todo/authorization"
	"github.com/kalaiselvan07/todo/pgdatabase"
	"github.com/kalaiselvan07/todo/todoservice"
)

func main() {
	// Initialize the database
	dbConfig := pgdatabase.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "mysecretpassword",
		DBName:   "test_db",
	}

	if err := pgdatabase.InitDB(dbConfig); err != nil {
		log.Fatalf("Error initializing the database: %v\n", err)
	}
	defer pgdatabase.CloseDB()

	// Create a new Gin router
	router := gin.Default()

	// Authentication route
	router.POST("/login", authentication.Login)

	// Todo API routes with token validation and authorization
	api := router.Group("/todos")
	api.Use(authentication.ValidateToken())

	api.GET("/", todoservice.GetTodos)
	api.GET("/:id", todoservice.GetTodoById)
	api.POST("/", authorization.Authorization("admin"), todoservice.AddTodo)
	api.PATCH("/:id", authorization.Authorization("admin"), todoservice.UpdateTodoById)

	// Start the Gin server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting the server: %v\n", err)
	}
}
