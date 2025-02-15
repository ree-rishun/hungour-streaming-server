package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode"

	"hungour-streaming-server/models"
	"hungour-streaming-server/repositories"
	"hungour-streaming-server/services"
)

var process models.Process

func ProcessController(w http.ResponseWriter, r *http.Request) {
	var err error
	start := time.Now()
	ctx := context.Background()
	pathParts := strings.Split(r.URL.Path, "/")
	conciergeId := pathParts[2]
	processId := pathParts[3]

	// Twilioからの音声認識結果を取得
	userMessage := r.PostFormValue("SpeechResult")
	log.Println(fmt.Sprintf("[%s/%s] user : %s", conciergeId, processId, userMessage))

	// ローカルに読み込まれていない場合のみFirestoreから取得
	if len(process.Messages) == 0 {
		process, err = repositories.GetProcessDocument(ctx, processId)
		if err != nil {
			log.Fatalf("Gemini Error: %s", err.Error())
		}
	}
	// Geminiで返信文章を考える
	replyText, err := services.GeminiRequest(ctx, userMessage, process.Messages)
	if err != nil {
		log.Fatalf("Gemini Error: %s", err.Error())
	}

	// 返信文章をFirestoreに格納
	process.Messages = append(
		process.Messages,
		models.Message{
			Role: "user",
			Text: userMessage,
		},
	)
	process.Messages = append(
		process.Messages,
		models.Message{
			Role: "model",
			Text: replyText,
		},
	)
	repositories.AddMessage(ctx, processId, process.Messages)

	log.Println(fmt.Sprintf("[%s/%s] model : %s", conciergeId, processId, replyText))

	w.Header().Set("Content-Type", "text/xml")
	w.WriteHeader(http.StatusOK)

	// 通話終了の場合
	if strings.Contains(removeFinishedIfFullWidth(replyText), "finished") {
		_, _ = w.Write([]byte(services.BuildFarewell()))
		return
	}

	log.Println("処理時間：", time.Since(start))
	// 返答
	_, _ = w.Write([]byte(services.BuildReply(replyText, conciergeId, processId)))
}

// 全角文字が含まれているか判定
func containsFullWidth(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Han, r) || unicode.Is(unicode.Hiragana, r) || unicode.Is(unicode.Katakana, r) || (r > 0xFF00 && r < 0xFFEF) {
			return true
		}
	}
	return false
}

// 条件に応じて `finished` を削除する
func removeFinishedIfFullWidth(s string) string {
	if strings.Contains(s, "finished") && containsFullWidth(s) {
		// "finished" を削除
		re := regexp.MustCompile(`finished`)
		return re.ReplaceAllString(s, "")
	}
	return s
}
