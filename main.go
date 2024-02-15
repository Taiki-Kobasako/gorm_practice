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
	"gorm.io/gorm/schema"
)

type advertiser struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type ItemList []advertiser

func dbTruncate(db *gorm.DB) error {
	tx := db.Begin()
	txDB := tx.Session(&gorm.Session{NewDB: true})
	txDB.Exec("SET FOREIGN_KEY_CHECKS=0")
	defer txDB.Exec("SET FOREIGN_KEY_CHECKS=1")
	truncateDB := txDB.Exec("TRUNCATE TABLE advertiser")
	if truncateDB.Error != nil {
		tx.Rollback()
		return fmt.Errorf("failed to truncate table: %w", truncateDB.Error)
	}
	tx.Commit()
	return nil
}

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
			Colorful:                  true,        // able to colorize log output
		},
	)

	// MySQLへの接続情報
	dsn := os.Getenv("DB_CONFIG")

	// GORMを使用してMySQLに接続
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",   // テーブル名のプレフィックス
			SingularTable: true, // テーブル名を複数形にしない
		},
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

	fmt.Printf("\n--Start Read--\n")
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

	fmt.Printf("\n--Start Insert--\n")
	// テーブルのデータを追加 C
	insertData := advertiser{
		ID:   "1",
		Name: "test",
	}
	// トランザクションを開始
	tx := db.Begin()
	result := db.Create(&insertData)
	// エラーチェック
	if result.Error != nil {
		//トランザクション
		tx.Rollback()
		panic("Failed to insert data: " + result.Error.Error())
	}
	fmt.Printf("Insert Result: %d\n", result.RowsAffected)
	// トランザクションを確定する
	tx.Commit()

	// テーブルのデータ数を取得
	db.Count(&TableCount)
	// テーブルのデータ数を表示
	fmt.Printf("TableCount: %d\n", TableCount)
	var TableList2 ItemList
	// テーブルのデータを取得
	db.Find(&TableList2)
	// 取得したデータを表示
	for _, user := range TableList2 {
		fmt.Printf("ID: %s, Name: %s\n", user.ID, user.Name)
	}

	fmt.Printf("\n--Start Delete--\n")
	// 挿入されたデータのIDと名前を使用してデータを削除 D
	var deleteResult *gorm.DB
	// トランザクションを開始
	tx = db.Begin()

	// dbオブジェクトを複製して新しいトランザクション用のdbオブジェクトを作成
	txDB := tx.Session(&gorm.Session{NewDB: true})
	// deleteResult = db.Where("name =?", "test3").Delete(&TableList)
	deleteResult = txDB.Delete(&advertiser{}, "id = ? AND name = ?", insertData.ID, insertData.Name)
	// エラーチェック
	if deleteResult.Error != nil {
		//トランザクション
		tx.Rollback()
		panic("Failed to delete data: " + deleteResult.Error.Error())
	}
	// トランザクションを確定する
	tx.Commit()
	fmt.Printf("Delete Result: %d\n", deleteResult.RowsAffected)
	fmt.Printf("Delete Data ID: %s, Name: %s\n", insertData.ID, insertData.Name)

	// テーブルのデータ数を取得
	db.Count(&TableCount)
	// テーブルのデータ数を表示
	fmt.Printf("TableCount: %d\n", TableCount)
	var TableList3 ItemList
	// テーブルのデータを取得
	db.Find(&TableList3)
	// 取得したデータを表示
	for _, user := range TableList3 {
		fmt.Printf("ID: %s, Name: %s\n", user.ID, user.Name)
	}

	fmt.Printf("\n--Start Truncate--\n")
	// テーブルのデータを全て削除
	err = dbTruncate(db)
	if err != nil {
		panic("Failed to truncate table: " + err.Error())
	}
	fmt.Printf("Truncate Table\n")

	// テーブルのデータ数を取得
	db.Count(&TableCount)
	// テーブルのデータ数を表示
	fmt.Printf("TableCount: %d\n", TableCount)
	var TableList4 ItemList
	// テーブルのデータを取得
	db.Find(&TableList4)
	// 取得したデータを表示
	for _, user := range TableList4 {
		fmt.Printf("ID: %s, Name: %s\n", user.ID, user.Name)
	}

}
