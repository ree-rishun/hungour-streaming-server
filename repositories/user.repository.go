package repositories


import (
	"context"

	"hungour-streaming-server/models"
	"hungour-streaming-server/services"
)

// プロセスの取得
func GetUserDocument(ctx context.Context, id string) (models.User, error) {
	var user models.User
	client, err := services.BuildApp(ctx)
	if err != nil {
		return user, err
	}
	defer client.Close()

	doc, err := client.Collection("users").Doc(id).Get(ctx)

	if err != nil {
		return user, err
	}

	data := doc.Data()
	user = models.User{
		ReserveName: data["reserve_name"].(string),
		Tel: data["tel"].(string),
	}

	return user, nil
}

