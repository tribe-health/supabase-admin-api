package optimizations

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/pkg/errors"
	"os"
	"reflect"
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

func generateSettings(settings interface{}) (*string, error) {
	var buffer bytes.Buffer
	val := reflect.ValueOf(&settings)
	for _, field := range reflect.VisibleFields(reflect.TypeOf(settings)) {
		serverFieldName := field.Tag.Get("conf")
		f := val.Elem().Elem().FieldByName(field.Name)
		if f.IsZero() {
			continue
		}
		switch f.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			buffer.WriteString(fmt.Sprintf("%s = %d\n", serverFieldName, f.Int()))
			break
		case reflect.Float64, reflect.Float32:
			buffer.WriteString(fmt.Sprintf("%s = %0.03f\n", serverFieldName, f.Float()))
			break
		case reflect.String:
			buffer.WriteString(fmt.Sprintf("%s = '%s'\n", serverFieldName, f.String()))
			break
		default:
			return nil, fmt.Errorf("unsupported type encountered for field %+v: %+v", f, f.Kind())
		}
	}
	return aws.String(buffer.String()), nil
}

func writeRecommendationsToFile(settings interface{}, destinationFilePath string) error {
	output, err := generateSettings(settings)
	if err != nil {
		return errors.Wrap(err, "couldn't serialize settings")
	}
	err = os.WriteFile(destinationFilePath, []byte(*output), 0644)
	if err != nil {
		return errors.Wrapf(err, "couldn't write recommendations to %s", destinationFilePath)
	}
	return nil
}

