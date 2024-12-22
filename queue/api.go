package queue

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var data map[string]Queue

type PublishMessageReq struct {
	Message string `json:"message"`
}

func InitilizeQueue(name string) {
	data[name] = Queue{Name: name}
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
	if _, ok := data[newQueue.Name]; ok {
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
	if _, ok := data[QueueName]; !ok {
		context.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "Queue named " + QueueName + " not found",
		})
		return
	}
	myq := data[QueueName]
	myq.Publish(newMessageReq.Message)
	data[QueueName] = myq
	context.IndentedJSON(http.StatusOK, gin.H{
		"message": "Message published to queue " + QueueName,
	})
}

func ConsumeMessage(context *gin.Context) {
	var QueueName string
	path := context.Request.URL.Path
	pathSegments := strings.Split(strings.TrimPrefix(path, "/"), "/")
	QueueName = pathSegments[0] // Get the first segment
	if _, ok := data[QueueName]; !ok {
		context.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "Queue named " + QueueName + " not found",
		})
		return
	}
	myq := data[QueueName]
	consumedMessage := myq.Consume()
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
	if data == nil {
		data = map[string]Queue{}
	}
	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
	router.POST("/queue", CreateQueue)
	fmt.Print("Queue server is running on port 6060...\n")
	router.Run(":6060")
}
