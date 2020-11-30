package api

import (
	"HVB_Stock_API/api/customError"
	"HVB_Stock_API/api/dto"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type stockService interface {
	CalculateEarning(key string, responseEntity dto.OutputDTO) (dto.OutputDTO, error)
}

type StockController struct {
	stockService stockService
}

func ProvideStockController(service stockService) *StockController {
	return &StockController{stockService: service}
}

func (stockController *StockController) CalculateEarning(req *gin.Context) {

	//fetch query strings
	share := req.Query("share")
	number := req.Query("number")
	singlePrice := req.Query("singlePrice")

	//throw bad request 400 if empty
	if share == "" || number == "" || singlePrice == "" {
		req.JSON(http.StatusBadRequest, "bad request")
	} else {

		//declare entities
		var responseEntity dto.OutputDTO
		var stockError *customError.ErrorStock

		//set query string share to entity.share
		responseEntity.Share = share

		//some germans might use "," as decimal point
		//catch possible error and parse string into float
		responseEntity.Numbers, _ = strconv.ParseFloat(strings.Replace(number, ",", ".", 1), 64)
		responseEntity.Price, _ = strconv.ParseFloat(strings.Replace(singlePrice, ",", ".", 1), 64)

		//fetch tokenkey
		key := req.GetHeader("x-rapidapi-key")

		//encapsulate controll-layer from service-layer
		returnEntity, err := stockController.stockService.CalculateEarning(key, responseEntity)

		//reponse to user based on err.code
		if errors.As(err, &stockError) {
			if stockError.Code == 200 {
				req.JSON(http.StatusOK, returnEntity)
			} else if stockError.Code == 204 {
				req.JSON(http.StatusNoContent, stockError.Text)
			} else if stockError.Code == 401 {
				req.JSON(http.StatusUnauthorized, stockError.Text)
			} else if stockError.Code == 403 {
				req.JSON(http.StatusForbidden, stockError.Text)
			} else if stockError.Code == 500 {
				req.JSON(http.StatusInternalServerError, stockError.Text)
			}
		}
	}
}

//Root Router to relative Path
func CreateRouter(router *gin.Engine, stockService stockService) {
	stockController := ProvideStockController(stockService)
	router.GET("HVB_Stock_API/calcGain", stockController.CalculateEarning)
}
