package main

import (
	"github.com/gin-gonic/gin"
	authn "github.com/kalaiselvan07/todo/authentication"
	authz "github.com/kalaiselvan07/todo/authorization"
	database "github.com/kalaiselvan07/todo/pgdatabase"
	todo "github.com/kalaiselvan07/todo/todoservice"
)

func main() {
	database.Init()
	defer database.DB.Close()
	router := setupRouter()
	router.Run("localhost:8080")
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	
	router.POST("/login", authn.Login)

	api := router.Group("/todos")
	{
		api.Use(authn.ValidateToken())

		api.GET("/", todo.GetTodos)
		api.GET("/:id", todo.GetTodoById)

		api.POST("/", authz.Authorization("admin"), todo.AddTodo)
		api.PATCH("/:id", authz.Authorization("admin"), todo.UpdateTodoById)
	}

	return router
}
