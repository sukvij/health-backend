package response

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ReturnErrorWithCode(ctx *gin.Context, err error, response *FinalResponse) {
	fmt.Println(errors.New("app_id are required"), err)
	if errors.New("app_id are required") == err {
		fmt.Println("case done....")
	}
	response.Error = &AppError{}
	switch err.Error() {
	case gorm.ErrRecordNotFound.Error():
		response.StatusCode = 204
		response.Error.Code = "204"
		response.Error.Message = "record not found check variables that you pass"
	case gorm.ErrDuplicatedKey.Error():
		response.StatusCode = 409
		response.Error.Code = "409"
		response.Error.Message = "duplicate key --> primary key constraint"
	case gorm.ErrForeignKeyViolated.Error():
		response.Error.Code = "23505"
		response.Error.Message = "foreign key violation."
	case gorm.ErrInvalidDB.Error():
		response.Error.Code = "252525"
		response.Error.Message = "database connection problem"
	case "app_id are required":
		fmt.Println("aa rha h kya")
		response.StatusCode = 400
		response.Error.Code = "400"
		response.Error.Message = "app_id missing"
	case "os_id are required":
		response.StatusCode = 400
		response.Error.Code = "400"
		response.Error.Message = "os_id missing"
	case "country_id are required":
		response.StatusCode = 400
		response.Error.Code = "400"
		response.Error.Message = "countryid missing"
	default:
		response.Error.Code = "898989"
		response.StatusCode = 500
		response.Error.Message = "pta nahi kya problem hain. 😔"
	}
	response.Error.Details = string(err.Error())
}
