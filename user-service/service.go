package userservice

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	DB *gorm.DB
}

func UserServiceRoute(engine *gin.Engine, db *gorm.DB) {
	controller := &Controller{DB: db}
	app := engine.Group("/user")
	app.POST("", controller.createUser)
	app.GET("/:id", controller.getUserById)
	app.GET("/login/:email", controller.getUserByEmail)
}

func (controller *Controller) createUser(ctx *gin.Context) {
	var user User
	ctx.ShouldBindJSON(&user)
	controller.DB.Create(&user)
	ctx.JSON(200, "succeed")
}

func (controller *Controller) getUserById(ctx *gin.Context) {
	x := ctx.Param("id")
	id, _ := strconv.Atoi(x)
	var user User
	err := controller.DB.Where("id = ?", id).First(&user).Error
	fmt.Println(id, err, user)
	ctx.JSON(200, &user)
}
func (controller *Controller) getUserByEmail(ctx *gin.Context) {
	email := ctx.Param("email")
	var user User
	err := controller.DB.Where("email = ?", email).First(&user).Error
	fmt.Println(email, err, user)
	ctx.JSON(200, &user)
}
