package traq

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	traq "github.com/traPtitech/go-traq"
)

const tokenEnvKey = "TRAQ_BOT_TOKEN"

type APIClient struct {
	client *traq.APIClient
	token  string
}

func NewTraqAPIClient() *APIClient {
	cfg := traq.NewConfiguration()
	cfg.HTTPClient = &http.Client{
		Transport: newETagCacheTransport(http.DefaultTransport),
	}

	return &APIClient{
		client: traq.NewAPIClient(cfg),
		token:  os.Getenv(tokenEnvKey),
	}
}

func (t *APIClient) authContext(ctx context.Context) context.Context {
	if t.token == "" {
		return ctx
	}
	return context.WithValue(ctx, traq.ContextAccessToken, t.token)
}

func (t *APIClient) GetGroupMembers(ctx context.Context, groupID string) ([]traq.UserGroupMember, error) {
	v, _, err := t.client.GroupApi.GetUserGroupMembers(t.authContext(ctx), groupID).Execute()
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (t *APIClient) GetUserTraqID(ctx context.Context, userUUID string) (string, error) {
	v, _, err := t.client.UserApi.GetUser(t.authContext(ctx), userUUID).Execute()
	if err != nil {
		return "", err
	}
	return v.Name, nil
}

func (t *APIClient) GetGroupName(ctx context.Context, groupID string) (string, error) {
	v, _, err := t.client.GroupApi.GetUserGroup(t.authContext(ctx), groupID).Execute()
	if err != nil {
		return "", err
	}
	return v.Name, nil
}

func (t *APIClient) GetUsers(ctx context.Context) ([]traq.User, error) {
	v, _, err := t.client.UserApi.GetUsers(t.authContext(ctx)).Execute()
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (t *APIClient) GetUsersByName(ctx context.Context, name string) ([]traq.User, error) {
	v, _, err := t.client.UserApi.GetUsers(t.authContext(ctx)).Name(name).Execute()
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (t *APIClient) GetGroups(ctx context.Context) ([]traq.UserGroup, error) {
	v, _, err := t.client.GroupApi.GetUserGroups(t.authContext(ctx)).Execute()
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (t *APIClient) GetStamps(ctx context.Context) ([]traq.StampWithThumbnail, error) {
	v, _, err := t.client.StampApi.GetStamps(t.authContext(ctx)).Execute()
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (t *APIClient) GetChannels(ctx context.Context) (*traq.ChannelList, error) {
	v, _, err := t.client.ChannelApi.GetChannels(t.authContext(ctx)).Execute()
	if err != nil {
		return nil, err
	}
	return v, nil
}

type cacheEntry struct {
	etag   string
	body   []byte
	header http.Header
}

type etagCacheTransport struct {
	base  http.RoundTripper
	mu    sync.RWMutex
	cache map[string]cacheEntry
}

func newETagCacheTransport(base http.RoundTripper) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	return &etagCacheTransport{
		base:  base,
		cache: map[string]cacheEntry{},
	}
}

func (t *etagCacheTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Method != http.MethodGet || !isETagTargetPath(req.URL.Path) {
		return t.base.RoundTrip(req)
	}

	cacheKey := req.URL.String()

	etag, hasCache := t.getETag(cacheKey)
	nextReq := req.Clone(req.Context())
	if hasCache && etag != "" {
		nextReq.Header.Set("If-None-Match", etag)
	}

	resp, err := t.base.RoundTrip(nextReq)
	if err != nil {
		// 通信失敗時はキャッシュがあればフォールバックする。
		if cachedResp, ok := t.buildCachedResponse(req, cacheKey); ok {
			return cachedResp, nil
		}
		return nil, err
	}

	if resp.StatusCode == http.StatusNotModified {
		resp.Body.Close()
		// 304のときは保存済みボディを200として返す。
		if cachedResp, ok := t.buildCachedResponse(req, cacheKey); ok {
			return cachedResp, nil
		}
		return resp, nil
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		resp.Body.Close()
		if cachedResp, ok := t.buildCachedResponse(req, cacheKey); ok {
			return cachedResp, nil
		}
		return resp, nil
	}

	if resp.StatusCode != http.StatusOK {
		return resp, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(body))
	resp.ContentLength = int64(len(body))

	if etag := resp.Header.Get("ETag"); etag != "" {
		t.setCache(cacheKey, cacheEntry{
			etag:   etag,
			body:   append([]byte(nil), body...),
			header: resp.Header.Clone(),
		})
	}

	return resp, nil
}

func (t *etagCacheTransport) getETag(key string) (string, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	entry, ok := t.cache[key]
	if !ok {
		return "", false
	}
	return entry.etag, true
}

func (t *etagCacheTransport) setCache(key string, entry cacheEntry) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.cache[key] = entry
}

func (t *etagCacheTransport) buildCachedResponse(req *http.Request, key string) (*http.Response, bool) {
	t.mu.RLock()
	entry, ok := t.cache[key]
	t.mu.RUnlock()
	if !ok {
		return nil, false
	}

	header := entry.header.Clone()
	if header == nil {
		header = make(http.Header)
	}
	header.Set("X-AnkeTo-Cache", "hit")

	return &http.Response{
		StatusCode:    http.StatusOK,
		Status:        "200 OK",
		Header:        header,
		Body:          io.NopCloser(bytes.NewReader(entry.body)),
		ContentLength: int64(len(entry.body)),
		Request:       req,
	}, true
}

func isETagTargetPath(path string) bool {
	trimmed := strings.TrimRight(path, "/")
	switch {
	case strings.HasSuffix(trimmed, "/users"):
		return true
	case strings.HasSuffix(trimmed, "/groups"):
		return true
	case strings.HasSuffix(trimmed, "/stamps"):
		return true
	case strings.HasSuffix(trimmed, "/channels"):
		return true
	default:
		return false
	}
}
