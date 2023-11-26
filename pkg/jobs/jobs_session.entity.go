package jobs

import "github.com/google/uuid"

type JobSession struct {
	SessionID uuid.UUID
	Jobs      []SessionJob
}
