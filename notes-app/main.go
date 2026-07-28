package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)


type note struct{
	ID int `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

}
var notes []note

func getNotes(c *gin.Context) {
	c.IndentedJSON(http.StatusOK,notes)
}

func addNotes(c *gin.Context) {
	var newNote note
	err := c.BindJSON(&newNote)
	if err != nil{
		c.IndentedJSON(http.StatusBadRequest,gin.H{"message":"Error adding the note"})
		return
	}
	notes = append(notes, newNote)
	c.IndentedJSON(http.StatusCreated,newNote)
}





func main() {
	router := gin.Default()
	router.GET("/notes",getNotes)
	router.POST("/notes",addNotes)
	// router.GET("/notes/:id",getNote)





}