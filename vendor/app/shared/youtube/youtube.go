package youtube

var (
	recap YouTube
)

type YouTube struct {
	ClientID string
}

func ReadConfig() YouTube {
	return recap
}
