package oxmies

import (
	"fmt"

	adapter "github.com/oxmies/oxmies/adapters"
)

var connectionManager *ConnectionManager

func init() {
	connectionManager = NewConnectionManager()
}

// GetManager returns the global DB manager.
func GetManager() *ConnectionManager {
	return connectionManager
}


type ConnectionManager struct {
	connections map[adapter.AdapterType]map[string]adapter.DBAdapter // adapterType -> key -> adapter
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[adapter.AdapterType]map[string]adapter.DBAdapter),
	}
}

func (m *ConnectionManager) Register(adapterType adapter.AdapterType, connectionName string, db adapter.DBAdapter) {
	if m.connections[adapterType] == nil {
		m.connections[adapterType] = make(map[string]adapter.DBAdapter)
	}
	m.connections[adapterType][connectionName] = db
}

func (m *ConnectionManager) GetConnection(adapterType adapter.AdapterType, connectionName string) (adapter.DBAdapter, bool) {
	db, ok := m.connections[adapterType][connectionName]
	return db, ok
}

// GetDB fetches a DBAdapter by key (or default).
func GetDB(connectionName string) (adapter.DBAdapter, error) {
	if connectionName == "" {
		connectionName = "default"
	}
	for _, conns := range connectionManager.connections {
		if db, ok := conns[connectionName]; ok {
			return db, nil
		}
	}
	return nil, fmt.Errorf("oxmies: connection '%s' not found", connectionName)
}
