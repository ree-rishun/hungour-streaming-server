package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"hungour-streaming-server/models"
	"hungour-streaming-server/repositories"
	"hungour-streaming-server/services"
)

func ProcessController(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	pathParts := strings.Split(r.URL.Path, "/")
	conciergeId := pathParts[2]
	processId := pathParts[3]

	// Twilioからの音声認識結果を取得
	userMessage := r.PostFormValue("SpeechResult")
	log.Println(fmt.Sprintf("[%s/%s] user : %s", conciergeId, processId, userMessage))

	// Geminiで返信文章を考える
	messages, err := repositories.GetProcessDocument(ctx, processId)
	if err != nil {
		log.Fatalf("Gemini Error: %s", err.Error())
	}
	replyText, err := services.GeminiRequest(ctx, userMessage, messages)
	if err != nil {
		log.Fatalf("Gemini Error: %s", err.Error())
	}

	// 返信文章をFirestoreに格納
	messages = append(
		messages,
		models.Message{
			Role: "user",
			Text: userMessage,
		},
	)
	messages = append(
		messages,
		models.Message{
			Role: "model",
			Text: replyText,
		},
	)
	repositories.AddMessage(ctx, processId, messages)

	log.Println(fmt.Sprintf("[%s/%s] model : %s", conciergeId, processId, replyText))

	w.Header().Set("Content-Type", "text/xml")
	w.WriteHeader(http.StatusOK)

	// 通話終了の場合
	if strings.Contains(replyText, "finished") {
		_, _ = w.Write([]byte(services.BuildFarewell()))
		return
	}

	// 返答
	_, _ = w.Write([]byte(services.BuildReply(replyText, conciergeId, processId)))
}

