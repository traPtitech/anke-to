package traq

import (
	"context"

	traq "github.com/traPtitech/go-traq"
)

const TOKEN = "/* your token */"

type APIClient struct {
	client *traq.APIClient
	auth   context.Context
}

func NewTraqAPIClient() *APIClient {
	return &APIClient{
		client: traq.NewAPIClient(traq.NewConfiguration()),
		auth:   context.WithValue(context.Background(), traq.ContextAccessToken, TOKEN),
	}
}

func (t *APIClient) GetGroupMembers(ctx context.Context, groupID string) ([]traq.UserGroupMember, error) {
	v, _, err := t.client.GroupApi.GetUserGroupMembers(ctx, groupID).Execute()
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (t *APIClient) GetUserTraqID(ctx context.Context, userUUID string) (string, error) {
	v, _, err := t.client.UserApi.GetUser(ctx, userUUID).Execute()
	if err != nil {
		return "", err
	}
	return v.Name, nil
}

func (t *APIClient) GetGroupName(ctx context.Context, groupID string) (string, error) {
	v, _, err := t.client.GroupApi.GetUserGroup(ctx, groupID).Execute()
	if err != nil {
		return "", err
	}
	return v.Name, nil
}
