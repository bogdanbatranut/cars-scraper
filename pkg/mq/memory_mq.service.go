package mq

import "carscraper/pkg/jobs"

type MemMQ struct {
	Jobs []jobs.SessionJob
}

func NewMemMQ() MemMQ {
	return MemMQ{}
}

func (mmq MemMQ) GetAJob() *jobs.SessionJob {
	j := mmq.Jobs[0]
	mmq.Jobs = mmq.Jobs[1:]
	return &j
}

func (mmq MemMQ) PutJob(j jobs.SessionJob) {
	mmq.Jobs = append(mmq.Jobs, j)
}

func (mmq MemMQ) PutJobs(js []jobs.SessionJob) {
	mmq.Jobs = append(mmq.Jobs, js...)
}
