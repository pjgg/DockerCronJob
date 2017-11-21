package connectors

import (
	"bytes"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var awsConnectorSingleton sync.Once
var awsConnectorImpl AwsConnector

type AwsConnector struct {
	Aws_secret_key string
	Aws_access_id  string
	Bucket_name    string
	credentials    *credentials.Credentials
	config         *aws.Config
	S3Client       *s3.S3
}

type AwsConnectorBehavior interface {
	Push(destPath, sourcePath string) error
}

func Instance(bucketName string) AwsConnectorBehavior {
	awsConnectorSingleton.Do(func() {
		var ok bool
		if awsConnectorImpl.Aws_secret_key, ok = os.LookupEnv("ORG_GRADLE_PROJECT_AWS_SECRET_KEY"); !ok {
			logrus.Fatal("missing ORG_GRADLE_PROJECT_AWS_SECRET_KEY")
		}

		if awsConnectorImpl.Aws_access_id, ok = os.LookupEnv("ORG_GRADLE_PROJECT_AWS_ACCESS_KEY"); !ok {
			logrus.Fatal("missing ORG_GRADLE_PROJECT_AWS_ACCESS_KEY")
		}

		awsConnectorImpl.Bucket_name = bucketName
		creds := credentials.NewStaticCredentials(awsConnectorImpl.Aws_access_id, awsConnectorImpl.Aws_secret_key, "")

		_, err := creds.Get()
		if err != nil {
			logrus.Fatal("bad credentials: %s", err)
		}
		awsConnectorImpl.credentials = creds
		awsConnectorImpl.config = aws.NewConfig().WithRegion("eu-central-1").WithCredentials(creds)
		awsConnectorImpl.S3Client = s3.New(session.New(), awsConnectorImpl.config)
	})

	return &awsConnectorImpl
}

func (self *AwsConnector) Push(destPath, sourcePath string) error {
	files, err := filepath.Glob(sourcePath)
	for _, elem := range files {
		logrus.Info("To s3, File: " + elem)
		file, err := os.Open(elem)
		if err != nil {
			logrus.Error("err opening file: %s", err)
		}
		defer file.Close()
		fileInfo, _ := file.Stat()
		size := fileInfo.Size()
		buffer := make([]byte, size) // read file content to buffer

		file.Read(buffer)
		fileBytes := bytes.NewReader(buffer)
		fileType := http.DetectContentType(buffer)
		path := destPath + filepath.Base(file.Name())

		if filepath.Ext(filepath.Base(file.Name())) == "css" {
			fileType = "text/css"
		}

		params := &s3.PutObjectInput{
			Bucket:        aws.String(self.Bucket_name),
			Key:           aws.String(path),
			Body:          fileBytes,
			ContentLength: aws.Int64(size),
			ContentType:   aws.String(fileType),
		}
		resp, err := self.S3Client.PutObject(params)
		if err != nil {
			logrus.Error("bad response: %s", err)
		}
		logrus.Info("response %s", awsutil.StringValue(resp))
	}

	return err
}
