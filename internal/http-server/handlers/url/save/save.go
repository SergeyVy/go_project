package save

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
	_ "golang.org/x/exp/slog"
	"net/http"
	_ "net/http"
	"url-shorter/internal/lib/api/response"
	"url-shorter/internal/lib/random"
)

type Request struct {
	URL   string `json:"url" validate:"required"`
	Alias string `json:"alias,omitempty"`
}
type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLength = 6

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.save.url.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to decode request", slog.String("error", err.Error()))

			render.JSON(w, r, response.Error("Failed to decode request"))

			return

		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("Failed to validate request", slog.String("error", err.Error()))

			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}
		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias) // было: erc
		if err != nil {
			log.Info("Failed to save url", slog.String("url", req.URL))
			render.JSON(w, r, response.Error("url already exists"))

			return
		}
		if err != nil {
			log.Error("Failed to save url", slog.String("url", req.URL))
			render.JSON(w, r, response.Error("failed to add url"))

			return
		}
		log.Info("url saved", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: response.OK(),
			Alias:    alias,
		})
	}
}
