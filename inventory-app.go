package main

import (
	"errors"

	"encoding/csv"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Item struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	UnitPrice float32 `json:"unit_price"`
}

func main() {
	// Initialize a new Gin router instance
	router := gin.New()

	// Connect to the database
	db, err := gorm.Open(sqlite.Open("inventory.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto Migrate the inventory table
	db.AutoMigrate(&Item{})

	// Define the CRUD handlers
	router.POST("/item", func(c *gin.Context) {
		var item Item

		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if item.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID field missing from JSON body"})
			return
		}

		if item.Name == "" || item.Quantity == 0 || item.UnitPrice == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing fields in JSON body"})
			return
		}

		// Check if item with requested ID already exists
		var existingItem Item
		result := db.First(&existingItem, item.ID)
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusConflict, gin.H{"error": "item with requested ID already exists"})
			return
		}

		db.Create(&item)
		c.JSON(http.StatusCreated, item)

	})

	router.GET("/item", func(c *gin.Context) {
		var items []Item
		db.Find(&items)
		c.JSON(http.StatusOK, items)
	})

	router.GET("/item/:id", func(c *gin.Context) {
		var item Item
		if err := db.First(&item, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}
		c.JSON(http.StatusOK, item)
	})

	router.DELETE("/item/:id", func(c *gin.Context) {
		var item Item
		if err := db.First(&item, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}
		db.Delete(&item)
		c.JSON(http.StatusOK, gin.H{"message": "Record deleted"})
	})

	router.PATCH("/item/:id", func(c *gin.Context) {
		var item Item
		if err := db.First(&item, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Save(&item)
		c.JSON(http.StatusOK, item)
	})

	// Define the CSV export handler
	router.GET("/item/csv", func(c *gin.Context) {
		var items []Item
		db.Find(&items)

		writer := csv.NewWriter(c.Writer)
		defer writer.Flush()

		// Write CSV header
		writer.Write([]string{"ID", "Name", "Quantity", "Unit Price"})

		// Write CSV data
		for _, item := range items {
			writer.Write([]string{
				strconv.Itoa(int(item.ID)),
				item.Name,
				strconv.Itoa(item.Quantity),
				strconv.FormatFloat(float64(item.UnitPrice), 'f', 2, 32),
			})
		}

		// Set Content-Type header to tell the browser to download the file
		c.Writer.Header().Set("Content-Type", "text/csv")
		c.Writer.Header().Set("Content-Disposition", "attachment;filename=inventory.csv")
	})
	router.Run(":8080")
}
