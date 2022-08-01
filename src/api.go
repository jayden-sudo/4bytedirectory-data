package src

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type (
	FourByteApi struct {
		baseUrl string
		client  *resty.Client
	}

	FourBytePageResponse struct {
		Results []FourByteSignature
	}

	FourByteSignature struct {
		Text string `json:"text_signature"`
		Hex  string `json:"hex_signature"`
	}
)

func NewFourByteApi() FourByteApi {
	return FourByteApi{
		baseUrl: "https://www.4byte.directory/api/v1",
		client:  resty.New(),
	}
}
func (a *FourByteApi) FetchPage(page int) (sigs []FourByteSignature, err error) {
	url := fmt.Sprintf("%s/signatures/?page=%d", a.baseUrl, page)

	resp, err := a.client.R().SetResult(FourBytePageResponse{}).Get(url)

	if err != nil {
		return sigs, err
	}

	if resp.StatusCode() != 200 {
		return sigs, errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode()))
	}

	return resp.Result().(*FourBytePageResponse).Results, nil
}
