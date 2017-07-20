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

// ManageFiles, given a userID and the map of playlist-videos, will create the user's temporary folder,
// create folders for each specific playlist and compress everything
func ManageFiles(userID string, videosMap map[model.Playlist][]model.Video) error {
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

	// then the 2 paths needed for creating the archive will be constructed
	userFolderPath := downloadsFolderPath + userID
	zipTargetPath := downloadsFolderPath + userID + ".zip"

	err = getZip(userFolderPath, zipTargetPath)

	return err
}

// getZip, given a directory to compress and a target path where to put the archive, will walk
// the directory tree recursively and compress everything
// TODO: need to find a workaround not being able to handle errors coming from deferred functions
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
func createPlaylistFolder(userID string, playlist model.Playlist, videos []model.Video) error {
	// first, we compute the full path of the playlist folder
	path := downloadsFolderPath + userID + "/" + playlist.Title

	var err error

	// next, create the folder with the name of the playlist with all the videos to be downloaded
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0700)
		if err != nil {
			return err
		}
	} else {
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
		err = os.MkdirAll(path, 0700)
		if err != nil {
			return err
		}
		return nil
	} else {
		return err
	}
}

// func cleanUp() error {
// 	path
// }
