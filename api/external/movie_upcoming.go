package main

import (
	"api/config"
	"fmt"
	"net/smtp"
	"os"

	_ "image/jpeg"
	"time"

	tmdb "github.com/cyruzin/golang-tmdb"
)

var confPath = "/go/api/config/movieConfig.yml"
var mailConfPath = "/go/api/config/mail.yml"

func main() {
	// GetMovierUpcoming: https://pkg.go.dev/github.com/cyruzin/golang-tmdb#Client.GetMovieUpcoming

	// config情報を取得
	confData, err := config.LoadConfigForYaml(confPath)
	if err != nil {
		fmt.Println(err)
		// return err
	}

	// tmbAPIインスタンス取得
	tmdbClient, err := tmdb.Init(confData.Settings.ApiKey)
	if err != nil {
		fmt.Println(err)
		// return err
	}

	// APIのオプションに値を設定
	options := map[string]string{
		"language": "ja",
		"region":   "JP",
	}

	// 今後公開予定の映画情報を取得
	movieUpcoming, err := tmdbClient.GetMovieUpcoming(options)
	if err != nil {
		fmt.Println(err)
		// return err
	}

	// 日時データ、レスポンス情報の初期化
	today := time.Now().Format("2006-01-02")
	after_days := time.Now().Add(24 * 30 * time.Hour).Format("2006-01-02")
	resp := make(map[string]interface{})
	var itemSlice movieUpcomingInfo

	// 処理実施日から1週間以内に公開予定の映画のみを抽出
	for i := 0; i < len(movieUpcoming.Results); i++ {
		movieItem := movieUpcoming.Results[i]
		if today <= movieItem.ReleaseDate && movieItem.ReleaseDate <= after_days {
			itemSlice.Results = append(itemSlice.Results, movieItem)
		} else if after_days < movieItem.ReleaseDate {
			continue
		}
	}

	// レスポンスデータ整形
	resp["page"] = movieUpcoming.Page
	resp["dates"] = movieUpcoming.Dates
	resp["total_pages"] = movieUpcoming.TotalPages
	resp["total_results"] = movieUpcoming.TotalResults

	// メール本文フォーマット
	mailText := ""
	for i := 0; i < len(itemSlice.Results); i++ {
		title := itemSlice.Results[i].Title
		releaseDate := itemSlice.Results[i].ReleaseDate
		baseUrl := confData.Settings.BaseUrl
		url := baseUrl + itemSlice.Results[i].PosterPath
		mailText += "\r\n" + "【タイトル】" + "\r\n" + title + "\r\n" + "【公開日】" + "\r\n" + releaseDate + "\r\n" + "【ポスターURL】" + "\r\n" + url + "\r\n"
	}

	// メール送信処理
	mailConfData, err := config.LoadConfigForYaml(mailConfPath)
	if err != nil {
		fmt.Println(err)
		// return err
	}

	from := mailConfData.Gmail.From
	to := mailConfData.Gmail.To
	_smtp := mailConfData.Gmail.Smtp
	port := mailConfData.Gmail.Port
	auth := smtp.PlainAuth("", from, mailConfData.Gmail.Password, _smtp)
	fromName := mailConfData.MovieUpcoming.FromName
	subject := mailConfData.MovieUpcoming.Subject
	msg := []byte("" +
		"From: " + fromName + "<" + from + ">\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		mailText +
		"\r\n" +
		"")
	address := _smtp + ":" + port
	sendMail := smtp.SendMail(address, auth, from, []string{to}, msg)
	if sendMail != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", sendMail)
		return
	}
}

type movieUpcomingInfo struct {
	Results []struct {
		PosterPath  string `json:"poster_path"`
		Adult       bool   `json:"adult"`
		Overview    string `json:"overview"`
		ReleaseDate string `json:"release_date"`
		Genres      []struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"genres"`
		ID               int64   `json:"id"`
		OriginalTitle    string  `json:"original_title"`
		OriginalLanguage string  `json:"original_language"`
		Title            string  `json:"title"`
		BackdropPath     string  `json:"backdrop_path"`
		Popularity       float32 `json:"popularity"`
		VoteCount        int64   `json:"vote_count"`
		Video            bool    `json:"video"`
		VoteAverage      float32 `json:"vote_average"`
	} `json:"results"`
}
