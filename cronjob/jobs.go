package cronjob

import "github.com/moonlightming/simple-docker-inside-webhook/dockercli"

const (
    CleanNoneTagImage = "CLEAN_NONE_TAG_IMAGE"
)

type CleanNoneTagImageJob struct {
}

func (c CleanNoneTagImageJob) Run() {
    dockercli.CleanNoneTagImage()
}
