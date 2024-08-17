package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

	_ "time/tzdata"

	"cloud.google.com/go/profiler"
)

func main() {
	// 環境変数からアプリバージョンを取得
	appVersion := os.Getenv("APP_VERSION")

	// profiling 設定
	cfg := profiler.Config{
		Service:        "app",
		ServiceVersion: appVersion,
		ProjectID:      "datadog-sandbox",
		DebugLogging:   true,
	}
	// Profiler initialization, best done as early as possible.
	if err := profiler.Start(cfg); err != nil {
		log.Fatal(err)
	}

	// Http server
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// リクエストのボディを取得します
		w.Write([]byte("Hello World!"))

		// ロジック
		count := calcTargetLogic(appVersion)
		fmt.Println("count: ", count)
	})

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func calcTargetLogic(appVersion string) (total int) {
	dummyData, err := read("./data/input.txt")
	if err != nil {
		fmt.Println(err.Error())
	}

	total = count(dummyData, appVersion)

	return total
}

func read(filename string) ([]int, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var dummyData []int
	for _, char := range string(content) {
		value, err := strconv.Atoi(string(char))
		if err != nil {
			return nil, fmt.Errorf("error converting character to int: %v", err)
		}
		if value != 0 && value != 1 {
			return nil, fmt.Errorf("invalid value in file: %d", value)
		}
		dummyData = append(dummyData, value)
	}

	return dummyData, nil
}

func count(dummyData []int, appVersion string) (total int) {
	switch appVersion {
	case "v1.0.0":
		n := len(dummyData)
		for i := 0; i < n; i++ {
			for j := 0; j < n-i-1; j++ {
				if dummyData[j] > dummyData[j+1] {
					dummyData[j], dummyData[j+1] = dummyData[j+1], dummyData[j]
				}
			}
		}
		index := sort.SearchInts(dummyData, 1)
		fmt.Println("index: ", index)

		return len(dummyData) - index

	case "v2.0.0":
		sort.Ints(dummyData)
		index := sort.SearchInts(dummyData, 1)
		fmt.Println("index: ", index)

		return len(dummyData) - index

	default:
		return 0
	}
}
