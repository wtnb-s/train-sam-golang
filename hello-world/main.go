package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	personID := request.PathParameters["personID"]
	personName := request.QueryStringParameters["personName"]
	old := 25

	person := PersonResponse{
		PersonID:   personID,
		PersonName: personName,
		Old:        old,
	}
	jsonBytes, _ := json.Marshal(person)

	return events.APIGatewayProxyResponse{
		Body:       string(jsonBytes),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}

type PersonResponse struct {
	PersonID   string `json:"personID"`
	PersonName string `json:"personName"`
	Old        int    `json:"old"`
}
