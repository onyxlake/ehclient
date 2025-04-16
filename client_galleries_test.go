package ehclient

import (
	"testing"
)

func TestSearchGalleriesWithLogin(t *testing.T) {
	client := InitTestClientWithLogin()
	t.Run("Compact", func(t *testing.T) {
		cursor, err := client.SearchGalleries(&SearchGalleriesOption{
			PreviewKind: Compact,
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(Json(cursor))
	})
	t.Run("Extended", func(t *testing.T) {
		cursor, err := client.SearchGalleries(&SearchGalleriesOption{
			PreviewKind: Extended,
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(Json(cursor))
	})
}
