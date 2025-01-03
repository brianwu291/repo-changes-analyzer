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
	GetContributors(ctx context.Context, owner, repo string) ([]*github.Contributor, error)
	GetCommits(ctx context.Context, owner, repo string, startDate, endDate time.Time) ([]*github.RepositoryCommit, error)
	ProcessCommitsConcurrently(ctx context.Context, owner, repo string, commits []*github.RepositoryCommit) (map[string]model.UserChanges, []error)
}

type githubRepository struct {
	client *github.Client
}

func NewGithubRepository(client *github.Client) GithubRepository {
	return &githubRepository{
		client: client,
	}
}

func (r *githubRepository) GetContributors(ctx context.Context, owner string, repo string) ([]*github.Contributor, error) {
	opts := &github.ListContributorsOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var allContributors []*github.Contributor
	for {
		contributors, resp, err := r.client.Repositories.ListContributors(ctx, owner, repo, opts)
		if err != nil {
			return nil, err
		}
		allContributors = append(allContributors, contributors...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return allContributors, nil
}

func (r *githubRepository) GetCommits(ctx context.Context, owner string, repo string, startDate time.Time, endDate time.Time) ([]*github.RepositoryCommit, error) {
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

func (r *githubRepository) ProcessCommitsConcurrently(ctx context.Context, owner string, repo string, commits []*github.RepositoryCommit) (map[string]model.UserChanges, []error) {
	commitsChan := make(chan *github.RepositoryCommit, len(commits))
	userCommitAnalysisChan := make(chan *model.CommitAnalysis, len(commits))
	errorsChan := make(chan *error, len(commits))
	var errors []error
	userChangesMap := make(map[string]model.UserChanges)

	workerSize := 5
	var wg sync.WaitGroup
	for i := 0; i < workerSize; i += 1 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for commit := range commitsChan {
				commitDetail, _, err := r.client.Repositories.GetCommit(ctx, owner, repo, commit.GetSHA(), nil)
				if err != nil {
					log.Printf("error getting commit detail for %s: %v\n", commit.GetSHA(), err)
					errorsChan <- &err
					continue
				}
				stats := commitDetail.GetStats()
				if stats == nil {
					continue
				}
				userCommitAnalysisChan <- &model.CommitAnalysis{
					Author:    *commit.Author.Login,
					Additions: stats.GetAdditions(),
					Deletions: stats.GetDeletions(),
				}
			}
		}()
	}
	go func() {
		defer close(commitsChan)
		for _, commit := range commits {
			if commit.Author == nil || commit.Author.Login == nil {
				continue
			}
			commitsChan <- commit
		}
	}()

	wg.Wait()
	close(userCommitAnalysisChan)
	close(errorsChan)

	commitCountsMap := r.getCommitCountsMap(commits)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range errorsChan {
			if err != nil {
				errors = append(errors, *err)
			}
		}
		for commitAnalysis := range userCommitAnalysisChan {
			if commitAnalysis != nil {
				current := userChangesMap[commitAnalysis.Author]
				current.Additions += commitAnalysis.Additions
				current.Deletions += commitAnalysis.Deletions
				current.Total = current.Additions + current.Deletions
				current.CommitCount = commitCountsMap[commitAnalysis.Author]
				userChangesMap[commitAnalysis.Author] = current
			}
		}
	}()
	wg.Wait()

	return userChangesMap, errors
}

func (r *githubRepository) getCommitCountsMap(commits []*github.RepositoryCommit) map[string]int {
	commitCountsMap := make(map[string]int)
	for _, commit := range commits {
		if commit.Author != nil && commit.Author.Login != nil {
			commitCountsMap[*commit.Author.Login] += 1
		}
	}
	return commitCountsMap
}
