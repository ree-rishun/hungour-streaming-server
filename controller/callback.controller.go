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
	messages, err := repositories.GetProcessDocument(ctx, processId)
	if err != nil {
		log.Fatalf("Gemini Error: %s", err.Error())
	}
	isReserved, err := services.IsReserved(ctx, messages)
	if err != nil {
		log.Fatalf("Gemini Error: %s", err.Error())
	}

	concierge, err := repositories.GetConciergeDocument(ctx, conciergeId)
	if err != nil {
		log.Fatalf("Gemini Error: %s", err.Error())
	}

	// 終了
	if isReserved {
		// 予約完了ステータスに
		repositories.UpdateConciergeDocument(ctx, conciergeId, "reserved", concierge.Cursor)

		// TODO: Podの削除処理

		return
	}

	// 次の予約を開始
	user, err := repositories.GetUserDocument(ctx, concierge.UserId)
	if err != nil {
		log.Fatalf("Gemini Error: %s", err.Error())
	}
	cursor := concierge.Cursor + 1

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
	services.StartCall(conciergeId, processId, "+819092244036")
	repositories.UpdateConciergeDocument(ctx, conciergeId, concierge.Status, cursor)
}

