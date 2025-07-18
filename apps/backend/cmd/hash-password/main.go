package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("使用方法: go run main.go <password>")
		fmt.Println("例: go run main.go password123")
		os.Exit(1)
	}

	password := os.Args[1]
	
	// bcryptでハッシュ化（コストファクター: 10）
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("パスワードのハッシュ化に失敗しました:", err)
	}

	fmt.Println("パスワード:", password)
	fmt.Println("ハッシュ値:", string(hash))
	fmt.Println("")
	fmt.Println("このハッシュ値をseed-dev-data.shスクリプトで使用してください。")
}