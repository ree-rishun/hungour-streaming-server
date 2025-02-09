package repositories

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"

	"hungour-streaming-server/models"
	"hungour-streaming-server/services"
)

// プロセスの取得
func GetProcessDocument(ctx context.Context, id string) ([]models.Message, error) {
	client, err := services.BuildApp(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	doc, err := client.Collection("processes").Doc(id).Get(ctx)

	if err != nil {
		return nil, err
	}

	var messages []models.Message

	data := doc.Data()
	ms, _ := data["messages"].([]interface{})
	for _, message := range ms {
		chatMap, _ := message.(map[string]interface{})
		role, _ := chatMap["role"].(string)
		text, _ := chatMap["text"].(string)

		messages = append(
			messages,
			models.Message{
				Role: role,
				Text: text,
			},
		)
	}

	return messages, nil
}

// プロセスの新規作成
func CreateNewProcess(
	ctx context.Context,
	id string,
	conciergeId string,
	shopName string,
	departureTime int64,
	partySize int64,
	seatType string,
	userName string,
	userTel string,
) error {
	client, err := services.BuildApp(ctx)
	if err != nil {
		return err
	}

	data := buildFirstData(
		conciergeId,
		shopName,
		departureTime,
		partySize,
		seatType,
		userName,
		userTel,
	)

	_, err = client.Collection("processes").Doc(id).Set(ctx, data)
	return err
}

// プロセスへメッセージ履歴の追加
func AddMessage(ctx context.Context, id string, messages []models.Message) error {
	client, err := services.BuildApp(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	docRef := client.Collection("processes").Doc(id)

	_, err = docRef.Update(ctx, []firestore.Update{
		{
			Path: "messages",
			Value: messages,
		},
	})
	return err
}

func buildFirstData(
	conciergeId string,
	shopName string,
	departureTime int64,
	partySize int64,
	seatType string,
	userName string,
	userTel string,
) models.Process {
	timestamp := time.Now()
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(jst)
	reserveTime := now.Add((departureTime + 3) * time.Minute)
	return models.Process{
		ConciergeId: conciergeId,
		Status: "created",
		Messages: []models.Message{
			models.Message{
				Role: "user",
				Text: fmt.Sprintf(`
あなたは日本語で話すAIアシスタントです。
これから、飲食店への電話予約を行います。
次のルールに従って、できるだけ自然にスムーズに予約を行ってください。
【基本ルール】
1. **自然な日本語**で話す。（できるだけ人間らしい口調を意識する）
2. **予約に必要な情報**を伝え、足りない情報は質問する。
3. **相手が理解しやすい話し方**を心がける。（短めの文）
4. **情報の確認**を忘れずに。（日付、時間、人数、名前、電話番号など）
5. もし予約ができなかった場合は場合は「承知しました。お忙しいところご対応いただき、ありがとうございました。」。
6. **最終的に予約が取れたかを確認し、礼儀正しく会話を終える**。
7. **店員の指示に従いながら柔軟に対応する**。
8. 予約の連絡を終えたタイミングで挨拶の代わりに finished と答える。
【予約のための基本情報】
- 店名: %s
- 希望時間: %s
- 人数: %d名
- 名前: %s
- 電話番号: %s`,
					shopName,
					now.Format("15:04"),
					partySize,
					userName,
					userTel,
				),
			},
			models.Message{
				Role: "model",
				Text: "承知しました。",
			},
		},
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}
}

