package taskrunner

import (
	"github.com/google/uuid"
)

type DesiredJob struct {
	TriggerID uuid.UUID
	Pattern   string
	Version   string
	Run       func()
}

type CronService interface {
	Start()
	Stop() error
	ReconcileJobs(desired []DesiredJob)
}
