package ehclient

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type GalleryDetail struct {
	Id         int
	Token      string
	Thumb      string
	Title      string
	TitleJpn   string
	Category   Category
	Uploader   string
	PostedAt   time.Time
	Parent     *IdTokenPair
	Visible    string
	Language   string
	FileSize   string
	Pages      int
	Favoriteds int
	Rating     float64
	Rateds     int
	Tags       []*TagGroup
	Versions   []*GalleryVersion
	Previews   []*PagePreview
	Comments   []*GalleryComment
}

type GalleryVersion struct {
	Id      int
	Token   string
	Title   string
	AddedAt time.Time
}

type PagePreview struct {
	Page     int
	Token    string
	FileName string
	Thumb    string
	Width    int
	Height   int
	Offset   int
}

type GalleryComment struct {
	PostedAt     time.Time
	LastEditedAt *time.Time
	User         string
	Score        int
	Content      []*html.Node
}

type GetGalleryOption struct {
	PreviewPage       int
	WithHiddenComment bool
}

func (c *Client) GetGallery(id int, token string, opt *GetGalleryOption) (*GalleryDetail, error) {
	u := fmt.Sprintf("https://%s/g/%d/%s/", string(c.opts.Endpoint), id, token)
	if opt != nil {
		values := make(url.Values)
		if opt.PreviewPage != 0 {
			values.Add("p", strconv.Itoa(opt.PreviewPage))
		}
		if opt.WithHiddenComment {
			values.Add("hc", "1")
		}
		query := values.Encode()
		if query != "" {
			u += "?" + query
		}
	}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return c.parser.parseGallery(doc)
}
