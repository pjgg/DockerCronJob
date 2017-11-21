package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/pjgg/DockerCronJob/connectors"
	"github.com/robfig/cron"
)

func main() {
	command := os.Getenv("COMMAND")
	arg := strings.Fields(os.Getenv("ARG"))
	cronExpr := os.Getenv("CRON_EXP")

	if isEmpty(command, cronExpr) {
		logrus.Fatal("missing required env variables as COMMAND, ARG, CRON_EXP")
	}

	logrus.Info("Cron job up!: " + cronExpr)

	var bucketName string
	var destPath string
	var sourcePath string
	var ok bool
	if bucketName, ok = os.LookupEnv("BUCKET_NAME"); !ok {
		logrus.Fatal("missing required env variables BUCKET_NAME")
	}

	if destPath, ok = os.LookupEnv("S3_DEST_PATH"); !ok {
		logrus.Fatal("missing required env variables S3_DEST_PATH")
	}

	if sourcePath, ok = os.LookupEnv("S3_SOURCE_PATH"); !ok {
		logrus.Fatal("missing required env variables S3_SOURCE_PATH")
	}

	s3Client := connectors.Instance(bucketName)
	job := &commandJob{command: command, arg: arg, s3Client: s3Client, destFolder: destPath, sourceFolder: sourcePath}

	c := cron.New()
	c.AddJob(cronExpr, job)
	go c.Start()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}

type commandJob struct {
	command      string
	arg          []string
	s3Client     connectors.AwsConnectorBehavior
	destFolder   string
	sourceFolder string
}

func (self commandJob) Run() {
	cmd := exec.Command(self.command, self.arg...)

	// open the out file for writing
	outfile, err := os.Create("/tmp/out.log")
	if err != nil {
		panic(err)
	}
	defer outfile.Close()

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	writer := bufio.NewWriter(outfile)

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	go io.Copy(writer, stdoutPipe)
	cmd.Wait()
	fmt.Printf("End.")

	self.s3Client.Push(self.destFolder, self.sourceFolder)
}

func isEmpty(params ...string) (empty bool) {

	for _, value := range params {
		if value == "" {
			empty = true
			break
		}
	}

	return
}
