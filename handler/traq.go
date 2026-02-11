package handler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/traPtitech/anke-to/openapi"
	traqAPI "github.com/traPtitech/anke-to/traq"
)

// (GET /traq/users)
func (h Handler) GetTraqUsers(ctx echo.Context) error {
	client := traqAPI.NewTraqAPIClient()
	users, err := client.GetUsers(ctx.Request().Context())
	if err != nil {
		ctx.Logger().Errorf("failed to get traq users: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get traq users: %w", err))
	}

	return ctx.JSON(http.StatusOK, users)
}

// (GET /traq/users/me)
func (h Handler) GetTraqUsersMe(ctx echo.Context) error {
	userID, err := h.Middleware.GetUserID(ctx)
	if err != nil {
		ctx.Logger().Errorf("failed to get userID: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
	}

	client := traqAPI.NewTraqAPIClient()
	users, err := client.GetUsersByName(ctx.Request().Context(), userID)
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

	return ctx.JSON(http.StatusOK, openapi.TraqMe{
		Id:   userID,
		Uuid: openapi_types.UUID(userUUID),
	})
}

// (GET /traq/groups)
func (h Handler) GetTraqGroups(ctx echo.Context) error {
	client := traqAPI.NewTraqAPIClient()
	groups, err := client.GetGroups(ctx.Request().Context())
	if err != nil {
		ctx.Logger().Errorf("failed to get traq groups: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get traq groups: %w", err))
	}

	return ctx.JSON(http.StatusOK, groups)
}

// (GET /traq/stamps)
func (h Handler) GetTraqStamps(ctx echo.Context) error {
	client := traqAPI.NewTraqAPIClient()
	stamps, err := client.GetStamps(ctx.Request().Context())
	if err != nil {
		ctx.Logger().Errorf("failed to get traq stamps: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get traq stamps: %w", err))
	}

	return ctx.JSON(http.StatusOK, stamps)
}

// (GET /traq/channels)
func (h Handler) GetTraqChannels(ctx echo.Context) error {
	client := traqAPI.NewTraqAPIClient()
	channels, err := client.GetChannels(ctx.Request().Context())
	if err != nil {
		ctx.Logger().Errorf("failed to get traq channels: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get traq channels: %w", err))
	}

	return ctx.JSON(http.StatusOK, channels)
}
