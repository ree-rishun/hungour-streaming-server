package controller

import (
	"context"
	"fmt"
	"log"
	"os"

	"hungour-streaming-server/repositories"
	"hungour-streaming-server/services"
)

func StartController() {
	ctx := context.Background()
	conciergeId := os.Getenv("CONCIERGE_ID")

	concierge, err := repositories.GetConciergeDocument(ctx, conciergeId)
	if err != nil {
		log.Fatalf("Gemini Error: %s", err.Error())
	}
	user, err := repositories.GetUserDocument(ctx, concierge.UserId)
	if err != nil {
		log.Fatalf("Gemini Error: %s", err.Error())
	}
	cursor := concierge.Cursor

	// TODO: 承認済みユーザのみ店舗に電話できるように変更
	toTel := user.Tel

	// プロセスを作成
	processId := fmt.Sprintf("%s-%s", conciergeId,concierge.ReserveList[cursor].Id)
	repositories.CreateNewProcess(
		ctx,
		processId,
		conciergeId,
		concierge.ReserveList[cursor].Name,
		concierge.DepartureTime,
		concierge.PartySize,
		concierge.SeatType,
		user.ReserveName,
		user.Tel,
	)
	// 予約の電話を開始
	services.StartCall(conciergeId, processId, toTel, concierge.ReserveList[cursor].Name)
}
