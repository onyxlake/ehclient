package ehclient

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	_regShowKey = regexp.MustCompile(`showkey="(\S+)"`)
	_regToken   = regexp.MustCompile(`startkey="(\S+)"`)
)

func (p *Parser) parsePage(doc *goquery.Document) (*Page, error) {
	var result Page

	if node := doc.FindMatcher(p.Matcher("script[type]")).Eq(1); node.Length() != 0 {
		script := node.Nodes[0].FirstChild.Data
		if submatch := _regShowKey.FindStringSubmatch(script); len(submatch) == 2 {
			result.ApiToken = submatch[1]
		} else {
			return nil, newNodeNotFoundError("page", "api_token")
		}
		if submatch := _regToken.FindStringSubmatch(script); len(submatch) == 2 {
			result.Token = submatch[1]
		} else {
			return nil, newNodeNotFoundError("page", "token")
		}
	} else {
		return nil, newNodeNotFoundError("page", "script")
	}

	navSelection := doc.FindMatcher(p.Single(".sn"))

	navItemsSelection := navSelection.Children()

	if node := navItemsSelection.Eq(1); node.Length() != 0 {
		href, exist := node.Attr("href")
		if !exist {
			return nil, newAttrNotFoundError("page", "prev", "href")
		}
		pair, err := p.parsePageTokenPairFromHref(href)
		if err != nil {
			return nil, newParserParseError("page", "prev", href, err)
		}
		result.Prev = pair
	} else {
		return nil, newNodeNotFoundError("page", "prev")
	}

	if node := navItemsSelection.Eq(3); node.Length() != 0 {
		href, exist := node.Attr("href")
		if !exist {
			return nil, newAttrNotFoundError("page", "next", "href")
		}
		pair, err := p.parsePageTokenPairFromHref(href)
		if err != nil {
			return nil, newParserParseError("page", "next", href, err)
		}
		result.Next = pair
	} else {
		return nil, newNodeNotFoundError("page", "next")
	}

	if node := navItemsSelection.Eq(2); node.Length() != 0 {
		pageStr := node.Children().First().Text()
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			return nil, newParserParseError("page", "page", pageStr, err)
		}
		result.Page = page
	} else {
		return nil, newNodeNotFoundError("page", "indicator")
	}

	if node := navSelection.Eq(0).Siblings().Eq(0); node.Length() != 0 {
		filename := node.Text()
		filename = filename[:strings.Index(filename, " ::")]
		result.FileName = filename
	} else {
		return nil, newNodeNotFoundError("page", "file_name")
	}

	if node := doc.FindMatcher(p.Single("#img")); node.Length() != 0 {
		url, exist := node.Attr("src")
		if !exist {
			return nil, newAttrNotFoundError("page", "url", "src")
		}
		result.Url = url
		onerror, exist := node.Attr("onerror")
		if !exist {
			return nil, newAttrNotFoundError("page", "url", "reload_token")
		}
		startIndex := strings.Index(onerror, "'")
		endIndex := strings.LastIndex(onerror, "'")
		result.ReloadToken = onerror[startIndex+1 : endIndex]
	} else {
		return nil, newNodeNotFoundError("page", "url")
	}

	if node := doc.FindMatcher(p.Single("#i6>div>a")).Last(); node.Length() != 0 {
		href, exist := node.Attr("href")
		if !exist {
			return nil, newAttrNotFoundError("page", "raw_url", "href")
		}
		result.RawUrl = href
	} else {
		return nil, newNodeNotFoundError("page", "raw_url")
	}

	return &result, nil
}

var (
	_regPrev = regexp.MustCompile(`<a\s*id="prev"[^>]+load_image\((\d)+, '(\S+)'\)`)
	_regNext = regexp.MustCompile(`<a\s*id="next"[^>]+load_image\((\d)+, '(\S+)'\)`)
	_regImg  = regexp.MustCompile(`<img\s*id="img"\s*src="(\S+)"[^>]+nl\('(\S+)'\)`)
	_regRaw  = regexp.MustCompile(`<a\s*href="(\S+)">Download`)
)

func (p *Parser) parsePageApi(resp *showPageResult) (*Page, error) {
	var result Page
	result.Page = resp.Page
	result.Token = resp.Token

	if submatch := _regPrev.FindStringSubmatch(resp.NavHtml); len(submatch) == 3 {
		page, err := strconv.Atoi(submatch[1])
		if err != nil {
			return nil, newParserParseError("page", "prev", submatch[1], err)
		}
		result.Prev = &PageTokenPair{
			Page:  page,
			Token: submatch[2],
		}
	} else {
		return nil, newNodeNotFoundError("page", "prev")
	}
	if submatch := _regNext.FindStringSubmatch(resp.NavHtml); len(submatch) == 3 {
		page, err := strconv.Atoi(submatch[1])
		if err != nil {
			return nil, newParserParseError("page", "next", submatch[1], err)
		}
		result.Next = &PageTokenPair{
			Page:  page,
			Token: submatch[2],
		}
	} else {
		return nil, newNodeNotFoundError("page", "next")
	}

	if fileName := strings.Split(resp.InfoHtml, " :: "); len(fileName) == 3 {
		result.FileName = strings.TrimLeft(fileName[0], "<div>")
	} else {
		return nil, newNodeNotFoundError("page", "file_name")
	}

	if submatch := _regImg.FindStringSubmatch(resp.ImageHtml); len(submatch) == 3 {
		result.Url = submatch[1]
		result.ReloadToken = submatch[2]
	} else {
		return nil, newNodeNotFoundError("page", "url")
	}

	if submatch := _regRaw.FindStringSubmatch(resp.FooterHtml); len(submatch) == 2 {
		result.RawUrl = submatch[1]
	} else {
		return nil, newNodeNotFoundError("page", "raw_url")
	}

	return &result, nil
}
