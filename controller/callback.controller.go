package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"hungour-streaming-server/repositories"
	"hungour-streaming-server/services"
)

func CallbackController(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	pathParts := strings.Split(r.URL.Path, "/")
	conciergeId := pathParts[2]
	processId := pathParts[3]

	// CallbackStatusの確認
	callStatus := r.FormValue("CallStatus")
	// callSid := r.FormValue("CallSid")

	if callStatus != "completed" {
		log.Println(fmt.Sprintf("[%s/%s] callStatus : %s", conciergeId, processId, callStatus))
	}

	// 会話履歴から予約完了したか確認
	process, err := repositories.GetProcessDocument(ctx, processId)
	if err != nil {
		log.Fatalf("Gemini Error: %s", err.Error())
	}
	messages := process.Messages
	isReserved, err := services.IsReserved(ctx, messages)
	if err != nil {
		log.Fatalf("Gemini Error: %s", err.Error())
	}

	concierge, err := repositories.GetConciergeDocument(ctx, conciergeId)
	if err != nil {
		log.Fatalf("Gemini Error: %s", err.Error())
	}

	user, err := repositories.GetUserDocument(ctx, concierge.UserId)
	if err != nil {
		log.Fatalf("Gemini Error: %s", err.Error())
	}

	// 予約完了
	if isReserved {
		// 予約完了ステータスに
		repositories.UpdateConciergeDocument(ctx, conciergeId, "reserved", concierge.Cursor, process.ReservedTime)

		// LINE送信処理
		services.SendLineMessage(
			user.LineId,
			"",
			"予約が完了しました",
			fmt.Sprintf("「%s」を%sに予約しました。詳細は投稿をご覧ください。", concierge.ReserveList[cursor].Name, process.ReservedTime.Format("15:04")),
			conciergeId,
		)

		// Podの削除処理
		services.DeletePod()

		return
	}

	// 全ての店が予約できなかった
	if concierge.Cursor >= 2 {
		// 予約完了ステータスに
		repositories.UpdateConciergeDocument(ctx, conciergeId, "failed", concierge.Cursor, process.ReservedTime)

		// LINE送信処理
		services.SendLineMessage(user.LineId, "", "予約できませんでした", "申し訳ございませんがお店に問い合わせたところ予約できませんでした。", conciergeId)

		// Podの削除処理
		services.DeletePod()

		return
	}

	// 次の予約を開始
	cursor := concierge.Cursor + 1

	// TODO: 承認済みユーザのみ店舗に電話できるように変更
	toTel := user.Tel

	repositories.CreateNewProcess(
		ctx,
		concierge.ReserveList[cursor].Id,
		conciergeId,
		concierge.ReserveList[cursor].Name,
		concierge.DepartureTime,
		concierge.PartySize,
		concierge.SeatType,
		user.ReserveName,
		user.Tel,
	)
	services.StartCall(conciergeId, processId, toTel, concierge.ReserveList[cursor].Name)
	repositories.UpdateConciergeDocument(ctx, conciergeId, concierge.Status, cursor, process.ReservedTime)
}

