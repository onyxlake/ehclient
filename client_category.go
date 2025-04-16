package ehclient

import "fmt"

type Category int

const (
	Doujinshi Category = 2
	Manga              = 4
	ArtistCG           = 8
	GameCG             = 16
	Western            = 512
	NonH               = 256
	ImageSet           = 32
	Cosplay            = 64
	AsianPorn          = 128
	Misc               = 1
)

func parseCategoryFromLabel(label string) (Category, error) {
	switch label {
	case "Doujinshi":
		return Doujinshi, nil
	case "Manga":
		return Manga, nil
	case "Artist CG":
		return ArtistCG, nil
	case "Game CG":
		return GameCG, nil
	case "Western":
		return Western, nil
	case "Non-H":
		return NonH, nil
	case "Image Set":
		return ImageSet, nil
	case "Cosplay":
		return Cosplay, nil
	case "Asian Porn":
		return AsianPorn, nil
	case "Misc":
		return Misc, nil
	default:
		return 0, fmt.Errorf("invalid category")
	}
}
