package ehclient

import "testing"

func TestGetTorrent(t *testing.T) {
	client := InitTestClient()
	torrents, err := client.GetTorrent(2905464, "fc361982ec")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(Json(torrents))
}
