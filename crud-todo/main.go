package main

import (
	"net/http"
	"errors"

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
func getTodo(c *gin.Context) {
	id := c.Param("id")
	todo, err := getTodoById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message":"todo not found"})
	}
	c.IndentedJSON(http.StatusOK,todo)


}
func getTodoById(id string)(*todo,error){
	for i,t:=range todos{
		if t.ID == id {
			return &todos[i],nil
		}
	}
	return nil, errors.New("todo not foud")
}
func toggleTodoStatus(c *gin.Context){
	id := c.Param("id")
	todo, err := getTodoById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message":"todo not found"})
	}
	todo.Completed = !todo.Completed
	c.IndentedJSON(http.StatusOK,todo)

}


func main() {
	router := gin.Default()
	router.GET("/todos",getTodos)
	router.GET("/todos/:id",getTodo)
	router.PATCH("/todos/:id",toggleTodoStatus)
	router.POST("/todos",addTodos)
	router.Run("localhost:9090")
}