package taskrunner

import (
	"log/slog"
	"sync"

	taskRunnerDomainTaskRunner "github.com/blocknextai/platform-api/internal/taskrunner/domain/taskrunner"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

type scheduledJob struct {
	entryID cron.EntryID
	version string
}

type CronService struct {
	cron          *cron.Cron
	triggerJobMap map[uuid.UUID]scheduledJob
	mutex         sync.Mutex
}

func NewCronService() taskRunnerDomainTaskRunner.CronService {
	cronInstance := cron.New(
		cron.WithSeconds(),
		cron.WithParser(cron.NewParser(
			cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow|cron.Descriptor,
		)),
	)

	return &CronService{
		cron:          cronInstance,
		triggerJobMap: make(map[uuid.UUID]scheduledJob),
	}
}

func (cs *CronService) Start() {
	cs.cron.Start()
	slog.Info("Cron service started", "component", "CronService")
}

func (cs *CronService) Stop() error {
	ctx := cs.cron.Stop()
	<-ctx.Done()
	return nil
}

func (cs *CronService) ReconcileJobs(desired []taskRunnerDomainTaskRunner.DesiredJob) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	desiredByID := make(map[uuid.UUID]taskRunnerDomainTaskRunner.DesiredJob, len(desired))
	for _, job := range desired {
		desiredByID[job.TriggerID] = job
	}

	for triggerID, scheduled := range cs.triggerJobMap {
		if _, ok := desiredByID[triggerID]; ok {
			continue
		}
		cs.cron.Remove(scheduled.entryID)
		delete(cs.triggerJobMap, triggerID)
		slog.Info("removed cron job for trigger",
			"component", "CronService",
			"trigger_id", triggerID)
	}

	for _, job := range desired {
		existing, ok := cs.triggerJobMap[job.TriggerID]
		if ok && existing.version == job.Version {
			continue
		}

		if ok {
			cs.cron.Remove(existing.entryID)
			delete(cs.triggerJobMap, job.TriggerID)
		}

		entryID, err := cs.cron.AddFunc(job.Pattern, job.Run)
		if err != nil {
			slog.Error("failed to schedule cron job for trigger",
				"component", "CronService",
				"trigger_id", job.TriggerID,
				"cron_pattern", job.Pattern,
				"error", err)
			continue
		}

		cs.triggerJobMap[job.TriggerID] = scheduledJob{entryID: entryID, version: job.Version}
		slog.Info("scheduled cron job for trigger",
			"component", "CronService",
			"trigger_id", job.TriggerID,
			"cron_pattern", job.Pattern,
			"entry_id", entryID)
	}
}
