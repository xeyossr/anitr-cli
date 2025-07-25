package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"time"

	"github.com/Masterminds/semver/v3"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
)

func fetchApi(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result interface{}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(respBody, &result)
	return result, err
}

func FetchUpdates() (msg string, err error) {
	data, err := fetchApi(githubApi)
	if err != nil {
		return "", err
	}

	latestVerStr := data.(map[string]interface{})["tag_name"].(string)
	currentVer, err := semver.NewVersion(CurrentVersion)
	if err != nil {
		return "", err
	}

	latestVer, err := semver.NewVersion(latestVerStr)
	if err != nil {
		return "", err
	}

	if !currentVer.LessThan(latestVer) {
		return "Zaten en son sürümdesiniz.", nil
	} else {
		return fmt.Sprintf("Yeni sürüm bulundu: %s -> %s", CurrentVersion, latestVerStr), err
	}
}

func Version() {
	fmt.Printf("anitr-cli %s\nLisans: GPL 3.0 (Özgür Yazılım)\nDestek ver: %s\n\nGo sürümü: %s\n", CurrentVersion, repoLink, runtime.Version())
}

func CheckUpdates() {
	msg, err := FetchUpdates()

	if err != nil {
		fmt.Println(ColorRed + "Güncelleme kontrolü sırasında bir hata oluştu!" + ColorReset)
		return
	}

	if msg != "Zaten en son sürümdesiniz." {
		fmt.Println(ColorCyan + msg + ColorReset)
		time.Sleep(2 * time.Second)
	}
}
