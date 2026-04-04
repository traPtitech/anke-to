package traq

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		message     string
		limit       int
		expected    []string
	}{
		{
			description: "limit以内のメッセージはそのまま返す",
			message:     "hello world",
			limit:       100,
			expected:    []string{"hello world"},
		},
		{
			description: "ちょうどlimitのメッセージはそのまま返す",
			message:     strings.Repeat("a", 10),
			limit:       10,
			expected:    []string{strings.Repeat("a", 10)},
		},
		{
			description: "改行で分割できる場合は改行単位で分割する",
			message:     "line1\nline2\nline3",
			limit:       11,
			expected:    []string{"line1\nline2", "line3"},
		},
		{
			description: "対象者行がlimitを超える場合はスペース単位で分割する",
			message:     "@user1 @user2 @user3 @user4 @user5",
			limit:       20,
			expected:    []string{"@user1 @user2 @user3", "@user4 @user5"},
		},
		{
			description: "ヘッダーと長い対象者行を含むメッセージを正しく分割する",
			message:     "#### 対象者\n@user1 @user2 @user3 @user4\n#### 回答リンク\nhttps://example.com",
			limit:       40,
			// "#### 対象者\n@user1 @user2 @user3 @user4" = 9+1+24=34 ≤ 40
			// "@user5..." 以降は元メッセージにないので次チャンクは残り
			expected: []string{
				"#### 対象者\n@user1 @user2 @user3 @user4",
				"#### 回答リンク\nhttps://example.com",
			},
		},
		{
			description: "日本語文字列もrune単位で正しく分割する",
			message:     "あいうえおかきくけこ さしすせそ",
			limit:       10,
			expected:    []string{"あいうえおかきくけこ", "さしすせそ"},
		},
		{
			description: "複数行でlimitをまたぐ場合に全チャンクがlimit以内に収まる",
			message:     "#### ヘッダー\n" + strings.Repeat("@user%d ", 300),
			limit:       2000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			result := splitMessage(tt.message, tt.limit)
			assert.NotEmpty(t, result)

			// 全チャンクがlimit以内であることを確認
			for _, chunk := range result {
				assert.LessOrEqual(t, len([]rune(chunk)), tt.limit, "chunk exceeded limit: %q", chunk)
			}

			// expectedが指定されている場合は内容も確認
			if tt.expected != nil {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
