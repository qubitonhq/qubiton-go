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

	qubiton "github.com/qubitonhq/qubiton-go"
)

func main() {
	apiKey := os.Getenv("QUBITON_API_KEY")
	if apiKey == "" {
		log.Fatal("Set QUBITON_API_KEY environment variable")
	}

	client := qubiton.NewClient(apiKey)
	ctx := context.Background()

	// ── Address Validation ──
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
		fmt.Printf("  Valid: %v, Score: %.2f\n", addr.IsValid, addr.ConfidenceScore)
	}

	// ── Tax ID Validation ──
	fmt.Println("\n=== Tax ID Validation ===")
	tax, err := client.ValidateTax(ctx, qubiton.TaxRequest{
		TaxNumber:   "12-3456789",
		TaxType:     "EIN",
		Country:     "US",
		CompanyName: "Example Corp",
	})
	if err != nil {
		handleError(err)
	} else {
		fmt.Printf("  Valid: %v, Type: %s, Name: %s\n", tax.IsValid, tax.TaxIDType, tax.RegisteredName)
	}

	// ── Bank Account Validation ──
	fmt.Println("\n=== Bank Account Validation ===")
	bank, err := client.ValidateBankAccount(ctx, qubiton.BankAccountRequest{
		BusinessEntityType: "Business",
		Country:            "US",
		BankAccountHolder:  "Example Corp",
		AccountNumber:      "123456789",
		BankCode:           "021000021",
	})
	if err != nil {
		handleError(err)
	} else {
		fmt.Printf("  Valid: %v, Bank: %s\n", bank.IsValid, bank.BankName)
	}

	// ── Email Validation ──
	fmt.Println("\n=== Email Validation ===")
	email, err := client.ValidateEmail(ctx, qubiton.EmailRequest{
		EmailAddress: "test@example.com",
	})
	if err != nil {
		handleError(err)
	} else {
		fmt.Printf("  Valid: %v, Deliverable: %v, Disposable: %v\n",
			email.IsValid, email.IsDeliverable, email.IsDisposable)
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
		fmt.Printf("  Valid: %v, Type: %s, Carrier: %s\n",
			phone.IsValid, phone.PhoneType, phone.Carrier)
	}

	// ── Sanctions Screening ──
	fmt.Println("\n=== Sanctions Screening ===")
	sanctions, err := client.CheckSanctions(ctx, qubiton.SanctionsRequest{
		CompanyName: "Example Trading Co",
		Country:     "US",
	})
	if err != nil {
		handleError(err)
	} else {
		fmt.Printf("  Matches: %v, Lists screened: %v\n",
			sanctions.HasMatches, sanctions.ScreenedLists)
	}

	// ── PEP Screening ──
	fmt.Println("\n=== PEP Screening ===")
	pep, err := client.ScreenPEP(ctx, qubiton.PEPRequest{
		Name:    "John Smith",
		Country: "US",
	})
	if err != nil {
		handleError(err)
	} else {
		fmt.Printf("  Matches: %v\n", pep.HasMatches)
	}

	// ── Business Registration ──
	fmt.Println("\n=== Business Registration Lookup ===")
	biz, err := client.LookupBusinessRegistration(ctx, qubiton.BusinessRegistrationRequest{
		CompanyName: "Apple Inc",
		Country:     "US",
		State:       "CA",
	})
	if err != nil {
		handleError(err)
	} else if len(biz.BusinessRegistrations) > 0 {
		r := biz.BusinessRegistrations[0]
		fmt.Printf("  Name: %s, ID: %s, Status: %s, Type: %s, Formed: %s\n",
			r.EntityName, r.RegistrationId, r.Status, r.BusinessEntityType, r.FormationDate)
	} else {
		fmt.Println("  No registrations found")
	}

	// ── ESG Score ──
	fmt.Println("\n=== ESG Score ===")
	esg, err := client.LookupESGScore(ctx, qubiton.ESGRequest{
		CompanyName: "Microsoft Corp",
		Country:     "US",
	})
	if err != nil {
		handleError(err)
	} else {
		fmt.Printf("  Overall: %.1f, E: %.1f, S: %.1f, G: %.1f\n",
			esg.OverallScore, esg.EnvironmentScore, esg.SocialScore, esg.GovernanceScore)
	}

	// ── Credit Analysis ──
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
		fmt.Printf("  Score: %d, Behavior: %s\n",
			credit.CreditScore, credit.PaymentBehavior)
	}

	// ── Beneficial Ownership ──
	fmt.Println("\n=== Beneficial Ownership ===")
	ubo, err := client.LookupBeneficialOwnership(ctx, qubiton.BeneficialOwnershipRequest{
		CompanyName: "Acme Holdings Ltd",
		CountryISO2: "GB",
	})
	if err != nil {
		handleError(err)
	} else {
		fmt.Printf("  Owners found: %d\n", len(ubo.Owners))
	}

	// ── Domain Security ──
	fmt.Println("\n=== Domain Security Report ===")
	domain, err := client.DomainSecurityReport(ctx, qubiton.DomainSecurityRequest{
		DomainName: "example.com",
	})
	if err != nil {
		handleError(err)
	} else {
		fmt.Printf("  Risk Score: %.1f, Threat Level: %s\n",
			domain.RiskScore, domain.ThreatLevel)
	}
}

func handleError(err error) {
	var apiErr *qubiton.ApiError
	if errors.As(err, &apiErr) {
		if apiErr.IsRateLimit() {
			fmt.Println("  Rate limited — retry later")
		} else if apiErr.IsAuthError() {
			fmt.Printf("  Auth error: %s\n", apiErr.Message)
		} else {
			fmt.Printf("  API error [%d]: %s\n", apiErr.StatusCode, apiErr.Message)
		}
	} else {
		fmt.Printf("  Error: %v\n", err)
	}
}
