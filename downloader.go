package synctube

// Downloader defines the interace that needs to be satisfied in order to be a qualified
// service downloader (e.g. YouTube, SoundCloud).
type Downloader interface {
	Download(url string) error
}

// YouTubeDownloader is the concrete implementation of the Downloader interface
// for YouTube.
type YouTubeDownloader struct {
	rateLimit int
	baseURL   string
}

// Download defines how a YouTube downloader, supplied with a URL for a video,
// will download that particular video.
func (d YouTubeDownloader) Download(url string) error {

	return nil
}
