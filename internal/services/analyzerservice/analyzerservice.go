package analyzerservice

import (
	"context"
	"fmt"
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
	// get all contributors
	contributors, err := s.repo.GetContributors(ctx, params.Owner, params.Repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get contributors: %w", err)
	}

	// user changes map with contributor info
	userChanges := make(map[string]model.UserChanges)
	for _, contributor := range contributors {
		if contributor.GetLogin() != "" {
			userChanges[contributor.GetLogin()] = model.UserChanges{
				Username:    contributor.GetLogin(),
				AvatarURL:   contributor.GetAvatarURL(),
				Deletions:   0,
				Total:       0,
				CommitCount: 0,
			}
		}
	}

	// get commits in date range
	commits, err := s.repo.GetCommits(ctx, params.Owner, params.Repo, params.StartDate, params.EndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits: %w", err)
	}

	// process commits and update statistics
	commitStats := s.repo.ProcessCommitsConcurrently(ctx, params.Owner, params.Repo, commits)

	// merge commit statistics
	for username, stats := range commitStats {
		if user, exists := userChanges[username]; exists {
			user.Additions = stats.Additions
			user.Deletions = stats.Deletions
			user.Total = stats.Total
			user.CommitCount = stats.CommitCount
			userChanges[username] = user
		} else {
			userChanges[username] = model.UserChanges{
				Username:    username,
				Additions:   stats.Additions,
				Deletions:   stats.Deletions,
				Total:       stats.Total,
				CommitCount: stats.CommitCount,
			}
		}
	}

	return userChanges, nil
}
