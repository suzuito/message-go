package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseURL(t *testing.T) {
	testCases := []struct {
		desc                  string
		inputText             string
		inputMessageEntities  MessageEntities
		expectMessageEntities MessageEntities
	}{
		{
			desc:                  "Empty inputText",
			inputText:             "",
			inputMessageEntities:  MessageEntities{},
			expectMessageEntities: MessageEntities{},
		},
		{
			desc:                  "Not empty inputText",
			inputText:             "abc def",
			inputMessageEntities:  MessageEntities{},
			expectMessageEntities: MessageEntities{},
		},
		{
			desc:                 "InputText includes URL",
			inputText:            "http://example.com/hoge",
			inputMessageEntities: MessageEntities{},
			expectMessageEntities: MessageEntities{
				URLs: []MessageEntityURL{
					{
						MessageEntity: MessageEntity{
							Indices: MessageEntityIndices{
								Begin: 0,
								End:   23,
							},
						},
						URL:        "http://example.com/hoge",
						DisplayURL: "http://example.com/hoge",
					},
				},
			},
		},
		{
			desc:                 "InputText includes URL",
			inputText:            "abchttps://example.com/hoge.json?hoge=fuga def",
			inputMessageEntities: MessageEntities{},
			expectMessageEntities: MessageEntities{
				URLs: []MessageEntityURL{
					{
						MessageEntity: MessageEntity{
							Indices: MessageEntityIndices{
								Begin: 3,
								End:   42,
							},
						},
						URL:        "https://example.com/hoge.json?hoge=fuga",
						DisplayURL: "https://example.com/hoge.json?hoge=fuga",
					},
				},
			},
		},
		{
			desc:                 "InputText includes URL",
			inputText:            "abc https://example.com/hoge.json?hoge=fuga def https://example.com/hoge.json?foo=bar ghi",
			inputMessageEntities: MessageEntities{},
			expectMessageEntities: MessageEntities{
				URLs: []MessageEntityURL{
					{
						MessageEntity: MessageEntity{
							Indices: MessageEntityIndices{
								Begin: 4,
								End:   43,
							},
						},
						URL:        "https://example.com/hoge.json?hoge=fuga",
						DisplayURL: "https://example.com/hoge.json?hoge=fuga",
					},
					{
						MessageEntity: MessageEntity{
							Indices: MessageEntityIndices{
								Begin: 48,
								End:   85,
							},
						},
						URL:        "https://example.com/hoge.json?foo=bar",
						DisplayURL: "https://example.com/hoge.json?foo=bar",
					},
				},
			},
		},
		{
			desc:                 "InputText includes URL",
			inputText:            "な https://example.com/あいう に https://example.com/えお ぬ",
			inputMessageEntities: MessageEntities{},
			expectMessageEntities: MessageEntities{
				URLs: []MessageEntityURL{
					{
						MessageEntity: MessageEntity{
							Indices: MessageEntityIndices{
								Begin: 4,
								End:   33,
							},
						},
						URL:        "https://example.com/あいう",
						DisplayURL: "https://example.com/あいう",
					},
					{
						MessageEntity: MessageEntity{
							Indices: MessageEntityIndices{
								Begin: 38,
								End:   64,
							},
						},
						URL:        "https://example.com/えお",
						DisplayURL: "https://example.com/えお",
					},
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ParseURL(tC.inputText, &tC.inputMessageEntities)
			assert.Equal(t, tC.expectMessageEntities, tC.inputMessageEntities)
		})
	}
}

func TestParseMention(t *testing.T) {
	testCases := []struct {
		desc                  string
		inputText             string
		inputMessageEntities  MessageEntities
		expectMessageEntities MessageEntities
	}{
		{
			desc:                  "Empty inputText",
			inputText:             "",
			inputMessageEntities:  MessageEntities{},
			expectMessageEntities: MessageEntities{},
		},
		{
			desc:                 "Empty inputText",
			inputText:            "@hoge",
			inputMessageEntities: MessageEntities{},
			expectMessageEntities: MessageEntities{
				Mentions: []MessageEntityMention{
					{
						MessageEntity: MessageEntity{
							Indices: MessageEntityIndices{
								Begin: 0,
								End:   5,
							},
						},
						Name: "@hoge",
					},
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ParseMention(tC.inputText, &tC.inputMessageEntities)
			assert.Equal(t, tC.expectMessageEntities, tC.inputMessageEntities)
		})
	}
}
