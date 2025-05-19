package ehclient

import (
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestTorrent(t *testing.T) {
	r := PrepareTestFile("torrent.html", "https://e-hentai.org/gallerytorrents.php?gid=2905464&t=fc361982ec")
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		panic(err)
	}
	parser := NewParser()
	result, err := parser.parseTorrent(doc)
	if err != nil {
		t.Error(err)
	}
	t.Log(Json(result))
}

func BenchmarkTorrent(b *testing.B) {
	r := PrepareTestFile("torrent.html", "https://e-hentai.org/gallerytorrents.php?gid=2905464&t=fc361982ec")
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		panic(err)
	}
	parser := NewParser()
	for i := 0; i < b.N; i++ {
		_, err := parser.parseTorrent(doc)
		if err != nil {
			b.Error(err)
		}
	}
}
