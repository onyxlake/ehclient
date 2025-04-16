package ehclient

import (
	"fmt"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	baseTimeLayout    = "2006-01-02 15:04"
	commentTimeLayout = "02 January 2006, 15:04"
)

func (p *Parser) parseBaseTime(str string) (time.Time, error) {
	return time.Parse(baseTimeLayout, str)
}

func (p *Parser) parseCommentTime(str string) (time.Time, error) {
	return time.Parse(commentTimeLayout, str)
}

var (
	regPx = regexp.MustCompile(`(-*\d+)px`)
)

func (p *Parser) parseRatingFromStyle(style string) (float64, error) {
	matches := regPx.FindAllStringSubmatch(style, 2)
	if len(matches) != 2 {
		return 0, fmt.Errorf("invalid rating style")
	}
	hOffset, err := strconv.ParseFloat(matches[0][1], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid h-offset (value: %s)", matches[0][1])
	}
	vOffset, err := strconv.ParseFloat(matches[1][1], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid v-offset (value: %s)", matches[1][1])
	}
	var result float64
	result = 5 - math.Abs(hOffset)/16
	if vOffset == -21 {
		result -= 0.5
	}
	return result, nil
}

type IdTokenPair struct {
	Id    int
	Token string
}

func (p *Parser) parseIdTokenPairFromHref(href string) (*IdTokenPair, error) {
	u, err := url.Parse(href)
	if err != nil {
		return nil, fmt.Errorf("invalid url")
	}
	segments := strings.Split((u.Path[1:]), "/")
	if len(segments) < 3 || segments[0] != "g" {
		return nil, fmt.Errorf("invalid href format")
	}
	id, err := strconv.Atoi(segments[1])
	if err != nil {
		return nil, fmt.Errorf("invalid id (value: %s)", segments[1])
	}
	result := IdTokenPair{
		Id:    id,
		Token: segments[2],
	}
	return &result, nil
}

var (
	regUrl = regexp.MustCompile(`url\((\S+)\)`)
)

func (p *Parser) parseUrlFromStyle(style string) (string, error) {
	submatch := regUrl.FindStringSubmatch(style)
	if submatch == nil || len(submatch) != 2 {
		return "", fmt.Errorf("invalid style")
	}
	return submatch[1], nil
}

type PageTokenPair struct {
	Page  int
	Token string
}

func (p *Parser) parsePageTokenPairFromHref(href string) (*PageTokenPair, error) {
	u, err := url.Parse(href)
	if err != nil {
		return nil, fmt.Errorf("invalid url")
	}
	segments := strings.Split((u.Path[1:]), "/")
	if len(segments) < 3 || segments[0] != "s" {
		return nil, fmt.Errorf("invalid href format")
	}
	token := segments[1]
	pageStr := segments[2]
	pageStr = pageStr[strings.Index(pageStr, "-")+1:]
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return nil, fmt.Errorf("invalid page (value: %s)", pageStr)
	}
	result := PageTokenPair{
		Page:  page,
		Token: token,
	}
	return &result, nil
}
