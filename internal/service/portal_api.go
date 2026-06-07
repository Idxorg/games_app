package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

// PortalAPI клиент для работы с корп порталом
type PortalAPI struct {
	client  *resty.Client
	baseURL string
	apiKey  string
}

// PortalUser данные пользователя из корп портала
type PortalUser struct {
	SID        string `json:"sid"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Gender     string `json:"gender"`
	Department string `json:"department"`
	Position   string `json:"position"`
	AvatarURL  string `json:"avatar_url"`
}

// PortalUserGroups группы пользователя
type PortalUserGroups struct {
	SID    string   `json:"sid"`
	Groups []string `json:"groups"`
}

// NewPortalAPI создает новый клиент Portal API
func NewPortalAPI(baseURL, apiKey string) *PortalAPI {
	return &PortalAPI{
		client:  resty.New().SetBaseURL(baseURL).SetTimeout(10 * time.Second),
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

// GetUser получает данные пользователя из корп портала
func (p *PortalAPI) GetUser(ctx context.Context, sid string) (*PortalUser, error) {
	resp, err := p.client.R().
		SetAuthToken(p.apiKey).
		SetResult(&PortalUser{}).
		Get("/api/users/" + sid)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("portal API error: %s", resp.Status())
	}

	user := resp.Result().(*PortalUser)
	return user, nil
}

// GetUserGroups получает группы пользователя
func (p *PortalAPI) GetUserGroups(ctx context.Context, sid string) ([]string, error) {
	resp, err := p.client.R().
		SetAuthToken(p.apiKey).
		SetResult(&PortalUserGroups{}).
		Get("/api/users/" + sid + "/groups")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("portal API error: %s", resp.Status())
	}

	groups := resp.Result().(*PortalUserGroups)
	return groups.Groups, nil
}

// HasAccess проверяет, есть ли у пользователя доступ к группе
func (p *PortalAPI) HasAccess(ctx context.Context, sid string, requiredGroup string) (bool, error) {
	groups, err := p.GetUserGroups(ctx, sid)
	if err != nil {
		return false, err
	}

	for _, g := range groups {
		if g == requiredGroup {
			return true, nil
		}
	}

	return false, nil
}
