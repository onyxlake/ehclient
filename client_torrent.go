package ehclient

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Torrent struct {
	Hash       string
	Name       string
	PostedAt   time.Time
	Uploader   string
	Size       string
	Seeds      int
	Peers      int
	Downloads  int
	IsOutdated bool
}

func (c *Client) GetTorrent(gid int, token string) ([]Torrent, error) {
	u := fmt.Sprintf("https://%s/gallerytorrents.php?gid=%d&t=%s", string(c.opts.Endpoint), gid, token)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return c.parser.parseTorrent(doc)
}

func (c *Client) BuildTorrentUrl(gid int, hash string) string {
	return fmt.Sprintf("https://ehtracker.org/get/%d/%s.torrent", gid, hash)
}
