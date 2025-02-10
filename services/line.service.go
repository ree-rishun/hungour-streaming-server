package services

import (
	"fmt"
	"log"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func SendLineMessage(toId string, imageUrl string, title string, message string, conciergeId string) error {
	// bot作成
	bot, err := linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
	)
	if err != nil {
		return err
	}

	pageUrl := fmt.Sprintf("%s/reserves/%s/status", CLIENT_BASE_URL, conciergeId)

	_, err = bot.PushMessage(toId, linebot.NewTemplateMessage(
		title,
		&linebot.ButtonsTemplate{
			ThumbnailImageURL: imageUrl,
			ImageAspectRatio:  linebot.ImageAspectRatioRectangle,
			ImageSize:         linebot.ImageSizeCover,
			ImageBackgroundColor: "#FFFFFF",
			Title: title,
			Text: message,
			DefaultAction: &linebot.URIAction{
				Label: "詳細確認",
				URI: pageUrl,
			},
			Actions: []linebot.TemplateAction{
				&linebot.URIAction{
					Label: "詳細確認",
					URI: pageUrl,
				},
			},
		},
	)).Do()
	if err != nil {
		return err
	} else {
		log.Println("LINEで予約完了メッセージを送信")
	}
	return nil
}