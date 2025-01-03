package model

type RepoAnalysisRequest struct {
	Owner     string `json:"owner" binding:"required"`
	Repo      string `json:"repo" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

type UserChanges struct {
	Additions int `json:"additions"`
	Deletions int `json:"deletions"`
	Total     int `json:"total"`
}

type AnalysisResponse struct {
	Repository  string                 `json:"repository"`
	TimeRange   string                 `json:"time_range"`
	UserChanges map[string]UserChanges `json:"user_changes"`
	Error       string                 `json:"error,omitempty"`
}

type CommitAnalysis struct {
	Author    string
	Additions int
	Deletions int
}
