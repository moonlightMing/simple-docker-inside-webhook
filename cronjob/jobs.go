package cronjob

import "fmt"

const (
	CleanNoneTagImage = "CLEAN_NONE_TAG_IMAGE"
)

type CleanNoneTagImageJob struct {
}

func (c CleanNoneTagImageJob) Run()  {
	fmt.Printf("job 'cleanNoneTagImage' start...\n")
}