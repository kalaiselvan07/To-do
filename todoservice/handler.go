package todoservice

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
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

func GetTodos(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, todos)
}

func GetTodosById(context *gin.Context) {
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

func AddTodos(context *gin.Context) {
	var newTodo todo
	if err := context.BindJSON(&newTodo); err != nil {
		return
	}
	todos = append(todos, newTodo)
	context.IndentedJSON(http.StatusCreated, newTodo)
}

func UpdateTodosById(context *gin.Context) {
	id := context.Param("id")
	newTodo, err := getTodo(id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "todo not found"})
		return
	}
	newTodo.Completed = !newTodo.Completed
	context.IndentedJSON(http.StatusOK, newTodo)
}
