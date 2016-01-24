package services

import (
	"bytes"
	"encoding/json"
	"strings"
	"time"

	"github.com/gophergala2016/chroniton/models"
	"github.com/gophergala2016/chroniton/utils"

	"github.com/google/go-github/github"
)

type GithubIntegration struct {
	client *github.Client
}

func NewGithubIntegration() GithubIntegration {
	return GithubIntegration{
		client: github.NewClient(nil),
	}
}

func (g GithubIntegration) HandleWebhookEvent(payload []byte) {
	var webhookPayload github.WebHookPayload
	payload = clearJSONRepoOrgField(payload)
	err := json.Unmarshal(payload, &webhookPayload)
	if err != nil {
		utils.Log.Error("Failed to unmarshal github WebHook payload. %v", err)
		return
	}

	repo := *webhookPayload.Repo.FullName
	var project models.Project
	err = utils.ORM.First(&project, "git_repo_name = ?", repo).Error
	if err != nil {
		utils.Log.Error("Failed to get assocaited chroniton project for push event: %v", err)
		return
	}

	if project.GitIntegrationName != "github" {
		utils.Log.Error("Project %v git integration (%v) doesn't match current one (github", project.Name, project.GitIntegrationName)
		return
	}

	var user models.User
	err = utils.ORM.First(&user, project.UserId).Error
	if err != nil {
		utils.Log.Error("Couldn't get associated user for project %v. %v", project.Name, err)
		return
	}

	g.computeTimeWorked(webhookPayload, project, user)
}

func (g GithubIntegration) computeTimeWorked(webhookPayload github.WebHookPayload,
	project models.Project,
	user models.User) {

	fullRepo := strings.Split(project.GitRepoName, "/")

	var branch string
	if strings.Contains(*webhookPayload.Ref, "refs/heads/") {
		branch = strings.Replace(*webhookPayload.Ref, "refs/heads/", "", -1)
	}

	for _, commit := range webhookPayload.Commits {
		duration := g.computeCommitTime(commit, fullRepo[0], fullRepo[1], branch, user)
		utils.Log.Debug("Duration for commit %v is %v", (*commit.ID)[0:8], duration)
	}
}

func (g GithubIntegration) computeCommitTime(commit github.WebHookCommit, owner, repo, branch string, user models.User) int64 {
	var start time.Time
	end := *commit.Timestamp
	commitDetails, _, err := g.client.Git.GetCommit(owner, repo, *commit.ID)
	if err != nil {
		utils.Log.Error("Failed to get commit data from github for project %v, %v", repo, err)
		return -1
	}

	if len(commitDetails.Parents) > 0 {
		parentCommitID := *commitDetails.Parents[0].SHA
		parentCommit, _, err := g.client.Git.GetCommit(owner, repo, parentCommitID)
		if err != nil {
			utils.Log.Error("Failed to get parent commit data from github for project %v, %v", repo, err)
			return -1
		}
		start = *parentCommit.Committer.Date
	} else {
		start = time.Unix(0, 0)
	}

	var heartbeats []models.Heartbeat
	err = utils.ORM.Where("project = ? AND branch = ? AND ? <= time AND time <= ?",
		repo, branch, start.Format(time.RFC3339), end.Format(time.RFC3339)).Order("time ASC").Find(&heartbeats).Error
	if err != nil {
		utils.Log.Error("Failed to get associated heartbeats for commit. %v", err)
		return -1
	}

	utils.Log.Debug("Length of heartbeats for commit: %v", len(heartbeats))
	return models.CalcTotalHeartbeatsDuration(user, heartbeats)
}

func clearJSONRepoOrgField(body []byte) []byte {
	// workaround for https://github.com/google/go-github/issues/131
	var o map[string]interface{}
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.UseNumber()
	dec.Decode(&o)
	if o != nil {
		repo := o["repository"]
		if repo != nil {
			if repo, ok := repo.(map[string]interface{}); ok {
				delete(repo, "organization")
			}
		}
	}
	b, _ := json.MarshalIndent(o, "", "  ")
	return b
}
