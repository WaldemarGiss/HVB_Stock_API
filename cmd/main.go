package main

import (
	"HVB_Stock_API/api"
	"HVB_Stock_API/internal/repository"
	"HVB_Stock_API/internal/service"
	"github.com/gin-gonic/gin"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	//if err := godotenv.Load(); err != nil {
	//	log.Print("No .env file found")
	//}
}

func main() {

	router := gin.Default()

	stockRepository := repository.ProvideStockRepository()

	stockService := service.ProvideStockService(stockRepository)

	api.CreateRouter(router, stockService)

	if err := router.Run("localhost:8080"); err != nil {
		panic(err)
	}

}
