package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
)

var db *gorm.DB
var err error

func initDB() (*gorm.DB, error) {
	// db, err := gorm.Open("sqlite3", "./test.db")
	db, err := gorm.Open("mysql", "root@tcp(127.0.0.1:3306)/madrasa?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Student - the model for student
type Student struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Address   string `json:"address"`
}

func main() {
	db, err = initDB()
	if err != nil {
		log.Panicln(err)
	}
	defer db.Close()

	db.AutoMigrate(&Student{})

	r := gin.Default()
	r.GET("/students", GetStudents)
	r.GET("/students/:id", GetStudent)
	r.POST("/students", CreateStudent)
	r.PUT("/students/:id", UpdateStudent)
	r.DELETE("/students/:id", DeleteStudent)

	fmt.Println("Server started...")
	r.Run(":5050")
}

// GetStudents - READ ALL
func GetStudents(c *gin.Context) {
	var students []Student

	if err := db.Find(&students).Error; err != nil {
		c.AbortWithStatus(404)
		log.Panicln(err)
	} else {
		c.JSON(200, students)
	}
}

// GetStudent - Find by ID
func GetStudent(c *gin.Context) {
	var foundStd Student

	id := c.Params.ByName("id")
	if err := db.Where("id = ?", id).First(&foundStd).Error; err != nil {
		c.AbortWithStatus(404)
		log.Panicln(err)
	} else {
		c.JSON(200, foundStd)
	}

}

// CreateStudent - create student
func CreateStudent(c *gin.Context) {
	var student Student
	c.BindJSON(&student)

	db.Create(&student)
	c.JSON(200, student)
}

// UpdateStudent - update student data
func UpdateStudent(c *gin.Context) {
	var student Student
	id := c.Param("id")
	if err := db.Find(&student, id).Error; err != nil {
		c.AbortWithStatus(404)
		log.Panicln(err)
	}

	c.BindJSON(&student)
	db.Save(&student)
	c.JSON(200, student)
}

// DeleteStudent - delete student
func DeleteStudent(c *gin.Context) {
	id := c.Param("id")
	var student Student

	err := db.Delete(&student, id).Error
	if err != nil {
		c.AbortWithStatus(404)
		log.Panicln(err)
	}
	c.JSON(200, gin.H{"message": "ID " + id + " is deleted"})
}
