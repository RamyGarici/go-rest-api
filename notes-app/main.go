package main

import (
	"net/http"
	"time"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)


type note struct{
	ID int `json:"id"`
	Title string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
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
		c.IndentedJSON(http.StatusBadRequest,gin.H{"message":"Title and Content required"})
		return
	}
	newNote.ID = len(notes)+1
	
	newNote.CreatedAt = time.Now()
	newNote.UpdatedAt = time.Now()
	notes = append(notes, newNote)
	c.IndentedJSON(http.StatusCreated,newNote)
}

func getNote(c *gin.Context) {
	id,err := strconv.Atoi(c.Param("id"))
	if err != nil{
		c.IndentedJSON(http.StatusBadRequest,gin.H{"message":"invalid id"})
		return
	}
	note, err := getNoteById(id)
	if err !=nil{
		c.IndentedJSON(http.StatusNotFound,gin.H{"message":"Note not found"})
		return
	}
	c.IndentedJSON(http.StatusOK,note)
}

func getNoteById(id int) (*note, error){
	for i,n := range notes{
		if n.ID == id{
			return &notes[i],nil
		}
	}
	err := errors.New("Note not found")
	return nil,err

}





func main() {
	router := gin.Default()
	router.GET("/notes",getNotes)
	router.POST("/notes",addNotes)
	router.GET("/notes/:id",getNote)

	router.Run(":8080")




}