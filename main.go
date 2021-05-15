package main

import (
	"crypto/sha256"
	"encoding/base32"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var urls map[string]string
var hostname = "localhost:8080"

func home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main website",
	})
}

func shorten(c *gin.Context) {
	urlString, err := c.GetQuery("url")
	if err == false {
		log.Printf("Couldnt bind to request %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Bad request",
		})
		return
	}

	hash := sha256.Sum256([]byte(urlString))
	id := strings.ToLower(base32.StdEncoding.EncodeToString(hash[:])[:10])
	log.Printf("new id: %v", id)
	urls[id] = urlString

	c.HTML(http.StatusOK, "shortened.html", gin.H{
		"url": "http://" + hostname + "/get?url=" + id,
	})
}

func redirect(c *gin.Context) {
	uriRequest, err := c.GetQuery("url")
	if err == false {
		log.Printf("Couldnt bind to request %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Bad request",
		})
		return
	}
	redirectUrl := urls[uriRequest]
	log.Printf("Requestd url: %s", redirectUrl)
	c.Redirect(http.StatusMovedPermanently, redirectUrl)
}

func main() {
	urls = make(map[string]string)
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", home)
	r.GET("/shorten", shorten)
	r.GET("/get", redirect)
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
