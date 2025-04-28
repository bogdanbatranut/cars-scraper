package results

type ProcessedSession struct {
	SessionID string
	Criterias *[]ProcessedCriteria
}

func NewProcessedSession(sessionID string) *ProcessedSession {
	pc := make([]ProcessedCriteria, 0)
	return &ProcessedSession{
		SessionID: sessionID,
		Criterias: &pc,
	}
}

func (ps *ProcessedSession) getCriteria(criteriaID uint) *ProcessedCriteria {
	for _, criteria := range *ps.Criterias {
		if criteria.CriteriaID == criteriaID {
			return &criteria
		}
	}
	pc := NewProcessedCriteria(criteriaID)
	*ps.Criterias = append(*ps.Criterias, *pc)
	return pc
}

func (ps ProcessedSession) isComplete() bool {
	if ps.Criterias == nil {
		return false
	}
	for _, criteria := range *ps.Criterias {
		if !criteria.isComplete() {
			return false
		}
	}
	return true
}
