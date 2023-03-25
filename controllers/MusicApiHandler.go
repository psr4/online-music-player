package controllers

import (
	"errors"
	"github.com/BeanWei/BWOnlineMusicPlayer/config"
	"github.com/BeanWei/BWOnlineMusicPlayer/service/media"
	"github.com/BeanWei/BWOnlineMusicPlayer/service/spider"
	ms "github.com/BeanWei/MusicSpider"
	"github.com/gin-gonic/gin"
	"github.com/wtolson/go-taglib"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func MusicApiHandler(c *gin.Context) {
	types := c.PostForm("types")
	switch types {
	case "search":
		count := c.PostForm("count")
		source := c.PostForm("source")
		pages := c.PostForm("pages")
		name := c.PostForm("name")
		args := map[string]int{"page": str2intSafe(pages), "limit": str2intSafe(count)}
		r := spider.Search(source, name, args)
		os.WriteFile("test.json", []byte(r["result"]), fs.ModePerm)
		c.JSON(200, gin.H{"data": searchFormat(source, r["result"])})
	case "url":
		id := c.PostForm("id")
		source := c.PostForm("source")
		c.JSON(200, gin.H{"url": ms.Downloadurl(source, id)["url"]})
	case "pic":
		id := c.PostForm("id")
		source := c.PostForm("source")
		c.JSON(200, gin.H{"url": songFormat(source, ms.Song(source, id)["result"])["cover_url"]})
	case "playlist":
		lid := c.PostForm("lid")
		source := c.PostForm("source")
		var r map[string]string
		if lid == "-1" {
			tracks, _ := media.ScanDir(config.MediaPath)
			c.JSON(200, gin.H{"data": map[string]interface{}{
				"source":            "local",
				"name":              "本地列表",
				"coverImgUrl":       "images/history.png",
				"creator_nickname":  "试试",
				"creator_avatarUrl": "images/history.png",
				"brief_desc":        "试试",
				"tags":              "抖店,抖店",
				"creat_time":        "11",
				"play_count":        "0",
				"subscribed_count":  "0",
				"tracks":            tracks,
			}})
		} else {
			r = ms.Playlist(source, lid)
			c.JSON(200, gin.H{"data": playlistFormat(source, r["result"])})
		}
	case "lyric":
		id := c.PostForm("id")
		source := c.PostForm("source")
		c.JSON(200, gin.H{"data": ms.Lyric(source, id)["lyric"]})
	case "userlist":
		offset := c.PostForm("offset")
		limit := c.PostForm("limit")
		uid := c.PostForm("uid")
		source := c.PostForm("source")
		args := map[string]int{"page": str2intSafe(offset), "limit": str2intSafe(limit)}
		r := ms.UserPlaylist(source, uid, args)
		c.JSON(200, gin.H{"data": userPlaylistFormat(source, r["result"])})
	case "download":
		url := c.PostForm("url")
		name := c.PostForm("name")
		artist := c.PostForm("artist")
		album := c.PostForm("album")
		ext := filepath.Ext(url)
		// Create the file to write to
		fpath := config.MediaPath + "/" + name + " - " + artist + ext
		if _, err := os.Stat(fpath); !errors.Is(err, os.ErrNotExist) {
			c.JSON(200, gin.H{"msg": "文件已存在!"})
			return
		}
		file, err := os.Create(fpath)
		if err != nil {
			c.JSON(200, gin.H{"msg": err})
			return
		}
		defer file.Close()

		// Get the data from the URL
		resp, err := http.Get(url)
		if err != nil {
			c.JSON(200, gin.H{"msg": err})
			return
		}
		defer resp.Body.Close()

		// Copy the data to the file
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			c.JSON(200, gin.H{"msg": err})
			return
		}

		item, err := taglib.Read(fpath)
		if err != nil {
			c.JSON(200, gin.H{"msg": err})
			return
		}
		defer item.Close()
		item.SetAlbum(album)
		item.SetArtist(artist)
		item.SetTitle(name)
		err = item.Save()
		if err != nil {
			c.JSON(200, gin.H{"msg": err})
			return
		}
		c.JSON(200, gin.H{"msg": "下载成功"})
	default:
	}
}

func str2intSafe(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	} else {
		return i
	}
}
