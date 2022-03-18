package models

import "github.com/jinzhu/gorm"

type Book struct {
	gorm.Model
	StockNumber int
	PageNumber  int
	Price       float64
	Name        string
	StockCode   string
	Isbn        string
	Author      string
}
