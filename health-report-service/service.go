package healthreportservice

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	DB *gorm.DB
}

func HealthReportServiceController(engine *gin.Engine, db *gorm.DB) {
	controller := &Controller{DB: db}
	app := engine.Group("/health-reports")
	app.POST("", controller.createHealthReport)
	app.GET("/:id", controller.getHealthReportsById)
}

func (controller *Controller) createHealthReport(ctx *gin.Context) {
	var healthReport HealthReport
	ctx.ShouldBindJSON(&healthReport)
	err := controller.DB.Create(&healthReport).Error
	ctx.JSON(200, err)
}

func (controller *Controller) getHealthReportsById(ctx *gin.Context) {
	x := ctx.Param("id")
	userId, _ := strconv.Atoi(x)
	var healthReport []HealthReport
	err := controller.DB.Where("user_id=?", userId).Order("updated_at desc").Find(&healthReport).Error
	fmt.Println(err)
	ctx.JSON(200, &healthReport)
}
