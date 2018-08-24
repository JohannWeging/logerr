package logerr_test

import (
	"os"

	"github.com/JohannWeging/logerr"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
)

func Example() {
	// make logrus play nice with go test
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{DisableTimestamp: true})

	err := job()
	if err != nil {
		fields := logerr.GetFields(err)
		log.WithFields(fields).Errorf("job failed: %s", err)
	}
	// Output: level=error msg="job failed: task failed: cause" jobID=0
}

func job() error {
	err := task()
	if err != nil {
		err = logerr.WithField(err, "jobID", "0")
		err = errors.Annotate(err, "task failed")
		return err
	}
	return nil
}

func task() error {
	return errors.New("cause")
}
