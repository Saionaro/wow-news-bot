package cacher

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
	"wow-news-bot/helpers"
	"wow-news-bot/types"
)

const (
	cacheFilePath       string = "./cache.json"
	cacheSyncPeriodMins int    = 2
)

var (
	sended = make(map[string]bool)
)

func LoadCache() {
	if _, err := os.Stat(cacheFilePath); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Cache not found")
			return
		}
	}
	cache, err := ioutil.ReadFile(cacheFilePath)
	helpers.Check(err)
	var objmap map[string]bool
	parseErr := json.Unmarshal(cache, &objmap)
	if parseErr != nil {
		fmt.Println("Cache file is broken!")
		return
	}
	sended = objmap
}

func syncCache() {
	fmt.Println("Starting cache sync...")
	cacheFile, err := os.OpenFile(cacheFilePath, os.O_RDWR|os.O_CREATE, 0666)
	helpers.Check(err)
	defer cacheFile.Close()
	jsonString, parseErr := json.Marshal(sended)
	helpers.Check(parseErr)
	n2, writeErr := cacheFile.Write(jsonString)
	helpers.Check(writeErr)
	fmt.Printf("Wrote %d bytes\n", n2)
}

func CheckExistence(hash string) bool {
	return sended[hash] || false
}

func MarkSended(item *types.NewsItem) {
	sended[item.Hash] = true
}

func CalcHash(item *types.NewsItem) string {
	hasher := md5.New()
	hasher.Reset()
	io.WriteString(hasher, item.Href)
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}

func StartSyncDeamon() {
	cacheSyncTicker := time.NewTicker(time.Duration(cacheSyncPeriodMins) * time.Minute)
	defer cacheSyncTicker.Stop()
	for {
		select {
		case <-cacheSyncTicker.C:
			go syncCache()
		}
	}
}