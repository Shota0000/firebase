package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/labstack/echo"
	"google.golang.org/api/option"
)

type User struct {
	Name    string `json:"name"`
	Age     string `json:"age"`
	Address string `json:"address"`
}

type Users *[]User

func main() {
	e := echo.New()
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World")
	// })
	// e.GET("/show", show)
	e.POST("/users", addUser)
	e.Logger.Fatal(e.Start(":1323"))
}

func firebaseInit(ctx context.Context) (*firestore.Client, error) {
	sa := option.WithCredentialsFile("path/to/serviceAccount.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return client, nil
}

func dataAdd(name string, age string, address string) ([]*User, error) {
	ctx := context.Background()
	client, err := firebaseInit(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// データ追加
	_, err = client.Collection("users").Doc(name).Set(ctx, map[string]interface{}{
		"age":     age,
		"address": address,
	})
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}
	//データ読み込み
	allData := client.Collection("users").Documents(ctx)
	// 全ドキュメント取得
	docs, err := allData.GetAll()
	if err != nil {
		log.Fatalf("Failed adding getAll: %v", err)
	}
	// 配列の初期化
	users := make([]*User, 0)
	for _, doc := range docs {
		u := new(User)
		// 構造体にFirestoreのデータをセット
		mapToStruct(doc.Data(), u)
		// ドキュメント名を取得してnameにセット
		u.Name = doc.Ref.ID
		// 配列に構造体をセット
		users = append(users, u)

	}

	defer client.Close()
	return users, err
}

func mapToStruct(m map[string]interface{}, val interface{}) error {
	tmp, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(tmp, val)
	if err != nil {
		return err
	}
	return nil
}

func addUser(c echo.Context) error {
	u := new(User)
	if error := c.Bind(u); error != nil {
		return error
	}
	users, error := dataAdd(u.Name, u.Age, u.Address)
	if error != nil {
		return error
	}
	//ここ修正必要かも
	return c.JSON(http.StatusOK, users)
}
