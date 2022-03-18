package main

import (
	"fmt"
	"hw2/models"
	"os"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

var booksRepo *models.BookDB

var db *gorm.DB

func init() {

	var err error

	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=* password=*")

	if err != nil {
		panic(err)
	}

	booksRepo = models.NewBookRepo(db)

	booksRepo.Setup()
	// booksRepo.InsertData()

}

var books models.Book

func main() {

	input := os.Args

	firstArg := strings.ToLower(input[1])

	switch firstArg {
	case "search":
		secondArg := strings.ToLower(input[2])
		bookResult := booksRepo.GetBookByName(secondArg)

		for _, v := range bookResult {
			fmt.Println(v.Name)
		}

	case "list":
		bookList := booksRepo.GetAllBook()
		for _, v := range bookList {
			fmt.Println(v.Name)
		}
	case "delete":
		secondArg, _ := strconv.Atoi(input[2])
		booksRepo.DeleteByID(secondArg)

	case "buy":
		secondArg, _ := strconv.Atoi(input[2])
		thirdArg, _ := strconv.Atoi(input[3])
		booksRepo.UpdateStock(secondArg, thirdArg, books)

	}

}
