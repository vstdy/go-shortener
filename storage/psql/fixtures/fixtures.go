package fixtures

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"runtime"
	"text/template"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"

	"github.com/vstdy0/go-project/storage/psql/schema"
)

// Fixtures keeps all fixture objects.
type Fixtures struct {
	URLS schema.URLS
}

func (f *Fixtures) appendURL(obj interface{}) error {
	dbObj, ok := obj.(*schema.URL)
	if !ok {
		return fmt.Errorf("url: type assert failed: %T", obj)
	}
	if dbObj == nil {
		return fmt.Errorf("url: type assert failed: nil")
	}
	f.URLS = append(f.URLS, *dbObj)

	return nil
}

// LoadFixtures load fixtures to DB and returns DB objects aggregate.
func LoadFixtures(ctx context.Context, db *bun.DB) (Fixtures, error) {
	type fixturesAppender struct {
		id     string
		append func(obj interface{}) error
	}

	fixtureManager := dbfixture.New(
		db,
		dbfixture.WithTemplateFuncs(template.FuncMap{
			"now": func() string {
				return time.Now().UTC().Format(time.RFC3339Nano)
			},
			"uuid": func() uuid.UUID {
				return uuid.New()
			},
		}),
	)

	err := fixtureManager.Load(
		ctx,
		os.DirFS(getFixturesDir()),
		"urls.yaml",
	)
	if err != nil {
		return Fixtures{}, fmt.Errorf("loading fixtures: %w", err)
	}

	fixtures := Fixtures{}
	appenders := []fixturesAppender{
		{id: "URL.link_1", append: fixtures.appendURL},
		{id: "URL.link_2", append: fixtures.appendURL},
	}
	for _, appender := range appenders {
		obj, err := fixtureManager.Row(appender.id)
		if err != nil {
			return Fixtures{}, fmt.Errorf("reading fixtures row (%s): %w", appender.id, err)
		}
		if obj == nil {
			return Fixtures{}, fmt.Errorf("reading fixtures row (%s): nil", appender.id)
		}
		if err := appender.append(obj); err != nil {
			return Fixtures{}, fmt.Errorf("appending fixtures row (%s): %w", appender.id, err)
		}
	}

	return fixtures, nil
}

// getFixturesDir returns current file directory.
func getFixturesDir() string {
	_, filePath, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}

	return filepath.Dir(filePath)
}
