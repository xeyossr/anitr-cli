package animecix

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/xeyossr/anitr-cli/internal"
)

var configAnimecix = internal.Config{
	BaseUrl:        "https://animecix.tv/",
	AlternativeUrl: "https://mangacix.net/",
	VideoPlayers:   []string{"tau-video.xyz", "sibnet"},
	HttpHeaders:    map[string]string{"Accept": "application/json", "User-Agent": "Mozilla/5.0"},
}

type VideoURL struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

type VideoResponse struct {
	URLs []VideoURL `json:"urls"`
}

func FetchAnimeSearchData(query string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%ssecure/search/%s?type=&limit=20", configAnimecix.BaseUrl, query)
	data, err := internal.GetJson(url, configAnimecix.HttpHeaders)

	if err != nil {
		return nil, err
	}

	m, ok := data.(map[string]interface{})

	if !ok {
		return nil, fmt.Errorf("data is not a map")
	}

	resultsRaw, exists := m["results"]
	if !exists {
		return nil, fmt.Errorf("'results' key not found")
	}

	resultsSlice, ok := resultsRaw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("'results' is not a slice")
	}

	var parsed []map[string]interface{}
	for _, item := range resultsSlice {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("'results'.item is not a map")
		}

		entry := map[string]interface{}{
			"name":           itemMap["name"],
			"id":             itemMap["id"],
			"type":           itemMap["type"],
			"title_type":     itemMap["title_type"],
			"original_title": itemMap["original_title"],
			"poster":         itemMap["poster"],
		}

		parsed = append(parsed, entry)
	}

	return parsed, nil
}

func FetchAnimeSeasonsData(id int) ([]int, error) {
	url := fmt.Sprintf("%ssecure/related-videos?episode=1&season=1&titleId=%d&videoId=637113", configAnimecix.AlternativeUrl, id)
	data, err := internal.GetJson(url, configAnimecix.HttpHeaders)
	if err != nil {
		return nil, err
	}

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data is not a map")
	}

	videosField, ok := dataMap["videos"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("'videos' key not found")
	}

	if len(videosField) == 0 {
		return nil, fmt.Errorf("no videos found")
	}

	video, ok := videosField[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("'videos'[0] is not a map")
	}

	title, ok := video["title"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("'title' key not found")
	}

	seasons, ok := title["seasons"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("'seasons' key not found")
	}

	if len(seasons) == 0 {
		return nil, fmt.Errorf("no seasons found")
	}

	count := len(seasons)
	indices := make([]int, count)
	for i := range indices {
		indices[i] = i
	}

	return indices, nil
}

func FetchAnimeEpisodesData(id int) ([]map[string]interface{}, error) {
	var episodes []map[string]interface{}
	seenEpisodes := make(map[string]bool)
	seasons, err := FetchAnimeSeasonsData(id)

	if err != nil {
		return nil, err
	}

	if len(seasons) == 0 {
		return nil, fmt.Errorf("no seasons available for this anime")
	}

	for _, seasonIndex := range seasons {
		url := fmt.Sprintf("%ssecure/related-videos?episode=1&season=%d&titleId=%d&videoId=637113", configAnimecix.AlternativeUrl, seasonIndex+1, id)
		data, err := internal.GetJson(url, configAnimecix.HttpHeaders)

		if err != nil {
			return nil, err
		}

		dataMap, ok := data.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("data is not a map")
		}

		videosRaw, ok := dataMap["videos"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("'videos' key not found")
		}

		for _, video := range videosRaw {
			video, ok := video.(map[string]interface{})

			if !ok {
				return nil, err
			}

			name, ok := video["name"].(string)
			if !ok {
				return nil, fmt.Errorf("name is not a string")
			}

			if !seenEpisodes[name] {
				episodeUrl, ok := video["url"].(string)

				if !ok {
					return nil, fmt.Errorf("url is not a string")
				}

				seasonNum := video["season_num"]
				episode := map[string]interface{}{"name": name, "url": episodeUrl, "season_num": seasonNum}
				episodes = append(episodes, episode)
				seenEpisodes[name] = true
			}
		}
	}

	return episodes, nil
}

func AnimeWatchApiUrl(Url string) ([]map[string]string, error) {
	watch_url := fmt.Sprintf("%s%s", configAnimecix.BaseUrl, Url)
	resp, err := http.Get(watch_url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	finalUrl := resp.Request.URL.String()
	parsedUrl, err := url.Parse(finalUrl)
	if err != nil {
		return nil, err
	}

	pathParts := strings.Split(parsedUrl.Path, "/")

	if len(pathParts) < 3 {
		return nil, fmt.Errorf("path format unexpected")
	}

	embedID := pathParts[2]

	queryParams := parsedUrl.Query()
	vid := queryParams.Get("vid")

	if len(configAnimecix.VideoPlayers) == 0 {
		return nil, fmt.Errorf("no video players configured")
	}
	apiUrl := fmt.Sprintf("https://%s/api/video/%s?vid=%s", configAnimecix.VideoPlayers[0], embedID, vid)

	response, err := http.Get(apiUrl)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var videoResp VideoResponse
	err = json.Unmarshal(body, &videoResp)
	if err != nil {
		return nil, err
	}

	results := []map[string]string{}
	for _, item := range videoResp.URLs {
		entry := map[string]string{
			"label": item.Label,
			"url":   item.URL,
		}

		results = append(results, entry)
	}

	return results, nil
}

func FetchTRCaption(seasonIndex, episodeIndex, id int) (string, error) {
	url := fmt.Sprintf("%ssecure/related-videos?episode=1&season=%d&titleId=%d&videoId=637113", configAnimecix.AlternativeUrl, seasonIndex+1, id)
	data, err := internal.GetJson(url, configAnimecix.HttpHeaders)
	if err != nil {
		return "", err
	}

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("data is not a map")
	}

	videosSlice, ok := dataMap["videos"].([]interface{})
	if !ok {
		return "", fmt.Errorf("'videos' key not found")
	}

	if episodeIndex >= len(videosSlice) || episodeIndex < 0 {
		return "", fmt.Errorf("episode index out of range")
	}

	video, ok := videosSlice[episodeIndex].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("episode not found")
	}

	captions, ok := video["captions"].([]interface{})
	if !ok {
		return "", fmt.Errorf("'captions' key not found")
	}

	for _, caption := range captions {
		captionMap, captionOk := caption.(map[string]interface{})
		if !captionOk {
			return "", fmt.Errorf("caption not found")
		}

		lang, langOk := captionMap["language"].(string)
		if langOk && lang == "tr" {
			return captionMap["url"].(string), nil
		}
	}

	if len(captions) == 0 {
		return "", fmt.Errorf("no captions found")
	}
	caption0, ok := captions[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("first caption is not a map")
	}
	captionUrl, ok := caption0["url"].(string)
	if !ok {
		return "", fmt.Errorf("caption url is not a string")
	}
	return captionUrl, nil
}

func AnimeMovieWatchApiUrl(id int) (map[string]interface{}, error) {
	Url := fmt.Sprintf("%ssecure/titles/%d?titleId=%d", configAnimecix.BaseUrl, id, id)

	client := &http.Client{}
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", configAnimecix.HttpHeaders["Accept"])
	req.Header.Set("User-Agent", configAnimecix.HttpHeaders["User-Agent"])
	req.Header.Set("x-e-h", "=.a")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, err
	}

	dataMap, ok := result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data is not a map")
	}

	titleMap, ok := dataMap["title"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("'title' key not found")
	}

	videosRaw, ok := titleMap["videos"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("'videos' key not found")
	}

	for _, video := range videosRaw {
		video, ok := video.(map[string]interface{})

		if !ok {
			return nil, fmt.Errorf("'videos' key unexpected format")
		}

		videoUrl, ok := video["url"].(string)

		if !ok {
			return nil, fmt.Errorf("'url' key not found")
		}

		client := &http.Client{}
		req, err := http.NewRequest("GET", videoUrl, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Accept", configAnimecix.HttpHeaders["Accept"])
		req.Header.Set("User-Agent", configAnimecix.HttpHeaders["User-Agent"])
		req.Header.Set("x-e-h", "=.a")

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		resp.Body.Close()

		finalUrl := resp.Request.URL.String()
		parsedUrl, err := url.Parse(finalUrl)
		if err != nil {
			return nil, err
		}

		pathParts := strings.Split(parsedUrl.Path, "/")

		if len(pathParts) < 3 {
			log.Printf("path format unexpected")
			continue
		}

		embedID := pathParts[2]
		queryParams := parsedUrl.Query()
		vid := queryParams.Get("vid")

		if len(configAnimecix.VideoPlayers) == 0 {
			return nil, fmt.Errorf("no video players configured")
		}
		apiUrl := fmt.Sprintf("https://%s/api/video/%s?vid=%s", configAnimecix.VideoPlayers[0], embedID, vid)

		response, err := http.Get(apiUrl)
		if err != nil {
			return nil, err
		}

		defer response.Body.Close()

		respBody, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		var videoResp VideoResponse
		err = json.Unmarshal(respBody, &videoResp)
		if err != nil {
			return nil, err
		}

		result := make(map[string]interface{})
		streams := make([]interface{}, 0)
		for _, item := range videoResp.URLs {
			entry := map[string]interface{}{
				"label": item.Label,
				"url":   item.URL,
			}

			streams = append(streams, entry)
		}

		result["video_streams"] = streams

		captions, ok := video["captions"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("'captions' key not found")
		}

		if len(captions) < 1 {
			result["caption_url"] = nil
			return result, nil
		}

		for _, caption := range captions {
			captionMap, captionOk := caption.(map[string]interface{})
			if !captionOk {
				return nil, fmt.Errorf("caption unexpected format")
			}

			lang, langOk := captionMap["language"].(string)
			if langOk && lang == "tr" {
				result["caption_url"] = captionMap["url"]
			} else {
				if len(captions) == 0 {
					return nil, fmt.Errorf("no captions found")
				}
				firstCaption, firstOk := captions[0].(map[string]interface{})
				if !firstOk {
					return nil, fmt.Errorf("first caption is not a map")
				}
				result["caption_url"] = firstCaption["url"]
			}

			return result, nil
		}
	}

	return nil, fmt.Errorf("AnimeMovieWatchApiUrl() panic")
}
