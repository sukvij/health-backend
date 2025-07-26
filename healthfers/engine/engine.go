package engine

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sukvij/health-checker/backend/healthfers/database"
	"github.com/sukvij/health-checker/backend/healthfers/logs"
	"github.com/sukvij/health-checker/backend/healthfers/tracing"
	"go.opentelemetry.io/otel/sdk/trace"
	"gorm.io/gorm"
)

type Engine struct {
	Ctx     *gin.Context
	DB      *gorm.DB
	Logs    *logs.AgreeGateLoager
	Tracker *trace.TracerProvider
	App     *gin.Engine
}

func MakeNewEngine() *Engine {
	engine := &Engine{}
	db, dbConnError := database.Connection()
	engine.DB = db
	if dbConnError != nil {
		fmt.Println("problem with database connections... in engine file")
		return nil
	}
	engine.Logs = logs.NewAgreeGateLogger()

	engine.Tracker = tracing.InitTracer()
	engine.App = gin.Default()
	return engine
}
