package tga

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

func (tga *TGA) GetTrainingDetails(code string) (*TrainingComponent, error) {

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
                <GetDetails xmlns="http://training.gov.au/services/2/">
                        <request>
                                <Code>` + code + `</Code>
                                <IncludeLegacyData>false</IncludeLegacyData>
                                <InformationRequested>
                                        <ShowClassifications>true</ShowClassifications>
                                        <ShowCompanionVolumeLinks>true</ShowCompanionVolumeLinks>
                                        <ShowCompletionMapping>true</ShowCompletionMapping>
                                        <ShowComponents>true</ShowComponents>
                                        <ShowContacts>true</ShowContacts>
                                        <ShowCurrencyPeriods>true</ShowCurrencyPeriods>
                                        <ShowDataManagers>true</ShowDataManagers>
                                        <ShowFiles>true</ShowFiles>
                                        <ShowIndustrySectors>true</ShowIndustrySectors>
                                        <ShowMappingInformation>true</ShowMappingInformation>
                                        <ShowOccupations>true</ShowOccupations>
                                        <ShowRecognitionManagers>true</ShowRecognitionManagers>
                                        <ShowReleases>true</ShowReleases>
                                        <ShowUnitGrid>true</ShowUnitGrid>
                                        <ShowUsageRecommendation>true</ShowUsageRecommendation>
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

	req, _ := http.NewRequest("POST", tga.Endpoint+"/TrainingComponentServiceV2.svc/Training", strings.NewReader(soapRequest))
	req.Header.Add("Content-Type", "text/xml")
	req.Header.Add("SOAPAction", "\"http://training.gov.au/services/2/ITrainingComponentService/GetDetails\"")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	tga.lastResponse = string(body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("HTTP Error: " + resp.Status)
	}

	respEnvelope := new(tSOAPEnvelope)
	//respEnvelope.Body = SOAPBody{Content: response}
	err = xml.Unmarshal(body, respEnvelope)
	if err != nil {
		fmt.Println("unmarshal failed:", err)
		return nil, err
	}

	return respEnvelope.Body.GetDetailsResponse.GetDetailsResult, nil
}

type tSOAPEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Body    tSOAPBody
}

type tSOAPBody struct {
	XMLName            xml.Name             `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	GetDetailsResponse *tGetDetailsResponse `xml:",omitempty"`
}

type tGetDetailsResponse struct {
	XMLName          xml.Name           `xml:"http://training.gov.au/services/2/ GetDetailsResponse"`
	GetDetailsResult *TrainingComponent `xml:"GetDetailsResult,omitempty"`
}

type TrainingComponent struct {
	Classifications           *ArrayOfClassification                  `xml:"Classifications,omitempty"`
	Code                      string                                  `xml:"Code,omitempty"`
	CompletionMapping         *ArrayOfNrtCompletion                   `xml:"CompletionMapping,omitempty"`
	ComponentType             *TrainingComponentTypes                 `xml:"ComponentType,omitempty"`
	Contacts                  *ArrayOfContact                         `xml:"Contacts,omitempty"`
	CreatedDate               *DateTimeOffset                         `xml:"CreatedDate,omitempty"`
	CurrencyPeriods           *ArrayOfNrtCurrencyPeriod               `xml:"CurrencyPeriods,omitempty"`
	CurrencyStatus            string                                  `xml:"CurrencyStatus,omitempty"`
	DataManagers              *ArrayOfDataManagerAssignment           `xml:"DataManagers,omitempty"`
	IndustrySectors           *ArrayOfTrainingComponentIndustrySector `xml:"IndustrySectors,omitempty"`
	IsConfidential            bool                                    `xml:"IsConfidential,omitempty"`
	IsLegacyData              bool                                    `xml:"IsLegacyData,omitempty"`
	IscOrganisationCode       string                                  `xml:"IscOrganisationCode,omitempty"`
	MappingInformation        *ArrayOfMapping                         `xml:"MappingInformation,omitempty"`
	Occupations               *ArrayOfTrainingComponentOccupation     `xml:"Occupations,omitempty"`
	ParentCode                string                                  `xml:"ParentCode,omitempty"`
	ParentTitle               string                                  `xml:"ParentTitle,omitempty"`
	RecognitionManagers       *ArrayOfRecognitionManagerAssignment    `xml:"RecognitionManagers,omitempty"`
	Releases                  *ArrayOfRelease                         `xml:"Releases,omitempty"`
	Restrictions              *ArrayOfNrtRestriction                  `xml:"Restrictions,omitempty"`
	ReverseMappingInformation *ArrayOfMapping                         `xml:"ReverseMappingInformation,omitempty"`
	Title                     string                                  `xml:"Title,omitempty"`
	UpdatedDate               *DateTimeOffset                         `xml:"UpdatedDate,omitempty"`
	UsageRecommendations      *ArrayOfUsageRecommendation             `xml:"UsageRecommendations,omitempty"`
}

/*
type ArrayOfClassification struct {
	Classification []*Classification `xml:"Classification,omitempty"`
}

type Classification struct {
	*AbstractDto
	PurposeCode string `xml:"PurposeCode,omitempty"`
	SchemeCode  string `xml:"SchemeCode,omitempty"`
	ValueCode   string `xml:"ValueCode,omitempty"`
}
*/

/*
type AbstractDto struct {
	ActionOnEntity *ActionOnEntity `xml:"ActionOnEntity,omitempty"`
	EndDate        time.Time       `xml:"EndDate,omitempty"`
	StartDate      time.Time       `xml:"StartDate,omitempty"`
}
*/

type ArrayOfNrtCompletion struct {
	NrtCompletion []*NrtCompletion `xml:"NrtCompletion,omitempty"`
}

type NrtCompletion struct {
	*AbstractDto
	Code        string `xml:"Code,omitempty"`
	IsMandatory bool   `xml:"IsMandatory,omitempty"`
}

/*
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
*/

/*
type Address struct {
	CountryCode   string `xml:"CountryCode,omitempty"`
	Line1         string `xml:"Line1,omitempty"`
	Line2         string `xml:"Line2,omitempty"`
	Postcode      string `xml:"Postcode,omitempty"`
	StateCode     string `xml:"StateCode,omitempty"`
	StateOverseas string `xml:"StateOverseas,omitempty"`
	Suburb        string `xml:"Suburb,omitempty"`
}
*/

type ArrayOfNrtCurrencyPeriod struct {
	NrtCurrencyPeriod []*NrtCurrencyPeriod `xml:"NrtCurrencyPeriod,omitempty"`
}

type NrtCurrencyPeriod struct {
	*AbstractDto
	Authority     string `xml:"Authority,omitempty"`
	EndComment    string `xml:"EndComment,omitempty"`
	EndReasonCode string `xml:"EndReasonCode,omitempty"`
}

/*
type ArrayOfDataManagerAssignment struct {
	DataManagerAssignment []*DataManagerAssignment `xml:"DataManagerAssignment,omitempty"`
}

type DataManagerAssignment struct {
	*AbstractDto
	Code string `xml:"Code,omitempty"`
}
*/

type ArrayOfTrainingComponentIndustrySector struct {
	TrainingComponentIndustrySector []*TrainingComponentIndustrySector `xml:"TrainingComponentIndustrySector,omitempty"`
}

type TrainingComponentIndustrySector struct {
	Code        string `xml:"Code,omitempty"`
	Description string `xml:"Description,omitempty"`
	ParentCode  string `xml:"ParentCode,omitempty"`
	Title       string `xml:"Title,omitempty"`
}

type ArrayOfMapping struct {
	Mapping []*Mapping `xml:"Mapping,omitempty"`
}

type Mapping struct {
	Code         string `xml:"Code,omitempty"`
	IsEquivalent bool   `xml:"IsEquivalent,omitempty"`
	MapsToCode   string `xml:"MapsToCode,omitempty"`
	MapsToTitle  string `xml:"MapsToTitle,omitempty"`
	Notes        string `xml:"Notes,omitempty"`
	Title        string `xml:"Title,omitempty"`
}

type ArrayOfTrainingComponentOccupation struct {
	TrainingComponentOccupation []*TrainingComponentOccupation `xml:"TrainingComponentOccupation,omitempty"`
}

type TrainingComponentOccupation struct {
	Code        string `xml:"Code,omitempty"`
	Description string `xml:"Description,omitempty"`
	Title       string `xml:"Title,omitempty"`
}

type ArrayOfRecognitionManagerAssignment struct {
	RecognitionManagerAssignment []*RecognitionManagerAssignment `xml:"RecognitionManagerAssignment,omitempty"`
}

type RecognitionManagerAssignment struct {
	*AbstractDto
	Code string `xml:"Code,omitempty"`
}

type ArrayOfRelease struct {
	Release []*Release `xml:"Release,omitempty"`
}

type Release struct {
	ApprovalProcess          string                             `xml:"ApprovalProcess,omitempty"`
	CompanionVolumeLinks     *ArrayOfReleaseCompanionVolumeLink `xml:"CompanionVolumeLinks,omitempty"`
	Components               *ArrayOfReleaseComponent           `xml:"Components,omitempty"`
	Currency                 string                             `xml:"Currency,omitempty"`
	Files                    *ArrayOfReleaseFile                `xml:"Files,omitempty"`
	IscApprovalDate          string                             `xml:"IscApprovalDate,omitempty"`
	MinisterialAgreementDate string                             `xml:"MinisterialAgreementDate,omitempty"`
	NqcEndorsementDate       string                             `xml:"NqcEndorsementDate,omitempty"`
	ReleaseDate              string                             `xml:"ReleaseDate,omitempty"`
	ReleaseNumber            string                             `xml:"ReleaseNumber,omitempty"`
	UnitGrid                 *ArrayOfUnitGridEntry              `xml:"UnitGrid,omitempty"`
}

type ArrayOfReleaseCompanionVolumeLink struct {
	ReleaseCompanionVolumeLink []*ReleaseCompanionVolumeLink `xml:"ReleaseCompanionVolumeLink,omitempty"`
}

type ReleaseCompanionVolumeLink struct {
	LinkNotes              string `xml:"LinkNotes,omitempty"`
	LinkText               string `xml:"LinkText,omitempty"`
	LinkUrl                string `xml:"LinkUrl,omitempty"`
	PublishedComponentType string `xml:"PublishedComponentType,omitempty"`
}

type ArrayOfReleaseComponent struct {
	ReleaseComponent []*ReleaseComponent `xml:"ReleaseComponent,omitempty"`
}

type ReleaseComponent struct {
	Code            string                  `xml:"Code,omitempty"`
	ReleaseCurrency string                  `xml:"ReleaseCurrency,omitempty"`
	ReleaseDate     string                  `xml:"ReleaseDate,omitempty"`
	ReleaseNumber   string                  `xml:"ReleaseNumber,omitempty"`
	Title           string                  `xml:"Title,omitempty"`
	Type            *TrainingComponentTypes `xml:"Type,omitempty"`
}

type ArrayOfReleaseFile struct {
	ReleaseFile []*ReleaseFile `xml:"ReleaseFile,omitempty"`
}

type ReleaseFile struct {
	RelativePath string `xml:"RelativePath,omitempty"`
	Size         int32  `xml:"Size,omitempty"`
}

type ArrayOfUnitGridEntry struct {
	UnitGridEntry []*UnitGridEntry `xml:"UnitGridEntry,omitempty"`
}

type UnitGridEntry struct {
	Code        string `xml:"Code,omitempty"`
	IsEssential bool   `xml:"IsEssential,omitempty"`
	Title       string `xml:"Title,omitempty"`
}

type ArrayOfNrtRestriction struct {
	NrtRestriction []*NrtRestriction `xml:"NrtRestriction,omitempty"`
}

type NrtRestriction struct {
	*AbstractDto
	Restriction string `xml:"Restriction,omitempty"`
}

type ArrayOfUsageRecommendation struct {
	UsageRecommendation []*UsageRecommendation `xml:"UsageRecommendation,omitempty"`
}

type UsageRecommendation struct {
	*AbstractDto
	State string `xml:"State,omitempty"`
}
