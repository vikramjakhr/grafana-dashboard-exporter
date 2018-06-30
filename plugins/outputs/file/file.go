package file

import (
	"github.com/vikramjakhr/grafana-dashboard-exporter/plugins/outputs"
	"github.com/vikramjakhr/grafana-dashboard-exporter"
	"fmt"
	"os"
	"strings"
	"io/ioutil"
	"log"
	"errors"
	"archive/zip"
	"path/filepath"
	"io"
)

type File struct {
	OutputDir    string `toml:"output_dir"`
	OutputFormat string `toml:"output_format"`
}

var sampleConfig = `
  output_dir = "<dir>" # default is /tmp/gde
  output_format = "zip" # zip, dir # default is zip
`

func (f *File) SampleConfig() string {
	return sampleConfig
}

func (f *File) Connect() error {
	if strings.Trim(f.OutputDir, " ") != "" {
		_, err := os.Stat(f.OutputDir)
		return err
	}
	of := strings.Trim(f.OutputFormat, " ")
	if !(strings.EqualFold(of, "file") || strings.EqualFold(of, "zip")) {
		return errors.New("E! File output_format can only be 'file' or 'zip' only")
	}
	return nil
}

func (f *File) Description() string {
	return "Send grafana json to specified directory"
}

func (f *File) Write(metric gde.Metric) error {
	if metric.Action() != "" {
		dir := "/tmp"

		if strings.Trim(f.OutputDir, " ") != "" {
			dir = strings.TrimRight(f.OutputDir, "/")
		}

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
		case gde.ActionZIP:
			fmt.Println(baseDir)
			if strings.EqualFold(f.OutputFormat, "zip") {
				err := zipit(baseDir, fmt.Sprintf("%s.zip", baseDir))
				if err != nil {
					log.Printf("E! Unable to create zip file. %v", err)
				}
				log.Printf("D! Clearing the temporary directory")
				removeDir(baseDir)
			}
			break
		}
	}
	return nil
}

func removeDir(dir string) {
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
	outputs.Add("file", func() gde.Output {
		return &File{}
	})
}
