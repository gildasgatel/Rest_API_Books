package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var books = map[int]book{}

type book struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", getBookById) // localhost:8080/books/2
	router.POST("/books", createBooks)
	router.PATCH("/checkout", checkout) // localhost:8080/checkout?id=2
	router.PATCH("/quantity", quantity) // localhost:8080/checkout?id=2&quantity=20
	router.DELETE("/books/:id", deleteBookById)
	router.PUT("/books", updateById)
	router.Run("localhost:8080")
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func createBooks(c *gin.Context) {
	var newBook book
	err := c.BindJSON(&newBook)
	handelError(err)
	newBook.Id = len(books) + 1
	books[newBook.Id] = newBook
	writeJsonFile()
}

func updateById(c *gin.Context) {
	var newBook book
	err := c.BindJSON(&newBook)
	handelError(err)
	books[newBook.Id] = newBook
	writeJsonFile()
}

func getBookById(c *gin.Context) {
	param := c.Param("id")
	c.IndentedJSON(http.StatusOK, bookById(param))
}

func deleteBookById(c *gin.Context) {
	param := c.Param("id")
	book := bookById(param)
	if book.Id == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not found"})
		return
	}
	delete(books, book.Id)
	writeJsonFile()
}

func bookById(param string) book {
	id, err := strconv.Atoi(param)
	handelError(err)
	return books[id]
}

func checkout(c *gin.Context) {
	param, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error missing parameters"})
		return
	}
	book := bookById(param)
	if book.Id == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not found"})
		return
	}
	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available"})
		return
	}
	book.Quantity -= 1
	books[book.Id] = book
	c.IndentedJSON(http.StatusOK, books[book.Id])
	writeJsonFile()
}

func quantity(c *gin.Context) {
	paramId, okI := c.GetQuery("id")
	paramQuantity, okQ := c.GetQuery("quantity")

	if !okI || !okQ {
		c.String(http.StatusBadRequest, "Error missing parameters")
		return
	}
	book := bookById(paramId)
	if book.Id == 0 {
		c.String(http.StatusBadRequest, "Book not found")
		return
	}
	quant, err := strconv.Atoi(paramQuantity)
	handelError(err)
	book.Quantity = quant
	books[book.Id] = book
	c.String(http.StatusOK, "Quantity updated")
	writeJsonFile()
}

// * * * TOOLS * * *

func writeJsonFile() {
	content, err := json.Marshal(books)
	handelError(err)
	err = ioutil.WriteFile(".file.json", content, 0644)
	handelError(err)
}
func handelError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// * * * INIT * * *

var booksInit = map[int]book{
	1: {Id: 1, Title: "test1", Author: "author1", Quantity: 1},
	2: {Id: 2, Title: "test2", Author: "author2", Quantity: 2},
	3: {Id: 3, Title: "test3", Author: "author3", Quantity: 3},
	4: {Id: 4, Title: "test4", Author: "author4", Quantity: 4},
}

func init() {
	file, err := ioutil.ReadFile(".file.json")
	if err != nil {
		content, err := json.Marshal(booksInit)
		handelError(err)
		ioutil.WriteFile(".file.json", content, 0644)
		initBooks(content)
	} else if len(file) > 0 {
		initBooks(file)
	}
}
func initBooks(file []byte) {
	err := json.Unmarshal(file, &books)
	handelError(err)
}
