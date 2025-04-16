package ehclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestParsePage(t *testing.T) {
	r := PrepareTestFile("page.html", "https://e-hentai.org/s/456e0aca7b/2905464-2")
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		panic(err)
	}
	parser := NewParser()
	result, err := parser.parsePage(doc)
	if err != nil {
		t.Error(err)
	}
	t.Log(Json(result))
}

func BenchmarkParsePage(b *testing.B) {
	r := PrepareTestFile("page.html", "https://e-hentai.org/s/456e0aca7b/2905464-2")
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		panic(err)
	}
	parser := NewParser()
	for i := 0; i < b.N; i++ {
		_, err := parser.parsePage(doc)
		if err != nil {
			b.Error(err)
		}
	}
}

func TestParsePageApi(t *testing.T) {
	r := preparePageApiFile()
	var resp showPageResult
	err := json.NewDecoder(r).Decode(&resp)
	if err != nil {
		t.Error(err)
		return
	}
	parser := NewParser()
	result, err := parser.parsePageApi(&resp)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(Json(result))
}

func BenchmarkParsePageApi(b *testing.B) {
	r := preparePageApiFile()
	var resp showPageResult
	err := json.NewDecoder(r).Decode(&resp)
	if err != nil {
		panic(err)
	}
	parser := NewParser()
	for i := 0; i < b.N; i++ {
		_, err := parser.parsePageApi(&resp)
		if err != nil {
			b.Error(err)
		}
	}
}

func preparePageApiFile() io.Reader {
	filepath := fmt.Sprintf("./testdata/%s", "page_api.json")
	if _, err := os.Stat("./testdata"); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			err := os.Mkdir("./testdata", os.ModePerm)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	if _, err := os.Stat(filepath); err == nil {
		file, err := os.Open(filepath)
		if err != nil {
			panic(err)
		}
		return file
	}
	client := InitTestClient()
	page, err := client.GetPage(2905464, 2, "456e0aca7b", nil)
	if err != nil {
		panic(err)
	}
	u := "https://api.e-hentai.org/api.php"
	payload := map[string]any{
		"method":  "showpage",
		"gid":     2905464,
		"page":    page.Next.Page,
		"imgkey":  page.Next.Token,
		"showkey": page.ApiToken,
	}
	body := &bytes.Buffer{}
	err = json.NewEncoder(body).Encode(payload)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", u, body)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	io.Copy(file, resp.Body)
	file.Seek(0, 0)
	return file
}
