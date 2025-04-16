package ehclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/cookiejar"
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/net/publicsuffix"
)

func InitTestClient() *Client {
	client := New(&ClientOptions{
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:137.0) Gecko/20100101 Firefox/137.0",
	})
	return client
}

func InitTestClientWithLogin() *Client {
	client := InitTestClient()
	memberId := os.Getenv("TEST_COOKIE_MEMBER_ID")
	passHash := os.Getenv("TEST_COOKIE_PASS_HASH")
	_, err := client.LoginWithCookie(memberId, passHash)
	if err != nil {
		panic(err)
	}
	return client
}

func TestLoginWithCookie(t *testing.T) {
	client := InitTestClient()
	memberId := os.Getenv("TEST_COOKIE_MEMBER_ID")
	passHash := os.Getenv("TEST_COOKIE_PASS_HASH")
	profile, err := client.LoginWithCookie(memberId, passHash)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(profile)
}

func Json(v any) string {
	bs, _ := json.MarshalIndent(v, "  ", "    ")
	return string(bs)
}

func PrepareTestFile(filename string, url string) io.Reader {
	filepath := fmt.Sprintf("./testdata/%s", filename)
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

	jar, _ := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	client := http.Client{
		Jar: jar,
	}

	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	var buffer bytes.Buffer
	w := io.MultiWriter(&buffer, file)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		panic(err)
	}
	file.Close()
	return &buffer
}
