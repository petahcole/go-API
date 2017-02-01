package main

import(
    "github.com/gin-gonic/gin"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"

    "net/http"
    "os"
)

type Book struct {
  ID    bson.ObjectId `json:"id" bson:"_id,omitempty"`
  Title       string
  Author      string
  Genre       string
}

func main() {
  lab := os.Getenv("MONGOLAB_URI")
  db := os.Getenv("MyBooks")

  session, err := mgo.Dial(lab)
  col:= session.DB(db).C("books")
  if err != nil {
    panic(err)
  }

  port := os.Getenv("PORT")
  if port == "" {
    port = "3000"
  }

  r := gin.Default()

  r.LoadHTMLGlob("*.html")

  r.Static("/public", "public")


  r.GET("/books", func(c *gin.Context) {
      var books []Book
      col.Find(nil).All(&books)
      c.JSON(http.StatusOK, gin.H{
			"books": books,
		})
  })

  r.POST("/new", func(c *gin.Context) {
    title := c.PostForm("title")
    author := c.PostForm("author")
    genre := c.PostForm("genre")
    err = col.Insert(&Book{Title: title, Author: author, Genre: genre})
    if err != nil {
      panic(err)
    }
    c.JSON(http.StatusOK, gin.H{
      "holy shit": "it worked",
    })
  })

  r.GET("/book/:id", func(c *gin.Context)  {
      result := Book{}
      title := c.Param("id")
      query := bson.M{"title": title}
      err := col.Find(query).One(&result)
      if err != nil {
        panic(err)
      }
      c.JSON(http.StatusOK, gin.H{
        "data": result,
      })
  })

  r.POST("/delete/:id", func(c *gin.Context)  {
    title := c.Param("id")
    err := col.Remove(bson.M{"title": title})
    if err != nil {
      panic(err)
    }
    c.JSON(http.StatusOK, gin.H{
      "deleted?": "yep",
    })
  })

  r.POST("/edit/:id", func(c *gin.Context)  {
    author := c.PostForm("author")
    genre := c.PostForm("genre")
    update := bson.M{
      "author": author,
      "genre": genre,
    }
    change := bson.M{"$set": update}
    title := c.Param("id")
    col.Update(bson.M{"title": title}, change)
    c.JSON(http.StatusOK, gin.H{
      "updated?": "totally",
    })
  })



  r.Run(":" + port)
}
