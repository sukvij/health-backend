package response

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	redisservice "github.com/sukvij/health-checker/backend/healthfers/redis-service"
)

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type Meta struct {
	Version      string `json:"version"`
	LatencyMs    int64  `json:"latencyMs"`
	Environment  string `json:"environment,omitempty"`
	TotalRecords int64  `json:"total_records"`
	PageNumber   int    `json:"page_number"`
	Query        string `json:"query"`
	// Pagination  *Pagination `json:"pagination,omitempty"`
	// RateLimit   *RateLimit  `json:"rateLimit,omitempty"`
}

// type Pagination struct {
// 	Total      int64 `json:"total"`
// 	Page       int   `json:"page"`
// 	PerPage    int   `json:"perPage"`
// 	TotalPages int   `json:"totalPages"`
// }

// type RateLimit struct {
// 	Limit     int   `json:"limit"`
// 	Remaining int   `json:"remaining"`
// 	Reset     int64 `json:"reset"`
// }

type FinalResponse struct {
	Success    bool      `json:"success"`
	StatusCode int       `json:"statusCode"`
	Data       any       `json:"articles,omitempty"`
	Error      *AppError `json:"error,omitempty"`
	Meta       *Meta     `json:"meta"`
}

func JSONResponse(ctx *gin.Context, err error, data interface{}, redisClinet *redis.Client, cacheKey string, totalRecords int64, queryDetails string) {
	response := &FinalResponse{Data: data, Meta: &Meta{TotalRecords: totalRecords, Query: RemoveExtraFromQuery(queryDetails)}}
	if err == nil {
		response.Success = true
		response.StatusCode = 200
		if data != nil {
			err1 := redisservice.SetValueToRedis(redisClinet, cacheKey, *response)
			fmt.Println(cacheKey, err1)
		}
		ctx.JSON(response.StatusCode, response)
		return
	}
	ReturnErrorWithCode(ctx, err, response)
	if response.StatusCode == 204 {
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	ctx.JSON(response.StatusCode, response)
}

func RemoveExtraFromQuery(query string) string {
	cleanedString := strings.ReplaceAll(query, "\t", " ")
	cleanedString = strings.ReplaceAll(cleanedString, "\n", "")
	return cleanedString
}
