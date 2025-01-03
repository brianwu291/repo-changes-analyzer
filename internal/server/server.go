package server

import (
	"net/http"

	"github.com/brianwu291/repo-changes-analyzer/config"
	analysishandler "github.com/brianwu291/repo-changes-analyzer/internal/handlers/analysishandler"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config  *config.Config
	router  *gin.Engine
	handler *analysishandler.AnalysisHandler
}

func NewServer(config *config.Config, handler *analysishandler.AnalysisHandler) *Server {
	server := &Server{
		config:  config,
		router:  gin.Default(),
		handler: handler,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	s.router.Use(corsMiddleware())

	api := s.router.Group("/api")
	{
		api.POST("/analyze", s.handler.HandleAnalysis)
	}

	api.GET("/ping",
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
}

func (s *Server) Start() error {
	return s.router.Run(":" + s.config.ServerPort)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
