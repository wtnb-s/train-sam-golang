package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// 環境変数設定
	endpoint := os.Getenv("DYNAMODB_ENDPOINT")

	// リクエストボディ取得
	reqBody := request.Body
	jsonBytes := ([]byte)(reqBody)
	presonReq := new(PersonRequest)
	if err := json.Unmarshal(jsonBytes, presonReq); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	personID := presonReq.PersonID
	personName := presonReq.PersonName
	testID := presonReq.TestID
	score := presonReq.Score
	passingMark := false
	if score >= 80 {
		passingMark = true
	}

	person := Person{
		PersonID:    personID,
		PersonName:  personName,
		TestID:      testID,
		Score:       score,
		PassingMark: passingMark,
	}
	// item を dynamodb attributeに変換
	av, err := dynamodbattribute.MarshalMap(person)
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{}, err
	}

	// Dynamodb接続設定
	session := session.Must(session.NewSession())
	config := aws.NewConfig().WithRegion("ap-northeast-1")
	if len(endpoint) > 0 {
		config = config.WithEndpoint(endpoint)
	}
	db := dynamodb.New(session, config)

	_, err = db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("Score"),
		Item:      av,
	})
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(jsonBytes),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}

type PersonRequest struct {
	PersonID   string `json:"personID"`
	PersonName string `json:"personName"`
	TestID     string `json:"testID"`
	Score      int    `json:"score"`
}

type Person struct {
	PersonID    string
	PersonName  string
	TestID      string
	Score       int
	PassingMark bool
}
