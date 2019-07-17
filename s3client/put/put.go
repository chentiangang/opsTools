package put

import (
	"fmt"
	"log"
	"os"
	"project/http-deploy/butterfly/command/s3client/com"
	"sync/atomic"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/urfave/cli"
)

type CustomReader struct {
	fp   *os.File
	size int64
	read int64
}

var verbose bool

func (r *CustomReader) Read(p []byte) (int, error) {
	return r.fp.Read(p)
}

func (r *CustomReader) ReadAt(p []byte, off int64) (int, error) {
	n, err := r.fp.ReadAt(p, off)
	if err != nil {
		return n, err
	}

	// Got the length have read( or means has uploaded), and you can construct your message
	atomic.AddInt64(&r.read, int64(n))

	// I have no idea why the read length need to be div 2,
	// maybe the request read once when Sign and actually send call ReadAt again
	// It works for me

	if verbose {
		log.Printf("total read:%d    progress:%d%%\n", r.read/2, int(float32(r.read*100/2)/float32(r.size)))
	}
	return n, err
}

func (r *CustomReader) Seek(offset int64, whence int) (int64, error) {
	return r.fp.Seek(offset, whence)
}

func Put(c *cli.Context) error {
	var (
		bucket, fileName, region, key string
		err                           error
	)
	verbose = c.Bool("verbose")
	if bucket, err = com.String("bucket", c); err != nil {
		return err
	}
	if fileName, err = com.String("src", c); err != nil {
		return err
	}

	region = c.String("region")
	if key, err = com.String("dest", c); err != nil {
		return err
	}

	credential := ""

	creds := credentials.NewSharedCredentials(credential, "default")
	if _, err := creds.Get(); err != nil {
		return cli.NewExitError(err, 1)
	}

	sess := session.New(&aws.Config{
		Region: aws.String(region),
	})

	file, err := os.Open(fileName)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	reader := &CustomReader{
		fp:   file,
		size: fileInfo.Size(),
	}

	uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024
		u.LeavePartsOnError = true
	})

	output, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:       aws.String(bucket),
		Key:          aws.String(key),
		Body:         reader,
		StorageClass: aws.String("STANDARD_IA"),
	})

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	if verbose {
		log.Println(output.Location)
	}
	fmt.Println(key)
	return nil
}
