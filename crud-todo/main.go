package main

import (
	"errors"
	"net/http"

	"sync"
	
    "time"
	"strings"
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
	Password string     `json:"password"`
}
var usersDB = make(map[string]string)
var mu      sync.RWMutex
var jwtSecret = []byte("secret_key")
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
		return
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
		return
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
		return
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
		return
	}
	usersDB[req.Username] = string(hashedPassword)
	c.IndentedJSON(http.StatusCreated,gin.H{"message":"Successfully Registered"})
}

func login(c *gin.Context) {
	var req user
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON"})
		return
	}
	mu.RLock()
	storedHash, exists := usersDB[req.Username]
	mu.RUnlock()
	if !exists {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Invalid Credentials"})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(storedHash),[]byte(req.Password))
	if err!=nil{
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Invalid Credentials"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{"username":req.Username,
"exp": time.Now().Add(time.Hour * 24).Unix()})
    
tokenString,err := token.SignedString(jwtSecret)
if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error Generating Token"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"token": tokenString})

}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Next()
	}
}



func main() {
	router := gin.Default()
	router.POST("/register",register)
	router.POST("/login", login)


	protected := router.Group("/")
	protected.Use(authMiddleware())
	{
	protected.GET("/todos",getTodos)
	protected.GET("/todos/:id",getTodo)
	protected.PATCH("/todos/:id",toggleTodoStatus)
	protected.POST("/todos",addTodos)
	}

	router.Run("localhost:9090")
}