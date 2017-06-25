package youtube

var (
	recap Info
)

type Info struct {
	ClientID string
}

func ReadConfig() Info {
	return recap
}
