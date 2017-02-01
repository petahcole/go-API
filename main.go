package main

import(
    "net/http"
    "os"
    "github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()

  r.LoadHTMLGlob("*.html")

  r.Static("/public", "public")

  port := os.Getenv("PORT")
  if port == "" {
    port = "3000"
  }

  r.GET("/", func(c *gin.Context) {
    c.HTML(http.StatusOK, "index.html", gin.H{
      "HelloMessage": "Yay Go!",
    })
  })

  r.Run(":" + port)
}
