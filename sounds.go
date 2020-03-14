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

func (s *Sound) Link() string {
	return fmt.Sprintf(
		"https://cdn.frustra.org/sounds/sound/%s/%s.mp3?id=%s",
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

func GetSounds(opts SoundFilterOptions, page int) []Sound {
	q := SoundQuery{
		SoundFilterOptions: opts,
		Page:               page,
		UnknownParameter:   "1",
	}

	v, _ := query.Values(q)
	url := fmt.Sprintf("http://www.portal2sounds.com/list.php?%s", v.Encode())

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
