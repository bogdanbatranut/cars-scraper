package logging

import "carscraper/pkg/adsdb"

func (repo LogsRepository) GetSessions() (*[]adsdb.SessionLog, error) {
	var sessions []adsdb.SessionLog
	tx := repo.db.Find(&sessions)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &sessions, nil
}

func (repo LogsRepository) GetSession(sessionID uint) (*adsdb.SessionLog, error) {
	var session adsdb.SessionLog
	tx := repo.db.Preload("CriteriaLogs.PageLogs").Find(&session, sessionID)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &session, nil
}

func (repo LogsRepository) GetPageLogsForCriteriaLogs(criteriaLogID uint) (*[]adsdb.PageLog, error) {
	var pageLogs []adsdb.PageLog
	tx := repo.db.Preload("CriteriaLog").Where("criteria_log_id = ?", criteriaLogID).Find(&pageLogs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &pageLogs, nil
}

func (repo LogsRepository) DeleteSession(sessionID uint) error {
	var existingSessionLog adsdb.SessionLog
	tx := repo.db.Preload("CriteriaLogs").Preload("PageLogs").Find(&existingSessionLog, sessionID)
	if tx.Error != nil {
		return tx.Error
	}

	pageLogs := existingSessionLog.PageLogs
	tx = repo.db.Delete(&pageLogs)
	if tx.Error != nil {
		return tx.Error
	}

	criterialLogs := existingSessionLog.CriteriaLogs
	tx = repo.db.Delete(&criterialLogs)
	if tx.Error != nil {
		return tx.Error
	}

	tx = repo.db.Delete(&existingSessionLog, sessionID)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
