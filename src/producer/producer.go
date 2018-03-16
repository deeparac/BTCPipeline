package producer

import (
    "github.com/Shopify/sarama"
    "github.com/jasonlvhit/gocron"

    "fmt"
    "strings"
    "io/ioutil"
    "net/http"
    "encoding/json"
)

/*
var (
    broker      = flag.String("broker", "192.168.99.103:9092", "The Kafka brokers to connect to")
    topic       = flag.String("topic", "btcpipe", "topic that's wrote into")
    verbose     = flag.Bool("verbose", false, "Turn on Sarama logging")
    queryUrl    = flag.String("queryUrl", "https://api.coindesk.com/v1/bpi/currentprice.json", "the url that be queried for btc price")
)
*/

type BTC_Info struct {
    Time struct {
        Updated string `json: time: updated`
        updatedISO string
        updateduk string
    } `json: time`
    diclaimer string
    chartName string
    Bpi struct {
        USD struct {
            code string
            symbol string
            Rate string `json: bpi: USD: rate`
            description string
            rate_float float32
        }
        GBP struct {
            code string
            symbol string
            rate string
            description string
            rate_float float32
        }
        EUR struct {
            code string
            symbol string
            rate string
            description string
            rate_float float32
        }
    } `json: bpi`
}

var config = sarama.NewConfig()

func InitProducer() {
    config.Producer.RequiredAcks = sarama.WaitForAll
    config.Producer.Retry.Max = 5
    config.Producer.Return.Successes = true
}

func ProducerProcess(broker string, topic string, queryUrl string, freqInSecond uint64) {
    brokerslc := []string{broker}
    producer, err := sarama.NewSyncProducer(brokerslc, config)
    if err != nil {
        panic(err)
    }

    defer func() {
        if err := producer.Close(); err != nil {
            panic(err)
        }
    }()

    s := gocron.NewScheduler()
    s.Every(freqInSecond).Seconds().Do(queryPrice, queryUrl, topic, producer)
    <-s.Start()
}

func queryPrice(url string, topic string, producer sarama.SyncProducer) error {
    r, err := http.Get(url)
    if err != nil {
        fmt.Println("Get fail")
        panic(err)
    }
    defer r.Body.Close()

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        fmt.Println("ReadAll fail")
        panic(err)
    }

    var parsedInfo BTC_Info
    if err := json.Unmarshal([]byte(body), &parsedInfo); err != nil {
        fmt.Println("Unmarshall fail")
        panic(err)
    }

    // parsing information
    price := parsedInfo.Bpi.USD.Rate
    timestamp := strings.Split(parsedInfo.Time.Updated, " ")[0:4]

    // produce to kafka
    msg_val := price + "@" + strings.Join(timestamp[:], "")
    msg := &sarama.ProducerMessage {
        Topic: topic,
        Value: sarama.StringEncoder(msg_val),
    }

    _, _, sendErr := producer.SendMessage(msg)
    if sendErr != nil {
        fmt.Println("SendMessage fail")
        panic(sendErr)
    }
    fmt.Println("s flag")
    fmt.Sprintf("Successfully sent BTC price (%s) at time (%s) in topic (%s)", price, timestamp, topic)
    return nil
}
