package ehclient

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (p *Parser) parseCompactGalleryPreview(parent *goquery.Selection) (*GalleryPreview, error) {
	var result GalleryPreview

	columnsSelection := parent.Children()

	column1Selection := columnsSelection.Eq(0)

	if selection := column1Selection.Children(); selection.Length() != 0 {
		categoryStr := selection.Text()
		category, err := parseCategoryFromLabel(categoryStr)
		if err != nil {
			return nil, newParserParseError("gallery_preview", "category", categoryStr, err)
		}
		result.Category = category
	} else {
		return nil, newNodeNotFoundError("gallery_preview", "category")
	}

	column2Selection := columnsSelection.Eq(1)

	if node := column2Selection.FindMatcher(p.Single(".glthumb img")); node.Length() != 0 {
		src, exist := node.Attr("src")
		if !exist {
			return nil, newAttrNotFoundError("gallery_preview", "thumb", "src")
		}
		if strings.HasPrefix(src, "data") {
			src, exist := node.Attr("data-src")
			if !exist {
				return nil, newAttrNotFoundError("gallery_preview", "thumb", "data-src")
			}
			result.Thumb = src
		} else {
			result.Thumb = src
		}
	} else {
		return nil, newNodeNotFoundError("gallery_preview", "thumb")
	}

	if node := column2Selection.FindMatcher(p.Single("div>div[id]")).Eq(0); node.Length() != 0 {
		postedAtStr := node.Text()
		postedAt, err := p.parseBaseTime(postedAtStr)
		if err != nil {
			return nil, newParserParseError("gallery_preview", "posted_at", postedAtStr, err)
		}
		result.PostedAt = postedAt
	} else {
		return nil, newNodeNotFoundError("gallery_preview", "posted_at")
	}

	if node := column2Selection.FindMatcher(p.Single("div>.ir")); node.Length() != 0 {
		ratingStyleStr, exist := node.Attr("style")
		if !exist {
			return nil, newAttrNotFoundError("gallery_preview", "rating", "style")
		}
		rating, err := p.parseRatingFromStyle(ratingStyleStr)
		if err != nil {
			return nil, newParserParseError("gallery_preview", "rating", ratingStyleStr, err)
		}
		result.Rating = rating
	} else {
		return nil, newNodeNotFoundError("gallery_preview", "rating")
	}

	if node := column2Selection.FindMatcher(p.Single("div>.gldown>a")); node.Length() != 0 {
		result.HasTorrent = true
	}

	column3Selection := columnsSelection.Eq(2)

	if node := column3Selection.FindMatcher(p.Single("a")); node.Length() != 0 {
		href, exist := node.Attr("href")
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

	if node := column3Selection.FindMatcher(p.Single(".glink")); node.Length() != 0 {
		result.Title = node.Text()
	} else {
		return nil, newNodeNotFoundError("gallery_preview", "title")
	}

	var tags []*TagGroup
	for _, s := range column3Selection.FindMatcher(p.Single("a>.glink~div")).Children().EachIter() {
		title, _ := s.Attr("title")
		i := strings.Index(title, ":")
		namespace := title[:i]
		value := title[i+1:]
		if l := len(tags); l != 0 && tags[l-1].Namespace == namespace {
			tags[l-1].Values = append(tags[l-1].Values, &TagValue{
				Value: value,
			})
		} else {
			tags = append(tags, &TagGroup{
				Namespace: namespace,
				Values: []*TagValue{
					{Value: value},
				},
			})
		}
	}
	result.Tags = tags

	column4ItemsSelection := columnsSelection.Eq(3).FindMatcher(p.Matcher("div"))

	if node := column4ItemsSelection.Eq(0); node.Length() != 0 {
		uploader := node.Text()
		result.Uploader = uploader
	} else {
		return nil, newNodeNotFoundError("gallery_preview", "uploader")
	}

	if node := column4ItemsSelection.Eq(1); node.Length() != 0 {
		pagesStr := node.Text()
		pagesStr = pagesStr[:strings.Index(pagesStr, " ")]
		pages, err := strconv.Atoi(pagesStr)
		if err != nil {
			return nil, newParserParseError("gallery_preview", "pages", pagesStr, err)
		}
		result.Pages = pages
	} else {
		return nil, newNodeNotFoundError("gallery_preview", "uploader")
	}

	return &result, nil
}
