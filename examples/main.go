// Example usage of the QubitOn Go SDK.
//
// Run: go run main.go
//
// Set your API key via environment variable:
//
//	export QUBITON_API_KEY=svm...
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	qubiton "github.com/qubitonhq/qubiton-go"
)

func main() {
	apiKey := os.Getenv("QUBITON_API_KEY")
	if apiKey == "" {
		log.Fatal("Set QUBITON_API_KEY environment variable")
	}

	client := qubiton.NewClient(apiKey)
	ctx := context.Background()

	// ── Address Validation (single-object response) ──
	fmt.Println("=== Address Validation ===")
	addr, err := client.ValidateAddress(ctx, qubiton.AddressRequest{
		AddressLine1: "1600 Pennsylvania Ave NW",
		City:         "Washington",
		State:        "DC",
		PostalCode:   "20500",
		Country:      "US",
	})
	if err != nil {
		handleError(err)
	} else {
		fmt.Printf("  Address: %s, %s, %s %s\n", addr.Address1, addr.City, addr.State, addr.PostalCode)
		if addr.IsResidential != nil {
			fmt.Printf("  Residential: %v\n", *addr.IsResidential)
		}
	}

	// ── Tax ID Validation (single-object response) ──
	fmt.Println("\n=== Tax ID Validation ===")
	tax, err := client.ValidateTax(ctx, qubiton.TaxRequest{
		IdentityNumber:     "12-3456789",
		IdentityNumberType: "EIN",
		Country:            "US",
		EntityName:         "Example Corp",
	})
	if err != nil {
		handleError(err)
	} else {
		nameMatch := "n/a"
		if tax.IsEntityNameMatch != nil {
			nameMatch = fmt.Sprintf("%v", *tax.IsEntityNameMatch)
		}
		fmt.Printf("  Valid: %v, Name Match: %s, Entity: %s, FileNo: %s\n",
			tax.TaxValid, nameMatch, tax.EntityName, tax.FileNumber)
	}

	// ── Tax Format Validation (offline regex+checksum, no provider call) ──
	fmt.Println("\n=== Tax Format Validation ===")
	taxFmt, err := client.ValidateTaxFormat(ctx, qubiton.TaxFormatRequest{
		TaxNumber: "12-3456789",
		TaxType:   "EIN",
		Country:   "US",
	})
	if err != nil {
		handleError(err)
	} else {
		fmt.Printf("  Valid: %v, Format: %v, Checksum: %v\n",
			taxFmt.IsValid, taxFmt.FormatMatch, taxFmt.ChecksumPass)
	}

	// ── Bank Account Validation (single-object response) ──
	fmt.Println("\n=== Bank Account Validation ===")
	bank, err := client.ValidateBankAccount(ctx, qubiton.BankAccountRequest{
		BankNumberType:     "ROUTING",
		Country:            "US",
		BusinessEntityType: "Business",
		BusinessName:       "Example Corp",
		AccountNumber:      "123456789",
		BankCode:           "021000021",
	})
	if err != nil {
		handleError(err)
	} else {
		fmt.Printf("  Bank: %s, Status: %s, Account: %s\n", bank.BankName, bank.ValidationStatus, bank.AccountHolder)
	}

	// ── Email Validation ──
	fmt.Println("\n=== Email Validation ===")
	email, err := client.ValidateEmail(ctx, qubiton.EmailRequest{
		EmailAddress: "test@example.com",
	})
	if err != nil {
		handleError(err)
	} else {
		fraud := -1
		if email.FraudScore != nil {
			fraud = *email.FraudScore
		}
		fmt.Printf("  Email: %s, Type: %s, FraudScore: %d\n",
			email.EmailAddress, email.EmailType, fraud)
	}

	// ── Phone Validation ──
	fmt.Println("\n=== Phone Validation ===")
	phone, err := client.ValidatePhone(ctx, qubiton.PhoneRequest{
		PhoneNumber: "+12025551234",
		Country:     "US",
	})
	if err != nil {
		handleError(err)
	} else {
		fraud := -1
		if phone.FraudScore != nil {
			fraud = *phone.FraudScore
		}
		fmt.Printf("  Phone: %s, FraudScore: %d\n",
			phone.FullPhoneNumber, fraud)
	}

	// ── Sanctions Screening (ARRAY response — one entry per matching list) ──
	fmt.Println("\n=== Sanctions Screening ===")
	sanctions, err := client.CheckSanctions(ctx, qubiton.SanctionsRequest{
		CompanyName: "Example Trading Co",
		Country:     "US",
	})
	if err != nil {
		handleError(err)
	} else {
		fmt.Printf("  Sources returned: %d\n", len(sanctions))
		for i, s := range sanctions {
			fmt.Printf("    [%d] match=%v score=%.2f desc=%s\n", i, s.IsMatch, s.Score, s.Description)
		}
	}

	// ── PEP Screening (ARRAY response — one entry per matching dataset) ──
	fmt.Println("\n=== PEP Screening ===")
	pep, err := client.ScreenPEP(ctx, qubiton.PEPRequest{
		Name:    "John Smith",
		Country: "US",
	})
	if err != nil {
		handleError(err)
	} else {
		fmt.Printf("  Sources returned: %d\n", len(pep))
		for i, src := range pep {
			fmt.Printf("    [%d] persons=%d organizations=%d\n",
				i, len(src.Persons), len(src.Organizations))
		}
	}

	// ── Business Registration (ARRAY of wrappers; each wrapper carries
	//    a list of registrations + validation metadata) ──
	fmt.Println("\n=== Business Registration Lookup ===")
	bizList, err := client.LookupBusinessRegistration(ctx, qubiton.BusinessRegistrationRequest{
		EntityName: "Apple Inc",
		Country:    "US",
		State:      "CA",
	})
	if err != nil {
		handleError(err)
	} else if len(bizList) > 0 && len(bizList[0].BusinessRegistrations) > 0 {
		w := bizList[0]
		r := w.BusinessRegistrations[0]
		fmt.Printf("  Found %d wrapper(s). First wrapper has %d registration(s); ValidationDescription=%q\n",
			len(bizList), len(w.BusinessRegistrations), w.ValidationDescription)
		fmt.Printf("    First registration: Name=%s ID=%s Status=%s Type=%s Formed=%s\n",
			r.EntityName, r.RegistrationId, r.Status, r.BusinessEntityType, r.FormationDate)
	} else {
		fmt.Println("  No registrations found")
	}

	// ── ESG Score (ARRAY response — one per provider) ──
	fmt.Println("\n=== ESG Score ===")
	esgList, err := client.LookupESGScore(ctx, qubiton.ESGRequest{
		CompanyName: "Microsoft Corp",
	})
	if err != nil {
		handleError(err)
	} else if len(esgList) > 0 {
		esg := esgList[0]
		fmt.Printf("  %d providers. First: Grade=%s, E=%d, S=%d, G=%d, Total=%d\n",
			len(esgList), esg.Grade, esg.EnvironmentScore, esg.SocialScore, esg.GovernanceScore, esg.Total)
	}

	// ── Credit Analysis (ARRAY response — one per provider) ──
	fmt.Println("\n=== Credit Analysis ===")
	credit, err := client.LookupCreditAnalysis(ctx, qubiton.CreditAnalysisRequest{
		CompanyName:  "Example Corp",
		AddressLine1: "123 Main St",
		City:         "New York",
		State:        "NY",
		Country:      "US",
	})
	if err != nil {
		handleError(err)
	} else {
		fmt.Printf("  %d analyses returned\n", len(credit))
		if len(credit) > 0 {
			ca := credit[0]
			fmt.Printf("    ApplicationId: %s, ExactMatchFound: %v, BureauCompanies: %d\n",
				ca.ApplicationId, ca.ExactMatchFound, len(ca.BureauCompanyList))
			if ca.ApplicationEcf != nil && ca.ApplicationEcf.DecisionOutcome != nil {
				fmt.Printf("    Outcome: %s\n", ca.ApplicationEcf.DecisionOutcome.Outcome)
			}
		}
	}

	// ── Beneficial Ownership (ARRAY response) ──
	fmt.Println("\n=== Beneficial Ownership ===")
	uboList, err := client.LookupBeneficialOwnership(ctx, qubiton.BeneficialOwnershipRequest{
		CompanyName: "Acme Holdings Ltd",
		CountryISO2: "GB",
	})
	if err != nil {
		handleError(err)
	} else {
		fmt.Printf("  %d hierarchies returned\n", len(uboList))
		if len(uboList) > 0 {
			tree := uboList[0].OwnershipTreeField
			if tree != nil {
				fmt.Printf("    First entity: %s, ubos=%d, nodes=%d\n",
					tree.NameField, len(tree.UltimateBeneficialOwners), len(tree.NodesField))
			} else {
				fmt.Printf("    Response code: %s, message: %s\n",
					uboList[0].ResponseCodeField, uboList[0].Message)
			}
		}
	}

	// ── Domain Security (single-object response) ──
	fmt.Println("\n=== Domain Security Report ===")
	domain, err := client.DomainSecurityReport(ctx, qubiton.DomainSecurityRequest{
		DomainName: "example.com",
	})
	if err != nil {
		handleError(err)
	} else {
		score := 0.0
		if domain.Score != nil {
			score = *domain.Score
		}
		fmt.Printf("  Domain: %s, Score: %.1f, Grade: %s\n",
			domain.DomainName, score, domain.Grade)
	}

	// ── Exchange Rates (path-param + array body — array response) ──
	fmt.Println("\n=== Exchange Rates ===")
	today := time.Now().UTC().Truncate(24 * time.Hour)
	rates, err := client.LookupExchangeRates(ctx, qubiton.ExchangeRateRequest{
		BaseCurrency: "USD",
		Dates:        []time.Time{today},
	})
	if err != nil {
		handleError(err)
	} else {
		fmt.Printf("  %d days returned\n", len(rates))
		if len(rates) > 0 {
			fmt.Printf("    %s: %d currency pairs\n", rates[0].Date, len(rates[0].ExchangeRates))
		}
	}

	// ── Reference: Supported Tax Formats (ARRAY response) ──
	fmt.Println("\n=== Supported Tax Formats ===")
	formats, err := client.GetSupportedTaxFormats(ctx)
	if err != nil {
		handleError(err)
	} else {
		fmt.Printf("  %d format entries\n", len(formats))
	}
}

func handleError(err error) {
	var apiErr *qubiton.ApiError
	if errors.As(err, &apiErr) {
		switch {
		case apiErr.IsRateLimit():
			fmt.Printf("  Rate limited — retry after %ds\n", apiErr.RetryAfter)
		case apiErr.IsAuthError():
			fmt.Printf("  Auth error: %s\n", apiErr.Message)
		case apiErr.IsServerError():
			fmt.Printf("  Server error [%d]: %s\n", apiErr.StatusCode, apiErr.Message)
		default:
			fmt.Printf("  API error [%d]: %s\n", apiErr.StatusCode, apiErr.Message)
		}
	} else {
		fmt.Printf("  Error: %v\n", err)
	}
}
