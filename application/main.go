package main

import (
	"log"
	hd "application/handler"
	"github.com/gin-gonic/gin"
)


func main() {

	if err := hd.InitDB("todos.db"); err != nil {
		log.Fatalf("Error initializing the database: %v\n", err)
	}
	defer hd.CloseDB()

	router := gin.Default()

	router.POST("/login", hd.Authenticate)

	api := router.Group("/todo")
	api.Use(hd.ValidateToken())

	api.GET("/", hd.GetTodos)
	api.GET("/:id", hd.GetTodoById)
	api.POST("/", hd.AuthorizeRequest("admin"), hd.AddTodo)
	api.PATCH("/:id", hd.AuthorizeRequest("admin"), hd.UpdateTodoById)

	if err := router.Run(":8087"); err != nil {
		log.Fatalf("Error starting the server: %v\n", err)
	}
}
