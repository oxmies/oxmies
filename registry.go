package oxmies

import (
	"database/sql"
	"fmt"
	"reflect"
	"sync"

	adapter "github.com/oxmies/oxmies/adapters"
	adaptersql "github.com/oxmies/oxmies/adapters/sql"
)

func InitSQL(cfg SQLConfig, connectionName string) error {
	dsn := cfg.DSN()
	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	ad := adaptersql.NewSQLAdapter(db, cfg.Debug)
	if cfg.Debug {
		fmt.Printf("✅ SQL Adapter '%s' initialized\n", connectionName)
	}

	GetManager().Register(adapter.SQL, connectionName, ad)

	for _, model := range cfg.OxmiesDbConfig.Models {
		// TODO: Auto-migrate logic here

		// Get the Model field using reflection
		val := reflect.ValueOf(model)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		// Look for embedded Model field
		modelField := val.FieldByName("Model")
		if !modelField.IsValid() {
			fmt.Printf("⚠️  oxmies: Model registration skipped: no embedded Model field found %+v\n", model)
			continue
		}

		baseModel := modelField.Interface().(Model)
		fmt.Printf("")
		if baseModel.ResourceName == "" {
			// If ResourceName is not set, use pluralized struct name
			baseModel.ResourceName = getTableName(val.Type())
			// Update the ResourceName in the original struct
			modelField.Set(reflect.ValueOf(baseModel))
		}

		RegisterModel(adapter.SQL, baseModel.ResourceName, model)
		fmt.Printf("ℹ️  Model registered: ResourceName=%s\n", baseModel.ResourceName)
	}
	return nil
}
type ModelMeta struct {
	ResourceName   string
	AdapterName    adapter.AdapterType
	ConnectionName string
}

type ModelRegistry struct {
	models map[string]ModelMeta
	mu     sync.RWMutex
}

var registry = &ModelRegistry{
	models: make(map[string]ModelMeta),
}

// RegisterModel registers a model under a specific resource name
func RegisterModel(adapterType adapter.AdapterType, resourceName string, model any) {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	registry.mu.Lock()
	defer registry.mu.Unlock()

	name := t.Name()
	registry.models[name] = ModelMeta{
		AdapterName:  adapterType,
		ResourceName: resourceName,
	}
}

// GetModelMeta fetches metadata by type name
func GetModelMeta(model any) (ModelMeta, error) {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	name := t.Name()

	registry.mu.RLock()
	defer registry.mu.RUnlock()

	meta, ok := registry.models[name]
	if !ok {
		return ModelMeta{}, fmt.Errorf("oxmies: model %s not registered", name)
	}
	return meta, nil
}

// derive table name
// NOTE: This is a naive pluralization. For proper pluralization, consider using a library or extend this logic.
func getTableName(t reflect.Type) string {
	name := t.Name()
	// Basic pluralization rules for demonstration
	if len(name) > 1 && name[len(name)-1] == 'y' && name[len(name)-2] not in []byte{'a', 'e', 'i', 'o', 'u'} {
		return fmt.Sprintf("%sies", name[:len(name)-1])
	}
	if len(name) > 1 && (name[len(name)-1] == 's' || name[len(name)-1] == 'x' || (len(name) > 2 && name[len(name)-2:] == "ch") || (len(name) > 2 && name[len(name)-2:] == "sh")) {
		return fmt.Sprintf("%ses", name)
	}
	return fmt.Sprintf("%ss", name) // fallback
}
