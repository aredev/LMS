package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"fmt"
	"reflect"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type User struct {
	gorm.Model
	Name string `gorm:"not null" json:"name" binding:"required,max=150"`
	Email string `gorm:"unique;not null" json:"email" binding:"required,email"`
	Password string `gorm:"not null" json:"password" binding:"required,min=6,max=150"`
}

type Parent struct {
	User
	Students []Student `gorm:"foreignkey:ParentRefer"`
}

type Student struct {
	User
	ParentRefer uint
	ClassRefer uint
}

type Teacher struct {
	ClassRefer uint
}

type Class struct {
	gorm.Model
	Name string `gorm:"unique;not null" json:"name" binding:"required"`
	Teachers []Teacher `gorm:"foreignkey:ClassRefer"`
	Students []Student `gorm:"foreignkey:ClassRefer"`
	Assignments []Assignment `gorm:"foreignkey:ClassRefer"`
}

type Assignment struct {
	gorm.Model
	Title string `gorm:"not null" json:"title" binding:"required"`
	Description string `json:"description"`
	Deadline time.Time `json:"deadline"`
	ClassRefer uint `binding:"required"`
	SubmittedAssignments []SubmittedAssignment `gorm:"foreignkey:AssignmentRefer"`
}

type Grade struct {
	gorm.Model
	//Score ``
	Comment string `json:"comment"`
}

type SubmittedAssignment struct {
	gorm.Model
	StudentRefer uint
	AssignmentRefer uint
}

func main(){
	db, err := gorm.Open("sqlite3", "lms.db")
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}

	govalidator.SetFieldsRequiredByDefault(true)

	for _, model := range []interface{}{
		User{}, Parent{}, Student{}, Class{}, Assignment{},
	} {
		if err := db.AutoMigrate(model).Error; err != nil {
			fmt.Println(err)
		}else{
			fmt.Println("Auto migrating ", reflect.TypeOf(model).Name(), " ... ")
		}
	}


	router := gin.Default()

	router.POST("/user", func(c *gin.Context){
		var user User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}else{
			c.JSON(http.StatusOK, user)
		}
	})

	router.GET("/users", func(c *gin.Context){
		var users []User

		if err := db.Find(&users); err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, "")
		}else{
			c.JSON(http.StatusOK, users)
		}
	})

	router.POST("/class", func(c *gin.Context){
		var class Class
		if err := c.BindJSON(&class); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}else{
			db.Create(&class)
			c.JSON(http.StatusCreated, class)
		}
	})

	router.GET("/class", func(c *gin.Context){
		var classes []Class

		db.Find(&classes)
		c.JSON(http.StatusOK, classes)
	})

	router.POST("/assignment", func(c * gin.Context){
		var assignment Assignment

		if err := c.BindJSON(&assignment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}else{
			c.JSON(http.StatusOK, assignment)
		}
	})

	router.Run()




}