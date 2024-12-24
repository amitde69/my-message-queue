package queue

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type QueueData struct {
	Queues map[string]Queue
	Mutex  sync.Mutex
}

var queueData QueueData

type PublishMessageReq struct {
	Message string `json:"message"`
}

func InitilizeQueue(name string) {
	queueData.Queues[name] = Queue{Name: name}
	router.POST("/"+name, PublishMessage)
	router.GET("/"+name, ConsumeMessage)
}

func CreateQueue(context *gin.Context) {
	var newQueue Queue

	if err := context.BindJSON(&newQueue); err != nil {
		fmt.Println(err)
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}
	queueData.Mutex.Lock()
	defer queueData.Mutex.Unlock()
	if _, ok := queueData.Queues[newQueue.Name]; ok {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Queue named " + newQueue.Name + " already exists",
		})
		return
	}
	InitilizeQueue(newQueue.Name)
	context.IndentedJSON(http.StatusOK, gin.H{
		"message": "Queue " + newQueue.Name + " created succcesfully",
	})
}

func PublishMessage(context *gin.Context) {
	var newMessageReq PublishMessageReq
	var QueueName string

	if err := context.BindJSON(&newMessageReq); err != nil {
		log.Print(err)
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}

	path := context.Request.URL.Path
	pathSegments := strings.Split(strings.TrimPrefix(path, "/"), "/")
	QueueName = pathSegments[0] // Get the first segment
	queueData.Mutex.Lock()
	defer queueData.Mutex.Unlock()
	if _, ok := queueData.Queues[QueueName]; !ok {
		context.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "Queue named " + QueueName + " not found",
		})
		return
	}
	myq := queueData.Queues[QueueName]
	myq.Publish(newMessageReq.Message)
	queueData.Queues[QueueName] = myq
	context.IndentedJSON(http.StatusOK, gin.H{
		"message": "Message published to queue " + QueueName,
	})
}

func ConsumeMessage(context *gin.Context) {
	var QueueName string
	path := context.Request.URL.Path
	header := context.GetHeader("ack")
	var ack bool
	if header != "" {
		ack, _ = strconv.ParseBool(header)
		// if err != nil {
		// 	context.IndentedJSON(http.StatusBadRequest, gin.H{
		// 		"message": "ack header doesnt contain a boolean value",
		// 	})
		// 	return
		// }
	} else {
		ack = true
	}

	pathSegments := strings.Split(strings.TrimPrefix(path, "/"), "/")
	QueueName = pathSegments[0] // Get the first segment
	queueData.Mutex.Lock()
	defer queueData.Mutex.Unlock()
	if _, ok := queueData.Queues[QueueName]; !ok {
		context.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "Queue named " + QueueName + " not found",
		})
		return
	}
	myq := queueData.Queues[QueueName]
	consumedMessage := myq.Consume(ack)
	if consumedMessage.Payload == "" {
		context.IndentedJSON(http.StatusNoContent, gin.H{
			"message": "No messages in queue " + QueueName,
		})
	}
	context.IndentedJSON(http.StatusOK, gin.H{
		"message": consumedMessage.Payload,
	})
}

var router *gin.Engine

func RunServer() {
	queueData.Mutex.Lock()
	if queueData.Queues == nil {
		queueData.Queues = map[string]Queue{}
	}
	queueData.Mutex.Unlock()
	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
	router.POST("/queue", CreateQueue)
	fmt.Print("Queue server is running on port 6060...\n")
	router.Run(":6060")
}
