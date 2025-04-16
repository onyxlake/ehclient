package ehclient

import "testing"

func TestParseRatingFromStyle(t *testing.T) {
	tests := []struct {
		style string
		want  float64
	}{
		{"background-position: -48px -21px;", 1.5},
		{"background-position: -48px -1px;", 2},
		{"background-position: -32px -21px;", 2.5},
		{"background-position: -32px -1px;", 3},
		{"background-position: -16px -21px;", 3.5},
		{"background-position: -16px -1px;", 4},
		{"background-position: -0px -21px;", 4.5},
		{"background-position: -0px -1px;", 5},
	}
	parser := NewParser()
	for _, test := range tests {
		got, err := parser.parseRatingFromStyle(test.style)
		if err != nil {
			t.Errorf("parseRatingFromStyle(%q) = %v, want %v", test.style, err, test.want)
			continue
		}
		if got != test.want {
			t.Errorf("parseRatingFromStyle(%q) = %v, want %v", test.style, got, test.want)
		}
	}
}

func BenchmarkParseRatingFromStyle(b *testing.B) {
	parser := NewParser()
	style := "background-position: -48px -21px;opacity: 1;"
	for i := 0; i < b.N; i++ {
		_, err := parser.parseRatingFromStyle(style)
		if err != nil {
			b.Error(err)
		}
	}
}

func TestParseIdTokenPairFromHref(t *testing.T) {
	tests := []struct {
		href string
		want IdTokenPair
	}{
		{"https://e-hentai.org/g/2905464/fc361982ec/", IdTokenPair{Id: 2905464, Token: "fc361982ec"}},
		{"https://e-hentai.org/g/3272585/93a7646c14/", IdTokenPair{Id: 3272585, Token: "93a7646c14"}},
	}
	parser := NewParser()
	for _, test := range tests {
		got, err := parser.parseIdTokenPairFromHref(test.href)
		if err != nil {
			t.Errorf("parseIdTokenPairFromHref(%q) returned an error: %v", test.href, err)
			continue
		}
		if got.Id != test.want.Id || got.Token != test.want.Token {
			t.Errorf("parseIdTokenPairFromHref(%q) = %v, want %v", test.href, got, test.want)
		}
	}
}

func BenchmarkParseIdTokenPairFromHref(b *testing.B) {
	parser := NewParser()
	href := "https://e-hentai.org/g/2905464/fc361982ec/"
	for i := 0; i < b.N; i++ {
		_, err := parser.parseIdTokenPairFromHref(href)
		if err != nil {
			b.Error(err)
		}
	}
}

func TestParsePageTokenPairFromHref(t *testing.T) {
	tests := []struct {
		href string
		want PageTokenPair
	}{
		{"https://e-hentai.org/s/456e0aca7b/2905464-2", PageTokenPair{Page: 2, Token: "456e0aca7b"}},
	}
	parser := NewParser()
	for _, test := range tests {
		got, err := parser.parsePageTokenPairFromHref(test.href)
		if err != nil {
			t.Errorf("parsePageTokenPairFromHref(%q) returned an error: %v", test.href, err)
			continue
		}
		if got.Page != test.want.Page || got.Token != test.want.Token {
			t.Errorf("parsePageTokenPairFromHref(%q) = %v, want %v", test.href, got, test.want)
		}
	}
}

func BenchmarkParsePageTokenPairFromHref(b *testing.B) {
	parser := NewParser()
	href := "https://e-hentai.org/s/456e0aca7b/2905464-2"
	for i := 0; i < b.N; i++ {
		_, err := parser.parsePageTokenPairFromHref(href)
		if err != nil {
			b.Error(err)
		}
	}
}
