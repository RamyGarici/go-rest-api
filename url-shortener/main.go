package main 

import (
    "github.com/jackc/pgx/v5"
    "github.com/gin-gonic/gin"
    "time"
    "net/http"
    "net/url"
)

type URL struct{
    ID int
    ShortCode string
    OriginalURL string
    CreatedAt time.Time

}
type request struct{
    Url string `json:"url"`
}

func shortenURL(c *gin.Context) {
    var request body
    err := c.ShouldBindJSON(&request)
	if err!=nil{
		c.IndentedJSON(http.StatusBadRequest,gin.H{"message":"Invalid body"})
        return
	}
    _,err = url.ParseRequestURI(request.Url)
    if err != nil{
        c.IndentedJSON(http.StatusBadRequest,gin.H{"message":"Invalid URL"})
        return
    }

    

}
 



func main() {
   router := gin.Default()
   router.POST("/shorten",shortenURL)
   router.Run(":8080")
}
