# QubitOn API — Go SDK

Go client for the [QubitOn API](https://www.qubiton.com). Validate, enrich, and assess business data across 250+ countries with 70+ APIs.

## Installation

```bash
go get github.com/qubitonhq/qubiton-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    qubiton "github.com/qubitonhq/qubiton-go"
)

func main() {
    client := qubiton.NewClient("svm...")

    resp, err := client.ValidateAddress(context.Background(), qubiton.AddressRequest{
        AddressLine1: "1600 Pennsylvania Ave NW",
        City:         "Washington",
        State:        "DC",
        PostalCode:   "20500",
        Country:      "US",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Valid: %v, Score: %.2f\n", resp.IsValid, resp.ConfidenceScore)
}
```

## Authentication

### API Key (recommended)

```go
client := qubiton.NewClient("svm...")
```

### OAuth2 Client Credentials

```go
client := qubiton.NewClient("",
    qubiton.WithOAuth2("client-id", "client-secret", ""),
)
```

## Methods

### Address & Contact Validation

| Method | Description |
|--------|-------------|
| `ValidateAddress` | Validate and standardize a postal address (249 countries) |
| `ValidateEmail` | Validate email deliverability and risk |
| `ValidatePhone` | Validate phone number against carrier databases |

### Tax Validation

| Method | Description |
|--------|-------------|
| `ValidateTax` | Validate a tax ID with live authority checks (60+ countries) |
| `ValidateTaxFormat` | Validate tax ID format via regex + checksum (193 countries, 242 types) |
| `GetSupportedTaxFormats` | List all supported country + tax type combinations |

### Bank Account Validation

| Method | Description |
|--------|-------------|
| `ValidateBankAccount` | Validate bank accounts (IBAN, SWIFT, routing — 180+ countries) |
| `ValidateBankPro` | Premium bank analytics with ownership verification |

### Business Intelligence

| Method | Description |
|--------|-------------|
| `LookupBusinessRegistration` | Look up official business registration records |
| `LookupBusinessClassification` | Look up NAICS/SIC business classification codes |
| `LookupDUNS` | Look up a DUNS number for company identification |
| `IdentifyGender` | Predict gender from a person's name |

### Corporate Structure

| Method | Description |
|--------|-------------|
| `LookupBeneficialOwnership` | Beneficial ownership for corporate transparency (49 countries) |
| `LookupCorporateHierarchy` | Corporate hierarchy and ownership structure (US) |
| `LookupHierarchy` | Parent-child company hierarchy |

### Compliance & Screening

| Method | Description |
|--------|-------------|
| `CheckSanctions` | Screen against 100+ sanctions lists (OFAC, EU, UN) |
| `ScreenPEP` | Politically Exposed Person screening |
| `CheckDirectors` | Check for disqualified directors |
| `CheckEPAProsecution` | Screen against EPA criminal prosecution records |
| `LookupEPAProsecution` | Look up EPA criminal prosecution details |
| `ValidatePeppol` | Validate Peppol participant ID (70+ ICD schemes) |
| `GetPeppolSchemes` | List all supported Peppol ICD schemes |

### Healthcare

| Method | Description |
|--------|-------------|
| `CheckHealthcareExclusion` | Screen against healthcare provider exclusion lists |
| `LookupHealthcareExclusion` | Look up healthcare provider exclusion details |
| `ValidateNPI` | Validate US National Provider Identifier |
| `ValidateMedpass` | Validate healthcare supplier via Medpass |

### Risk & Financial

| Method | Description |
|--------|-------------|
| `CheckBankruptcyRisk` | Check bankruptcy filings and proceedings |
| `LookupCreditScore` | Commercial credit score and financial stability |
| `LookupCreditAnalysis` | Comprehensive credit analysis and payment behavior |
| `LookupFailRate` | Payment failure rate and risk classification |
| `AssessEntityRisk` | Entity fraud risk and adverse media assessment |
| `AnalyzePaymentTerms` | Payment terms optimization and early-pay discounts |
| `LookupExchangeRates` | Currency exchange rates for specific dates |

### ESG & Cybersecurity

| Method | Description |
|--------|-------------|
| `LookupESGScore` | ESG (Environmental, Social, Governance) scores |
| `DomainSecurityReport` | Domain cybersecurity and threat intelligence |
| `CheckIPQuality` | IP address quality and fraud risk |

### Industry Specific

| Method | Description |
|--------|-------------|
| `LookupDOTCarrier` | USDOT/FMCSA motor carrier safety data |
| `ValidateIndiaIdentity` | Indian identity validation (Driver License, Voter ID) |
| `ValidateCertification` | Business certification validation (MBE, WBE, DBE) |
| `LookupCertification` | Business certification lookup |

### SAP Ariba

| Method | Description |
|--------|-------------|
| `LookupAribaSupplier` | SAP Ariba supplier profile lookup by ANID |
| `ValidateAribaSupplier` | SAP Ariba supplier profile validation by ANID |

## Examples

### Tax ID Validation

```go
resp, err := client.ValidateTax(ctx, qubiton.TaxRequest{
    TaxNumber:   "12-3456789",
    TaxType:     "EIN",
    Country:     "US",
    CompanyName: "Acme Corp",
})
fmt.Printf("Valid: %v, Registered: %s\n", resp.IsValid, resp.RegisteredName)
```

### Bank Account Validation

```go
resp, err := client.ValidateBankAccount(ctx, qubiton.BankAccountRequest{
    BusinessEntityType: "Business",
    Country:            "US",
    BankAccountHolder:  "Acme Corp",
    AccountNumber:      "123456789",
    BankCode:           "021000021", // routing number
})
fmt.Printf("Valid: %v, Bank: %s\n", resp.IsValid, resp.BankName)
```

### Email Validation

```go
resp, err := client.ValidateEmail(ctx, qubiton.EmailRequest{
    EmailAddress: "contact@acme.com",
})
fmt.Printf("Valid: %v, Disposable: %v\n", resp.IsValid, resp.IsDisposable)
```

### Sanctions Screening

```go
resp, err := client.CheckSanctions(ctx, qubiton.SanctionsRequest{
    CompanyName: "Acme Trading Co",
    Country:     "US",
})
fmt.Printf("Matches: %v, Lists: %v\n", resp.HasMatches, resp.ScreenedLists)
```

### Beneficial Ownership

```go
resp, err := client.LookupBeneficialOwnership(ctx, qubiton.BeneficialOwnershipRequest{
    CompanyName: "Acme Holdings Ltd",
    CountryISO2: "GB",
})
fmt.Printf("Owners: %d\n", len(resp.Owners))
```

### ESG Score

```go
resp, err := client.LookupESGScore(ctx, qubiton.ESGRequest{
    CompanyName: "Tesla Inc",
    Country:     "US",
})
fmt.Printf("Overall: %.1f, E: %.1f, S: %.1f, G: %.1f\n",
    resp.OverallScore, resp.EnvironmentScore, resp.SocialScore, resp.GovernanceScore)
```

## Error Handling

```go
resp, err := client.ValidateAddress(ctx, req)
if err != nil {
    var apiErr *qubiton.ApiError
    if errors.As(err, &apiErr) {
        if apiErr.IsRateLimit() {
            log.Println("Rate limited, retry later")
        } else if apiErr.IsAuthError() {
            log.Println("Auth failed:", apiErr.Message)
        } else {
            log.Printf("API error [%d]: %s\n", apiErr.StatusCode, apiErr.Message)
        }
    }
}
```

All response structs include a `Raw map[string]interface{}` field with the full API response for fields not covered by the typed struct.

## Options

```go
client := qubiton.NewClient("svm...",
    qubiton.WithBaseURL("https://custom-api.example.com"),
    qubiton.WithTimeout(15 * time.Second),
    qubiton.WithHTTPClient(customClient),
)
```

## Requirements

- Go 1.22+
- Zero external dependencies (stdlib only)

## MCP Protocol Support

This API is available as a native [Model Context Protocol](https://modelcontextprotocol.io) (MCP) server for Claude, ChatGPT, and other AI agents.

| Category | Count | Description |
|----------|-------|-------------|
| MCP Tools | 37+ | 1:1 mapped to API endpoints — same auth, rate limits, and plan access |
| MCP Prompts | 20 | Multi-tool workflow templates (onboarding, compliance, risk, payment) |
| MCP Resources | 7 | Reference datasets (tool inventory, risk categories, country coverage) |

- [MCP Manifest](https://mcp.qubiton.com/.well-known/mcp.json)

## Getting an API Key

1. Sign up free at [www.qubiton.com](https://www.qubiton.com/auth/register) — 100 API calls/month, no credit card
2. Navigate to Dashboard → API Keys
3. Copy your API key (starts with `svm`)

## License

MIT — Copyright (c) 2026 apexanalytix
