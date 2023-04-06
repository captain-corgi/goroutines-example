package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	_ "github.com/godror/godror"
)

type Data struct {
	Col1 string `json:"col1"`
	Col2 int    `json:"col2"`
}

func getData(wg *sync.WaitGroup, db *sql.DB, offset int, limit int, dataChan chan []Data) {
	defer wg.Done()

	// Thay đổi câu truy vấn của bạn
	query := fmt.Sprintf("SELECT * FROM yourtable OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", offset, limit)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var data []Data

	// Thay đổi cấu trúc dữ liệu của bạn
	var col1 string
	var col2 int

	for rows.Next() {
		err := rows.Scan(&col1, &col2)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, Data{Col1: col1, Col2: col2})
	}

	dataChan <- data
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	// Thay đổi thông tin kết nối cơ sở dữ liệu của bạn
	db, err := sql.Open("godror", `user="yourusername" password="yourpassword" connectString="yourhost:yourport/yourSID"`)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var wg sync.WaitGroup
	dataChan := make(chan []Data)

	// Thay đổi số lượng goroutine và số lượng dòng dữ liệu được trả về cho mỗi goroutine
	numGoroutines := 10
	rowsPerGoroutine := 1000000000 / numGoroutines

	for i := 0; i < numGoroutines; i++ {
		offset := i * rowsPerGoroutine
		wg.Add(1)
		go getData(&wg, db, offset, rowsPerGoroutine, dataChan)
	}

	go func() {
		wg.Wait()
		close(dataChan)
	}()

	var data []Data

	for d := range dataChan {
		data = append(data, d...)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func main() {
	http.HandleFunc("/", handlePost)

	// Thay đổi cổng mà máy chủ lắng nghe
	log.Fatal(http.ListenAndServe(":8080", nil))
}
