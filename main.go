package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Junkes887/web-crawler/artifacts"
	"github.com/Junkes887/web-crawler/connection"
	"github.com/PuerkitoBio/goquery"
	"github.com/segmentio/kafka-go"
)

func main() {
	con := connection.ConnectionKafka()
	consume(con)
}

func consume(con *kafka.Reader) {
	for {
		menssage, err := con.ReadMessage(context.Background())
		artifacts.HandlerError(err)
		links := searchUrl(string(menssage.Value))
		topic := "topic-links"
		partition := 0

		conn, err := kafka.DialLeader(context.Background(), "tcp", "kafka:29092", topic, partition)
		if err != nil {
			log.Fatal("Erro DialLeader Kafka: \n ", err)
		}

		for link := range links {
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			fmt.Printf("Link encontrado: %s \n" + link)
			_, err = conn.WriteMessages(
				kafka.Message{Value: []byte(link)},
			)
			if err != nil {
				log.Fatal("Erro WriteMessage Kafka: \n ", err)
			}
		}
	}
}

func searchUrl(url string) map[string]int {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	artifacts.HandlerError(err)

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:58.0) Gecko/20100101 Firefox/58.0")
	res, err := client.Do(req)
	artifacts.HandlerError(err)
	return findLinks(res)
}

func findLinks(res *http.Response) map[string]int {
	links := make(map[string]int)
	doc, err := goquery.NewDocumentFromReader(res.Body)
	artifacts.HandlerError(err)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if len(href) >= 4 && href[0:4] == "http" {
			links[href] = i
		}
	})
	return links
}
