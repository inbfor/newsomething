package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

func main() {
	data, err := http.Get("http://prometheus-web.dc.wildberries.ru:9090/api/v1/query?query=max%20by%20(instance)(up{role=~%22.*basket.*%22})")
	if err != nil {
		fmt.Println("cant get data")
	}

	defer data.Body.Close()
	body, _ := ioutil.ReadAll(data.Body)
	var config Config

	json.Unmarshal(body, &config)

	for _, item := range config.Data.Result {
		fmt.Println(item.Metric.Instance)
		fmt.Println(DecodeReplica(item.Metric.Instance))
	}
}

func DecodeReplica(replica string) string {
	userStorage := regexp.MustCompile(`^user-storage-\d\d\w?(-\d)?.\w\w.wb.ru$`)

	ordersBasket := regexp.MustCompile(`^(catalog-)?mediabasket-orders-(basket-)?(baskets-)?\d\d\w?(-\d)?.\w\w.wb.ru$ `)

	mediaBasketFirst := regexp.MustCompile(`^mediabasket-nsk-cdn-\d\d\w?(-\d)?.\w\w.wb.ru$`)
	mediaBasketSecond := regexp.MustCompile(`^(catalog-)?mediabasket-(basket-)?(baskets-)?\d\d\w?(-\d)?.\w\w.wb.ru$`)
	mediaBasketThird := regexp.MustCompile(`^basket-\d\d\w?(-\d)?.\w\w.wb.ru$`)

	digitalBasketFirst := regexp.MustCompile(`^(catalog-)?mediabasket-digital-(basket-)?(baskets-)?\d\d\w?(-\d)?.\w\w.wb.ru$`)
	digitalBasketSecond := regexp.MustCompile(`^digital-basket-\d\d\w?(-\d)?.\w\w.wb.ru$`)

	switch {
	case userStorage.MatchString(replica):
		return "user-storage"
	case ordersBasket.MatchString(replica):
		return "orders-basket"
	case mediaBasketFirst.MatchString(replica):
		return "media-basket"
	case mediaBasketSecond.MatchString(replica):
		return "media-basket"
	case mediaBasketThird.MatchString(replica):
		return "media-basket"
	case digitalBasketFirst.MatchString(replica):
		return "digital-basket"
	case digitalBasketSecond.MatchString(replica):
		return "digital-basket"
	default:
		return "no"
	}
}

type Config struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

type Data struct {
	ResultType string `json:"resultType"`
	Result     []struct {
		Metric struct {
			Instance string `json:"instance"`
			Value    string `json:"value"`
		} `json:"metric"`
	} `json:"result"`
}
