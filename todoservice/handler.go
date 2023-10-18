package todoservice

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kalaiselvan07/todo/pgdatabase"
)

type Todo struct {
	Id        int    `json:"id"`
	Item      string `json:"item"`
	Completed bool   `json:"completed"`
}

func GetTodos(c *gin.Context) {
	rows, err := pgdatabase.DB.Query("SELECT id, item, completed FROM todos")
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve todos"})
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.Id, &t.Item, &t.Completed); err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve todos"})
			return
		}
		todos = append(todos, t)
	}

	c.IndentedJSON(http.StatusOK, todos)
}

func GetTodoById(c *gin.Context) {
	id := c.Param("id")

	var t Todo
	err := pgdatabase.DB.QueryRow("SELECT id, item, completed FROM todos WHERE id = $1", id).Scan(&t.Id, &t.Item, &t.Completed)
	if err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "todo not found"})
		} else {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve todo"})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, t)
}

func AddTodo(c *gin.Context) {
	var t Todo
	if err := c.BindJSON(&t); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	// Assuming id is auto-incremented in the database
	err := pgdatabase.DB.QueryRow("INSERT INTO todos(item, completed) VALUES($1, $2) RETURNING id", t.Item, t.Completed).Scan(&t.Id)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to add todo"})
		return
	}

	c.IndentedJSON(http.StatusCreated, t)
}

func UpdateTodoById(c *gin.Context) {
	id := c.Param("id")

	var t Todo
	err := pgdatabase.DB.QueryRow("UPDATE todos SET completed = NOT completed WHERE id = $1 RETURNING id, item, completed", id).Scan(&t.Id, &t.Item, &t.Completed)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to update todo"})
		return
	}

	c.IndentedJSON(http.StatusOK, t)
}
