package di

import (
	"context"

	"github.com/brianwu291/repo-changes-analyzer/config"
	analysishandler "github.com/brianwu291/repo-changes-analyzer/internal/handlers/analysishandler"
	githubrepository "github.com/brianwu291/repo-changes-analyzer/internal/repos/githubrepository"
	"github.com/brianwu291/repo-changes-analyzer/internal/server"
	analyzerservice "github.com/brianwu291/repo-changes-analyzer/internal/services/analyzerservice"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

type Container struct {
	Config       *config.Config
	HTTPServer   *server.Server
	GithubClient *github.Client
}

func NewDI(config *config.Config) *Container {
	tokenSecure := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GithubToken},
	)
	tokenClient := oauth2.NewClient(context.Background(), tokenSecure)
	githubClient := github.NewClient(tokenClient)

	repoAnalyzer := githubrepository.NewGithubRepository(githubClient)

	analyzerService := analyzerservice.NewAnalyzerService(repoAnalyzer)

	analysisHandler := analysishandler.NewAnalysisHandler(analyzerService)

	httpServer := server.NewServer(config, analysisHandler)

	return &Container{
		Config:       config,
		HTTPServer:   httpServer,
		GithubClient: githubClient,
	}
}
