package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kalaiselvan07/todo/token"
)

type todo struct {
	Id        string `json:"id"`
	Item      string `json:"item"`
	Completed bool   `json:"completed"`
}

var todos = []todo{
	{Id: "1", Item: "read book", Completed: false},
	{Id: "2", Item: "play football", Completed: false},
	{Id: "3", Item: "listen music", Completed: false},
}

func getTodos(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, todos)
}

func getTodosById(context *gin.Context) {
	newTodo, err := getTodo(context.Param("id"))
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "todo not found"})
		return
	}
	context.IndentedJSON(http.StatusOK, newTodo)
}

func getTodo(id string) (*todo, error) {
	for i, t := range todos {
		if t.Id == id {
			return &todos[i], nil
		}
	}
	return nil, errors.New("todo not found")
}

func addTodos(context *gin.Context) {
	var newTodo todo
	if err := context.BindJSON(&newTodo); err != nil {
		return
	}
	todos = append(todos, newTodo)
	context.IndentedJSON(http.StatusCreated, newTodo)
}

func updateTodosById(context *gin.Context) {
	id := context.Param("id")
	newTodo, err := getTodo(id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "todo not found"})
		return
	}
	newTodo.Completed = !newTodo.Completed
	context.IndentedJSON(http.StatusOK, newTodo)
}

func login(context *gin.Context) {
	var loginObj models.LoginRequest
	if err := context.ShouldBindJSON(&loginObj); err != nil {
		log.Fatal("Bad request")
	}

	var claims = &models.JwtClaims{}
	claims.CompanyId = "companyId"
	claims.Username = loginObj.Username
	claims.Roles = []int{1, 2}

	var tokenCreationTime = time.Now().UTC()
	var expirationTime = tokenCreationTime.Add(time.Duration(10) * time.Minute)
	tokenString, err := token.GenrateToken(claims, expirationTime)

	if err != nil {
		log.Fatal("Token not created")
	}

	context.AbortWithStatusJSON(http.StatusOK, models.Response{
		Data:    tokenString,
		Status:  http.StatusOK,
		Message: "Token created",
	})
}

func main() {
	router := gin.Default()
	router.POST("/login", login)
	router.GET("/todos", getTodos)
	router.GET("/todos/:id", getTodosById)
	router.POST("/todos", addTodos)
	router.PATCH("/todos/:id", updateTodosById)
	router.Run("localhost:8080")
}
