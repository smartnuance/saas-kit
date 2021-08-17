package auth

import (
	"database/sql"

	_ "github.com/golang-migrate/migrate/v4"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

//lint:file-ignore U1000 Ignore unused fields, gorm model fields are used to create migrations
type (
	Instance struct {
		gorm.Model
		Name string
		url  string
	}

	User struct {
		gorm.Model
		Name        string
		Email       *string
		Username    sql.NullString
		ActivatedAt sql.NullTime
	}

	Profile struct {
		gorm.Model
		UserID     int
		User       User
		InstanceID int
		Instance   Instance
		Roles      pq.StringArray `gorm:"type:text[]"`
	}
)

// // createTodo add a new todo
// func createTodo(c *gin.Context) {
// 	completed, _ := strconv.Atoi(c.PostForm("completed"))
// 	todo := todoModel{Title: c.PostForm("title"), Completed: completed}
// 	db.Save(&todo)
// 	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Todo item created successfully!", "resourceId": todo.ID})
// }

// // fetchAllTodo fetch all todos
// func fetchAllTodo(c *gin.Context) {
// 	var todos []todoModel
// 	var _todos []transformedTodo

// 	db.Find(&todos)

// 	if len(todos) <= 0 {
// 		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
// 		return
// 	}

// 	//transforms the todos for building a good response
// 	for _, item := range todos {
// 		completed := false
// 		if item.Completed == 1 {
// 			completed = true
// 		} else {
// 			completed = false
// 		}
// 		_todos = append(_todos, transformedTodo{ID: item.ID, Title: item.Title, Completed: completed})
// 	}
// 	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todos})
// }

// // fetchSingleTodo fetch a single todo
// func fetchSingleTodo(c *gin.Context) {
// 	var todo todoModel
// 	todoID := c.Param("id")

// 	db.First(&todo, todoID)

// 	if todo.ID == 0 {
// 		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
// 		return
// 	}

// 	completed := false
// 	if todo.Completed == 1 {
// 		completed = true
// 	} else {
// 		completed = false
// 	}

// 	_todo := transformedTodo{ID: todo.ID, Title: todo.Title, Completed: completed}
// 	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todo})
// }

// // updateTodo update a todo
// func updateTodo(c *gin.Context) {
// 	var todo todoModel
// 	todoID := c.Param("id")

// 	db.First(&todo, todoID)

// 	if todo.ID == 0 {
// 		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
// 		return
// 	}

// 	db.Model(&todo).Update("title", c.PostForm("title"))
// 	completed, _ := strconv.Atoi(c.PostForm("completed"))
// 	db.Model(&todo).Update("completed", completed)
// 	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo updated successfully!"})
// }

// // deleteTodo remove a todo
// func deleteTodo(c *gin.Context) {
// 	var todo todoModel
// 	todoID := c.Param("id")

// 	db.First(&todo, todoID)

// 	if todo.ID == 0 {
// 		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
// 		return
// 	}

// 	db.Delete(&todo)
// 	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo deleted successfully!"})
// }
