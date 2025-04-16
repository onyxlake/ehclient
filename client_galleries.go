package ehclient

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type PreviewKind string

var (
	Minimal     PreviewKind = "m"
	MinimalPlus PreviewKind = "p"
	Compact     PreviewKind = "l"
	Extended    PreviewKind = "e"
	Thumbnail   PreviewKind = "t"
)

type GalleryCursor struct {
	Kind     PreviewKind
	Prev     string
	Next     string
	Previews []*GalleryPreview
}

type GalleryPreview struct {
	Id         int
	Token      string
	Thumb      string
	Title      string
	Category   Category
	PostedAt   time.Time
	Rating     float64
	Uploader   string
	Pages      int
	HasTorrent bool
	Tags       []*TagGroup
}

type SearchGalleriesOption struct {
	Keywords                      string
	Categories                    int
	IncludeExpunged               bool
	MustHaveTorrent               bool
	PagesRangeStart               int
	PagesRangeEnd                 int
	MinimunRating                 int
	DisableCustomFilterOfLanguage bool
	DisableCustomFilterOfUploader bool
	DisableCustomFilterOfTags     bool
	PreviewKind                   PreviewKind
	Next                          int

	RawQuery string
}

func (c *Client) SearchGalleries(opt *SearchGalleriesOption) (*GalleryCursor, error) {
	var query string
	if opt.RawQuery != "" {
		query = opt.RawQuery
	} else {
		values := make(url.Values)
		if opt.Next != 0 {
			values.Add("next", strconv.Itoa(opt.Next))
		}
		if opt.Categories != 0 {
			values.Add("f_cats", strconv.Itoa(opt.Categories))
		}
		if opt.Keywords != "" {
			values.Add("f_search", opt.Keywords)
		}
		if opt.IncludeExpunged {
			values.Add("f_sh", "on")
		}
		if opt.MustHaveTorrent {
			values.Add("f_sto", "on")
		}
		if opt.PagesRangeStart != 0 && opt.PagesRangeEnd != 0 {
			values.Add("f_spf", strconv.Itoa(opt.PagesRangeStart))
			values.Add("f_spt", strconv.Itoa(opt.PagesRangeEnd))
		}
		if opt.MinimunRating != 0 {
			values.Add("f_srdd", strconv.Itoa(opt.MinimunRating))
		}
		if opt.DisableCustomFilterOfLanguage {
			values.Add("f_sfl", "on")
		}
		if opt.DisableCustomFilterOfUploader {
			values.Add("f_sfu", "on")
		}
		if opt.DisableCustomFilterOfTags {
			values.Add("f_sft", "on")
		}
		if opt.PreviewKind != "" {
			values.Add("inline_set", "dm_"+string(opt.PreviewKind))
		}
		query = values.Encode()
	}

	u := "https://" + string(c.opts.Endpoint)
	if query != "" {
		u += "?" + query
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
	return c.parser.parseGalleries(doc)
}
