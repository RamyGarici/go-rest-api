package main

import (
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type todo struct{
	ID        string    `json:"id"`
	Item       string   `json:"item"`
	Completed  bool     `json:"completed"`
}

type user struct{
	Username string     `json:"username"`
	Password string     `json:"-"`
}
var usersDB = make(map[string]string)
var mu      sync.RWMutex

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

func register(c *gin.Context) {
	var req user
	if err := c.BindJSON(&req); err!=nil{
		c.IndentedJSON(http.StatusBadRequest,gin.H{"message":"Invalid data"})
		return
	}
	if req.Username =="" || req.Password ==""{
		c.IndentedJSON(http.StatusBadRequest,gin.H{"message":"Empty Body"})
	}
	mu.Lock()
	defer mu.Unlock()
	if _,ok := usersDB[req.Username];ok{
			c.IndentedJSON(http.StatusConflict,gin.H{"message":"User already exists."})
			return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password),14)
	if err!=nil{
		c.IndentedJSON(http.StatusInternalServerError,gin.H{"message":"Error when hashing the password"})
	}
	usersDB[req.Username] = string(hashedPassword)
	c.IndentedJSON(http.StatusCreated,gin.H{"message":"Successfully Registered"})
}





func main() {
	router := gin.Default()
	router.GET("/todos",getTodos)
	router.GET("/todos/:id",getTodo)
	router.PATCH("/todos/:id",toggleTodoStatus)
	router.POST("/todos",addTodos)
	router.Run("localhost:9090")
}