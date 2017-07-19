package file_manager

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nclandrei/synctube/model"
)

const (
	downloadsFolderPath string = "tmp/"
)

func ManageFiles(userID string, videosMap map[string][]model.Video) error {
	// first, we create the user's folder
	err := createUserFolder(userID)
	if err != nil {
		log.Fatalf("Error in creating the user's temporary folder: %v", err.Error())
	}

	// then, we create all playlist folders inside the user's one
	for playlist, videos := range videosMap {
		err = createPlaylistFolder(userID, playlist, videos)
		if err != nil {
			log.Fatalf("Error in creating the playlist with ID %v for user %v: %v", playlist, userID, err.Error())
		}
	}
	return err
}

func getZip(dir, target string) error {
	zipfile, err := os.Create(target)

	if err != nil {
		return err
	}

	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(dir)

	if err != nil {
		log.Fatalf("Error reading directory: %v", err.Error())
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(dir)
	}

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, dir))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

// createPlaylistFolder creates a new folder with the name=playlistName
// which will contain all songs synchronized for that user on that playlist
func createPlaylistFolder(userID string, playlistName string, videos []model.Video) error {
	// first, we compute the full path of the playlist folder
	fullPath := downloadsFolderPath + userID + "/" + playlistName

	// next, create the folder with the name of the playlist with all the videos to be downloaded
	err := exec.Command("bash", "-c", "mkdir", fullPath).Run()
	if err != nil {
		return err
	}

	for _, video := range videos {
		// next, add all mp3s inside the folderName folder
		err = exec.Command("bash", "-c", "mv "+video.Title+".mp3").Run()
	}

	return err
}

// createUserFolder creates a folder named after the user's ID that will hold the zip with synced songs
func createUserFolder(userID string) error {
	path := downloadsFolderPath + userID
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(downloadsFolderPath, 0700)
		return nil
	} else {
		return err
	}
}

//func CleanUp() error {
//	// TODO: add code to remove zip/mp3 files after being downloaded by the user
//}
