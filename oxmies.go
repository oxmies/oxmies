package oxmies

import (
	"fmt"
)

/**
* Initialize the ORM with the given configuration.
*/
func Initialize(cfg map[string]any){
	if cfg == nil {
		panic("oxmies: configuration cannot be nil")
	}
	for connectionName, config := range cfg {
		switch c := config.(type) {
		case SQLConfig:
			// Passing the config to InitSQL function with connection name
			InitSQL(c, connectionName)
		case RedisConfig:
			// Placeholder: Initialize Redis adapter here
		// Add other adapters as needed
		default:
			panic(fmt.Sprintf("Unsupported config type for adapter '%s'", connectionName))
		}
	}
}