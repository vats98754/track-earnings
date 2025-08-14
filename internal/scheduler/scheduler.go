package scheduler

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

type Job func(ctx context.Context) error

type Scheduler struct {
	jobs []scheduledJob
}

type scheduledJob struct {
	name string
	fn   Job
	intv time.Duration
}

func New() *Scheduler { return &Scheduler{} }

func (s *Scheduler) Every(name string, d time.Duration, fn Job) {
	s.jobs = append(s.jobs, scheduledJob{name: name, fn: fn, intv: d})
}

func (s *Scheduler) Start(ctx context.Context) {
	for _, j := range s.jobs {
		j := j
		go func() {
			t := time.NewTicker(j.intv)
			defer t.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-t.C:
					if err := j.fn(ctx); err != nil {
						log.Error().Str("job", j.name).Err(err).Msg("job error")
					}
				}
			}
		}()
	}
}
