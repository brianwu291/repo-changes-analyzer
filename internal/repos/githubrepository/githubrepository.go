package githubrepository

import (
	"context"
	"log"
	"sync"
	"time"

	model "github.com/brianwu291/repo-changes-analyzer/internal/models"

	github "github.com/google/go-github/v45/github"
)

type GithubRepository interface {
	GetCommits(ctx context.Context, owner, repo string, startDate, endDate time.Time) ([]*github.RepositoryCommit, error)
	ProcessCommitsConcurrently(ctx context.Context, owner, repo string, commits []*github.RepositoryCommit) map[string]model.UserChanges
}

type githubRepository struct {
	client *github.Client
}

func NewGithubRepository(client *github.Client) GithubRepository {
	return &githubRepository{
		client: client,
	}
}

func (r *githubRepository) GetCommits(ctx context.Context, owner, repo string, startDate, endDate time.Time) ([]*github.RepositoryCommit, error) {
	opts := &github.CommitsListOptions{
		Since: startDate,
		Until: endDate,
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var allCommits []*github.RepositoryCommit
	for {
		commits, resp, err := r.client.Repositories.ListCommits(ctx, owner, repo, opts)
		if err != nil {
			return nil, err
		}
		allCommits = append(allCommits, commits...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return allCommits, nil
}

func (r *githubRepository) ProcessCommitsConcurrently(ctx context.Context, owner, repo string, commits []*github.RepositoryCommit) map[string]model.UserChanges {
	var wg sync.WaitGroup
	workerSize := 5
	buffer := min(len(commits), workerSize * 5)
	userChangesChan := make(chan model.CommitAnalysis, buffer)

	workerPool := make(chan struct{}, workerSize) // limit concurrent api requests
	for _, commit := range commits {
		if commit.Author == nil || commit.Author.Login == nil {
			continue
		}

		wg.Add(1)
		go func(commit *github.RepositoryCommit) {
			defer wg.Done()
			workerPool <- struct{}{}        // acquire a worker
			defer func() { <-workerPool }() // release a worker

			// get commit details
			commitDetail, _, err := r.client.Repositories.GetCommit(ctx, owner, repo, commit.GetSHA(), nil)
			if err != nil {
				log.Printf("error getting commit detail for %s: %v", commit.GetSHA(), err)
				return
			}

			userChangesChan <- model.CommitAnalysis{
				Author:    *commit.Author.Login,
				Additions: commitDetail.GetStats().GetAdditions(),
				Deletions: commitDetail.GetStats().GetDeletions(),
			}
		}(commit)
	}

	// close channel when all goroutines are done
	go func() {
		wg.Wait()
		close(userChangesChan)
	}()

	// aggregate results
	userChanges := make(map[string]model.UserChanges)
	for change := range userChangesChan {
		current := userChanges[change.Author]
		current.Additions += change.Additions
		current.Deletions += change.Deletions
		userChanges[change.Author] = current
	}

	return userChanges
}
