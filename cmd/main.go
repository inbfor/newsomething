package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/BurntSushi/toml"
)

func main() {
	var conf = map[string]interface{}{
		"user_storage": map[string]string{},
		"mediabasket":  map[string]string{},
	}

	data, err := http.Get("http://prometheus-web.dc.wildberries.ru:9090/api/v1/query?query=max%20by%20(instance)(up{role=~%22.*basket.*%22})")
	if err != nil {
		fmt.Println("cant get data")
	}

	defer data.Body.Close()
	body, _ := ioutil.ReadAll(data.Body)
	var config Config

	json.Unmarshal(body, &config)

	f, err := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("cant open file")
	}

	for _, item := range config.Data.Result {
		replica := item.Metric.Instance
		conf[DecodeReplica(replica)] = replica
	}

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(config); err != nil {
		log.Fatal(err)
	}

	f.Write(buf.Bytes())
}

func DecodeReplica(replica string) string {
	replica = replica + ".wb.ru"
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
