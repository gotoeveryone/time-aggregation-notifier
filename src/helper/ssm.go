package helper

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func GetParameter(c ssm.Client, name string) (*string, error) {
	o, err := c.GetParameter(context.TODO(), &ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: true,
	})
	if err != nil {
		return nil, err
	}
	return o.Parameter.Value, nil
}
