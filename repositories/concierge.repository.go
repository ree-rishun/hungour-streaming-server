package repositories

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"

	"hungour-streaming-server/models"
	"hungour-streaming-server/services"
)

// コンシェルジュ（予約情報）の取得
func GetConciergeDocument(ctx context.Context, id string) (models.Concierge, error) {
	var res models.Concierge
	client, err := services.BuildApp(ctx)
	if err != nil {
		return res, err
	}
	defer client.Close()

	doc, err := client.Collection("concierges").Doc(id).Get(ctx)
	if err != nil {
		return res, err
	}

	data := doc.Data()
	rawList, ok := data["reserve_list"].([]interface{})
	if !ok {
		return res, fmt.Errorf("reserve_list is not an array")
	}
	var reserveList []models.ReserveList
	for _, item := range rawList {
		obj, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		// 手動で `models.ReserveList` に変換
		reserve := models.ReserveList{
			Id:    obj["id"].(string),
			Name:  obj["name"].(string),
			Tel: obj["tel"].(string),
		}
		reserveList = append(reserveList, reserve)
	}
	res = models.Concierge{
		ReserveList: reserveList,
		UserId: data["user_id"].(string),
		Status: data["status"].(string),
		Cursor:	data["cursor"].(int64),
		CreatedAt: data["created_at"].(time.Time),
		UpdatedAt: data["updated_at"].(time.Time),
	}
	return res, nil
}

// 更新
func UpdateConciergeDocument(ctx context.Context, id string, status string, cursor int64) error {
	client, err := services.BuildApp(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	docRef := client.Collection("concierges").Doc(id)

	timestamp := time.Now()
	_, err = docRef.Update(ctx, []firestore.Update{
		firestore.Update{
			Path: "status",
			Value: status,
		},
		firestore.Update{
			Path: "cursor",
			Value: cursor,
		},
		firestore.Update{
			Path: "updated_at",
			Value: timestamp,
		},
	})
	return err
}
