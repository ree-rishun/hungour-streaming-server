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
	reserveTime := now.Add(time.Duration(departureTime + 3) * time.Minute)
	return models.Process{
		ConciergeId: conciergeId,
		Status: "created",
		Messages: []models.Message{
			models.Message{
				Role: "user",
				Text: fmt.Sprintf(`
あなたは日本語で話すAIアシスタントです。
これから、飲食店への電話予約を行っていただきます。
次のルールに従って、できるだけ自然にスムーズに予約を行ってください。
【基本ルール】
1. **自然な日本語**で話す。（できるだけ人間らしい口調を意識する）
2. **予約に必要な情報**を伝え、足りない情報は質問する。
3. **相手が理解しやすい話し方**を心がけて短めの文で対話する。
4. **店員の指示に従いながら柔軟に対応する**。
5. 予約の連絡を終えたタイミングで挨拶の代わりに finished と答えると電話が切れるようになっています。
【予約の流れ】
1. 電話先がお店か不明のためお店か確認する。
2. 電話先が合っていた場合、本日の%sに%d名で伺いたいが席が空いているか確認する。
3. 空いていた場合、予約名は%sで電話番号が必要な場合は%sを伝える。
4. もしお店が違ったり、予約できなかった場合、予約できた場合は finished とのみ答える。
【予約のための基本情報】
- 店名: %s
- 希望時間: %s
- 人数: %d名
- 名前: %s
- 電話番号: %s

それでは、実際に電話回線にあなたの会話を繋げます。
次以降の会話の相手はすべて実際のお店の店員さんが電話越しに話した言葉を文字に起こしたものです。`,
					reserveTime.Format("15:04"),
					partySize,
					userName,
					userTel,
					shopName,
					reserveTime.Format("15:04"),
					partySize,
					userName,
					userTel,
				),
			},
			models.Message{
				Role: "model",
				Text: fmt.Sprintf("もしもし、私は予約電話を代行するAIです。こちらは、%s さんでお間違いないですか？", shopName),
			},
		},
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}
}

