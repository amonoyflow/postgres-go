package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	guuid "github.com/google/uuid"
)

var dbConnect *pg.DB

// Todo model
type Todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Completed string    `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateTodoTable function
func CreateTodoTable(db *pg.DB) error {
	opts := &orm.CreateTableOptions{
		IfNotExists: true,
	}

	createError := db.CreateTable(&Todo{}, opts)

	if createError != nil {
		log.Printf("Error while creating todo table, Reason: %v\n", createError)
		return createError
	}

	log.Printf("Todo table created")
	return nil
}

// InitiateDB function
func InitiateDB(db *pg.DB) {
	dbConnect = db
}

// GetAllTodos function
func GetAllTodos(c *gin.Context) {
	var todos []Todo
	err := dbConnect.Model(&todos).Select()

	if err != nil {
		log.Printf("Error while getting all todos, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "All Todos",
		"data":    todos,
	})
	return
}

// CreateTodo function
func CreateTodo(c *gin.Context) {
	var todo Todo
	c.BindJSON(&todo)

	title := todo.Title
	body := todo.Body
	completed := todo.Completed
	id := guuid.New().String()

	insertError := dbConnect.Insert(&Todo{
		ID:        id,
		Title:     title,
		Body:      body,
		Completed: completed,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if insertError != nil {
		log.Printf("Error while inserting new todo into db, Reason: %v\n", insertError)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Todo created successfullt",
	})
	return
}

// GetSingleTodo function
func GetSingleTodo(c *gin.Context) {
	todoID := c.Param("todoId")
	todo := Todo{ID: todoID}

	err := dbConnect.Select(todo)

	if err != nil {
		log.Printf("Error while getting single todo, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Todo not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Single todo",
		"data":    todo,
	})
	return
}

// EditTodo function
func EditTodo(c *gin.Context) {
	todoID := c.Param("todoId")
	var todo Todo
	c.BindJSON(&todo)
	completed := todo.Completed

	_, err := dbConnect.Model(&Todo{}).Set("completed = ?", completed).Where("id = ?", todoID).Update()

	if err != nil {
		log.Printf("Error, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Todo edited successfully",
	})
	return
}

// DeleteTodo function
func DeleteTodo(c *gin.Context) {
	todoID := c.Param("todoId")
	todo := &Todo{ID: todoID}

	err := dbConnect.Delete(todo)

	if err != nil {
		log.Printf("Error while deleting a single todo, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Todo deleted successfully",
	})
	return
}
