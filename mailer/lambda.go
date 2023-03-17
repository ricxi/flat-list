package mailer

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

type lambdaHandler func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

// HandleSendActivationEmail is a lambda function that sends an email to a user.
// It can be accessed through AWS API Gateway
func SendActivationEmail(es *EmailService) lambdaHandler {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		var data UserActivationData

		if err := json.Unmarshal([]byte(request.Body), &data); err != nil {
			log.Println(err)
			return events.APIGatewayProxyResponse{}, err
		}

		if err := es.SendActivationEmail(data); err != nil {
			log.Println(err)
			return events.APIGatewayProxyResponse{}, err
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
		}, nil
	}
}
