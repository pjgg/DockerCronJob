# DockerCronJob

docker pull pjgg/dockercronjob:0.0.1-SNAPSHOT

docker run -e "CRON_EXP=@every 30s" -e "COMMAND=gradle" -e "ARG=-version" -w "/usr/local/bin/" pjgg/dockercronjob:0.0.1-SNAPSHOT dockerCronJob