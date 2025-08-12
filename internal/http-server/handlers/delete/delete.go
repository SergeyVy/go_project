// internal/http-server/handlers/url/delete/delete.go
package delete

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"

	"url-shorter/internal/storage"
)

type Deleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, d Deleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(slog.String("op", "handlers.url.delete.New"))

		alias := chi.URLParam(r, "alias") // или r.PathValue("alias") если используешь Chi 5.1+
		if alias == "" {
			http.Error(w, "alias is required", http.StatusBadRequest)
			return
		}

		if err := d.DeleteURL(alias); err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				http.Error(w, "alias not found", http.StatusNotFound)
				return
			}
			log.Error("failed to delete", slog.String("alias", alias), slog.Any("err", err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent) // 204
	}
}
