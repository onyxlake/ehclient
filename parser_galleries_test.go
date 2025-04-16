package ehclient

import (
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestParseGalleries(t *testing.T) {
	parser := NewParser()
	t.Run("Fixed-Compact", func(t *testing.T) {
		file := PrepareTestFile("galleries_fixed_compact.html", "https://e-hentai.org/?f_search=-o:%22ai+generated%22&next=3293142&seek=2025-01-01&inline_set=dm_l")
		doc, err := goquery.NewDocumentFromReader(file)
		if err != nil {
			panic(err)
		}
		result, err := parser.parseGalleries(doc)
		if err != nil {
			t.Error(err)
		}
		t.Log(Json(result))
	})
	t.Run("Fixed-Extended", func(t *testing.T) {
		file := PrepareTestFile("galleries_fixed_extended.html", "https://e-hentai.org/?f_search=-o:%22ai+generated%22&next=3293142&seek=2025-01-01&inline_set=dm_e")
		doc, err := goquery.NewDocumentFromReader(file)
		if err != nil {
			panic(err)
		}
		result, err := parser.parseGalleries(doc)
		if err != nil {
			t.Error(err)
		}
		t.Log(Json(result))
	})
}

func BenchmarkParseGalleries(b *testing.B) {
	b.Run("Fixed-Compact", func(b *testing.B) {
		file := PrepareTestFile("galleries_fixed_compact.html", "https://e-hentai.org/?f_search=-o:%22ai+generated%22&next=3293142&seek=2025-01-01&inline_set=dm_l")
		doc, err := goquery.NewDocumentFromReader(file)
		if err != nil {
			panic(err)
		}
		parser := NewParser()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = parser.parseGalleries(doc)
		}
	})
	b.Run("Fixed-Extended", func(b *testing.B) {
		file := PrepareTestFile("galleries_fixed_extended.html", "https://e-hentai.org/?f_search=-o:%22ai+generated%22&next=3293142&seek=2025-01-01&inline_set=dm_e")
		doc, err := goquery.NewDocumentFromReader(file)
		if err != nil {
			panic(err)
		}
		parser := NewParser()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = parser.parseGalleries(doc)
		}
	})
}
