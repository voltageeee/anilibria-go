// This module is basically a tiny Golang wrapper for Anilibria V3 API.
// Shot out to guys from anilibria for doing a gigantic work!
// Credits go to: https://github.com/anilibria, especially
// to the authors of https://github.com/anilibria/docs/blob/master/api_v3.md

package anilibria

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Wraps the franchise
type FranchiseWrapper struct {
	Franchise Franchise `json:"franchise"`
	Releases  []Release `json:"releases"`
}

// Wraps the posters
type PostersWrapper struct {
	Small    Poster `json:"small"`
	Medium   Poster `json:"medium"`
	Original Poster `json:"original"`
}

// The poster object
type Poster struct {
	URL           string `json:"url"`
	RawBase64File any    `json:"raw_base_64_file"`
}

// The franchise object
type Franchise struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// The release object
type Release struct {
	ID      int64  `json:"id"`
	Code    string `json:"code"`
	Ordinal int    `json:"ordinal"`
	Names   Names  `json:"names"`
}

// The names object
type Names struct {
	Ru          string `json:"ru"`
	En          string `json:"en"`
	Alternative string `json:"alternative"`
}

// The status object
type Status struct {
	String string `json:"status"`
	Code   int    `json:"code"`
}

// The object that represents the anime's type
type AnimeType struct {
	FullString string `json:"full_string"`
	Code       int64  `json:"code"`
	String     string `json:"string"`
	Episodes   any    `json:"episodes"`
	Length     int    `json:"length"`
}

// The object that represents the anime's team
type AnimeTeam struct {
	Voice      []string `json:"voice"`
	Translator []string `json:"translator"`
	Editing    []string `json:"editing"`
	Decor      []string `json:"decor"`
	Timing     []string `json:"timing"`
}

// The season object
type Season struct {
	String  string `json:"string"`
	Code    int    `json:"code"`
	Year    int32  `json:"year"`
	WeekDay int    `json:"week_day"`
}

// The blocked object
type Blocked struct {
	Blocked bool `json:"blocked"`
	Bakanim bool `json:"bakanim"`
}

// I'm tired of documenting my code.
// You are on your own now.
type Episodes struct {
	First  int    `json:"first"`
	Last   int    `json:"last"`
	String string `json:"string"`
}

type Skips struct {
	Opening any `json:"opening"`
	Ending  any `json:"ending"`
}

type HLS struct {
	FHD string `json:"fhd"`
	HD  string `json:"hd"`
	SD  string `json:"sd"`
}

type Episode struct {
	EpisodeNumber    int    `json:"episode"`
	Name             any    `json:"name"`
	UUID             string `json:"uuid"`
	CreatedTimestamp int64  `json:"created_timestamp"`
	Preview          any    `json:"preview"`
	Skips            Skips  `json:"skips"`
	HLS              HLS    `json:"hls"`
}

type Rutube struct {
}

type Player struct {
	AlternativePlayer any                `json:"alternative_player"`
	Host              string             `json:"host"`
	Episodes          Episodes           `json:"episodes"`
	EpisodeList       map[string]Episode `json:"list"`
	Rutube            Rutube             `json:"rutube"`
}

type Quality struct {
	String     string `json:"string"`
	Type       string `json:"type"`
	Resolution string `json:"resolution"`
	Ecoder     string `json:"encoder"`
	LQAudio    any    `json:"lq_audio"`
}

type Torrent struct {
	TorrentId         int64    `json:"torrent_id"`
	Episodes          Episodes `json:"episodes"`
	Quality           Quality  `json:"quality"`
	Leechers          int16    `json:"leechers"`
	Seeders           int16    `json:"seeders"`
	Downloads         int64    `json:"downloads"`
	TotalSize         int64    `json:"total_size"`
	SizeString        string   `json:"size_string"`
	URL               string   `json:"url"`
	Magnet            string   `json:"magnet"`
	UploadedTimestamp int64    `json:"uploaded_timestamp"`
	Hash              string   `json:"hash"`
	Metadata          any      `json:"metadata"`
	RawBase64File     any      `json:"raw_base_64_file"`
}

type TorrentWrapper struct {
	Episodes Episodes  `json:"episodes"`
	List     []Torrent `json:"list"`
}

type Anime struct {
	Id          int64              `json:"id"`
	Code        string             `json:"code"`
	Names       Names              `json:"names"`
	Franchises  []FranchiseWrapper `json:"franchises"`
	Announce    string             `json:"announce"`
	Status      Status             `json:"status"`
	Posters     PostersWrapper     `json:"posters"`
	Updated     int64              `json:"updated"`
	LastChange  int64              `json:"last_change"`
	Type        AnimeType          `json:"type"`
	Genres      []string           `json:"genres"`
	Team        AnimeTeam          `json:"team"`
	Season      Season             `json:"season"`
	Description string             `json:"description"`
	InFavorites int64              `json:"in_favorites"`
	Blocked     Blocked            `json:"blocked"`
	Player      Player             `json:"player"`
	Torrents    TorrentWrapper     `json:"torrents"`
}

type Pagination struct {
	Pages        int   `json:"pages"`
	CuuretnPage  int   `json:"current_page"`
	ItemsPerPage int   `json:"items_per_page"`
	TotalItems   int32 `json:"total_items"`
}

type Updates struct {
	List       []Anime    `json:"list"`
	Pagination Pagination `json:"pagination"`
}

type SearchResponse struct {
	List []Anime `json:"list"`
}

type Schedule struct {
	Day  int     `json:"day"`
	List []Anime `json:"list"`
}

const RouteURL string = "https://api.anilibria.tv/v3/"

// anilibria.Search returns an array of Anime objects that we get by making a request to:
// https://api.anilibria.tv/v3/title/search?search=keywords
// And filtering the response using "filters" array of strings. May as well be empty.
func Search(keywords string, filters []string) ([]Anime, error) {
	requestLink := fmt.Sprintf(RouteURL+"title/search?search=%s&filter=%s", keywords, strings.Join(filters, ","))
	// We must do this terribleness as the response we get doesn't return a simple array of Anime objects
	// But it does, in fact, return a "list" object that contains... Yeah, an array of Anime objects :/
	// TODO: There must be a better way of doing this.
	var result SearchResponse

	response, err := http.Get(requestLink)
	if err != nil {
		return nil, fmt.Errorf("error occurred while we were trying to send the request: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error occurred while we were trying to read the response: %s", err)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error occurred while we were trying to unmarshal the response: %s", err)
	}

	return result.List, nil
}

// anilibria.Random returns a random Anime object
// And filtering the response using "filters" array of strings. May as well be empty.
func Random(filters []string) (*Anime, error) {
	requestLink := fmt.Sprintf(RouteURL+"title/random?filter=%s", strings.Join(filters, ","))
	var result Anime

	response, err := http.Get(requestLink)
	if err != nil {
		return nil, fmt.Errorf("error occurred while we were trying to send the request: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error occurred while we were trying to read the response: %s", err)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error occurred while we were trying to unmarshal the response: %s", err)
	}

	return &result, nil
}

// anilibria.GetTitle returns an Anime object by making a request to:
// https://api.anilibria.tv/v3/title?id/code=id/code
func GetTitle(searchByID bool, id string, code string, filters []string) (*Anime, error) {
	var result Anime
	var requestLink string

	// TODO: There MUST be a better way of doing this.
	if searchByID {
		requestLink = fmt.Sprintf(RouteURL+"title?id=%s", id)
	} else {
		requestLink = fmt.Sprintf(RouteURL+"title?id=%s", code)
	}

	response, err := http.Get(requestLink)
	if err != nil {
		return nil, fmt.Errorf("error occurred while we were trying to send the request: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error occurred while we were trying to read the response: %s", err)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error occurred while we were trying to unmarshal the response: %s", err)
	}

	return &result, nil
}

// anilibria.GetTitleList returns an array of Anime objects by making a request to:
// https://api.anilibria.tv/v3/title?list?id_list/code_list=filter
// And filtering the response using "filters" array of strings. May as well be empty.
func GetTitleList(searchByID bool, ids []string, codes []string, filters []string) ([]Anime, error) {
	var result []Anime
	var requestLink string

	// TODO: There MUST be a better way of doing this.
	if searchByID {
		requestLink = fmt.Sprintf(RouteURL+"title/list?id_list=%s&filters=%s", strings.Join(ids, ","), strings.Join(filters, ","))
	} else {
		requestLink = fmt.Sprintf(RouteURL+"title/list?code_list=%s&filter=%s", strings.Join(codes, ","), strings.Join(filters, ","))
	}

	response, err := http.Get(requestLink)
	if err != nil {
		return nil, fmt.Errorf("error occurred while we were trying to send the request: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error occurred while we were trying to read the response: %s", err)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error occurred while we were trying to unmarshal the response: %s", err)
	}

	return result, nil
}

// anilibria.GetUpdates returns an array of Anime objects (no more then "limit" value)
// That got any updates since "since" timestamp (in the UNIX format)
func GetUpdates(limit int, filters []string, since time.Time) (Updates, error) {
	var result Updates

	requestLink := fmt.Sprintf(RouteURL+"title/updates?since=%d&filter=%s&limit=%d", since.Unix(), strings.Join(filters, ","), limit)

	response, err := http.Get(requestLink)
	if err != nil {
		return Updates{}, fmt.Errorf("error occurred while we were trying to send the request: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Updates{}, fmt.Errorf("error occurred while we were trying to read the response: %s", err)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return Updates{}, fmt.Errorf("error occurred while we were trying to unmarshal the response: %s", err)
	}

	return result, nil
}

// anilibria.GetChanges returns an array of Anime objects (no more then "limit" value)
// That got any changes since "since" timestamp (in the UNIX format)
// Don't ask how is this different from anilibria.GetUpdates. Even the response structure is the same.
// The response when testing wasn't any different from anilibria.GetUpdates.
func GetChanges(limit int, filters []string, since time.Time) (Updates, error) {
	var result Updates

	requestLink := fmt.Sprintf(RouteURL+"title/changes?since=%d&filter=%s&limit=%d", since.Unix(), strings.Join(filters, ","), limit)

	response, err := http.Get(requestLink)
	if err != nil {
		return Updates{}, fmt.Errorf("error occurred while we were trying to send the request: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Updates{}, fmt.Errorf("error occurred while we were trying to read the response: %s", err)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return Updates{}, fmt.Errorf("error occurred while we were trying to unmarshal the response: %s", err)
	}

	return result, nil
}

// anilibria.GetSchedule returns an array of Schdule objects that contain an array of Anime objects
// And a day of the week when the title will be updated (0 is monday and 6 is sunday)
// We pass days as an array of strings ({"0", "1"}) because it's easier to format it lmao
func GetSchedule(filters []string, days []string) ([]Schedule, error) {
	var result []Schedule

	requestLink := fmt.Sprintf(RouteURL+"title/schedule?filter=%s&days=%s", strings.Join(filters, ","), strings.Join(days, ","))

	response, err := http.Get(requestLink)
	if err != nil {
		return nil, fmt.Errorf("error occurred while we were trying to send the request: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error occurred while we were trying to read the response: %s", err)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error occurred while we were trying to unmarshal the response: %s", err)
	}

	return result, nil
}

// anilibria.GetFranchises will return an array of FranchiseWrapper objects that contain
// an array of Release objects and a Franchise object.
// We get the franchise of the title by passing the title's ID to this function.
func GetFranchises(filters []string, id int64) ([]FranchiseWrapper, error) {
	var result []FranchiseWrapper

	requestLink := fmt.Sprintf(RouteURL+"title/franchises?filter=%s&id=%d", strings.Join(filters, ","), id)

	response, err := http.Get(requestLink)
	if err != nil {
		return nil, fmt.Errorf("error occurred while we were trying to send the request: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error occurred while we were trying to read the response: %s", err)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error occurred while we were trying to unmarshal the response: %s", err)
	}

	return result, nil
}
