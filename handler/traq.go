package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/traPtitech/anke-to/openapi"
	traq "github.com/traPtitech/go-traq"
)

// (GET /traq/users)
func (h Handler) GetTraqUsers(ctx echo.Context) error {
	users, err := h.TraqClient.GetUsers(ctx.Request().Context())
	if err != nil {
		ctx.Logger().Errorf("failed to get traq users: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get traq users: %w", err))
	}

	traqUsers := make(openapi.TraqUsers, 0, len(users))
	for _, user := range users {
		userUUID, err := parseOpenAPIUUID(user.Id)
		if err != nil {
			ctx.Logger().Errorf("invalid traq user uuid: %s", user.Id)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("invalid traq user uuid: %w", err))
		}

		traqUsers = append(traqUsers, openapi.TraqUser{
			Id:   userUUID,
			Name: user.Name,
		})
	}

	return ctx.JSON(http.StatusOK, traqUsers)
}

// (GET /traq/users/me)
func (h Handler) GetTraqUsersMe(ctx echo.Context) error {
	userID, err := h.Middleware.GetUserID(ctx)
	if err != nil {
		ctx.Logger().Errorf("failed to get userID: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
	}

	users, err := h.TraqClient.GetUsersByName(ctx.Request().Context(), userID)
	if err != nil {
		ctx.Logger().Errorf("failed to get traq user by name: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get traq user by name: %w", err))
	}
	if len(users) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "traq user not found")
	}

	// 最初のユーザーを使用（GetUsersByNameは名前でフィルタリング済み）
	user := users[0]
	userUUID, err := uuid.Parse(user.Id)
	if err != nil {
		ctx.Logger().Errorf("invalid user uuid: %s", user.Id)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("invalid user uuid: %w", err))
	}

	return ctx.JSON(http.StatusOK, openapi.TraqUser{
		Id:   userUUID,
		Name: user.Name,
	})
}

// (GET /traq/groups)
func (h Handler) GetTraqGroups(ctx echo.Context) error {
	groups, err := h.TraqClient.GetGroups(ctx.Request().Context())
	if err != nil {
		ctx.Logger().Errorf("failed to get traq groups: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get traq groups: %w", err))
	}

	traqGroups := make(openapi.TraqGroups, 0, len(groups))
	for _, group := range groups {
		groupUUID, err := parseOpenAPIUUID(group.Id)
		if err != nil {
			ctx.Logger().Errorf("invalid traq group uuid: %s", group.Id)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("invalid traq group uuid: %w", err))
		}

		traqGroups = append(traqGroups, openapi.TraqGroup{
			Id:   groupUUID,
			Name: group.Name,
		})
	}

	return ctx.JSON(http.StatusOK, traqGroups)
}

// (GET /traq/stamps)
func (h Handler) GetTraqStamps(ctx echo.Context) error {
	stamps, err := h.TraqClient.GetStamps(ctx.Request().Context())
	if err != nil {
		ctx.Logger().Errorf("failed to get traq stamps: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get traq stamps: %w", err))
	}

	traqStamps := make(openapi.TraqStamps, 0, len(stamps))
	for _, stamp := range stamps {
		stampUUID, err := parseOpenAPIUUID(stamp.Id)
		if err != nil {
			ctx.Logger().Errorf("invalid traq stamp uuid: %s", stamp.Id)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("invalid traq stamp uuid: %w", err))
		}

		traqStamps = append(traqStamps, openapi.TraqStamp{
			Id:   stampUUID,
			Name: stamp.Name,
		})
	}

	return ctx.JSON(http.StatusOK, traqStamps)
}

// (GET /traq/channels)
func (h Handler) GetTraqChannels(ctx echo.Context) error {
	channels, err := h.TraqClient.GetChannels(ctx.Request().Context())
	if err != nil {
		ctx.Logger().Errorf("failed to get traq channels: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get traq channels: %w", err))
	}
	if channels == nil {
		return ctx.JSON(http.StatusOK, openapi.TraqChannels{})
	}

	traqChannels, err := mapTraqChannels(channels.Public)
	if err != nil {
		ctx.Logger().Errorf("failed to map traq channels: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to map traq channels: %w", err))
	}

	return ctx.JSON(http.StatusOK, traqChannels)
}

func parseOpenAPIUUID(raw string) (openapi_types.UUID, error) {
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return openapi_types.UUID{}, err
	}
	return openapi_types.UUID(parsed), nil
}

func mapTraqChannels(channels []traq.Channel) (openapi.TraqChannels, error) {
	pathByID := make(map[string]string, len(channels))
	channelByID := make(map[string]traq.Channel, len(channels))
	for _, channel := range channels {
		channelByID[channel.Id] = channel
	}

	var buildPath func(channel traq.Channel) string
	buildPath = func(channel traq.Channel) string {
		if cachedPath, ok := pathByID[channel.Id]; ok {
			return cachedPath
		}

		current := "/" + channel.Name
		parentID := channel.ParentId.Get()
		if parentID != nil {
			if parent, ok := channelByID[*parentID]; ok {
				current = strings.TrimRight(buildPath(parent), "/") + current
			}
		}

		pathByID[channel.Id] = current
		return current
	}

	traqChannels := make(openapi.TraqChannels, 0, len(channels))
	for _, channel := range channels {
		channelUUID, err := parseOpenAPIUUID(channel.Id)
		if err != nil {
			return nil, err
		}
		path := buildPath(channel)
		traqChannels = append(traqChannels, openapi.TraqChannel{
			Id:   channelUUID,
			Name: channel.Name,
			Path: &path,
		})
	}
	return traqChannels, nil
}
