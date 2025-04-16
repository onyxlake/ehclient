package ehclient

import "testing"

func TestGetGallery(t *testing.T) {
	client := InitTestClient()
	t.Run("Full Attribute", func(t *testing.T) {
		gallery, err := client.GetGallery(3274096, "82691ad863", nil)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(Json(gallery))
	})
}
