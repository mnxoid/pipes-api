package main

type (
	Integration struct {
		ID         string  `json:"id"`
		Name       string  `json:"name"`
		Link       string  `json:"link"`
		Image      string  `json:"image"`
		Pipes      []*Pipe `json:"pipes"`
		AuthURL    string  `json:"auth_url,omitempty"`
		AuthType   string  `json:"auth_type,omitempty"`
		Authorized bool    `json:"authorized"`
	}
)

var pipeNames = map[string]string{
	"users":       "Users",
	"tasks":       "Tasks",
	"todos":       "Todos",
	"clients":     "Clients",
	"projects":    "Projects",
	"todolists":   "Todo lists",
	"timeentries": "Time entries export",
}

func workspaceIntegrations(workspaceID int) ([]Integration, error) {
	authorizations, err := loadAuthorizations(workspaceID)
	if err != nil {
		return nil, err
	}
	workspacePipes, err := loadPipes(workspaceID)
	if err != nil {
		return nil, err
	}
	pipeStatuses, err := loadPipeStatuses(workspaceID)
	if err != nil {
		return nil, err
	}

	var integrations []Integration
	for _, integration := range availableIntegrations {
		integration.AuthURL = oAuth2URL(integration.ID)
		integration.Authorized = authorizations[integration.ID]
		var pipes []*Pipe
		for _, integrationPipe := range integration.Pipes {
			key := pipesKey(integration.ID, integrationPipe.ID)
			pipe := workspacePipes[key]
			if pipe == nil {
				pipe = integrationPipe
			}
			pipe.Premium = integrationPipe.Premium
			pipe.PipeStatus = pipeStatuses[key]
			pipe.serviceID = integration.ID
			pipe.workspaceID = workspaceID
			pipe.key = key
			pipes = append(pipes, pipe)
		}
		integration.Pipes = pipes
		integrations = append(integrations, *integration)
	}
	return integrations, nil
}
