package traq

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func makeResponse(statusCode int, body string, headers map[string]string) *http.Response {
	resp := &http.Response{
		StatusCode: statusCode,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
	for k, v := range headers {
		resp.Header.Set(k, v)
	}
	return resp
}

func TestIsETagTargetPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		path        string
		expected    bool
	}{
		{
			description: "/users suffix returns true",
			path:        "/api/v3/users",
			expected:    true,
		},
		{
			description: "/users with trailing slash returns true",
			path:        "/api/v3/users/",
			expected:    true,
		},
		{
			description: "/groups suffix returns true",
			path:        "/api/v3/groups",
			expected:    true,
		},
		{
			description: "/groups with trailing slash returns true",
			path:        "/api/v3/groups/",
			expected:    true,
		},
		{
			description: "/stamps suffix returns true",
			path:        "/api/v3/stamps",
			expected:    true,
		},
		{
			description: "/stamps with trailing slash returns true",
			path:        "/api/v3/stamps/",
			expected:    true,
		},
		{
			description: "/channels suffix returns true",
			path:        "/api/v3/channels",
			expected:    true,
		},
		{
			description: "/channels with trailing slash returns true",
			path:        "/api/v3/channels/",
			expected:    true,
		},
		{
			description: "bare /users returns true",
			path:        "/users",
			expected:    true,
		},
		{
			description: "specific user ID path returns false",
			path:        "/api/v3/users/abc123",
			expected:    false,
		},
		{
			description: "specific group ID path returns false",
			path:        "/api/v3/groups/abc123",
			expected:    false,
		},
		{
			description: "unrelated path returns false",
			path:        "/api/v3/messages",
			expected:    false,
		},
		{
			description: "empty path returns false",
			path:        "",
			expected:    false,
		},
		{
			description: "root path returns false",
			path:        "/",
			expected:    false,
		},
		{
			description: "path containing users but not as suffix returns false",
			path:        "/api/v3/users/123/tags",
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			result := isETagTargetPath(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewETagCacheTransport(t *testing.T) {
	t.Parallel()

	t.Run("nil base uses http.DefaultTransport", func(t *testing.T) {
		t.Parallel()
		transport := newETagCacheTransport(nil)
		require.NotNil(t, transport)
	})

	t.Run("non-nil base is used", func(t *testing.T) {
		t.Parallel()
		base := roundTripFunc(func(_ *http.Request) (*http.Response, error) {
			return nil, nil
		})
		transport := newETagCacheTransport(base)
		require.NotNil(t, transport)
	})
}

func TestETagCacheTransport_RoundTrip_NonGET(t *testing.T) {
	t.Parallel()

	called := 0
	base := roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		called++
		return makeResponse(http.StatusOK, "ok", nil), nil
	})

	transport := newETagCacheTransport(base).(*etagCacheTransport)

	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "http://example.com/api/v3/users", nil)
			resp, err := transport.RoundTrip(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
	assert.Equal(t, len(methods), called, "base should be called once per non-GET method")
}

func TestETagCacheTransport_RoundTrip_NonCacheablePath(t *testing.T) {
	t.Parallel()

	callCount := 0
	base := roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		callCount++
		return makeResponse(http.StatusOK, "data", nil), nil
	})

	transport := newETagCacheTransport(base).(*etagCacheTransport)
	req := httptest.NewRequest(http.MethodGet, "http://example.com/api/v3/messages", nil)

	resp, err := transport.RoundTrip(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	req2 := httptest.NewRequest(http.MethodGet, "http://example.com/api/v3/messages", nil)
	resp2, err := transport.RoundTrip(req2)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp2.StatusCode)
	assert.Equal(t, 2, callCount)
}

func TestETagCacheTransport_RoundTrip_FirstRequestCachesETag(t *testing.T) {
	t.Parallel()

	const etag = `"abc123"`
	const body = `[{"id":"u1"}]`

	callCount := 0
	var receivedIfNoneMatch string
	base := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		callCount++
		receivedIfNoneMatch = r.Header.Get("If-None-Match")
		return makeResponse(http.StatusOK, body, map[string]string{"ETag": etag}), nil
	})

	transport := newETagCacheTransport(base).(*etagCacheTransport)
	req := httptest.NewRequest(http.MethodGet, "http://example.com/api/v3/users", nil)

	resp, err := transport.RoundTrip(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 1, callCount)
	assert.Empty(t, receivedIfNoneMatch, "first request must not send If-None-Match")

	req2 := httptest.NewRequest(http.MethodGet, "http://example.com/api/v3/users", nil)
	_, err = transport.RoundTrip(req2)
	require.NoError(t, err)
	assert.Equal(t, etag, receivedIfNoneMatch, "second request must send stored ETag as If-None-Match")
}

func TestETagCacheTransport_RoundTrip_304UsesCachedBody(t *testing.T) {
	t.Parallel()

	const etag = `"v1"`
	const cachedBody = `["user1","user2"]`

	callNum := 0
	base := roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		callNum++
		if callNum == 1 {
			return makeResponse(http.StatusOK, cachedBody, map[string]string{"ETag": etag}), nil
		}
		return makeResponse(http.StatusNotModified, "", nil), nil
	})

	transport := newETagCacheTransport(base).(*etagCacheTransport)

	req1 := httptest.NewRequest(http.MethodGet, "http://example.com/api/v3/users", nil)
	_, err := transport.RoundTrip(req1)
	require.NoError(t, err)

	req2 := httptest.NewRequest(http.MethodGet, "http://example.com/api/v3/users", nil)
	resp, err := transport.RoundTrip(req2)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, cachedBody, string(bodyBytes))
}

func TestETagCacheTransport_RoundTrip_304NoCacheFallthrough(t *testing.T) {
	t.Parallel()

	base := roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		return makeResponse(http.StatusNotModified, "", nil), nil
	})

	transport := newETagCacheTransport(base).(*etagCacheTransport)
	req := httptest.NewRequest(http.MethodGet, "http://example.com/api/v3/users", nil)
	resp, err := transport.RoundTrip(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotModified, resp.StatusCode)
}

func TestETagCacheTransport_RoundTrip_5xxFallback(t *testing.T) {
	t.Parallel()

	const etag = `"v2"`
	const cachedBody = `["stamp1"]`

	callNum := 0
	base := roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		callNum++
		if callNum == 1 {
			return makeResponse(http.StatusOK, cachedBody, map[string]string{"ETag": etag}), nil
		}
		return makeResponse(http.StatusInternalServerError, "error", nil), nil
	})

	transport := newETagCacheTransport(base).(*etagCacheTransport)

	req1 := httptest.NewRequest(http.MethodGet, "http://example.com/api/v3/stamps", nil)
	_, err := transport.RoundTrip(req1)
	require.NoError(t, err)

	req2 := httptest.NewRequest(http.MethodGet, "http://example.com/api/v3/stamps", nil)
	resp, err := transport.RoundTrip(req2)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "5xx should fall back to cached 200")

	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, cachedBody, string(bodyBytes))
}

func TestETagCacheTransport_RoundTrip_5xxNoCacheFallthrough(t *testing.T) {
	t.Parallel()

	base := roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		return makeResponse(http.StatusServiceUnavailable, "unavailable", nil), nil
	})

	transport := newETagCacheTransport(base).(*etagCacheTransport)
	req := httptest.NewRequest(http.MethodGet, "http://example.com/api/v3/users", nil)
	resp, err := transport.RoundTrip(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
}

func TestETagCacheTransport_RoundTrip_NetworkErrorWithCache(t *testing.T) {
	t.Parallel()

	const etag = `"v3"`
	const cachedBody = `["group1"]`
	networkErr := errors.New("connection refused")

	callNum := 0
	base := roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		callNum++
		if callNum == 1 {
			return makeResponse(http.StatusOK, cachedBody, map[string]string{"ETag": etag}), nil
		}
		return nil, networkErr
	})

	transport := newETagCacheTransport(base).(*etagCacheTransport)

	req1 := httptest.NewRequest(http.MethodGet, "http://example.com/api/v3/groups", nil)
	_, err := transport.RoundTrip(req1)
	require.NoError(t, err)

	req2 := httptest.NewRequest(http.MethodGet, "http://example.com/api/v3/groups", nil)
	resp, err := transport.RoundTrip(req2)
	require.NoError(t, err, "network error should be suppressed when cache is available")
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, cachedBody, string(bodyBytes))
}

func TestETagCacheTransport_RoundTrip_NetworkErrorNoCache(t *testing.T) {
	t.Parallel()

	networkErr := errors.New("connection refused")
	base := roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		return nil, networkErr
	})

	transport := newETagCacheTransport(base).(*etagCacheTransport)
	req := httptest.NewRequest(http.MethodGet, "http://example.com/api/v3/users", nil)
	resp, err := transport.RoundTrip(req)
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, networkErr)
}

func TestETagCacheTransport_RoundTrip_NonOKNonSpecialStatus(t *testing.T) {
	t.Parallel()

	base := roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		return makeResponse(http.StatusUnauthorized, "unauthorized", nil), nil
	})

	transport := newETagCacheTransport(base).(*etagCacheTransport)
	req := httptest.NewRequest(http.MethodGet, "http://example.com/api/v3/users", nil)
	resp, err := transport.RoundTrip(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestETagCacheTransport_RoundTrip_200WithoutETagNotCached(t *testing.T) {
	t.Parallel()

	callCount := 0
	var receivedIfNoneMatch string
	base := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		callCount++
		receivedIfNoneMatch = r.Header.Get("If-None-Match")
		return makeResponse(http.StatusOK, "data", nil), nil
	})

	transport := newETagCacheTransport(base).(*etagCacheTransport)

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "http://example.com/api/v3/users", nil)
		_, err := transport.RoundTrip(req)
		require.NoError(t, err)
	}

	assert.Empty(t, receivedIfNoneMatch, "no ETag stored, so If-None-Match must never be sent")
	assert.Equal(t, 2, callCount)
}

func TestETagCacheTransport_RoundTrip_CacheHitHeader(t *testing.T) {
	t.Parallel()

	const etag = `"hdr-test"`

	callNum := 0
	base := roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		callNum++
		if callNum == 1 {
			return makeResponse(http.StatusOK, "body", map[string]string{"ETag": etag}), nil
		}
		return makeResponse(http.StatusNotModified, "", nil), nil
	})

	transport := newETagCacheTransport(base).(*etagCacheTransport)

	req1 := httptest.NewRequest(http.MethodGet, "http://example.com/api/v3/channels", nil)
	_, _ = transport.RoundTrip(req1)

	req2 := httptest.NewRequest(http.MethodGet, "http://example.com/api/v3/channels", nil)
	resp, err := transport.RoundTrip(req2)
	require.NoError(t, err)
	assert.Equal(t, "hit", resp.Header.Get("X-AnkeTo-Cache"))
}

func TestETagCacheTransport_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	const etag = `"concurrent"`
	const body = `["a","b"]`

	base := roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		return makeResponse(http.StatusOK, body, map[string]string{"ETag": etag}), nil
	})

	transport := newETagCacheTransport(base).(*etagCacheTransport)

	const goroutines = 20
	done := make(chan struct{}, goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer func() { done <- struct{}{} }()
			req := httptest.NewRequest(http.MethodGet, "http://example.com/api/v3/users", nil)
			resp, err := transport.RoundTrip(req)
			if err == nil && resp != nil {
				resp.Body.Close()
			}
		}()
	}
	for i := 0; i < goroutines; i++ {
		<-done
	}
}
