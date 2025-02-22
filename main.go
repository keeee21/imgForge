package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
	"github.com/rs/cors"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// CORS設定
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// OPTIONSリクエストを処理（CORS用）
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// ファイル取得
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		log.Println("ファイル取得エラー:", err)
		http.Error(w, "ファイルを取得できません", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 画像フォーマット判定
	format := strings.ToLower(fileHeader.Header.Get("Content-Type"))
	log.Println("アップロードされたファイルの形式:", format)

	var img image.Image
	if format == "image/png" {
		img, err = png.Decode(file)
	} else if format == "image/jpeg" || format == "image/jpg" {
		img, err = jpeg.Decode(file)
	} else {
		log.Println("対応していない画像形式:", format)
		http.Error(w, "対応していない画像形式です (JPEG / PNG のみ対応)", http.StatusUnsupportedMediaType)
		return
	}

	if err != nil {
		log.Println("画像デコードエラー:", err)
		http.Error(w, "画像のデコードに失敗しました", http.StatusInternalServerError)
		return
	}

	// クライアントからリサイズ幅・高さを取得（デフォルト: 300x300）
	width, err := strconv.Atoi(r.FormValue("width"))
	if err != nil || width <= 0 {
		width = 300
	}

	height, err := strconv.Atoi(r.FormValue("height"))
	if err != nil || height <= 0 {
		height = 300
	}

	log.Printf("リサイズ: %dx%d\n", width, height)

	// 画像リサイズ
	resizedImg := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

	// バッファに書き込む
	var buf bytes.Buffer
	if format == "image/png" {
		err = png.Encode(&buf, resizedImg)
		w.Header().Set("Content-Type", "image/png")
	} else {
		err = jpeg.Encode(&buf, resizedImg, nil)
		w.Header().Set("Content-Type", "image/jpeg")
	}

	if err != nil {
		log.Println("画像エンコードエラー:", err)
		http.Error(w, "画像のエンコードに失敗しました", http.StatusInternalServerError)
		return
	}

	// クライアントに返す
	w.Write(buf.Bytes())
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/upload", uploadHandler).Methods("POST", "OPTIONS")

	// CORSミドルウェア
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	port := "8080"
	fmt.Println("サーバー起動: http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
