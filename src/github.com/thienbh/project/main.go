package main

import (
	// "time"
	// "context"
	// "os/signal"
	// "os"
	// "net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"	
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

type Note struct {
	gorm.Model
	Title string
	Status bool
}

type Paging struct {
	Page uint `form:"p"`
	Limit uint `form:"l"`
}

func main() {
	db, _ := gorm.Open("mysql", "root:123456789@/golang?charset=utf8&parseTime=True&loc=Local")
	db.AutoMigrate(&Note{

	})
	//defer db.Close()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(201, "pong")
	})

	r.POST("/note", func(c *gin.Context) {
		var note Note
		c.ShouldBind(&note)
		db.Create(&note)
		// c.JSON(200, c.JSON)
		c.JSON(200, note)
	})

	r.PUT("/note/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var note Note
		// tim theo id
		db.Where("id = ?", id).First(&note)
		c.ShouldBind(&note)
		note.ID = uint(id) // convert string -> uint truoc
		// update vo db
		err := db.Model(&note).Update(&note).Error
		if err != nil {
			c.AbortWithStatus(404)
			return
		}

		c.JSON(200, note)
	})

	r.DELETE("/note/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var note Note
		c.ShouldBind(&note)
		note.ID = uint(id) // convert string -> uint truoc
		// update vo db
		err := db.Model(&note).Delete(id).Error
		if err != nil {
			c.AbortWithStatus(404)
			return
		}

		c.JSON(200, note)
	})

	r.GET("/note", func(c *gin.Context) {
		var pager Paging
		c.ShouldBindQuery(&pager)

		// offset := pager.GetOffSet(pager)
		// limit := pager.GetLimit(pager)

		offset := 1
		limit := 1

		var notes []Note
		db.Offset(offset).Limit(limit).Find(&notes)

		c.JSON(200, notes)
	})

	// srv := &http.Server{
	// 	Addr: ":8080",
	// 	Handler: r,
	// }

	// go func() {
	// 	srv.ListenAndServe()
	// }()

	// quit := make(chan os.Signal)
	// signal.Notify(quit, os.Interrupt)
	// <-quit

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// srv.Shutdown(ctx)
	r.Run(":8081") 
	// r.Run() // listen and serve on 0.0.0.0:8080
}

//curl -XPOST -H "Content-type: application/json" http://localhost:8080/note -d '{"title": "test2", "status": false}'