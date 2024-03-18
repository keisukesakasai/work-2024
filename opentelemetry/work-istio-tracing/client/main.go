package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	logging "client/internal/log"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	Service               = "Population-Query"
	SleepEnv              = "SLEEP"
	ServerAddressEnv      = "SERVER_ADDRESS"
	FruitServerAddressEnv = "FRUIT_SERVER_ADDRESS"
)

func main() {
	url := os.Getenv(ServerAddressEnv)
	// fruitUrl := os.Getenv(FruitServerAddressEnv)
	sleep, _ := strconv.Atoi(os.Getenv(SleepEnv))

	// Start main task.
	for range time.Tick(time.Duration(sleep) * time.Second) {
		ctx := context.Background()
		rand.NewSource(time.Now().UnixNano())
		prefecture := prefectures[rand.Intn(len(prefectures))]
		queryUrl := fmt.Sprintf(url+"?pref=%s", prefecture)

		// queryPopulation
		err := queryPopulation(ctx, queryUrl, prefecture)
		if err != nil {
			log.Fatal(err)
		}

		// queryFruit
		/*
			err = queryFruit(ctx, fruitUrl)
			if err != nil {
				log.Fatal(err)
			}
		*/
	}
}

func queryPopulation(ctx context.Context, url string, pref string) (err error) {
	logger := logging.GetLoggerFromCtx(ctx)
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	err = func() error {
		req, _ := http.NewRequestWithContext(
			ctx,
			"GET",
			url,
			nil,
		)

		logger.Infoln(pref + "の人口をクエリ...")
		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		body, _ := io.ReadAll(res.Body)
		_ = res.Body.Close()

		logger.Infof(pref+"の人口は %s 人です。", body)

		/*
			projectID := "my-datastore-project"
			clientDatastore, err := datastore.NewClient(ctx, projectID)
			if err != nil {
				logger.Error(err)
			}
			data := &Data{
				Prefecture: pref,
				Population: string(body),
				Date:       time.Now(),
			}
			k := datastore.NameKey("Data", uuid.New().String(), nil)
			if _, err := clientDatastore.Put(ctx, k, data); err != nil {
				logger.Error(err)
			}
			logger.Infof("Datastore への Put をしました。", data)
		*/

		return err
	}()

	return err
}

/*
func queryFruit(ctx context.Context, url string) (err error) {
	logger := logging.GetLoggerFromCtx(ctx)
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	err = func() error {
		req, _ := http.NewRequestWithContext(
			ctx,
			"GET",
			url,
			nil,
		)

		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		body, _ := io.ReadAll(res.Body)
		_ = res.Body.Close()

		logger.Infof("フルーツは %s です。", body)

		return err
	}()

	return err
}
*/

var prefectures = []string{
	"Hokkaido", "Aomori-Ken", "Iwate-Ken", "Miyagi-Ken", "Akita-Ken",
	"Yamagata-Ken", "Fukushima-Ken", "Ibaraki-Ken", "Tochigi-Ken", "Gunma-Ken",
	"Saitama-Ken", "Chiba-Ken", "Tokyo-To", "Kanagawa-Ken", "Niigata-Ken",
	"Toyama-Ken", "Ishikawa-Ken", "Fukui-Ken", "Yamanashi-Ken", "Nagano-Ken",
	"Gifu-Ken", "Shizuoka-Ken", "Aichi-Ken", "Mie-Ken", "Shiga-Ken",
	"Kyoto-Fu", "Osaka-Fu", "Hyogo-Ken", "Nara-Ken", "Wakayama-Ken",
	"Tottori-Ken", "Shimane-Ken", "Okayama-Ken", "Hiroshima-Ken", "Yamaguchi-Ken",
	"Tokushima-Ken", "Kagawa-Ken", "Ehime-Ken", "Kochi-Ken", "Fukuoka-Ken",
	"Saga-Ken", "Nagasaki-Ken", "Kumamoto-Ken", "Oita-Ken", "Miyazaki-Ken",
	"Kagoshima-Ken", "Okinawa-Ken",
}

type Data struct {
	Prefecture string
	Population string
	Date       time.Time
}
