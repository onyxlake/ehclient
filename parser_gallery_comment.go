package ehclient

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (p *Parser) parseGalleryComments(root *goquery.Document) ([]*GalleryComment, error) {
	var comments []*GalleryComment
	for _, s := range root.FindMatcher(p.Matcher("#cdiv>.c1")).EachIter() {
		var comment GalleryComment

		if node := s.FindMatcher(p.Single(".c3")); node.Length() != 0 {
			text := node.Text()
			i := strings.Index(text, "by:")
			postedAtStr := text[:i-1]
			userStr := text[i+4:]
			postedAtStr = strings.TrimLeft(postedAtStr, "Posted on ")
			postedAt, err := p.parseCommentTime(postedAtStr)
			if err != nil {
				return nil, newParserParseError("gallery_comment", "posted_at", postedAtStr, err)
			}
			comment.PostedAt = postedAt
			comment.User = strings.TrimSpace(userStr)
		}

		if node := s.FindMatcher(p.Single(".c5>span")); node.Length() != 0 {
			scoreStr := node.Text()
			score, err := strconv.Atoi(scoreStr)
			if err != nil {
				return nil, newParserParseError("gallery_comment", "score", scoreStr, err)
			}
			comment.Score = score
		} else {
			// not .c5 but have .c4 means this is a uploader comment
			if s.FindMatcher(p.Single(".c4")).Length() == 0 {
				return nil, newNodeNotFoundError("gallery_comment", "score")
			}
		}

		if node := s.FindMatcher(p.Single(".c6")); node.Length() != 0 {
			comment.Content = node.Nodes[0]
		}

		if node := s.FindMatcher(p.Single(".c8>strong")); node.Length() != 0 {
			lastEditedAtStr := node.Text()
			lastEditedAt, err := p.parseCommentTime(lastEditedAtStr)
			if err != nil {
				return nil, newParserParseError("gallery_comment", "last_edited_at", lastEditedAtStr, err)
			}
			comment.LastEditedAt = &lastEditedAt
		}

		comments = append(comments, &comment)
	}
	return comments, nil
}
