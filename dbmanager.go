package main

type DatabaseManager struct {
	DBS []*Database
}

func (m DatabaseManager) GetDB(id int) *Database {
	return m.DBS[id%len(m.DBS)]
}

/*func (m DatabaseManager) GetVisitsFilter(id, filters) ([]UserVisits, error) {
	return []UserVisits{}, nil
}*/
/*
func NewDatabaseManager(count int) (*DatabaseManager, error) {
	m := DatabaseManager{
		DBS: make(Database, count),
	}
	for i := 0; i < COUNT_DB; i++ {
		DB, err := DatabaseInit()
		if err != nil {
			return nil, err
		}
		m.DBS[i] = DB
	}

	return &m, nil
}
*/
