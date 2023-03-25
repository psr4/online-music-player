package media

import (
	"fmt"
	"github.com/dhowden/tag"
	"os"
	"path/filepath"
	"strings"
)

func ScanDir(dir string) ([]map[string]string, error) {
	music := []map[string]string{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		switch filepath.Ext(path) {
		case ".mp3":
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			meta, err := tag.ReadFrom(file)
			if err != nil {
				return nil
			}

			var pic string = "/static/images/history.png"
			if meta.Picture() != nil {
				//pic = "data:image/" + meta.Picture().Ext + ";base64," + base64.StdEncoding.EncodeToString(meta.Picture().Data)
			}

			music = append(music, map[string]string{
				"album_name":  meta.Album(),
				"artist_name": meta.Artist(),
				"song_name":   meta.Title(),
				"source_url":  "/music/" + strings.Replace(path, filepath.Dir(path)+"/", "", -1),
				"url":         "/music/" + strings.Replace(path, filepath.Dir(path)+"/", "", -1),
				"pic":         pic,
			})

		}

		return nil
	})
	if err != nil {
		return music, err
	}
	fmt.Println(music)
	return music, err
}
