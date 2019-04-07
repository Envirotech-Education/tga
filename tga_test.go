package tga

import (
	"fmt"
	"testing"
)

func TestTga(t *testing.T) {

	// Test basic operation
	{
		tga := TGA{"https://ws.sandbox.training.gov.au/Deewr.Tga.Webservices/", "WebService.Read", "Asdf098"}

		// Check this random email address is not throttled
		o, err := tga.GetDetails("90525")
		if err != nil {
			t.Fatalf("tga.GetDetails() failed: %v", err)
		}
		if o == nil {
			t.Fatalf("tga.GetDetails() should return object")
		}
		if o.Codes != nil {
			fmt.Println(o.Codes.OrganisationCode[0].Code)
		}
		if o.Contacts != nil {
			for _, r := range o.Contacts.Contact {
				fmt.Println(r.FirstName, r.LastName, r.Email, r.RoleCode, r.Phone, r.JobTitle, r.GroupName)
			}
			fmt.Printf("%d Contacts\n", len(o.Contacts.Contact))
		} else {
			fmt.Println("Contacts not set")
		}
		if o.TradingNames != nil {
			fmt.Println(o.TradingNames.TradingName[0].Name)
		}
		if o.Locations != nil {
			for _, l := range o.Locations.OrganisationLocation {
				fmt.Println(l.Address.Line1, l.Address.Suburb, l.Address.StateCode)
			}
			fmt.Printf("%d Locations\n", len(o.Locations.OrganisationLocation))
		} else {
			fmt.Println("Locations not set")
		}
		if o.ResponsibleLegalPersons != nil {
			for _, p := range o.ResponsibleLegalPersons.ResponsibleLegalPerson {
				fmt.Println(" Person:", p.Name)
			}
			fmt.Printf("%d ResponsibleLegalPersons\n", len(o.ResponsibleLegalPersons.ResponsibleLegalPerson))
		} else {
			fmt.Println("ResponsibleLegalPersons not set")
		}
		if o.Roles != nil {
			for _, r := range o.Roles.Role {
				fmt.Println(r.Abbreviation, r.Code, r.Description)
			}
			fmt.Printf("%d Roles\n", len(o.Roles.Role))
		} else {
			fmt.Println("Roles not set")
		}

	}

}
