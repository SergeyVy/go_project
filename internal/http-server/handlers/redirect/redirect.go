// internal/http-server/handlers/redirect/redirect.go
package redirect

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"

	"url-shorter/internal/storage"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, getter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(slog.String("op", "handlers.url.redirect.New"))

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			http.Error(w, "alias is required", http.StatusBadRequest)
			return
		}

		target, err := getter.GetURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				http.Error(w, "alias not found", http.StatusNotFound)
				return
			}
			log.Error("failed to get url", slog.String("alias", alias), slog.Any("err", err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, target, http.StatusFound) // 302
	}
}
