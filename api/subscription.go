package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/mishvets/WeatherForecast/db/sqlc"
)

type subscribeRequest struct {
	Email     string           `json:"email" binding:"required,email"`
	City      string           `json:"city" binding:"required"`
	Frequency db.FrequencyEnum `json:"frequency" binding:"required,frequency"`
}

type subscribeResponse struct {
	ID        int64            `json:"id"`
	Email     string           `json:"email"`
	City      string           `json:"city"`
	Frequency db.FrequencyEnum `json:"frequency"`
	CreatedAt time.Time        `json:"created_at"`
}

func (server *Server) subscribe(ctx *gin.Context) {
	var req subscribeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	} else {
		fmt.Println(req)
	}

	subscription, err := server.store.CreateSubscription(ctx, db.CreateSubscriptionParams(req))
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation", "":
				ctx.JSON(http.StatusConflict, errorResponse(err))
				return
			}
		}
		// ctx.JSON(http.StatusInternalServerError, errorResponse(err)) // TODO: the provided swagger file doesn't contain this type of error
		return
	}

	// TODO: send confirmation email
	resp := subscribeResponse{
		ID:        subscription.ID,
		Email:     subscription.Email,
		City:      subscription.City,
		Frequency: subscription.Frequency,
		CreatedAt: subscription.CreatedAt,
	}
	ctx.JSON(http.StatusOK, resp)
}
