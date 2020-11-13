package model

import (
	"regexp"

	"mvdan.cc/xurls/v2"
)

// Message ...
type Message struct {
	Text     string
	Entities MessageEntities
}

// MessageEntities ...
type MessageEntities struct {
	URLs     []MessageEntityURL
	Mentions []MessageEntityMention
	Hashtags []MessageEntityHash
}

// MessageEntityIndices ...
type MessageEntityIndices struct {
	Begin int
	End   int
}

// MessageEntity ...
type MessageEntity struct {
	Indices MessageEntityIndices
}

// MessageEntityURL ...
type MessageEntityURL struct {
	MessageEntity
	URL        string
	DisplayURL string
}

// MessageEntityMention ...
type MessageEntityMention struct {
	MessageEntity
	Name string
}

// MessageEntityHash ...
type MessageEntityHash struct {
	MessageEntity
	Text string
}

// ParseURL ...
func ParseURL(text string, me *MessageEntities) {
	matcher := xurls.Strict()
	indices := matcher.FindAllStringIndex(text, -1)
	for _, index := range indices {
		uri := text[index[0]:index[1]]
		m := MessageEntityURL{
			MessageEntity: MessageEntity{
				Indices: MessageEntityIndices{
					Begin: index[0],
					End:   index[1],
				},
			},
			URL:        uri,
			DisplayURL: uri,
		}
		me.URLs = append(me.URLs, m)
	}
}

// ParseMention ...
func ParseMention(text string, me *MessageEntities) {
	matcher := regexp.MustCompile("(@.+)")
	indices := matcher.FindAllStringIndex(text, -1)
	for _, index := range indices {
		mention := text[index[0]:index[1]]
		m := MessageEntityMention{
			MessageEntity: MessageEntity{
				Indices: MessageEntityIndices{
					Begin: index[0],
					End:   index[1],
				},
			},
			Name: mention,
		}
		me.Mentions = append(me.Mentions, m)
	}
}
