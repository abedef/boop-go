package boop

import (
	"context"
	"encoding/json"
	"junk/boop-server/pgdb"
	"net/http"
	"time"
)

type Bean struct {
	Name  string
	Value float32
}

type BeanSummary struct {
	Totals []Bean
	Boops  []string
}

type Boop struct {
	ID      int32
	Text    string
	Created time.Time
}

func MakeHttpHandler(fn MyHandlerFunc, ctx context.Context, queries *pgdb.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := MyResponseWriter{w: w}
		err := fn(&rw, r, ctx, queries)
		if err != nil {
			// Note that this is now the only place where we use
			// the original writer's error reporting capabilities
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

type MyResponseWriter struct {
	w http.ResponseWriter
}

type MyHandlerFunc func(w *MyResponseWriter, r *http.Request, ctx context.Context, queries *pgdb.Queries) error

func (rw *MyResponseWriter) WriteHeader(statusCode int) error {
	rw.w.WriteHeader(statusCode)
	return nil
}
func (rw *MyResponseWriter) WriteString(str string) error {
	_, err := rw.w.Write([]byte(str))
	return err
}
func (rw *MyResponseWriter) WriteJSON(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = rw.w.Write(data)
	return err
}
