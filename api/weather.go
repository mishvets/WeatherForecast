package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mishvets/WeatherForecast/internal/errs"
)

type getWeatherRequest struct {
	City string `form:"city" binding:"required,min=1"`
}

func (server *Server) getWeather(ctx *gin.Context) {
	var req getWeatherRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	weatherRes, err := server.weatherService.GetWeatherForCity(ctx, req.City)
	if err != nil {
		if errors.Is(err, errs.CityNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, weatherRes)
}
