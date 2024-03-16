package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/johnldev/4-deploy-cloud-run/internal/utils"
)

type ICepService interface {
	GetCep(zipcode string) (string, error)
}
type cepService struct {
	ctx context.Context
}

func (s cepService) GetCep(zipcode string) (string, error) {

	var cepForCdn string = zipcode[:5] + "-" + zipcode[5:]
	var city string
	cdnUrl := fmt.Sprintf("https://cdn.apicep.com/file/apicep/%s.json", cepForCdn)
	viaCepUrl := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", strings.Replace(zipcode, "-", "", 1))

	resultViaCep := make(chan []byte)
	resultCdn := make(chan []byte)
	defer close(resultViaCep)
	defer close(resultCdn)

	go func() {
		defer utils.PanicRecovery()
		response, _ := utils.RequestWithContext(s.ctx, cdnUrl)
		resultCdn <- response
	}()

	go func() {
		defer utils.PanicRecovery()
		response, _ := utils.RequestWithContext(s.ctx, viaCepUrl)
		resultViaCep <- response
	}()

	select {
	case result := <-resultCdn:
		fmt.Println("Response from cdn")
		response := struct {
			City string `json:"city"`
		}{}
		json.Unmarshal(result, &response)
		city = response.City
	case result := <-resultViaCep:
		fmt.Println("Response from viacep")
		response := struct {
			Localidade string `json:"localidade"`
		}{}
		json.Unmarshal(result, &response)
		city = response.Localidade
	case <-s.ctx.Done():
		fmt.Println("Timeout on request")
		fmt.Println(s.ctx.Err())
	}

	if city == "" {
		return "", fmt.Errorf("can not find zipcode")
	}

	return city, nil
}

func NewCepService(ctx context.Context) cepService {
	return cepService{ctx: ctx}
}
