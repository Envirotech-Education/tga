package tga

import ()

type TGA struct {
	Endpoint     string // https://ws.sandbox.training.gov.au/Deewr.Tga.Webservices/
	username     string // WebService.Read
	password     string
	lastResponse string
}

func (tga *TGA) LastSoapResponse() string {
	return tga.lastResponse
}

func NewTGA(endpoint, username, password string) *TGA {
	return &TGA{endpoint, username, password, ""}

}
