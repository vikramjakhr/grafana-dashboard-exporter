package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"archive/zip"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/vikramjakhr/grafana-dashboard-exporter"
	"github.com/vikramjakhr/grafana-dashboard-exporter/plugins/outputs"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type S3 struct {
	Bucket       string `toml:"bucket"`
	AccessKey    string `toml:"access_key"`
	SecretKey    string `toml:"secret_key"`
	Region       string `toml:"region"`
	BucketPrefix string `toml:"bucket_prefix"`
	OutputFormat string `toml:"output_format"`
}

var sampleConfig = `
  bucket = "<bucket-name>" # required
  access_key = "$ACCESS_KEY" # required
  secret_key = "$SECRET_KEY" # required
  region = "s3-bucket-region"
  bucketPrefix = "<prefix>"
  output_format = "zip" # zip, dir
`

func (f *S3) SampleConfig() string {
	return sampleConfig
}

func (f *S3) Connect() error {
	of := strings.Trim(f.OutputFormat, " ")
	if !(strings.EqualFold(of, "dir") || strings.EqualFold(of, "zip")) {
		return errors.New("E! S3 output_format can only be 'file' or 'zip' only")
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(f.Region),
		Credentials: credentials.NewStaticCredentials(f.AccessKey, f.SecretKey, ""),
	})

	// Create S3 service client
	svc := s3.New(sess)
	_, err = svc.ListObjects(&s3.ListObjectsInput{Bucket: aws.String(f.Bucket)})
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to list items in bucket %q, %v", f.Bucket, err))
	}
	return nil
}

func (f *S3) Description() string {
	return "Send grafana json to s3"
}

func (f *S3) Write(metric gde.Metric) error {
	if metric.Action() != "" {
		dir := "/tmp/gde"

		baseDir := fmt.Sprintf("%s/%s", dir, metric.Dir())
		dir = fmt.Sprintf("%s/%ss/", baseDir, string(metric.Type()))

		switch metric.Action() {
		case gde.ActionCreate:

			if _, err := os.Stat(dir); os.IsNotExist(err) {
				err = os.MkdirAll(dir, 0774)
				if err != nil {
					log.Printf("E! Unable to create direcotry. %v", err)
					return err
				}
			}

			switch metric.Type() {
			case gde.TypeDatasource:
				filename := fmt.Sprintf("%s%s.json", dir, strings.Replace(metric.Title(), " ", "", -1))
				err := ioutil.WriteFile(filename, metric.Content(), 0644)
				if err != nil {
					log.Printf("E! Unable to create file. %v", err)
					return err
				}
				break
			case gde.TypeDashboard:
				filename := fmt.Sprintf("%s%s.json", dir, strings.Replace(metric.Title(), " ", "", -1))
				err := ioutil.WriteFile(filename, metric.Content(), 0644)
				if err != nil {
					log.Printf("E! Unable to create file. %v", err)
					return err
				}
				break
			}

			break
		case gde.ActionFinish:
			if strings.EqualFold(f.OutputFormat, "zip") {
				zipFileName := fmt.Sprintf("%s.zip", baseDir)
				err := zipit(baseDir, zipFileName)
				if err != nil {
					removeDir(dir)
					log.Printf("E! Unable to create zip file. %v", err)
					return err
				} else {
					sess, err := f.makeSession()
					if err != nil {
						removeDir(dir)
						return errors.New(fmt.Sprintf("E! failed to create aws session, %v", err))
					}
					err = uploadFileToS3(sess, f.Bucket, f.BucketPrefix, zipFileName)
					if err != nil {
						removeDir(dir)
						return errors.New(fmt.Sprintf("E! Failed to upload data to %s/%s, %s\n",
							f.Bucket, zipFileName, err))
					}
					log.Printf("D! %s uploaded to s3", zipFileName)
					removeDir(dir)
				}
			}
			if strings.EqualFold(f.OutputFormat, "dir") {
				sess, err := f.makeSession()
				if err != nil {
					removeDir(dir)
					return errors.New(fmt.Sprintf("E! failed to create aws session, %v", err))
				}
				err = uploadDirToS3(sess, f.Bucket, f.BucketPrefix, baseDir)
				if err != nil {
					removeDir(dir)
					return errors.New(fmt.Sprintf("E! Failed to upload data to %s/%s, %s\n",
						f.Bucket, baseDir, err))
				}
				log.Printf("D! %s uploaded to s3", baseDir)
				removeDir(dir)
			}
			break
		}
	}
	return nil
}

func (f *S3) makeSession() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(f.Region),
		Credentials: credentials.NewStaticCredentials(f.AccessKey, f.SecretKey, ""),
	})
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func uploadDirToS3(sess *session.Session, bucketName string, bucketPrefix string, dirPath string) error {
	fileList := []string{}
	filepath.Walk(dirPath, func(path string, f os.FileInfo, err error) error {
		if isDirectory(path) {
			// Do nothing
			return nil
		} else {
			fileList = append(fileList, path)
			return nil
		}
	})

	for _, file := range fileList {
		err := uploadFileToS3(sess, bucketName, bucketPrefix, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func isDirectory(path string) bool {
	fd, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	switch mode := fd.Mode(); {
	case mode.IsDir():
		return true
	case mode.IsRegular():
		return false
	}
	return false
}

func uploadFileToS3(sess *session.Session, bucketName string, bucketPrefix string, filePath string) error {
	log.Printf("D! uploading %s to S3", filePath)
	// An s3 service
	s3Svc := s3.New(sess)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var key string
	fileDirectory, _ := filepath.Abs(filePath)
	key = bucketPrefix + "/" + strings.TrimPrefix(fileDirectory, "/tmp/gde")

	// Upload the file to the s3 given bucket
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName), // Required
		Key:    aws.String(key),        // Required
		Body:   file,
	}
	_, err = s3Svc.PutObject(params)
	if err != nil {
		return err
	}
	return nil
}

func removeDir(dir string) {
	log.Printf("D! Clearing the directory: %s", dir)
	err := os.RemoveAll(dir)
	if err != nil {
		log.Printf("E! Unable to remove directory: %s. %v", dir, err)
	}
}

func zipit(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

func init() {
	outputs.Add("s3", func() gde.Output {
		return &S3{}
	})
}
