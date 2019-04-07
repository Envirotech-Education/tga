package tga

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

type TGA struct {
	Endpoint string // https://ws.sandbox.training.gov.au/Deewr.Tga.Webservices/
	username string // WebService.Read
	password string
}

// CheckMoodlePassword completes the moodle signin form and checks the response from moodle for indicators of signin success or failure.
func (tga *TGA) GetDetails(code string) (*Organisation, error) {

	soapRequest := `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
 xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
 <soapenv:Header>
      <wsse:Security soapenv:mustUnderstand="1" xmlns:wsse="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd" xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">
         <wsse:UsernameToken wsu:Id="UsernameToken-1">
            <wsse:Username>` + tga.username + `</wsse:Username>
            <wsse:Password Type="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText">` + tga.password + `</wsse:Password>
    </wsse:UsernameToken>
      </wsse:Security>
    
 </soapenv:Header>
 <soapenv:Body>
                <GetDetails xmlns="http://training.gov.au/services/7/">
                        <request>
                                <Code>` + code + `</Code>
                                <IncludeLegacyData>false</IncludeLegacyData>
                                <InformationRequested>
                                        <ShowCodes>true</ShowCodes>
                                        <ShowContacts>true</ShowContacts>
                                        <ShowDataManagers>true</ShowDataManagers>
                                        <ShowExplicitScope>true</ShowExplicitScope>
                                        <ShowImplicitScope>true</ShowImplicitScope>
                                        <ShowLocations>true</ShowLocations>
                                        <ShowOrganisatoinRoles>true</ShowOrganisatoinRoles>
                                        <ShowRegistrationManagers>true</ShowRegistrationManagers>
                                        <ShowRegistrationPeriods>true</ShowRegistrationPeriods>
                                        <ShowResponsibleLegalPersons>true</ShowResponsibleLegalPersons>
                                        <ShowRestrictions>true</ShowRestrictions>
                                        <ShowRtoClassifications>true</ShowRtoClassifications>
                                        <ShowRtoDeliveryNotification>true</ShowRtoDeliveryNotification>
                                        <ShowTradingNames>true</ShowTradingNames>
                                        <ShowUrls>true</ShowUrls>
                                </InformationRequested>
                        </request>
                </GetDetails>
 </soapenv:Body>
</soapenv:Envelope>`

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}

	req, _ := http.NewRequest("POST", tga.Endpoint+"/OrganisationServiceV7.svc/Organisation", strings.NewReader(soapRequest))
	req.Header.Add("Content-Type", "text/xml")
	req.Header.Add("SOAPAction", "\"http://training.gov.au/services/7/IOrganisationService/GetDetails\"")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	//b := string(body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("HTTP Error: " + resp.Status)
	}

	respEnvelope := new(SOAPEnvelope)
	//respEnvelope.Body = SOAPBody{Content: response}
	err = xml.Unmarshal(body, respEnvelope)
	if err != nil {
		fmt.Println("unmarshal failed:", err)
		return nil, err
	}

	return respEnvelope.Body.GetDetailsResponse.GetDetailsResult, nil
}

type ActionOnEntity string

const (
	ActionOnEntityNone   ActionOnEntity = "None"
	ActionOnEntityUpdate ActionOnEntity = "Update"
	ActionOnEntityDelete ActionOnEntity = "Delete"
	ActionOnEntityAdd    ActionOnEntity = "Add"
)

type SOAPEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Body    SOAPBody
}

type SOAPBody struct {
	XMLName            xml.Name            `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	GetDetailsResponse *GetDetailsResponse `xml:",omitempty"`
}

type GetDetailsResponse struct {
	XMLName          xml.Name      `xml:"http://training.gov.au/services/7/ GetDetailsResponse"`
	GetDetailsResult *Organisation `xml:"GetDetailsResult,omitempty"`
}

type Organisation struct {
	Codes                   *ArrayOfOrganisationCode       `xml:"Codes,omitempty"`
	Contacts                *ArrayOfContact                `xml:"Contacts,omitempty"`
	CreatedDate             *DateTimeOffset                `xml:"CreatedDate,omitempty"`
	DataManagers            *ArrayOfDataManagerAssignment  `xml:"DataManagers,omitempty"`
	IsLegacyData            bool                           `xml:"IsLegacyData,omitempty"`
	Locations               *ArrayOfOrganisationLocation   `xml:"Locations,omitempty"`
	ResponsibleLegalPersons *ArrayOfResponsibleLegalPerson `xml:"ResponsibleLegalPersons,omitempty"`
	Roles                   *ArrayOfRole                   `xml:"Roles,omitempty"`
	TradingNames            *ArrayOfTradingName            `xml:"TradingNames,omitempty"`
	UpdatedDate             *DateTimeOffset                `xml:"UpdatedDate,omitempty"`
	Urls                    *ArrayOfUrl                    `xml:"Urls,omitempty"`
}

type ArrayOfOrganisationCode struct {
	OrganisationCode []*OrganisationCode `xml:"OrganisationCode,omitempty"`
}

type OrganisationCode struct {
	*AbstractDto

	Code string `xml:"Code,omitempty"`
}

type AbstractDto struct {
	ActionOnEntity *ActionOnEntity `xml:"ActionOnEntity,omitempty"`
	EndDate        string          `xml:"EndDate,omitempty"`
	StartDate      string          `xml:"StartDate,omitempty"`
}

type ArrayOfContact struct {
	Contact []*Contact `xml:"Contact,omitempty"`
}

type Contact struct {
	*AbstractDto
	Email            string   `xml:"Email,omitempty"`
	Fax              string   `xml:"Fax,omitempty"`
	FirstName        string   `xml:"FirstName,omitempty"`
	GroupName        string   `xml:"GroupName,omitempty"`
	JobTitle         string   `xml:"JobTitle,omitempty"`
	LastName         string   `xml:"LastName,omitempty"`
	Mobile           string   `xml:"Mobile,omitempty"`
	OrganisationName string   `xml:"OrganisationName,omitempty"`
	Phone            string   `xml:"Phone,omitempty"`
	PostalAddress    *Address `xml:"PostalAddress,omitempty"`
	RoleCode         string   `xml:"RoleCode,omitempty"`
	Title            string   `xml:"Title,omitempty"`
	TypeCode         string   `xml:"TypeCode,omitempty"`
}

type Address struct {
	CountryCode   string  `xml:"CountryCode,omitempty"`
	Latitude      float64 `xml:"Latitude,omitempty"`
	Line1         string  `xml:"Line1,omitempty"`
	Line2         string  `xml:"Line2,omitempty"`
	Longitude     float64 `xml:"Longitude,omitempty"`
	Postcode      string  `xml:"Postcode,omitempty"`
	StateCode     string  `xml:"StateCode,omitempty"`
	StateOverseas string  `xml:"StateOverseas,omitempty"`
	Suburb        string  `xml:"Suburb,omitempty"`
}

type ArrayOfDataManagerAssignment struct {
	DataManagerAssignment []*DataManagerAssignment `xml:"DataManagerAssignment,omitempty"`
}

type DataManagerAssignment struct {
	*AbstractDto

	Code string `xml:"Code,omitempty"`
}

type ArrayOfOrganisationLocation struct {
	OrganisationLocation []*OrganisationLocation `xml:"OrganisationLocation,omitempty"`
}

type OrganisationLocation struct {
	*AbstractDto
	Address *Address `xml:"Address,omitempty"`
}

type ArrayOfResponsibleLegalPerson struct {
	ResponsibleLegalPerson []*ResponsibleLegalPerson `xml:"ResponsibleLegalPerson,omitempty"`
}
type ResponsibleLegalPerson struct {
	*AbstractDto
	Abns *ArrayOfstring `xml:"Abns,omitempty"`
	Acn  string         `xml:"Acn,omitempty"`
	Name string         `xml:"Name,omitempty"`
}

type ArrayOfRole struct {
	Role []*Role `xml:"Role,omitempty"`
}

type Role struct {
	*AbstractDto
	Abbreviation string `xml:"Abbreviation,omitempty"`
	Code         int32  `xml:"Code,omitempty"`
	Description  string `xml:"Description,omitempty"`
}

type ArrayOfTradingName struct {
	TradingName []*TradingName `xml:"TradingName,omitempty"`
}

type TradingName struct {
	*AbstractDto
	Name string `xml:"Name,omitempty"`
}

type ArrayOfUrl struct {
	Url []*Url `xml:"Url,omitempty"`
}

type Url struct {
	*AbstractDto
	Link string `xml:"Link,omitempty"`
}

type DateTimeOffset struct {
	DateTime      time.Time `xml:"DateTime,omitempty"`
	OffsetMinutes int16     `xml:"OffsetMinutes,omitempty"`
}

type ArrayOfstring struct {
	String []string `xml:"string,omitempty"`
}
