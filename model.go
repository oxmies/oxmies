package oxmies

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	adapter "github.com/oxmies/oxmies/adapters"
)

// Model is the base struct for all models.
// Embed this in your concrete structs to get simple CRUD methods:
//
//	type User struct {
//	    oxmies.Model
//	    ID   int    `orm:"primary_key,column:id"`
//	    Name string `orm:"column:name"`
//	}
//
// It looks up the adapter by the embedded Model.DBKey (defaults to "default").
type Model struct {
	ResourceName string `orm:"-"`
	DBKey        string `orm:"-"` // optional: if you have multiple DB connections
}

// helper to resolve adapter via connectionManager using model registry
func adapterForModel(m any) (adapter.DBAdapter, error) {
	if m == nil {
		return nil, errors.New("oxmies: nil model")
	}

	t := reflect.TypeOf(m)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	name := t.Name()

	modelMeta, ok := registry.models[name]
	if !ok {
		return nil, errors.New("oxmies: model not registered")
	}

	ad, ok := connectionManager.GetConnection(modelMeta.AdapterName, modelMeta.ConnectionName)
	if !ok {
		return nil, errors.New("oxmies: adapter not found for model")
	}

	return ad, nil
}

// Insert inserts the current model. The concrete adapter is responsible for
// any DB-specific behavior (returning ids, handling zero-values, etc.).
func (m *Model) Insert(ctx context.Context) error {
	if m == nil {
		return errors.New("oxmies: cannot Insert on nil model")
	}

	ad, err := adapterForModel(m)
	if err != nil {
		return err
	}
	return ad.Insert(ctx, m)
}

// Update updates the current model. Returns an error when called on nil.
func (m *Model) Update(ctx context.Context) error {
	if m == nil {
		return errors.New("oxmies: cannot Update on nil model")
	}

	ad, err := adapterForModel(m)
	if err != nil {
		return err
	}
	return ad.Update(ctx, m)

// FindByID populates the model with the record for the given id.
func (m *Model) FindByID(ctx context.Context, id any) error {
	if m == nil {
		return errors.New("oxmies: cannot FindByID on nil model")
	}

	ad, err := adapterForModel(m)
	ad, err := adapterForModel(m)
	if err != nil {
		return err
	}

	return ad.FindByID(ctx, m, id)
// Delete removes the record identified by the model's primary key.
func (m *Model) Delete(ctx context.Context) error {
	if m == nil {
		return errors.New("oxmies: cannot Delete on nil model")
	}

	ad, err := adapterForModel(m)
	ad, err := adapterForModel(m)
	ad, err := adapterForModel(m)
	if err != nil {
		return err
	}
	return ad.Delete(ctx, m)
