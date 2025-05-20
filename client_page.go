package ehclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Page struct {
	GalleryId   int
	Page        int
	Token       string
	FileName    string
	Url         string
	RawUrl      string
	Prev        *PageTokenPair
	Next        *PageTokenPair
	ReloadToken string
	ApiToken    string
}

type GetPageOptions struct {
	ApiToken    string
	ReloadToken string
}

type showPageResult struct {
	Page              int    `json:"p"`
	Uri               string `json:"s"`
	NavHtml           string `json:"n"`
	InfoHtml          string `json:"i"`
	Token             string `json:"k"`
	ImageHtml         string `json:"i3"`
	BackHtml          string `json:"i5"`
	FooterHtml        string `json:"i6"`
	ReloadTokenPrefix string `json:"si"`
	Widht             string `json:"x"`
	Height            string `json:"y"`
}

func (c *Client) GetPage(id int, page int, token string, opt *GetPageOptions) (*Page, error) {
	if opt == nil || opt.ApiToken == "" {
		u := fmt.Sprintf("https://%s/s/%s/%d-%d", string(c.opts.Endpoint), token, id, page)

		if opt != nil && opt.ReloadToken != "" {
			u += fmt.Sprintf("?ul=%s", opt.ReloadToken)
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

		return c.parser.parsePage(doc)
	} else {
		u := "https://api.e-hentai.org/api.php"
		payload := map[string]any{
			"method":  "showpage",
			"gid":     id,
			"page":    page,
			"imgkey":  token,
			"showkey": opt.ApiToken,
		}
		body := &bytes.Buffer{}
		err := json.NewEncoder(body).Encode(payload)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest("POST", u, body)
		if err != nil {
			return nil, err
		}

		resp, err := c.Do(req)
		if err != nil {
			return nil, err
		}

		var result showPageResult
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return nil, err
		}

		return c.parser.parsePageApi(&result)
	}
}
