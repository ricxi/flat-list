package mailer

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type lambdaHandler func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

// HandleSendActivationEmail is a lambda function that sends an email to a user.
// It is accessed through an AWS API Gateway.
func SendActivationEmail(mailerService *Service) lambdaHandler {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		var data ActivationEmailData

		if err := json.Unmarshal([]byte(request.Body), &data); err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		if err := mailerService.sendActivationEmail(data); err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
		}, nil
	}
}
