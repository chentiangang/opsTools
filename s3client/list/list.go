package list

import (
	"fmt"
	"project/http-deploy/butterfly/command/s3client/com"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/urfave/cli"
)

func List(c *cli.Context) error {
	bucket, err := com.String("bucket", c)
	if err != nil {
		return err
	}

	sess := session.Must(session.NewSession())
	svc := s3.New(sess, &aws.Config{
		Region: aws.String("ap-northeast-1"),
	})
	i := 0
	err = svc.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: &bucket,
	}, func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {
		fmt.Println("Page,", i)
		i++

		for _, obj := range p.Contents {
			fmt.Println("Object:", *obj.Key)
		}
		return true
	})
	if err != nil {
		err = fmt.Errorf("failed to list objects %s", err)
		return cli.NewExitError(err, 1)
	}
	return nil
}
