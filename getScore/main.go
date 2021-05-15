package main

import (
	"encoding/json"
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
	// クエリパラメータ取得
	personID := request.PathParameters["personID"]
	testID := request.QueryStringParameters["testID"]

	// Dynamodb接続設定
	session := session.Must(session.NewSession())
	config := aws.NewConfig().WithRegion("ap-northeast-1")
	if len(endpoint) > 0 {
		config = config.WithEndpoint(endpoint)
	}
	db := dynamodb.New(session, config)

	param, err := db.Query(&dynamodb.QueryInput{
		TableName: aws.String("Score"),
		ExpressionAttributeNames: map[string]*string{
			"#PersonID":    aws.String("PersonID"),
			"#TestID":      aws.String("TestID"),
			"#PersonName":  aws.String("PersonName"),
			"#Score":       aws.String("Score"),
			"#PassingMark": aws.String("PassingMark"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":personID": {
				S: aws.String(personID),
			},
			":testID": {
				S: aws.String(testID),
			},
		},
		KeyConditionExpression: aws.String("#PersonID=:personID AND #TestID=:testID"),
		ProjectionExpression:   aws.String("#PersonID, #TestID, #PersonName, #Score, #PassingMark"),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	scores := make([]*ScoreRes, 0)
	if err := dynamodbattribute.UnmarshalListOfMaps(param.Items, &scores); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	jsonBytes, _ := json.Marshal(scores)

	return events.APIGatewayProxyResponse{
		Body:       string(jsonBytes),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}

type ScoreRes struct {
	PersonID    string `json:"personID"`
	PersonName  string `json:"personName"`
	TestID      string `json:"testID"`
	Score       int    `json:"score"`
	PassingMark bool   `json:"PassingMark"`
}
