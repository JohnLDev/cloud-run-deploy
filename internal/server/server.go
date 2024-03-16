package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/johnldev/4-deploy-cloud-run/internal/services"
	usecases "github.com/johnldev/4-deploy-cloud-run/internal/useCases"

	"github.com/go-playground/validator/v10"
)

type Input struct {
	Zipcode string `json:"zipcode" validate:"required,len=8,number"`
}

func StartServer() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		validate := validator.New(validator.WithRequiredStructEnabled())
		input := Input{}
		// Decode the request body into the input struct
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		validateErr := validate.Struct(input)
		if validateErr != nil {
			http.Error(w, "invalid zipcode", http.StatusBadRequest)
			return
		}

		response, err := usecases.NewGetTemperatureUseCase(services.NewCepService(r.Context()), services.NewWeatherService(r.Context())).Execute(r.Context(), input.Zipcode)
		if err != nil {
			if err.Error() == "can not find zipcode" {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"temp_C": %.2f, "temp_F": %.2f, "temp_K": %.2f}`, response.Celcius, response.Fahrenheit, response.Kelvin)))
	})
	slog.Info("Server started at :8080")
	http.ListenAndServe(":8080", r)
}
