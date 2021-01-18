package model

import (
	"mime"
	"net/http"
	"regexp"
	"strconv"

	"github.com/otiai10/opengraph/v2"
	"golang.org/x/xerrors"
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
	MediaType     string // RFC-1521
	ContentLength int64
	URL           string
	DisplayURL    string
	OpenGraph     *OpenGraph
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
			URL:           uri,
			DisplayURL:    uri,
			MediaType:     "",
			OpenGraph:     nil,
			ContentLength: 0,
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

// OpenGraph ...
type OpenGraph struct {
	Title       string
	Type        string
	Image       []OpenGraphImage
	URL         string
	Audio       []OpenGraphAudio
	Description string
	Determiner  string
	Locale      string
	LocaleAlt   []string
	SiteName    string
	Video       []OpenGraphVideo
}

// OpenGraphImage ...
type OpenGraphImage struct {
	URL       string
	SecureURL string
	Type      string
	Width     int
	Height    int
	Alt       string
}

// OpenGraphAudio ....
type OpenGraphAudio struct {
	URL       string
	SecureURL string
	Type      string
}

// OpenGraphVideo ...
type OpenGraphVideo struct {
	URL       string
	SecureURL string
	Type      string
	Width     int
	Height    int
	Duration  int
}

func newOpenGraphFromOpenGraphGo(ogp *opengraph.OpenGraph) *OpenGraph {
	return &OpenGraph{
		Title:       ogp.Title,
		Type:        ogp.Type,
		Image:       newOpenGraphImageFromOpenGraphGo(ogp.Image),
		URL:         ogp.URL,
		Audio:       newOpenGraphAudioFromOpenGraphGo(ogp.Audio),
		Description: ogp.Description,
		Locale:      ogp.Locale,
		LocaleAlt:   ogp.LocaleAlt,
		SiteName:    ogp.SiteName,
		Video:       newOpenGraphVideoFromOpenGraphGo(ogp.Video),
	}
}

func newOpenGraphImageFromOpenGraphGo(image []opengraph.Image) []OpenGraphImage {
	r := []OpenGraphImage{}
	for _, i := range image {
		r = append(r, OpenGraphImage{
			URL:       i.URL,
			SecureURL: i.SecureURL,
			Type:      i.Type,
			Width:     i.Width,
			Height:    i.Height,
			Alt:       i.Alt,
		})
	}
	return r
}

func newOpenGraphAudioFromOpenGraphGo(audio []opengraph.Audio) []OpenGraphAudio {
	r := []OpenGraphAudio{}
	for _, i := range audio {
		r = append(r, OpenGraphAudio{
			URL:       i.URL,
			SecureURL: i.SecureURL,
			Type:      i.Type,
		})
	}
	return r
}

func newOpenGraphVideoFromOpenGraphGo(video []opengraph.Video) []OpenGraphVideo {
	r := []OpenGraphVideo{}
	for _, i := range video {
		r = append(r, OpenGraphVideo{
			URL:       i.URL,
			SecureURL: i.SecureURL,
			Type:      i.Type,
			Width:     i.Width,
			Height:    i.Height,
			Duration:  i.Duration,
		})
	}
	return r
}

func fetchOpenGraph(url string) (*OpenGraph, error) {
	ogp, err := opengraph.Fetch(url)
	if err != nil {
		return nil, xerrors.Errorf("Cannot fetch ogp : %w", err)
	}
	return newOpenGraphFromOpenGraphGo(ogp), nil
}

func HeadURL(cli *http.Client, me *MessageEntities) []error {
	errs := []error{}
	for i := range me.URLs {
		res, err := cli.Head(me.URLs[i].URL)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		contentType := ""
		contentType = res.Header.Get("Content-Type")
		if contentType == "" {
			contentType = res.Header.Get("content-type")
		}
		if contentType == "" {
			contentType = res.Header.Get("Content-type")
		}
		if contentType == "" {
			contentType = res.Header.Get("content-Type")
		}
		mediaType, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			errs = append(errs, err)
		} else {
			me.URLs[i].MediaType = mediaType
		}

		contentLength := ""
		contentLength = res.Header.Get("Content-Length")
		if contentLength == "" {
			contentLength = res.Header.Get("content-length")
		}
		if contentLength == "" {
			contentLength = res.Header.Get("Content-length")
		}
		if contentLength == "" {
			contentLength = res.Header.Get("content-Length")
		}
		contentLengthInt, err := strconv.ParseInt(contentLength, 10, 64)
		if err != nil {
			errs = append(errs, err)
		} else {
			me.URLs[i].ContentLength = contentLengthInt
		}
	}
	return errs
}

func FetchOpenGraph(cli *http.Client, me *MessageEntities) []error {
	errs := []error{}
	for i := range me.URLs {
		if me.URLs[i].MediaType != "text/html" {
			continue
		}
		ogp, err := fetchOpenGraph(me.URLs[i].URL)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		me.URLs[i].OpenGraph = ogp
	}
	return errs
}
