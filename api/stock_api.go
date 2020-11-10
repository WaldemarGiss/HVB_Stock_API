package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type stockService interface {
	Calculate(tenant string) float64
}

type StockAPI struct {
	stockService stockService
}

func ProvideStockAPI(service stockService) *StockAPI {
	return &StockAPI{stockService: service}
}

func (stockAPI *StockAPI) Calculate(req *gin.Context) {
	//TODO: DIE TENANTS ABFANGEN
	tenant := req.Param("vicci_tenant")
	ret := stockAPI.stockService.Calculate(tenant)
	if ret != 0 {
		req.JSON(http.StatusOK, ret)
	} else {
		req.JSON(http.StatusBadRequest, ret)
	}
}

func CreateRouter(router *gin.Engine, stockService stockService) {
	stockController := ProvideStockAPI(stockService)
	router.GET("/tenant/:vicci_tenant/carlines", stockController.Calculate)
}
