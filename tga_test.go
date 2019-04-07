package tga

import (
	"fmt"
	"testing"
	"time"
)

func TestTgaTraining(t *testing.T) {

	// Test basic operation
	{
		tga := TGA{"https://ws.sandbox.training.gov.au/Deewr.Tga.Webservices/", "WebService.Read", "Asdf098", ""}

		// Check this random email address is not throttled
		tc, err := tga.GetTrainingDetails("CUACMP511")
		if err != nil {
			t.Fatalf("tga.GetTrainingDetails() failed: %v", err)
		}
		if tc == nil {
			t.Fatalf("tga.GetTrainingDetails() should return object")
		}
		if tc.Code != "" {
			fmt.Println("Code:", tc.Code)
		}

		fmt.Println()
		fmt.Println("Classifications:")
		if tc.Contacts != nil {
			for i, r := range tc.Classifications.Classification {
				if i == 5 {
					fmt.Println("     ....")
					break
				}
				fmt.Println("  -", r.PurposeCode, r.SchemeCode, r.ValueCode)
			}
			fmt.Printf("%d Classifications\n", len(tc.Classifications.Classification))
		} else {
			fmt.Println("  - Classifications not set")
		}

		fmt.Println()
		fmt.Println("Contacts:")
		if tc.Contacts != nil {
			for i, r := range tc.Contacts.Contact {
				if i == 5 {
					fmt.Println("     ....")
					break
				}
				fmt.Println("  -", r.FirstName, r.LastName, r.Email, r.RoleCode, r.Phone, r.JobTitle, r.GroupName)
			}
			fmt.Printf("%d Contacts\n", len(tc.Contacts.Contact))
		} else {
			fmt.Println("  - Contacts not set")
		}

		fmt.Println()
		fmt.Println("Completion Mapping:")
		if tc.CompletionMapping != nil {
			for _, p := range tc.CompletionMapping.NrtCompletion {
				fmt.Println("  -", p.Code, p.IsMandatory)
			}
			fmt.Printf("    %d CompletionMapping\n", len(tc.CompletionMapping.NrtCompletion))
		} else {
			fmt.Println("CompletionMapping not set")
		}

		fmt.Println()
		fmt.Println("MappingInformation:")
		if tc.MappingInformation != nil {
			for _, p := range tc.MappingInformation.Mapping {
				fmt.Println("  -", p.Code, p.Title, p.IsEquivalent, p.MapsToCode, p.MapsToTitle, p.Notes)
			}
			fmt.Printf("    %d MappingInformation\n", len(tc.MappingInformation.Mapping))
		} else {
			fmt.Println("MappingInformation not set")
		}

	}

}

func TestTgaOrganisation(t *testing.T) {

	// Test basic operation
	{
		tga := TGA{"https://ws.sandbox.training.gov.au/Deewr.Tga.Webservices/", "WebService.Read", "Asdf098", ""}

		// Check this random email address is not throttled
		o, err := tga.GetOrganisationDetails("90525")
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
			fmt.Println("Contacts")
			for i, r := range o.Contacts.Contact {
				if i == 5 {
					fmt.Println("     ....")
					break
				}
				fmt.Println("  -", r.FirstName, r.LastName, r.Email, r.RoleCode, r.Phone, r.JobTitle, r.GroupName)
			}
			fmt.Printf("%d Contacts\n", len(o.Contacts.Contact))
		} else {
			fmt.Println("  - Contacts not set")
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
		fmt.Println()
		fmt.Println("Responsible Legal Persons:")
		if o.ResponsibleLegalPersons != nil {
			for _, p := range o.ResponsibleLegalPersons.ResponsibleLegalPerson {
				fmt.Println("  -", p.Name)
			}
			fmt.Printf("    %d ResponsibleLegalPersons\n", len(o.ResponsibleLegalPersons.ResponsibleLegalPerson))
		} else {
			fmt.Println("ResponsibleLegalPersons not set")
		}
		if o.Roles != nil {
			for _, r := range o.Roles.Role {
				fmt.Println("  -", r.Abbreviation, r.Code, r.Description)
			}
			fmt.Printf("    %d Roles\n", len(o.Roles.Role))
		} else {
			fmt.Println("Roles not set")
		}
		fmt.Println()
		fmt.Println("Classifications:")
		if o.Classifications != nil {
			for _, r := range o.Classifications.Classification {
				fmt.Println("  - ", r.PurposeCode, r.SchemeCode, r.ValueCode)
			}
			fmt.Printf("%d Roles\n", len(o.Classifications.Classification))
		} else {
			fmt.Println(" - Classifications not set")
		}

		fmt.Println()
		fmt.Println("Scopes:")
		currentYear := time.Now().Year()
		if o.Scopes != nil {
			for i, r := range o.Scopes.Scope {
				if i == 30000 {
					fmt.Println("     ....")
					break
				}
				if r.End().Year() >= currentYear {
					fmt.Println("  - ", r.StartDate, r.EndDate, r.ExtentCode, r.IsImplicit, r.IsRefused, r.NrtCode, r.TrainingComponentType)
				}
			}
			fmt.Printf("%d scopes\n", len(o.Scopes.Scope))
		} else {
			fmt.Println(" - Scopes not set")
		}

		//fmt.Println(tga.LastSoapResponse())
	}

}
