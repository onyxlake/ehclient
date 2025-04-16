package ehclient

import (
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/andybalholm/cascadia"
)

type Parser struct {
	matcherCache sync.Map
}

func NewParser() *Parser {
	return &Parser{
		matcherCache: sync.Map{},
	}
}

func (p *Parser) loadMatcher(key string) (goquery.Matcher, bool) {
	if matcher, exist := p.matcherCache.Load(key); exist {
		return matcher.(goquery.Matcher), true
	} else {
		return nil, false
	}
}

func (p *Parser) storeMatcher(key string, matcher goquery.Matcher) {
	p.matcherCache.Store(key, matcher)
}

func (p *Parser) Matcher(sel string) goquery.Matcher {
	if matcher, exist := p.loadMatcher(sel); exist {
		return matcher
	}
	matcher := cascadia.MustCompile(sel)
	p.storeMatcher(sel, matcher)
	return matcher
}

func (p *Parser) Single(sel string) goquery.Matcher {
	matcher := p.Matcher(sel)
	return goquery.SingleMatcher(matcher)
}
