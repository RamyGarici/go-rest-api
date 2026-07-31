package main

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type URL struct{
    ID int
    ShortCode string
    OriginalURL string
    CreatedAt time.Time

}

type request struct{
    URL string `json:"url"`
}
type apiConfig struct {
    DB *pgx.Conn
}

func(cfg *apiConfig) shortenURL(c *gin.Context) {
    var request request
    err := c.ShouldBindJSON(&request)
	if err!=nil{
		c.IndentedJSON(http.StatusBadRequest,gin.H{"message":"Invalid body"})
        return
	}
    _,err = url.ParseRequestURI(request.URL)
    if err != nil{
        c.IndentedJSON(http.StatusBadRequest,gin.H{"message":"Invalid URL"})
        return
    }
    var id int
    err = cfg.DB.QueryRow(
        c.Request.Context(),
        "INSERT INTO urls(original_url) VALUES ($1) RETURNING id",
        request.URL,
    ).Scan(&id)

    if err!=nil{
        c.JSON(http.StatusInternalServerError,gin.H{"message":"Database error"})
        return
    }
    result := convertToBase62(id)
   _, err = cfg.DB.Exec(
    c.Request.Context(),
    "UPDATE urls SET short_code = $1 WHERE id = $2",
    result,
    id,
)
if err!=nil{
        c.JSON(http.StatusInternalServerError,gin.H{"message":"Database error"})
        return
    }
   c.IndentedJSON(http.StatusOK,gin.H{"short_url": "http://localhost:8080/"+result})

    

}
 func convertToBase62(id int)string{

    chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    result := ""
    for id > 0 {
        remainder := id % 62
       
        char := chars[remainder]
         result = string(char) + result
        id = id/62


}
    return result
 }

func (cfg *apiConfig) redirectURL(c *gin.Context) {
    code := c.Param("code")
    var OriginalURL string
    err:= cfg.DB.QueryRow(
        c.Request.Context(),
        "SELECT original_url FROM urls WHERE short_code = $1",
        code,
    ).Scan(&OriginalURL)
    if err!=nil{
        c.JSON(http.StatusNotFound,gin.H{"message":"URL not found"})
        return
    }
    c.Redirect(http.StatusMovedPermanently,OriginalURL)
}

   




func main() {
 connStr := os.Getenv("DATABASE_URL")
  conn, err := pgx.Connect(context.Background(),connStr)
 if err!=nil{
    return
 }
 cfg := apiConfig{
    DB: conn,
 }

 defer conn.Close(context.Background())








   router := gin.Default()
   router.POST("/shorten",cfg.shortenURL)
   router.GET("/:code",cfg.redirectURL)
   router.Run(":8080")
}
