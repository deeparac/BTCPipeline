package main

import (
    _ "fmt"
    "flag"
    _ "os"
    _ "os/signal"
    _ "strings"
    _ "io/ioutil"
    _ "net/http"
    _ "encoding/json"

    "./producer"
    "./consumer"
)

var (
    // kafka arguments
    broker         = flag.String("broker", "192.168.99.103:9092", "Kafka brokers to connect to")
    topic          = flag.String("topic", "btcpipe", "topic that's read from")
    verbose        = flag.Bool("verbose", false, "Turn on Sarama logging")

    // mongodb arguments
    mongo_host     = flag.String("mongo_host", "192.168.99.104:27017", "MongoDB address")
    dbName         = flag.String("db", "btcbd", "MongoDB Database Name")
    collectionName = flag.String("collection", "btc_price", "MongoDB Collection Name")

    // query arguments
    freqInSecond   = flag.Uint64("query_frequency", 2, "query frequency to server and send to kafka")
    queryUrl       = flag.String("queryUrl", "https://api.coindesk.com/v1/bpi/currentprice.json", "the url that be queried for btc price")
)

func init() {
    flag.Parse()
    producer.InitProducer()
    consumer.InitConsumer(*mongo_host)
}

func main() {
    producer.ProducerProcess(*broker, *topic, *queryUrl, *freqInSecond)
    consumer.PersistProcess(*broker, *topic, *dbName, *collectionName)
}
