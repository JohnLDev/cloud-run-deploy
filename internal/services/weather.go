package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/johnldev/4-deploy-cloud-run/internal/utils"
)

const (
	weatherApiUrl = "http://api.weatherapi.com/v1/current.json"
)

type ITemperatureResponse struct {
	Current struct {
		Celcius    float64 `json:"temp_c"`
		Fahrenheit float64 `json:"temp_f"`
	} `json:"current"`
}

type IWeatherService interface {
	GetTemperatureByCity(city string) (*ITemperatureResponse, error)
}

type WeatherService struct {
	Ctx context.Context
}

func (s WeatherService) GetTemperatureByCity(city string) (*ITemperatureResponse, error) {
	token := os.Getenv("WEATHER_API_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("weather api token not found")
	}

	result := ITemperatureResponse{}

	response, err := utils.RequestWithContext(s.Ctx, fmt.Sprintf("%s?q=%s&key=%s", weatherApiUrl, city, token))
	if err != nil {
		return nil, err
	}

	json.Unmarshal(response, &result)
	return &result, nil
}

func NewWeatherService(ctx context.Context) WeatherService {
	return WeatherService{
		Ctx: ctx,
	}
}
