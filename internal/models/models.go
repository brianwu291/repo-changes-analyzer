package model

type RepoAnalysisRequest struct {
	Owner     string `json:"owner" binding:"required"`
	Repo      string `json:"repo" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

type UserChanges struct {
	Username    string `json:"username"`
	Additions   int    `json:"additions"`
	Deletions   int    `json:"deletions"`
	Total       int    `json:"total"`
	CommitCount int    `json:"commit_count"`
	AvatarURL   string `json:"avatar_url,omitempty"`
}

type AnalysisResponse struct {
	Repository  string                 `json:"repository"`
	TimeRange   string                 `json:"time_range"`
	UserChanges map[string]UserChanges `json:"user_changes"`
	Error       string                 `json:"error,omitempty"`
}

type CommitAnalysis struct {
	Author      string
	Additions   int
	Deletions   int
	CommitCount int
	AvatarURL   string
}
