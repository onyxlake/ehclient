package ehclient

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (p *Parser) parseGalleries(doc *goquery.Document) (*GalleryCursor, error) {
	var result GalleryCursor

	selectNavSelection := doc.FindMatcher(p.Single(".searchnav"))

	if node := selectNavSelection.Children().Last().FindMatcher(p.Single("option[selected]")); node.Length() != 0 {
		kind, exist := node.Attr("value")
		if !exist {
			return nil, newAttrNotFoundError("galleries", "preview_type", "value")
		}
		result.Kind = PreviewKind(kind)
	} else {
		return nil, newNodeNotFoundError("galleries", "preview_type")
	}

	if node := selectNavSelection.FindMatcher(p.Single("#uprev")); node.Length() != 0 {
		if href, exist := node.Attr("href"); exist {
			href = href[strings.Index(href, "?")+1:]
			result.Prev = href
		}
	} else {
		return nil, newNodeNotFoundError("galleries", "prev")
	}

	if node := selectNavSelection.FindMatcher(p.Single("#unext")); node.Length() != 0 {
		if href, exist := node.Attr("href"); exist {
			href = href[strings.Index(href, "?")+1:]
			result.Next = href
		}
	} else {
		return nil, newNodeNotFoundError("galleries", "next")
	}

	var (
		previews []*GalleryPreview
		preview  *GalleryPreview
		err      error
	)
	if result.Kind != Thumbnail {
		for i, selection := range doc.FindMatcher(p.Matcher(".itg>tbody>tr")).EachIter() {
			if (result.Kind != Extended && i == 0) || selection.FindMatcher(p.Single(".itd")).Length() != 0 {
				continue
			}

			switch result.Kind {
			case Compact:
				preview, err = p.parseCompactGalleryPreview(selection)
			case Extended:
				preview, err = p.parseExpendedGalleryPreview(selection)
			default:

			}
			if err != nil {
				return nil, err
			}
			previews = append(previews, preview)
		}
	} else {
		// for i, selection := range doc.Find(".itg.gld>div").EachIter() {

		// }
	}
	result.Previews = previews

	return &result, nil
}
