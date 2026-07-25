package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type todo struct{
	ID        string    `json:"id"`
	Item       string   `json:"item"`
	Completed  bool     `json:"completed"`
}

var todos = []todo {
	{ID:"1",Item:"Read Book", Completed: false},
	{ID:"2",Item:"Clean Room", Completed: false},
	{ID:"3",Item:"Record video", Completed: false},
}

func getTodos(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, todos)

}
func addTodos(c *gin.Context) {
	var newTodo todo
	if err := c.BindJSON(&newTodo);err!=nil{
		return
	}
	todos = append(todos,newTodo)
	c.IndentedJSON(http.StatusCreated,newTodo)
	
}


func main() {
	router := gin.Default()
	router.GET("/todos",getTodos)
	router.POST("/todos",addTodos)
	router.Run("localhost:9090")
}