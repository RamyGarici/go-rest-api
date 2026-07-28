package main
import("github.com/gin-gonic/gin"
"net/http"
"time")


type note struct{
	ID int `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

}







func main() {




}