# QubitOn API — Go SDK

Go client for the [QubitOn API](https://www.qubiton.com). Validate, enrich, and assess business data across 250+ countries with 45 APIs.

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

    fmt.Printf("City: %s, State: %s, Postal: %s\n", resp.City, resp.State, resp.PostalCode)
}
```

## Authentication

### API Key (recommended)

```go
client := qubiton.NewClient("svm...")
```

The SDK sends the lowercase `apikey` header expected by the QubitOn API.

### OAuth (key + secret → JWE access token)

```go
client := qubiton.NewClient("",
    qubiton.WithOAuth("client-key", "client-secret", ""),
)
```

The `tokenURL` argument may be empty to use the default `{baseURL}/api/oauth/token`. The SDK caches the token until 30 seconds before its declared expiry, dedupes concurrent fetches via singleflight, and transparently refreshes on a stale 401.

## Method index

| Group | Method | Returns |
|-------|--------|---------|
| Address & contact | `ValidateAddress` | `*AddressResponse` |
|                   | `ValidateEmail`   | `*EmailResponse` |
|                   | `ValidatePhone`   | `*PhoneResponse` |
| Tax               | `ValidateTax`            | `*TaxResponse` |
|                   | `ValidateTaxFormat`      | `*TaxFormatResponse` |
|                   | `GetSupportedTaxFormats` | `[]SupportedTaxFormat` |
| Bank              | `ValidateBankAccount` | `*BankAccountResponse` |
|                   | `ValidateBankPro`     | `*BankProResponse` |
| Business intelligence | `LookupBusinessRegistration`   | `[]BusinessRegistrationResponse` |
|                       | `LookupBusinessClassification` | `[]BusinessClassificationResponse` |
|                       | `LookupDUNS`                   | `[]DUNSResponse` |
|                       | `IdentifyGender`               | `*GenderResponse` |
| Corporate structure | `LookupBeneficialOwnership` | `[]BeneficialOwnershipResponse` |
|                     | `LookupCorporateHierarchy`  | `[]CorporateHierarchyResponse` |
|                     | `LookupHierarchy`           | `*HierarchyResponse` |
| Compliance & screening | `CheckSanctions`        | `[]SanctionsResponse` |
|                        | `ScreenPEP`             | `[]PEPResponse` |
|                        | `CheckDirectors`        | `*DirectorsResponse` |
|                        | `CheckEPAProsecution`   | `[]EPAResponse` |
|                        | `LookupEPAProsecution`  | `[]EPAResponse` |
|                        | `ValidatePeppol`        | `*PeppolResponse` |
|                        | `GetPeppolSchemes`      | `[]SupportedPeppolScheme` |
| Healthcare | `CheckHealthcareExclusion`   | `[]HealthcareExclusionResponse` |
|            | `LookupHealthcareExclusion`  | `*HealthcareExclusionLookupResponse` |
|            | `ValidateNPI`                | `*NPIResponse` |
|            | `ValidateMedpass`            | `*MedpassResponse` |
| Risk & financial | `LookupRisk`             | `[]RiskResponse` (Social/Governance/Environmental adverse media) |
|                  | `CheckBankruptcyRisk`    | `*BankruptcyResponse` (routes to /api/risk/riskcontrol) |
|                  | `LookupCreditScore`      | `*CreditScoreResponse` (routes to /api/risk/riskcontrol) |
|                  | `LookupFailRate`         | `*FailRateResponse` (routes to /api/risk/riskcontrol) |
|                  | `LookupCreditAnalysis`   | `[]CreditAnalysisResponse` |
|                  | `AssessEntityRisk`       | `*EntityRiskResponse` |
|                  | `AnalyzePaymentTerms`    | `*PaymentTermsResponse` |
|                  | `LookupExchangeRates`    | `[]ExchangeRateResponse` |
| ESG & cybersecurity | `LookupESGScore`        | `[]ESGResponse` |
|                     | `DomainSecurityReport`  | `*DomainSecurityResponse` |
|                     | `CheckIPQuality`        | `*IPQualityResponse` |
| Industry specific | `LookupDOTCarrier`       | `[]DOTCarrierResponse` |
|                   | `ValidateIndiaIdentity`  | `*IndiaIdentityResponse` |
|                   | `ValidateCertification`  | `*CertificationResponse` |
|                   | `LookupCertification`    | `[]CertificationResponse` |
| SAP Ariba | `LookupAribaSupplier`    | `[]AribaSupplierResponse` |
|           | `ValidateAribaSupplier`  | `*AribaSupplierResponse` |
| Misc      | `CheckCallbackStatus`    | `*CheckStatusResponse` (bulk callback status by `CallBackID`) |
|           | `ScreenContinuous`       | `*ContinuousScreeningResponse` (server returns 501 today; field shape is best-guess) |

## Examples

### Sanctions screening (array response)

```go
results, err := client.CheckSanctions(ctx, qubiton.SanctionsRequest{
    CompanyName: "Acme Trading Co",
    Country:     "US",
})
if err != nil { /* ... */ }

for _, src := range results {
    fmt.Printf("source=%s match=%v score=%.2f\n", src.Description, src.IsMatch, src.Score)
}
```

### Tax ID validation (single object)

```go
tax, err := client.ValidateTax(ctx, qubiton.TaxRequest{
    EntityName:         "Acme Corp",
    IdentityNumber:     "12-3456789",
    IdentityNumberType: "EIN",
    Country:            "US",
})
if err != nil { /* ... */ }
fmt.Printf("Valid: %v, Entity name: %s, FileNumber: %s\n", tax.TaxValid, tax.EntityName, tax.FileNumber)
if tax.IsEntityNameMatch != nil {
    fmt.Printf("Name match: %v\n", *tax.IsEntityNameMatch)
}
```

### Bank account validation

```go
bank, err := client.ValidateBankAccount(ctx, qubiton.BankAccountRequest{
    BankNumberType:     "ROUTING",
    Country:            "US",
    BusinessEntityType: "Business",
    BankAccountHolder:  "Acme Corp",
    AccountNumber:      "123456789",
    BankCode:           "021000021",
})
if err != nil { /* ... */ }
fmt.Printf("Status: %s, Bank: %s\n", bank.ValidationStatus, bank.BankName)
```

### Email validation

```go
email, err := client.ValidateEmail(ctx, qubiton.EmailRequest{EmailAddress: "ap@acme.com"})
if err != nil { /* ... */ }
fmt.Printf("Type: %s\n", email.EmailType)
if email.FraudScore != nil {
    fmt.Printf("Fraud score: %d\n", *email.FraudScore)
}
```

### Beneficial ownership (array response)

```go
hierarchies, err := client.LookupBeneficialOwnership(ctx, qubiton.BeneficialOwnershipRequest{
    CompanyName: "Acme Holdings Ltd",
    CountryISO2: "GB",
})
if err != nil { /* ... */ }
for _, h := range hierarchies {
    if h.OwnershipTreeField != nil {
        fmt.Printf("entity=%s ubos=%d nodes=%d\n",
            h.OwnershipTreeField.NameField,
            len(h.OwnershipTreeField.UltimateBeneficialOwners),
            len(h.OwnershipTreeField.NodesField))
    } else {
        fmt.Printf("response=%s message=%s\n", h.ResponseCodeField, h.Message)
    }
}
```

### Currency exchange rates (path-param + array body)

```go
import "time"

today := time.Now().UTC().Truncate(24 * time.Hour)
rates, err := client.LookupExchangeRates(ctx, qubiton.ExchangeRateRequest{
    BaseCurrency: "USD",
    Dates:        []time.Time{today},
})
```

## Error handling

```go
resp, err := client.ValidateAddress(ctx, req)
if err != nil {
    var apiErr *qubiton.ApiError
    if errors.As(err, &apiErr) {
        switch {
        case apiErr.IsRateLimit():
            fmt.Printf("rate limited; retry after %ds\n", apiErr.RetryAfter)
        case apiErr.IsAuthError():
            fmt.Println("auth failed:", apiErr.Message)
        case apiErr.IsServerError():
            fmt.Println("server error:", apiErr.Message)
        default:
            fmt.Printf("api error [%d]: %s\n", apiErr.StatusCode, apiErr.Message)
        }
    }
}
```

Sentinel errors are available for `errors.Is` matching: `ErrAuth`, `ErrRateLimit`, `ErrServerError`, `ErrNotFound`, `ErrValidation` — 5 sentinels covering 8 status-code mappings (401/403 → `ErrAuth`, 404 → `ErrNotFound`, 429 → `ErrRateLimit`, 400/422 → `ErrValidation`, 5xx → `ErrServerError`).

The full parsed JSON body is preserved on `ApiError.Raw` for endpoints that return structured ProblemDetails (validation errors are folded into `Message` automatically). When `RetryAfter` is set, `Error()` appends `(retry after Ns)` so the hint surfaces in logs without a type assertion.

## Options

```go
client := qubiton.NewClient("svm...",
    qubiton.WithBaseURL("https://custom-api.example.com"),
    qubiton.WithHTTPClient(customClient),                   // optional: bring your own *http.Client
    qubiton.WithTimeout(15 * time.Second),                  // apply AFTER WithHTTPClient — clones, never mutates the caller's hc
    qubiton.WithOAuth("key", "secret", ""),                 // tokenURL "" → uses {baseURL}/api/oauth/token
)
```

The default HTTP client uses a tuned `http.Transport` (`MaxIdleConnsPerHost: 20`, `MaxConnsPerHost: 100`, HTTP/2 enabled, `IdleConnTimeout: 90s`). `MaxConnsPerHost: 100` is the hard ceiling on simultaneous connections per host — tune up for highly concurrent batch loads against one QubitOn host.

### Option ordering: `WithTimeout` and `WithHTTPClient`

- `WithHTTPClient(hc)` then `WithTimeout(d)`: `WithTimeout` shallow-clones `hc`, sets the new timeout on the clone, and uses the clone — `hc` itself is never mutated. The Transport pointer is shared, so the SDK reuses the caller's connection pool. This is the recommended ordering.
- `WithTimeout(d)` then `WithHTTPClient(hc)`: `WithHTTPClient` overrides — the earlier `WithTimeout` is lost. Apply `WithTimeout` last when you need both.

## Retry & resilience

- Exponential backoff with ±25% jitter, floored at 50ms and capped at 30s.
- Retries on 5xx, 408, 429 (respecting `Retry-After`), and network errors.
- 4xx responses (other than 408/429) are NOT retried.
- Single transparent retry on 401 when OAuth is configured (token refreshed; does not consume an attempt slot).
- All retries respect `ctx` cancellation.

## Requirements

- Go 1.22+ (uses `math/rand/v2` for concurrency-safe jitter).
- Zero external dependencies (stdlib only).

## MCP Protocol Support

This API is also available as a native [Model Context Protocol](https://modelcontextprotocol.io) (MCP) server for Claude, ChatGPT, and other AI agents.

| Category | Count | Description |
|----------|-------|-------------|
| MCP Tools | 39 | 1:1 mapped to API endpoints — same auth, rate limits, and plan access |
| MCP Prompts | 20 | Multi-tool workflow templates (onboarding, compliance, risk, payment) |
| MCP Resources | 7 | Reference datasets (tool inventory, risk categories, country coverage) |

- [MCP Manifest](https://mcp.qubiton.com/.well-known/mcp.json)

## Getting an API key

1. Sign up free at [www.qubiton.com](https://www.qubiton.com/auth/register) — 100 API calls/month, no credit card.
2. Navigate to Dashboard → API Keys.
3. Copy your API key (starts with `svm`).

## License

MIT — Copyright (c) 2026 apexanalytix
