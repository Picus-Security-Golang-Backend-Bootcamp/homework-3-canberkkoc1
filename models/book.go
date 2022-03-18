package models

import (
	"encoding/csv"
	"hw2/helper"
	"log"
	"os"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

type BookDB struct {
	db *gorm.DB
}

func (g *BookDB) InsertData() {

	var books []Book

	file, err := os.Open("/media/canberk/hdd1/HW/HW-3/booklist.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	text, err := csv.NewReader(file).ReadAll()

	if err != nil {

		log.Fatal(err)

	}

	for _, record := range text {

		books = append(books, Book{

			StockNumber: helper.RandomNumber(1, 100),
			PageNumber:  helper.RandomNumber(1, 100),
			Price:       helper.RandomFloat(1, 100),
			Name:        record[3],
			StockCode:   helper.RandomString(5),
			Isbn:        helper.RandomString(5),
			Author:      record[6],
		})

	}

	for _, v := range books[1:] {
		g.db.Create(&v)
	}

}

func NewBookRepo(db *gorm.DB) *BookDB {
	return &BookDB{
		db: db,
	}
}

func (g *BookDB) Setup() {
	g.db.AutoMigrate(&Book{})

}

func (g *BookDB) GetAllBook() []Book {

	var bookList []Book
	g.db.Find(&bookList)

	return bookList

}
func (g *BookDB) DeleteByID(id int) {
	var books []Book

	var n []int

	g.db.Model(&books).Pluck("id", &n)

	g.db.Unscoped().Delete(&books, id)

	isDeleted := helper.CheckSlice(n, id)

	if !isDeleted {
		panic("id not found")
	}

}

func (g *BookDB) GetBookByName(name string) []Book {

	var book []Book
	g.db.Where(" Name LIKE ?", "%"+name+"%").Find(&book)

	return book

}

func (g *BookDB) UpdateStock(id, stock int, book Book) {

	var n []int

	g.db.Model(&book).Pluck("stock_number", &n)

	stoc_num := n[id-1]

	if stoc_num <= 0 || stoc_num < stock {
		panic("not in stock")
	}

	newStock := stoc_num - stock

	g.db.Model(&book).Where("id = ?", id).Update("stock_number", newStock)

}
