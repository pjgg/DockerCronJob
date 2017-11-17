package main

import (
	"os"
	"os/exec"
	"os/signal"

	"github.com/Sirupsen/logrus"
	"github.com/robfig/cron"
)

func main() {
	command := os.Getenv("COMMAND")
	arg := os.Getenv("ARG")
	cronExpr := os.Getenv("CRON_EXP")

	if isEmpty(command, cronExpr) {
		logrus.Fatal("missing required env variables as COMMAND, ARG, CRON_EXP")
	}

	logrus.Info("Cron job up!: " + cronExpr)
	job := &commandJob{command: command, arg: arg}

	c := cron.New()
	c.AddJob(cronExpr, job)
	go c.Start()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}

type commandJob struct {
	command string
	arg     string
}

func (self commandJob) Run() {
	cmd := exec.Command(self.command, self.arg)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Error("unexpected error " + err.Error())
	} else {
		logrus.Info(string(stdoutStderr[:]))
	}

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
