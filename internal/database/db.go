package database

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"

	"github.com/rprtr258/xerr"
)

type entity interface {
	ID() string
}

func entityName[E entity]() string {
	var e E
	return reflect.TypeOf(e).Name()
}

func read[E entity](r DB, filter func(E) bool) ([]E, error) {
	bytes, err := os.ReadFile(filepath.Join(r.dir, entityName[E]()))
	if err != nil {
		return nil, xerr.NewWM(err, "can't open table file",
			xerr.Field("entity", entityName[E]()))
	}

	var all []E
	if err := json.Unmarshal(bytes, &all); err != nil {
		return nil, xerr.NewW(err)
	}

	res := make([]E, 0, len(all))
	for _, entity := range all {
		if filter(entity) {
			res = append(res, entity)
		}
	}

	return res, nil
}

func write[E entity](r DB, entities []E) error {
	bytes, err := json.Marshal(entities)
	if err != nil {
		return xerr.NewW(err)
	}

	if err := os.WriteFile(filepath.Join(r.dir, entityName[E]()), bytes, 0644); err != nil {
		return xerr.NewW(err)
	}

	return nil
}
