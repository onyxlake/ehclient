package ehclient

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (p *Parser) parseTagsFromTable(root *goquery.Selection) ([]*TagGroup, error) {
	var tags []*TagGroup
	for _, s := range root.EachIter() {
		var tagGroup TagGroup

		columnsSelection := s.Children()

		if node := columnsSelection.Eq(0); node.Length() != 0 {
			namespace := node.Text()
			namespace = strings.TrimSuffix(namespace, ":")
			tagGroup.Namespace = namespace
		}

		for _, ss := range columnsSelection.Eq(1).Children().EachIter() {
			var tagValue TagValue

			class, exist := ss.Attr("class")
			if !exist {
				return nil, newAttrNotFoundError("tag_table", "tag", "class")
			}
			tagValue.IsWeak = class == "gtl"

			tagValue.Value = ss.Text()

			tagGroup.Values = append(tagGroup.Values, &tagValue)
		}

		tags = append(tags, &tagGroup)
	}
	return tags, nil
}
