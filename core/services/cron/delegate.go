package cron

import (
	"context"

	"github.com/pkg/errors"

	"github.com/goplugin/pluginv3.0/v2/core/logger"
	"github.com/goplugin/pluginv3.0/v2/core/services/job"
	"github.com/goplugin/pluginv3.0/v2/core/services/pipeline"
)

type Delegate struct {
	pipelineRunner pipeline.Runner
	lggr           logger.Logger
}

var _ job.Delegate = (*Delegate)(nil)

func NewDelegate(pipelineRunner pipeline.Runner, lggr logger.Logger) *Delegate {
	return &Delegate{
		pipelineRunner: pipelineRunner,
		lggr:           lggr,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.Cron
}

func (d *Delegate) BeforeJobCreated(spec job.Job)              {}
func (d *Delegate) AfterJobCreated(spec job.Job)               {}
func (d *Delegate) BeforeJobDeleted(spec job.Job)              {}
func (d *Delegate) OnDeleteJob(context.Context, job.Job) error { return nil }

// ServicesForSpec returns the scheduler to be used for running cron jobs
func (d *Delegate) ServicesForSpec(ctx context.Context, spec job.Job) (services []job.ServiceCtx, err error) {
	if spec.CronSpec == nil {
		return nil, errors.Errorf("services.Delegate expects a *jobSpec.CronSpec to be present, got %v", spec)
	}

	cron, err := NewCronFromJobSpec(spec, d.pipelineRunner, d.lggr)
	if err != nil {
		return nil, err
	}

	return []job.ServiceCtx{cron}, nil
}
