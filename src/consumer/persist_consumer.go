package consumer

import (
    "gopkg.in/mgo.v2"
    _ "gopkg.in/mgo.v2/bson"

    "os"
    "os/signal"
    "syscall"
    "fmt"
    "strings"

    "github.com/Shopify/sarama"
    _ "github.com/jasonlvhit/gocron"
)

/*
var (
    broker         = flag.String("broker", "192.168.99.103:9092", "Kafka brokers to connect to")
    topic          = flag.String("topic", "btcpipe", "topic that's read from")
    verbose        = flag.Bool("verbose", false, "Turn on Sarama logging")

    mongo_host     = flag.String("mongo_host", "192.168.99.104:27071", "MongoDB address")
    dbName         = flag.String("db", "btcbd", "MongoDB Database Name")
    collectionName = flag.String("collection", "btc_price", "MongoDB Collection Name")
)
*/

type BTC struct {
    Timestamp string
    Rate      string
}

var session *mgo.Session
var config = sarama.NewConfig()

func InitConsumer(mongo_host string) {
    config.Consumer.Return.Errors = true

    session, err:= mgo.Dial(mongo_host)
    if err != nil {
        fmt.Println("Failed to connect to MongoDB")
        panic(err)
    }
    session.SetMode(mgo.Monotonic, true)
    mongoClose()
}

func mongoClose() {
    ch := make(chan os.Signal, 1)
    signal.Notify(ch, os.Interrupt)
    signal.Notify(ch, syscall.SIGTERM)
    signal.Notify(ch, syscall.SIGKILL)
    go func() {
        <-ch
        session.Close()
        os.Exit(0)
    }()
}

func PersistProcess(broker string, topic string, db string, collection string) {
   brokerslc := []string{broker}

   master, err := sarama.NewConsumer(brokerslc, config)
   if err != nil {
       fmt.Println("NewConsumer fail")
       panic(err)
   }
    defer func() {
        if err := master.Close(); err != nil {
            fmt.Println("Fail to close consumer!")
        }
    }()

    consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetOldest)
    if err != nil {
        fmt.Println("ConsumerPartition fail")
        panic(err)
    }

    signals := make(chan os.Signal, 1)
    signal.Notify(signals, os.Interrupt)

    msgCount := 0

    doneCh := make(chan struct{})
    go func() {
        for {
            select {
                case err := <-consumer.Errors():
                    fmt.Println(err)
                case msg := <-consumer.Messages():
                    msgCount++
                    fmt.Println("Received messages", string(msg.Key), string(msg.Value))
                    persist_data(msg, session, db, collection)
                case <-signals:
                    fmt.Println("Interrupt!")
                    doneCh <- struct{}{}
            }
        }
    }()

    <-doneCh
    fmt.Println("Done processing", msgCount, "messages")
}

func persist_data(msg *sarama.ConsumerMessage, session *mgo.Session, db string, collection string) {
    c := session.DB(db).C(collection)

    parsed := strings.Split(string(msg.Value), "@")
    statement := BTC{parsed[0], parsed[1]}

    err := c.Insert(&statement)
    if err != nil {
        fmt.Println("DB Insert fail")
        panic(err)
    }

    fmt.Println("Successful persist one record", statement)
}
