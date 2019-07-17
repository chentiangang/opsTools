package get

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"github.com/chentiangang/opsTools/s3/com"
	"sync/atomic"

	"github.com/chentiangang/xlog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/urfave/cli"
)

type progressWriter struct {
	written int64
	writer  io.WriterAt
	size    int64
}

func (pw *progressWriter) WriteAt(p []byte, off int64) (int, error) {
	atomic.AddInt64(&pw.written, int64(len(p)))

	percentageDownloaded := float32(pw.written*100) / float32(pw.size)

	fmt.Printf("File size:%d downloaded:%d percentage:%.2f%%\r", pw.size, pw.written, percentageDownloaded)

	return pw.writer.WriteAt(p, off)
}

func byteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func getFileSize(svc *s3.S3, bucket string, prefix string) (filesize int64, error error) {
	params := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(prefix),
	}

	resp, err := svc.HeadObject(params)
	if err != nil {
		return 0, err
	}

	return *resp.ContentLength, nil
}

func Get(c *cli.Context) error {
	bucket, err := com.String("bucket", c)
	if err != nil {
		return err
	}

	key, err := com.String("file", c)
	if err != nil {
		return err
	}
	filename := filepath.Base(key)
	region := c.String("region")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		xlog.LogDebug("%s", err)
		return cli.NewExitError(err, 1)
	}

	s3Client := s3.New(sess)
	downloader := s3manager.NewDownloader(sess)
	size, err := getFileSize(s3Client, bucket, key)
	if err != nil {
		xlog.LogDebug("%s", err)
		return cli.NewExitError(err, 1)
	}

	log.Println("Starting download, size:", byteCountDecimal(size))
	cwd, err := os.Getwd()
	if err != nil {
		xlog.LogDebug("%s", err)
		return cli.NewExitError(err, 1)
	}

	temp, err := ioutil.TempFile(cwd, "getObjWithProgress-tmp-")
	if err != nil {
		xlog.LogDebug("%s", err)
		return cli.NewExitError(err, 1)
	}
	tempfileName := temp.Name()

	writer := &progressWriter{writer: temp, size: size, written: 0}
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	if _, err := downloader.Download(writer, params); err != nil {
		log.Printf("Download failed! Deleting tempfile: %s", tempfileName)
		os.Remove(tempfileName)
		xlog.LogDebug("%s", err)
		return cli.NewExitError(err, 1)
	}

	if err := temp.Close(); err != nil {
		xlog.LogDebug("%s", err)
		return cli.NewExitError(err, 1)
	}

	if err := os.Rename(temp.Name(), filename); err != nil {
		xlog.LogDebug("%s", err)
		return cli.NewExitError(err, 1)
	}

	fmt.Println()
	log.Println("File downloaded! Avaliable at:", filename)
	return nil
}
