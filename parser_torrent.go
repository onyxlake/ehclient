package ehclient

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (p *Parser) parseTorrent(doc *goquery.Document) ([]Torrent, error) {
	var result []Torrent

	for _, selection := range doc.FindMatcher(p.Matcher("form>div")).EachIter() {
		var torrent Torrent

		trs := selection.FindMatcher(p.Matcher("tr"))
		tds := trs.Eq(0).FindMatcher(p.Matcher("td"))

		if node := tds.Eq(0); node.Length() != 0 {
			postedAtStr := node.Nodes[0].LastChild.FirstChild.Data
			posted, err := p.parseBaseTime(postedAtStr)
			if err != nil {
				return nil, newParserParseError("torrent", "postedAt", postedAtStr, err)
			}
			torrent.PostedAt = posted

			if len(node.Nodes[0].LastChild.Attr) != 0 {
				torrent.IsOutdated = true
			}
		} else {
			return nil, newNodeNotFoundError("torrent", "postedAt")
		}

		if node := tds.Eq(1); node.Length() != 0 {
			size := node.Nodes[0].LastChild.Data
			torrent.Size = strings.TrimSpace(size)
		} else {
			return nil, newNodeNotFoundError("torrent", "size")
		}

		if node := tds.Eq(3); node.Length() != 0 {
			seedsStr := node.Nodes[0].LastChild.Data
			seedsStr = strings.TrimSpace(seedsStr)
			seeds, err := strconv.Atoi(seedsStr)
			if err != nil {
				return nil, newParserParseError("torrent", "seeds", seedsStr, err)
			}
			torrent.Seeds = seeds
		} else {
			return nil, newNodeNotFoundError("torrent", "seeds")
		}

		if node := tds.Eq(4); node.Length() != 0 {
			peersStr := node.Nodes[0].LastChild.Data
			peersStr = strings.TrimSpace(peersStr)
			peers, err := strconv.Atoi(peersStr)
			if err != nil {
				return nil, newParserParseError("torrent", "peers", peersStr, err)
			}
			torrent.Peers = peers
		} else {
			return nil, newNodeNotFoundError("torrent", "peers")
		}

		if node := tds.Eq(5); node.Length() != 0 {
			downloadsStr := node.Nodes[0].LastChild.Data
			downloadsStr = strings.TrimSpace(downloadsStr)
			downloads, err := strconv.Atoi(downloadsStr)
			if err != nil {
				return nil, newParserParseError("torrent", "downloads", downloadsStr, err)
			}
			torrent.Downloads = downloads
		} else {
			return nil, newNodeNotFoundError("torrent", "downloads")
		}

		if node := trs.Eq(1).FindMatcher(p.Single("td")).First(); node.Length() != 0 {
			uploader := node.Nodes[0].LastChild.Data
			torrent.Uploader = strings.TrimSpace(uploader)
		} else {
			return nil, newNodeNotFoundError("torrent", "uploader")
		}

		if node := trs.Eq(2).FindMatcher(p.Single("a")); node.Length() != 0 {
			name := node.Text()
			name = strings.TrimSpace(name)
			torrent.Name = name

			if href, exist := node.Attr("href"); exist {
				href = href[strings.LastIndex(href, "/")+1 : strings.LastIndex(href, ".")]
				torrent.Hash = href
			} else {
				return nil, newAttrNotFoundError("torrent", "hash", "href")
			}
		} else {
			return nil, newNodeNotFoundError("torrent", "hash")
		}

		result = append(result, torrent)
	}

	return result, nil
}
