package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TableItem struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type ItemList []TableItem

func main() {
	// .envファイルの読み込み
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Fail Load env File!: %v", err)
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)

	// MySQLへの接続情報
	dsn := os.Getenv("DB_CONFIG")

	// GORMを使用してMySQLに接続
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("Failed to connect to database")
	}

	// テーブル名を環境変数から取得
	TabaleName := os.Getenv("TABLENAME")
	// テーブルのデータを取得
	db = db.Table(TabaleName)
	// 変数の定義
	var TableCount int64

	// テーブルのデータ数を取得
	db.Count(&TableCount)
	// テーブルのデータ数を表示
	fmt.Printf("TableCount: %d\n", TableCount)

	// テーブルのデータを取得 R
	var TableList ItemList
	db.Find(&TableList)
	// 取得したデータを表示
	for _, user := range TableList {
		fmt.Printf("ID: %s, Name: %s\n", user.ID, user.Name)
	}

	fmt.Printf("--Start Insert--\n")
	// テーブルのデータを追加 C
	insertData := TableItem{
		ID:   "103",
		Name: "test3",
	}
	// トランザクションを開始
	// tx := db.Begin()
	result := db.Create(&insertData)
	// エラーチェック
	if result.Error != nil {
		//トランザクション
		// tx.Rollback()
		panic("Failed to insert data: " + result.Error.Error())
	}
	fmt.Printf("Insert Result: %d\n", result.RowsAffected)
	// トランザクションを確定する
	// tx.Commit()

	// テーブルのデータ数を取得
	db.Count(&TableCount)
	// テーブルのデータ数を表示
	fmt.Printf("TableCount: %d\n", TableCount)
	// テーブルのデータを取得
	db.Find(&TableList)
	// 取得したデータを表示
	for _, user := range TableList {
		fmt.Printf("ID: %s, Name: %s\n", user.ID, user.Name)
	}

	fmt.Printf("--Start Delete--\n")
	// 挿入されたデータのIDと名前を使用してデータを削除 D
	var deleteResult *gorm.DB
	// トランザクションを開始
	// tx = db.Begin()
	deleteResult = db.Delete(&TableItem{}, "id = ? AND name = ?", insertData.ID, insertData.Name)
	// エラーチェック
	if deleteResult.Error != nil {
		//トランザクション
		// tx.Rollback()
		panic("Failed to delete data: " + deleteResult.Error.Error())
	}
	fmt.Printf("Delete Result: %d\n", deleteResult.RowsAffected)
	fmt.Printf("Delete Data ID: %s, Name: %s\n", insertData.ID, insertData.Name)
	// トランザクションを確定する
	// tx.Commit()

	// テーブルのデータ数を取得
	db.Count(&TableCount)
	// テーブルのデータ数を表示
	fmt.Printf("TableCount: %d\n", TableCount)

	// テーブルのデータを取得
	db.Find(&TableList)
	// 取得したデータを表示
	for _, user := range TableList {
		fmt.Printf("ID: %s, Name: %s\n", user.ID, user.Name)
	}

}
