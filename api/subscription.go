package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	db "github.com/mishvets/WeatherForecast/db/sqlc"
	"github.com/mishvets/WeatherForecast/worker"
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
	}

	arg := db.SubscribeTxParams{
		CreateSubscriptionParams: db.CreateSubscriptionParams(req),
		AfterCreate: func(subscription db.Subscription) error {
			taskPayload := &worker.PayloadSendVerifyEmail{
				Email: subscription.Email,
			}
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(5 * time.Second),
			}
			return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
		},
	}

	txResult, err := server.store.SubscribeTx(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation", "":
				ctx.JSON(http.StatusConflict, errorResponse(err))
				return
			}
		}
		// ctx.JSON(http.StatusInternalServerError, errorResponse(err)) // TODO: the provided swagger file doesn't contain this type of error
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, createSubscribeResp(txResult.Subscription))
}

func createSubscribeResp(s db.Subscription) subscribeResponse {
	return subscribeResponse{
		ID:        s.ID,
		Email:     s.Email,
		City:      s.City,
		Frequency: s.Frequency,
		CreatedAt: s.CreatedAt,
	}
}
