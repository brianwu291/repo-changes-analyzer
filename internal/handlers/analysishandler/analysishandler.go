package analysishandlers

import (
	"net/http"
	"time"

	model "github.com/brianwu291/repo-changes-analyzer/internal/models"
	analyzerservice "github.com/brianwu291/repo-changes-analyzer/internal/services/analyzerservice"

	"github.com/gin-gonic/gin"
)

type AnalysisHandler struct {
	analyzerService analyzerservice.AnalyzerService
}

func NewAnalysisHandler(service analyzerservice.AnalyzerService) *AnalysisHandler {
	return &AnalysisHandler{
		analyzerService: service,
	}
}

func (h *AnalysisHandler) HandleAnalysis(c *gin.Context) {
	var req model.RepoAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.AnalysisResponse{
			Error: "Invalid request format: " + err.Error(),
		})
		return
	}

	startDate, endDate, err := h.validateDates(req.StartDate, req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.AnalysisResponse{
			Error: err.Error(),
		})
		return
	}

	result, err := h.analyzerService.AnalyzeRepository(c.Request.Context(), analyzerservice.AnalysisParams{
		Owner:     req.Owner,
		Repo:      req.Repo,
		StartDate: startDate,
		EndDate:   endDate,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.AnalysisResponse{
			Error: "Analysis failed: " + err.Error(),
		})
		return
	}

	response := model.AnalysisResponse{
		Repository:  req.Owner + "/" + req.Repo,
		TimeRange:   req.StartDate + " to " + req.EndDate,
		UserChanges: result,
	}

	c.JSON(http.StatusOK, response)
}

func (h *AnalysisHandler) validateDates(startStr, endStr string) (time.Time, time.Time, error) {
	startDate, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endDate, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return startDate, endDate, nil
}
