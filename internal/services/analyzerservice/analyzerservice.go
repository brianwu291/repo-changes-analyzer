package analyzerservice

import (
	"context"
	"time"

	model "github.com/brianwu291/repo-changes-analyzer/internal/models"
	repos "github.com/brianwu291/repo-changes-analyzer/internal/repos/githubrepository"
)

type AnalysisParams struct {
	Owner     string
	Repo      string
	StartDate time.Time
	EndDate   time.Time
}

type AnalyzerService interface {
	AnalyzeRepository(ctx context.Context, params AnalysisParams) (map[string]model.UserChanges, error)
}

type analyzerService struct {
	repo repos.GithubRepository
}

func NewAnalyzerService(repo repos.GithubRepository) AnalyzerService {
	return &analyzerService{
		repo: repo,
	}
}

func (s *analyzerService) AnalyzeRepository(ctx context.Context, params AnalysisParams) (map[string]model.UserChanges, error) {
	commits, err := s.repo.GetCommits(ctx, params.Owner, params.Repo, params.StartDate, params.EndDate)
	if err != nil {
		return nil, err
	}

	results := s.repo.ProcessCommitsConcurrently(ctx, params.Owner, params.Repo, commits)

	userChanges := make(map[string]model.UserChanges)
	for user, changes := range results {
		userChanges[user] = model.UserChanges{
			Additions: changes.Additions,
			Deletions: changes.Deletions,
			Total:     changes.Additions + changes.Deletions,
		}
	}

	return userChanges, nil
}
