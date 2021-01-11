package handler

import (
	"api/config"

	_ "image/jpeg"
	"net/http"
	"time"

	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/labstack/echo"
)

func MovieUpcoming(c echo.Context) error {
	// GetMovierUpcoming: https://pkg.go.dev/github.com/cyruzin/golang-tmdb#Client.GetMovieUpcoming

	// config情報を取得
	confData, err := config.LoadConfigForYaml()
	if err != nil {
		return c.JSON(http.StatusCreated, err)
	}

	// tmbAPIインスタンス取得
	tmdbClient, err := tmdb.Init(confData.Settings.ApiKey)
	if err != nil {
		return c.JSON(http.StatusCreated, err)
	}

	// APIのオプションに値を設定
	options := map[string]string{
		"language": "ja",
		"region":   "JP",
	}

	// 今後公開予定の映画情報を取得
	movieUpcoming, err := tmdbClient.GetMovieUpcoming(options)
	if err != nil {
		return c.JSON(http.StatusCreated, err)
	}

	// 日時データ、レスポンス情報の初期化
	today := time.Now().Format("2006-01-02")
	after_six_days := time.Now().Add(24 * 6 * time.Hour).Format("2006-01-02")
	resp := make(map[string]interface{})
	var itemSlice []interface{}

	// 処理実施日から1週間以内に公開予定の映画のみを抽出
	for i := 0; i < len(movieUpcoming.Results); i++ {
		movieItem := movieUpcoming.Results[i]
		if today <= movieItem.ReleaseDate && movieItem.ReleaseDate <= after_six_days {
			itemSlice = append(itemSlice, movieItem)
		} else if after_six_days < movieItem.ReleaseDate {
			continue
		}
	}

	// レスポンスデータ整形
	resp["page"] = movieUpcoming.Page
	resp["dates"] = movieUpcoming.Dates
	resp["total_pages"] = movieUpcoming.TotalPages
	resp["total_results"] = movieUpcoming.TotalResults
	resp["results"] = itemSlice

	return c.JSON(http.StatusCreated, resp)
}
