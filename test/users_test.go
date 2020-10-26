package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsersMe(t *testing.T) {
	t.Parallel()

	rec := createRecorder(userUnAuthorized, methodGet, makePath("/users/me"), typeNone, "")

	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	u := userOne
	expectBody := fmt.Sprintf(
`{
  "traqID": "%s"
}`, u)
	rec = createRecorder(u, methodGet, makePath("/users/me"), typeNone, "")

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, expectBody, rec.Body.String())
}
