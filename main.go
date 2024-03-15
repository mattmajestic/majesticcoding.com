package main

import (
    "github.com/gin-gonic/gin"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "net/http"
)

type User struct {
    gorm.Model
    Name string
}

func main() {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    // Migrate the schema
    db.AutoMigrate(&User{})

    router := gin.Default()

    router.LoadHTMLGlob("templates/*")

    router.GET("/", func(c *gin.Context) {
        var users []User
        db.Find(&users)
        c.HTML(http.StatusOK, "index.tmpl", gin.H{
            "users": users,
        })
    })

    router.Run(":8080")
}