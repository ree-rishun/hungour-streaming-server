package services

import (
	"fmt"
	"log"
	"os"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func StartCall(conciergeId string, processId string, toTel string, shopName string) {
	// Pod名の取得
	podName, _ := os.Hostname()

	// Twilio APIクライアントを作成
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: os.Getenv("TWILIO_ACCOUNT_SID"),
		Password: os.Getenv("TWILIO_AUTH_TOKEN"),
	})

	twiml := fmt.Sprintf(`
		<Response>
			<Say language="ja-JP">%s</Say>
			<Gather input="speech" language="ja-JP" action="%s" speechTimeout="auto" />
		</Response>
	`, fmt.Sprintf("もしもし、私は予約電話を代行するエーアイです。こちらは、%s さんでお間違いないですか？", shopName), buildWebhookUrl(conciergeId, processId))

	log.Println(twiml)

	params := &openapi.CreateCallParams{}
	params.SetTo(toTel)
	params.SetFrom(os.Getenv("TWILIO_TEL_FROM"))
	params.SetTwiml(twiml)
	params.SetStatusCallback(fmt.Sprintf("%s/callback/%s/%s?pod=%s", os.Getenv("API_URL"), conciergeId, processId, podName))
	params.SetStatusCallbackEvent([]string{"completed"})

	// 発信
	resp, err := client.Api.CreateCall(params)
	if err != nil {
		log.Fatalf("通話の作成に失敗しました: %s", err.Error())
	}

	log.Printf("通話を開始しました！Call SID: %s\n", *resp.Sid)
}

func BuildReply(replyText string, conciergeId string, processId string) string {
	return fmt.Sprintf(`
		<Response>
		<Say language="ja-JP">%s</Say>
		<Gather input="speech" language="ja-JP" action="%s" speechTimeout="auto" />
		</Response>
	`, replyText, buildWebhookUrl(conciergeId, processId))
}

func BuildFarewell() string {
	return`
<Response>
    <Say language="ja-JP">それでは失礼いたします。</Say>
    <Hangup/>
</Response>
`
}

func buildWebhookUrl(conciergeId string, processId string) string {
	// Pod名の取得
	podName, _ := os.Hostname()

	return fmt.Sprintf("%s/process/%s/%s?pod=%s", os.Getenv("API_URL"), conciergeId, processId, podName)
}
