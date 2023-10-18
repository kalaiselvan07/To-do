package main

import (
	"github.com/gin-gonic/gin"
	authn "github.com/kalaiselvan07/todo/authentication"
	authz "github.com/kalaiselvan07/todo/authorization"
	todo "github.com/kalaiselvan07/todo/todoservice"
)

func main() {

	router := gin.Default()
	router.POST("/login", authn.Login)

	api := router.Group("/todos")
	api.Use(authn.ValidateToken())

	api.GET("/", todo.GetTodos)
	api.GET("/:id", todo.GetTodosById)

	api.POST("/", authz.Authorization("admin"), todo.AddTodos)
	api.PATCH("/:id", authz.Authorization("admin"), todo.UpdateTodosById)

	router.Run("localhost:8080")
}
