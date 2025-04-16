package ehclient

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/publicsuffix"
)

type UserProfile struct {
	DisplayName string
	Avatar      string
}

func (c *Client) LoginWithCookie(memberId string, hashPass string) (*UserProfile, error) {
	domain := "." + string(c.opts.Endpoint)
	expires := time.Now().Add(c.opts.CookieExpires)
	cookies := []*http.Cookie{
		{
			Domain:  domain,
			Name:    "ipb_member_id",
			Value:   memberId,
			Path:    "/",
			Expires: expires,
		},
		{
			Domain:  domain,
			Name:    "ipb_pass_hash",
			Value:   hashPass,
			Path:    "/",
			Expires: expires,
		},
	}
	u, _ := url.Parse("https://" + string(c.opts.Endpoint))

	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, err
	}
	c.httpc.Jar = jar

	c.httpc.Jar.SetCookies(u, cookies)

	userId, err := c.getCurrentUserId()
	if err != nil {
		return nil, err
	}

	return c.GetUserProfile(userId)
}

func (c *Client) getCurrentUserId() (int, error) {
	req, err := http.NewRequest("GET", "https://forums.e-hentai.org/", nil)
	if err != nil {
		return 0, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return 0, err
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return 0, err
	}
	if node := doc.Find("#userlinks>.home>b>a"); node.Length() != 0 {
		if href, exist := node.Attr("href"); exist {
			u, err := url.Parse(href)
			if err != nil {
				return 0, newParserParseError("forums", "user_url", href, err)
			}
			userIdStr := u.Query().Get("showuser")
			userId, err := strconv.Atoi(userIdStr)
			if err != nil {
				return 0, newParserParseError("forums", "user_id", userIdStr, err)
			}
			return userId, nil
		}
		return 0, newAttrNotFoundError("forums", "current_user", "href")
	}
	return 0, newNodeNotFoundError("forums", "current_user")
}

func (c *Client) GetUserProfile(userId int) (*UserProfile, error) {
	u := fmt.Sprintf("https://forums.e-hentai.org/index.php?showuser=%d", userId)
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

	var result UserProfile

	if node := doc.Find("#profilename"); node.Length() != 0 {
		result.DisplayName = node.Text()
	} else {
		return nil, newNodeNotFoundError("profile", "display_name")
	}

	if node := doc.Find("#profilename~div>img"); node.Length() != 0 {
		result.Avatar, _ = node.Attr("src")
	}

	return &result, nil
}
