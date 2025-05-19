package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
			case "unique_violation":
				ctx.JSON(http.StatusConflict, errorResponse(err))
				return
			}
		}
		// ctx.JSON(http.StatusInternalServerError, errorResponse(err)) // TODO: the provided swagger file doesn't contain this type of error
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, createSubscribeResp(txResult))
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

type uriTokenRequest struct {
	Token string `json:"token" uri:"token" binding:"required,min=36,max=36"`
}

func (server *Server) confirm(ctx *gin.Context) {
	var req uriTokenRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	uuid, err := uuid.Parse(req.Token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ConfirmSubscriptionTxParams{
		ConfirmSubscriptionParams: db.ConfirmSubscriptionParams{
		Token:     uuid,
		Confirmed: true,
	},
		AfterConfirm: func(subscription db.Subscription) error {
			taskPayload := &worker.PayloadGetWeatherData{
				ID: subscription.ID,
				Token: subscription.Token,
				Email: subscription.Email,
				City: subscription.City,
			}
			opts := []asynq.Option{
				asynq.MaxRetry(10),
			}
			return server.taskDistributor.DistributeTaskGetWeatherData(ctx, taskPayload, opts...)
		},
	}

	_, err = server.store.ConfirmSubscriptionTx(ctx, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, req)
}

func (server *Server) unsubscribe(ctx *gin.Context) {
	var req uriTokenRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	uuid, err := uuid.Parse(req.Token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = server.store.DeleteSubscription(ctx, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, req)
}
