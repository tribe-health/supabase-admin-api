package optimizations

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/pkg/errors"
)

func GetInstanceType(ctx context.Context) (*InstanceType, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't obtain ec2metadata config")
	}

	client := imds.NewFromConfig(cfg)
	document, err := client.GetInstanceIdentityDocument(ctx, &imds.GetInstanceIdentityDocumentInput{})
	if err != nil {
		return nil, errors.Wrap(err, "couldn't obtain metadata")
	}
	return &document.InstanceType, nil
}

type InstanceType = string
