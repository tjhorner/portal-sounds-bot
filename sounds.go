package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query"
)

type Sound struct {
	ID        string `json:"id"`
	Views     string `json:"views"`
	Downloads string `json:"downs"`
	ValveID   string `json:"valveid"`
	Who       string `json:"who"`
	Text      string `json:"text"`
	Folder    string `json:"folder"`
}

var cdnLinkMap = map[string]string{
	"www":      "sound",
	"dlc":      "sound_dlc1",
	"dlc2":     "sound_dlc2",
	"p2music":  "music_portal2",
	"p1":       "sound_portal1",
	"p1music":  "music_portal1",
	"tf2":      "sound_tf2",
	"tf2music": "tf2_songs",
}

func (s *Sound) Link(subdomain string) string {
	folder := cdnLinkMap[subdomain]
	return fmt.Sprintf(
		"https://cdn.frustra.org/sounds/%s/%s/%s.mp3?id=%s",
		folder,
		s.Folder,
		strings.ReplaceAll(s.ValveID, ".", "/"),
		s.ID,
	)
}

type SoundFilterOptions struct {
	Quote string `url:"quote"`
	Who   string `url:"who"`
}

type SoundQuery struct {
	SoundFilterOptions
	Page             int    `url:"page"`
	UnknownParameter string `url:"s"`
}

func GetSounds(opts SoundFilterOptions, subdomain string, page int) []Sound {
	q := SoundQuery{
		SoundFilterOptions: opts,
		Page:               page,
		UnknownParameter:   "1",
	}

	v, _ := query.Values(q)
	url := fmt.Sprintf("http://%s.portal2sounds.com/list.php?%s", subdomain, v.Encode())

	req, _ := http.NewRequest("POST", url, nil)
	resp, _ := http.DefaultClient.Do(req)

	var res []Sound

	// Parse the results. The first 3 items in the array are ints.
	// 0: number of pages
	// 1: this page
	// 2: results this page

	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	dec.Decode(&res)

	return res[3:]
}
