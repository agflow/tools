package sql

import (
	"embed"
	"io/fs"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/pkg/errors"

	"github.com/agflow/tools/log"
)

// mustReadQueries is an utility method for MustReadSQL for sql files in queries subdirectory
// Reads given folder and expects a SQL file named for each field of destination.
// Then sets the content of file into that field.
func mustReadQueries(v reflect.Value, queriesFS embed.FS, subfolder string, base ...string) {
	dir := subfolder
	if len(base) > 0 {
		dir = filepath.Join(base[0], subfolder)
	}

	files, err := fs.ReadDir(queriesFS, dir)
	if err != nil {
		panic(err)
	}

	got := make([]string, len(files))
	for i, f := range files {
		got[i] = strings.TrimSuffix(f.Name(), ".sql")
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	names := make([]string, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		names[i] = v.Type().Field(i).Name
	}
	setQueries(v, queriesFS, dir)
}

func setQueries(v reflect.Value, queriesFS embed.FS, dir string) {
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		f := filepath.Join(dir, field.Name+".sql")
		query, err := fs.ReadFile(queriesFS, filepath.Clean(f))
		if err != nil {
			log.Fatal(errors.Wrapf(err, "can't read file %q", f))
		}
		v.FieldByName(field.Name).SetString(string(query))
	}
}

// mustReadSQLFiles fills queries into Query variable.
func mustReadSQLFiles(root string, queriesFS embed.FS, query interface{}) {
	q := reflect.ValueOf(query).Elem()
	v := reflect.Indirect(reflect.ValueOf(query))
	for i := 0; i < v.NumField(); i++ {
		f := v.Type().Field(i)
		mustReadQueries(q.FieldByName(f.Name), queriesFS,
			filepath.Join(root, f.Tag.Get("sql")))
	}
}

// MustLoadSQLQueries enables pgu to use queries
// that are located in the default sql directory
func MustLoadSQLQueries(queriesFS embed.FS, query interface{}) {
	f, err := fs.ReadDir(queriesFS, ".")
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range f {
		mustReadSQLFiles(r.Name(), queriesFS, query)
	}
}
