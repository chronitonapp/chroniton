package services

type gitService interface {
	HandleWebhookEvent(payload []byte)
}

var registeredGitServices = map[string]gitService{}
