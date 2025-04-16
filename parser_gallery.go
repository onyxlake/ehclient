package ehclient

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (p *Parser) parseGallery(doc *goquery.Document) (*GalleryDetail, error) {
	var result GalleryDetail

	if node := doc.FindMatcher(p.Single(".g3>a")); node.Length() != 0 {
		href, exist := node.Attr("href")
		if !exist {
			return nil, newAttrNotFoundError("gallery_detail", "link", "href")
		}
		pair, err := p.parseIdTokenPairFromHref(href)
		if err != nil {
			return nil, newParserParseError("gallery_detail", "link", href, err)
		}
		result.Id = pair.Id
		result.Token = pair.Token
	} else {
		return nil, newNodeNotFoundError("gallery_detail", "link")
	}

	if node := doc.FindMatcher(p.Single("#gn")); node.Length() != 0 {
		result.Title = node.Text()
	} else {
		return nil, newNodeNotFoundError("gallery_detail", "title")
	}

	if node := doc.FindMatcher(p.Single("#gj")); node.Length() != 0 {
		result.Title = node.Text()
	} else {
		return nil, newNodeNotFoundError("gallery_detail", "title_jpn")
	}

	if node := doc.FindMatcher(p.Single("#gd1>div")); node.Length() != 0 {
		style, exist := node.Attr("style")
		if !exist {
			return nil, newAttrNotFoundError("gallery_preview", "thumb", "style")
		}
		thumb, err := p.parseUrlFromStyle(style)
		if err != nil {
			return nil, newParserParseError("gallery_preview", "thumb", style, err)
		}
		result.Thumb = thumb
	} else {
		return nil, newNodeNotFoundError("gallery_detail", "thumb")
	}

	infoSelection := doc.FindMatcher(p.Single("#gd3"))

	if node := infoSelection.FindMatcher(p.Single("#gdc>div")); node.Length() != 0 {
		categoryStr := node.Text()
		category, err := parseCategoryFromLabel(categoryStr)
		if err != nil {
			return nil, newParserParseError("gallery_detail", "category", categoryStr, err)
		}
		result.Category = category
	} else {
		return nil, newNodeNotFoundError("gallery_detail", "category")
	}

	if node := infoSelection.FindMatcher(p.Single("#gdn>a")).Eq(0); node.Length() != 0 {
		result.Uploader = node.Text()
	} else {
		return nil, newNodeNotFoundError("gallery_detail", "uploader")
	}

	infoRowsSelection := infoSelection.FindMatcher(p.Matcher("#gdd tr"))

	if node := infoRowsSelection.Eq(0).Children().Eq(1); node.Length() != 0 {
		postedAtStr := node.Text()
		postedAt, err := p.parseBaseTime(postedAtStr)
		if err != nil {
			return nil, newParserParseError("gallery_detail", "posted_at", postedAtStr, err)
		}
		result.PostedAt = postedAt
	} else {
		return nil, newNodeNotFoundError("gallery_detail", "posted_at")
	}

	if node := infoRowsSelection.Eq(1).Children().Eq(1).Find("a"); node.Length() != 0 {
		href, exist := node.Attr("href")
		if !exist {
			return nil, newAttrNotFoundError("gallery_detail", "parent", "href")
		}
		pair, err := p.parseIdTokenPairFromHref(href)
		if err != nil {
			return nil, newParserParseError("gallery_detail", "parent", href, err)
		}
		result.Parent = pair
	}

	if node := infoRowsSelection.Eq(2).Children().Eq(1); node.Length() != 0 {
		result.Visible = node.Text()
	} else {
		return nil, newNodeNotFoundError("gallery_detail", "visible")
	}

	if node := infoRowsSelection.Eq(3).Children().Eq(1); node.Length() != 0 {
		result.Language = node.Text()
	} else {
		return nil, newNodeNotFoundError("gallery_detail", "language")
	}

	if node := infoRowsSelection.Eq(4).Children().Eq(1); node.Length() != 0 {
		result.FileSize = node.Text()
	} else {
		return nil, newNodeNotFoundError("gallery_detail", "file_size")
	}

	if node := infoRowsSelection.Eq(5).Children().Eq(1); node.Length() != 0 {
		pagesStr := node.Text()
		pagesStr = pagesStr[:strings.Index(pagesStr, " ")]
		pages, err := strconv.Atoi(pagesStr)
		if err != nil {
			return nil, newParserParseError("gallery_detail", "pages", pagesStr, err)
		}
		result.Pages = pages
	} else {
		return nil, newNodeNotFoundError("gallery_detail", "pages")
	}

	if node := infoRowsSelection.Eq(6).Children().Eq(1); node.Length() != 0 {
		favoritedsStr := node.Text()
		favoritedsStr = favoritedsStr[:strings.Index(favoritedsStr, " ")]
		favoriteds, err := strconv.Atoi(favoritedsStr)
		if err != nil {
			return nil, newParserParseError("gallery_detail", "favoriteds", favoritedsStr, err)
		}
		result.Favoriteds = favoriteds
	} else {
		return nil, newNodeNotFoundError("gallery_detail", "favoriteds")
	}

	ratingGroupSelection := infoSelection.FindMatcher(p.Single("#gdr"))

	if node := ratingGroupSelection.FindMatcher(p.Single("#rating_count")); node.Length() != 0 {
		ratedsStr := node.Text()
		rateds, err := strconv.Atoi(ratedsStr)
		if err != nil {
			return nil, newParserParseError("gallery_detail", "rateds", ratedsStr, err)
		}
		result.Rateds = rateds
	} else {
		return nil, newNodeNotFoundError("gallery_detail", "rateds")
	}

	if node := ratingGroupSelection.FindMatcher(p.Single("#rating_label")); node.Length() != 0 {
		ratingStr := node.Text()
		ratingStr = ratingStr[strings.Index(ratingStr, " ")+1:]
		rating, err := strconv.ParseFloat(ratingStr, 64)
		if err != nil {
			return nil, newParserParseError("gallery_detail", "rating", ratingStr, err)
		}
		result.Rating = rating
	} else {
		return nil, newNodeNotFoundError("gallery_detail", "rating")
	}

	tags, err := p.parseTagsFromTable(doc.FindMatcher(p.Matcher("#taglist tr")))
	if err != nil {
		return nil, err
	}
	result.Tags = tags

	var versions []*GalleryVersion
	for _, s := range doc.FindMatcher(p.Matcher("#gnd a")).EachIter() {
		var version GalleryVersion

		href, exist := s.Attr("href")
		if !exist {
			return nil, newAttrNotFoundError("gallery_version", "link", "href")
		}
		pair, err := p.parseIdTokenPairFromHref(href)
		if err != nil {
			return nil, newParserParseError("gallery_version", "id_token_pair", href, err)
		}
		version.Id = pair.Id
		version.Token = pair.Token

		version.Title = s.Text()

		addedAtStr := s.Nodes[0].NextSibling.Data
		addedAtStr = strings.TrimPrefix(addedAtStr, ", added ")
		addedAt, err := p.parseBaseTime(addedAtStr)
		if err != nil {
			return nil, newParserParseError("gallery_version", "added_at", addedAtStr, err)
		}
		version.AddedAt = addedAt

		versions = append(versions, &version)
	}
	result.Versions = versions

	var previews []*PagePreview
	for _, s := range doc.FindMatcher(p.Matcher("#gdt>a")).EachIter() {
		var preview PagePreview

		href, exist := s.Attr("href")
		if !exist {
			return nil, newAttrNotFoundError("page_preview", "link", "href")
		}
		pair, err := p.parsePageTokenPairFromHref(href)
		if err != nil {
			return nil, newParserParseError("page_preview", "page_token_pair", href, err)
		}
		preview.Page = pair.Page
		preview.Token = pair.Token

		child := s.Children()

		title, exist := child.Attr("title")
		if !exist {
			return nil, newAttrNotFoundError("page_preview", "thumb", "title")
		}
		title = title[strings.Index(title, ":")+2:]
		preview.FileName = title

		style, exist := child.Attr("style")
		if !exist {
			return nil, newAttrNotFoundError("page_preview", "thumb", "style")
		}
		submatchs := regPx.FindAllStringSubmatch(style, 3)
		if len(submatchs) != 3 {
			return nil, newParserParseError("page_preview", "infos", style, fmt.Errorf("invalid style"))
		}

		widthStr := submatchs[0][1]
		width, err := strconv.Atoi(widthStr)
		if err != nil {
			return nil, newParserParseError("page_preview", "width", widthStr, err)
		}
		preview.Width = width

		heightStr := submatchs[1][1]
		height, err := strconv.Atoi(heightStr)
		if err != nil {
			return nil, newParserParseError("page_preview", "height", heightStr, err)
		}
		preview.Height = height

		offsetStr := submatchs[2][1]
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return nil, newParserParseError("page_preview", "offset", offsetStr, err)
		}
		if offset < 0 {
			offset = -offset
		}
		preview.Offset = offset

		thumb, err := p.parseUrlFromStyle(style)
		if err != nil {
			return nil, newParserParseError("page_preview", "thumb", style, err)
		}
		preview.Thumb = thumb

		previews = append(previews, &preview)
	}
	result.Previews = previews

	comments, err := p.parseGalleryComments(doc)
	if err != nil {
		return nil, err
	}
	result.Comments = comments

	return &result, nil
}
