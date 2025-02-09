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

func StartController(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	pathParts := strings.Split(r.URL.Path, "/")
	conciergeId := pathParts[2]

	concierge, err := repositories.GetConciergeDocument(ctx, conciergeId)
	if err != nil {
		log.Fatalf("Gemini Error: %s", err.Error())
	}
	user, err := repositories.GetUserDocument(ctx, concierge.UserId)
	if err != nil {
		log.Fatalf("Gemini Error: %s", err.Error())
	}
	cursor := concierge.Cursor

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
	services.StartCall(conciergeId, processId, "+819092244036")
}
