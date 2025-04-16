package ehclient

import (
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestParseGallery(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		r := PrepareTestFile("gallery_simple.html", "https://e-hentai.org/g/2905464/fc361982ec/")
		doc, err := goquery.NewDocumentFromReader(r)
		if err != nil {
			panic(err)
		}
		parser := NewParser()
		result, err := parser.parseGallery(doc)
		if err != nil {
			t.Error(err)
		}
		t.Log(Json(result))
	})
	t.Run("Complex", func(t *testing.T) {
		r := PrepareTestFile("gallery_complex.html", "https://e-hentai.org/g/3272585/93a7646c14/")
		doc, err := goquery.NewDocumentFromReader(r)
		if err != nil {
			panic(err)
		}
		parser := NewParser()
		result, err := parser.parseGallery(doc)
		if err != nil {
			t.Error(err)
		}
		t.Log(Json(result))
	})
}

func BenchmarkParseGallery(b *testing.B) {
	b.Run("Simple", func(b *testing.B) {
		r := PrepareTestFile("gallery_simple.html", "https://e-hentai.org/g/2905464/fc361982ec/")
		doc, err := goquery.NewDocumentFromReader(r)
		if err != nil {
			panic(err)
		}
		parser := NewParser()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = parser.parseGallery(doc)
		}
	})
	b.Run("Complex", func(b *testing.B) {
		r := PrepareTestFile("gallery_complex.html", "https://e-hentai.org/g/3272585/93a7646c14/")
		doc, err := goquery.NewDocumentFromReader(r)
		if err != nil {
			panic(err)
		}
		parser := NewParser()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = parser.parseGallery(doc)
		}
	})
}
