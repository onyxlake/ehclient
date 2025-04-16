package ehclient

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (p *Parser) parseExpendedGalleryPreview(parent *goquery.Selection) (*GalleryPreview, error) {
	var result GalleryPreview

	linkSelection := parent.FindMatcher(p.Single(".gl1e>div>a"))
	if linkSelection.Length() != 0 {
		href, exist := linkSelection.Attr("href")
		if !exist {
			return nil, newAttrNotFoundError("gallery_preview", "link", "href")
		}
		pair, err := p.parseIdTokenPairFromHref(href)
		if err != nil {
			return nil, newParserParseError("gallery_preview", "link", href, err)
		}
		result.Id = pair.Id
		result.Token = pair.Token
	} else {
		return nil, newNodeNotFoundError("gallery_preview", "link")
	}

	thumbSelection := linkSelection.Children().First()
	if thumbSelection.Length() != 0 {
		src, exist := thumbSelection.Attr("src")
		if !exist {
			return nil, newAttrNotFoundError("gallery_preview", "thumb", "src")
		}
		result.Thumb = src
	} else {
		return nil, newNodeNotFoundError("gallery_preview", "thumb")
	}

	column2Selection := parent.FindMatcher(p.Single(".gl3e")).Children()

	if node := column2Selection.Eq(0); node.Length() != 0 {
		categoryStr := node.Text()
		category, err := parseCategoryFromLabel(categoryStr)
		if err != nil {
			return nil, newParserParseError("gallery_preview", "category", categoryStr, err)
		}
		result.Category = category
	} else {
		return nil, newNodeNotFoundError("gallery_preview", "category")
	}

	if node := column2Selection.Eq(1); node.Length() != 0 {
		postedAtStr := node.Text()
		postedAt, err := p.parseBaseTime(postedAtStr)
		if err != nil {
			return nil, newParserParseError("gallery_preview", "posted_at", postedAtStr, err)
		}
		result.PostedAt = postedAt
	} else {
		return nil, newNodeNotFoundError("gallery_preview", "posted_at")
	}

	if node := column2Selection.Eq(2); node.Length() != 0 {
		style, exist := node.Attr("style")
		if !exist {
			return nil, newAttrNotFoundError("gallery_preview", "rating", "style")
		}
		rating, err := p.parseRatingFromStyle(style)
		if err != nil {
			return nil, newParserParseError("gallery_preview", "rating", style, err)
		}
		result.Rating = rating
	} else {
		return nil, newNodeNotFoundError("gallery_preview", "rating")
	}

	if node := column2Selection.Eq(3); node.Length() != 0 {
		uploader := node.Text()
		result.Uploader = uploader
	} else {
		return nil, newNodeNotFoundError("gallery_preview", "uploader")
	}

	if node := column2Selection.Eq(4); node.Length() != 0 {
		pagesStr := node.Text()
		pagesStr = pagesStr[:strings.Index(pagesStr, " ")]
		pages, err := strconv.Atoi(pagesStr)
		if err != nil {
			return nil, newParserParseError("gallery_preview", "pages", pagesStr, err)
		}
		result.Pages = pages
	} else {
		return nil, newNodeNotFoundError("gallery_preview", "pages")
	}

	torrentSelection := column2Selection.Eq(5).Children().First()
	if torrentSelection.Length() != 0 && torrentSelection.Nodes[0].Data == "a" {
		result.HasTorrent = true
	}

	column4Selection := parent.FindMatcher(p.Single(".gl4e"))

	if node := column4Selection.FindMatcher(p.Single(".glink")); node.Length() != 0 {
		title := node.Text()
		result.Title = title
	} else {
		return nil, newNodeNotFoundError("gallery_preview", "title")
	}

	var tags []*TagGroup
	for _, s := range column4Selection.FindMatcher(p.Single("tbody")).Children().EachIter() {
		var tagGroup TagGroup

		columns := s.Children()

		if node := columns.Eq(0); node.Length() != 0 {
			namespace := node.Text()
			namespace = strings.TrimSuffix(namespace, ":")
			tagGroup.Namespace = namespace
		}

		for _, ss := range columns.Eq(1).Children().EachIter() {
			var tagValue TagValue

			class, exist := ss.Attr("class")
			if !exist {
				return nil, newAttrNotFoundError("gallery_preview", "tag", "class")
			}
			tagValue.IsWeak = class == "gtl"

			tagValue.Value = ss.Text()

			tagGroup.Values = append(tagGroup.Values, &tagValue)
		}

		tags = append(tags, &tagGroup)
	}
	result.Tags = tags

	return &result, nil
}
