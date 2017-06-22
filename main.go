package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
)

var (
	repeat int
)

func repeatHandler(c *gin.Context) {
	var buffer bytes.Buffer
	for i := 0; i < repeat; i++ {
		buffer.WriteString("Hello from Go!\n")
	}
	c.String(http.StatusOK, buffer.String())
}

func main() {
	var err error
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	tStr := os.Getenv("REPEAT")
	repeat, err = strconv.Atoi(tStr)
	if err != nil {
		log.Printf("Error converting $REPEAT to an int: %q - Using default\n", err)
		repeat = 5
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/mark", func(c *gin.Context) {
		c.String(http.StatusOK, string(blackfriday.MarkdownBasic([]byte("**hi!**"))))
	})

	router.GET("/repeat", repeatHandler)

	router.Run(":" + port)
}