package main

import (
	"strings"

	"github.com/labstack/echo/v4"
)

type RouteConfig struct {
	path        string
	method      string
	middlewares []echo.MiddlewareFunc
	isGroup     bool
}

type MiddlewareSwitcher struct {
	routeConfigs []RouteConfig
}

func NewMiddlewareSwitcher() *MiddlewareSwitcher {
	return &MiddlewareSwitcher{
		routeConfigs: []RouteConfig{},
	}
}

func (m *MiddlewareSwitcher) AddGroupConfig(grouppath string, middlewares ...echo.MiddlewareFunc) {
	m.routeConfigs = append(m.routeConfigs, RouteConfig{
		path:        grouppath,
		middlewares: middlewares,
		isGroup:     true,
	})
}

func (m *MiddlewareSwitcher) AddRouteConfig(path string, method string, middlewares ...echo.MiddlewareFunc) {
	m.routeConfigs = append(m.routeConfigs, RouteConfig{
		path:        path,
		method:      method,
		middlewares: middlewares,
		isGroup:     false,
	})
}

func (m *MiddlewareSwitcher) IsWithinGroup(groupPath string, path string) bool {
	if !strings.HasPrefix(path, groupPath) {
		return false
	}
	return len(groupPath) == len(path) || path[len(groupPath)] == '/'
}

func (m *MiddlewareSwitcher) FindMiddlewares(path string, method string) []echo.MiddlewareFunc {
	var matchedMiddlewares []echo.MiddlewareFunc

	for _, config := range m.routeConfigs {
		if config.isGroup && m.IsWithinGroup(config.path, path) {
			matchedMiddlewares = append(matchedMiddlewares, config.middlewares...)
		}
		if !config.isGroup && config.path == path && config.method == method {
			matchedMiddlewares = append(matchedMiddlewares, config.middlewares...)
		}
	}

	return matchedMiddlewares
}

func (m *MiddlewareSwitcher) ApplyMiddlewares(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Path()
		method := c.Request().Method

		middlewares := m.FindMiddlewares(path, method)

		for _, mw := range middlewares {
			if err := mw(next)(c); err != nil {
				return err
			}
		}

		return next(c)
	}
}
