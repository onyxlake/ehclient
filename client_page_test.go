package ehclient

import "testing"

func TestGetPage(t *testing.T) {
	client := InitTestClient()
	t.Run("By Web", func(t *testing.T) {
		result, err := client.GetPage(2905464, 2, "456e0aca7b", nil)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(Json(result))
		t.Run("By Api", func(t *testing.T) {
			next := result.Next
			result, err := client.GetPage(2905464, next.Page, next.Token, &GetPageOptions{
				ApiToken: result.ApiToken,
			})
			if err != nil {
				t.Error(err)
				return
			}
			t.Log(Json(result))
		})
	})
}
