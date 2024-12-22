package queue

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type ServerConnection struct {
	URL string
}

type ServerReponse struct {
	Message string `json:"message"`
}

var client http.Client

func (c *ServerConnection) CreateQueue(name string) bool {
	rawBody := map[string]string{
		"name": name,
	}

	jsonValue, _ := json.Marshal(rawBody)
	bodyReader := bytes.NewReader(jsonValue)
	req, _ := http.NewRequest("POST", c.URL+"/queue", bodyReader)
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return false
	}
	defer resp.Body.Close()
	return true
}

func (c *ServerConnection) Publish(q string, m string) bool {
	rawBody := map[string]string{
		"message": m,
	}

	jsonValue, _ := json.Marshal(rawBody)
	bodyReader := bytes.NewReader(jsonValue)
	req, _ := http.NewRequest("POST", c.URL+"/"+q, bodyReader)
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return false
	}
	body, _ := io.ReadAll(resp.Body)
	var bodyJson ServerReponse
	json.Unmarshal(body, &bodyJson)
	defer resp.Body.Close()
	return true
}

func (c *ServerConnection) Consume(q string) chan string {
	consumer := make(chan string)
	go func() {
		for {
			req, _ := http.NewRequest("GET", c.URL+"/"+q, nil)
			resp, err := client.Do(req)
			if err != nil {
				log.Print(err)
				break
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				var bodyJson ServerReponse
				err := json.Unmarshal(body, &bodyJson)
				if err != nil {
					log.Printf("Failed decoding json message when consuming from queue "+q+" %s", err)
				}
				consumer <- bodyJson.Message
			} else if resp.StatusCode == http.StatusNotFound {
				log.Printf("Failed consuming message since queue " + q + " not found")
			}
			time.Sleep(time.Second)
		}
		close(consumer)
	}()
	return consumer
}
