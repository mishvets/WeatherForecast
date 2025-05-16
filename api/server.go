package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/mishvets/WeatherForecast/db/sqlc"
	"github.com/mishvets/WeatherForecast/worker"
)

type Server struct {
	store           db.Store
	router          *gin.Engine
	taskDistributor worker.TaskDistributor
}

// New server creates a new HTTP server and setup routing
func NewServer(store db.Store, taskDistributor worker.TaskDistributor) *Server {
	server := &Server{
		store: store,
		taskDistributor: taskDistributor,
	}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("frequency", validFrequency)
	}

	// router.GET("/weather", server.getWeather)

	router.POST("/subscribe", server.subscribe)
	// router.GET("/confirm/:token", server.confirm) // TODO: add city at this step
	// router.GET("/unsubscribe/:token", server.unsubscribe)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
