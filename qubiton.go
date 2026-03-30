// Package qubiton provides a Go client for the QubitOn API.
//
// 70+ APIs for validating, enriching, and assessing business data across 250+ countries.
//
// Usage:
//
//	client := qubiton.NewClient("svm...")
//	resp, err := client.ValidateAddress(ctx, qubiton.AddressRequest{
//	    AddressLine1: "123 Main St", City: "New York", State: "NY",
//	    PostalCode: "10001", Country: "US",
//	})
package qubiton

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

const (
	version        = "0.2.0"
	defaultBaseURL = "https://api.qubiton.com"
	maxRetries     = 3
)

// Client is the QubitOn API client.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	oauth      *oauth2TokenManager
}

// Option configures a Client.
type Option func(*Client)

// WithBaseURL overrides the default API base URL.
func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = url }
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.httpClient.Timeout = d }
}

// WithHTTPClient provides a custom HTTP client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// WithOAuth2 enables OAuth2 client credentials authentication.
func WithOAuth2(clientID, clientSecret, tokenURL string) Option {
	return func(c *Client) {
		if tokenURL == "" {
			tokenURL = c.baseURL + "/api/oauth/token"
		}
		c.oauth = newOAuth2TokenManager(clientID, clientSecret, tokenURL, c.httpClient)
	}
}

// NewClient creates a new QubitOn API client.
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (map[string]interface{}, error) {
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}
	}

	url := c.baseURL + path
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		var reqBody io.Reader
		if bodyBytes != nil {
			reqBody = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
		if err != nil {
			return nil, fmt.Errorf("create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "qubiton-go-sdk/"+version)

		if c.apiKey != "" {
			req.Header.Set("X-Api-Key", c.apiKey)
		} else if c.oauth != nil {
			token, err := c.oauth.getToken()
			if err != nil {
				return nil, err
			}
			req.Header.Set("Authorization", "Bearer "+token)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			if attempt < maxRetries-1 {
				sleep(ctx, backoff(attempt))
				continue
			}
			break
		}

		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == 429 {
			if attempt < maxRetries-1 {
				sleep(ctx, backoff(attempt))
				continue
			}
			return nil, &ApiError{StatusCode: 429, Message: "Rate limit exceeded"}
		}

		if resp.StatusCode >= 500 {
			lastErr = &ApiError{StatusCode: resp.StatusCode, Message: string(respBody)}
			if attempt < maxRetries-1 {
				sleep(ctx, backoff(attempt))
				continue
			}
			return nil, lastErr
		}

		if resp.StatusCode == 401 || resp.StatusCode == 403 {
			msg := "Authentication failed"
			var raw map[string]interface{}
			if json.Unmarshal(respBody, &raw) == nil {
				if m, ok := raw["message"].(string); ok {
					msg = m
				}
			}
			return nil, &ApiError{StatusCode: resp.StatusCode, Message: msg, Raw: raw}
		}

		if resp.StatusCode >= 400 {
			msg := "Request failed"
			var raw map[string]interface{}
			if json.Unmarshal(respBody, &raw) == nil {
				if m, ok := raw["message"].(string); ok {
					msg = m
				}
			}
			return nil, &ApiError{StatusCode: resp.StatusCode, Message: msg, Raw: raw}
		}

		var result map[string]interface{}
		if err := json.Unmarshal(respBody, &result); err != nil {
			return nil, fmt.Errorf("decode response: %w", err)
		}
		return result, nil
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, &ApiError{StatusCode: 0, Message: "request failed after retries"}
}

// doGet performs a GET request with no body.
func (c *Client) doGet(ctx context.Context, path string) (map[string]interface{}, error) {
	return c.doRequest(ctx, "GET", path, nil)
}

// ── Address Validation ────────────────────────────────────────────────────

// ValidateAddress validates and standardizes a postal address across 249 countries.
func (c *Client) ValidateAddress(ctx context.Context, req AddressRequest) (*AddressResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/address/validate", req)
	if err != nil {
		return nil, err
	}
	var resp AddressResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// ── Tax Validation ────────────────────────────────────────────────────────

// ValidateTax validates a tax identification number across 60+ countries with live authority checks.
func (c *Client) ValidateTax(ctx context.Context, req TaxRequest) (*TaxResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/tax/validate", req)
	if err != nil {
		return nil, err
	}
	var resp TaxResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// ValidateTaxFormat validates tax ID format using regex and checksum algorithms
// for 193 countries and 242 tax types.
func (c *Client) ValidateTaxFormat(ctx context.Context, req TaxFormatRequest) (*TaxFormatResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/tax/format-validate", req)
	if err != nil {
		return nil, err
	}
	var resp TaxFormatResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// ── Bank Account Validation ───────────────────────────────────────────────

// ValidateBankAccount validates bank accounts across 180+ countries (IBAN, SWIFT, routing numbers).
func (c *Client) ValidateBankAccount(ctx context.Context, req BankAccountRequest) (*BankAccountResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/bank/validate", req)
	if err != nil {
		return nil, err
	}
	var resp BankAccountResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// ValidateBankPro performs premium bank analytics with ownership verification and confidence scoring.
func (c *Client) ValidateBankPro(ctx context.Context, req BankProRequest) (*BankProResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/bankaccount/pro/validate", req)
	if err != nil {
		return nil, err
	}
	var resp BankProResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// ── Email & Phone Validation ──────────────────────────────────────────────

// ValidateEmail validates an email address for deliverability and risk.
func (c *Client) ValidateEmail(ctx context.Context, req EmailRequest) (*EmailResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/email/validate", req)
	if err != nil {
		return nil, err
	}
	var resp EmailResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// ValidatePhone validates a phone number against carrier databases.
func (c *Client) ValidatePhone(ctx context.Context, req PhoneRequest) (*PhoneResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/phone/validate", req)
	if err != nil {
		return nil, err
	}
	var resp PhoneResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// ── Business Registration ─────────────────────────────────────────────────

// LookupBusinessRegistration looks up official business registration records.
func (c *Client) LookupBusinessRegistration(ctx context.Context, req BusinessRegistrationRequest) (*BusinessRegistrationResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/businessregistration/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp BusinessRegistrationResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// ── Peppol ────────────────────────────────────────────────────────────────

// ValidatePeppol validates a Peppol participant ID against 70+ ISO 6523 ICD schemes.
func (c *Client) ValidatePeppol(ctx context.Context, req PeppolRequest) (*PeppolResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/peppol/validate", req)
	if err != nil {
		return nil, err
	}
	var resp PeppolResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// ── Sanctions & Compliance ────────────────────────────────────────────────

// CheckSanctions screens an entity against 100+ global sanctions lists (OFAC, EU, UN).
func (c *Client) CheckSanctions(ctx context.Context, req SanctionsRequest) (*SanctionsResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/prohibited/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp SanctionsResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// ScreenPEP screens against Politically Exposed Person databases.
func (c *Client) ScreenPEP(ctx context.Context, req PEPRequest) (*PEPResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/pep/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp PEPResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// CheckDirectors checks for disqualified directors.
func (c *Client) CheckDirectors(ctx context.Context, req DirectorsRequest) (*DirectorsResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/disqualifieddirectors/validate", req)
	if err != nil {
		return nil, err
	}
	var resp DirectorsResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// ── EPA Prosecution ───────────────────────────────────────────────────────

// CheckEPAProsecution screens against EPA criminal prosecution records.
func (c *Client) CheckEPAProsecution(ctx context.Context, req EPARequest) (*EPAResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/epa/validate", req)
	if err != nil {
		return nil, err
	}
	var resp EPAResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// LookupEPAProsecution looks up EPA criminal prosecution details.
func (c *Client) LookupEPAProsecution(ctx context.Context, req EPARequest) (*EPAResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/epa/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp EPAResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// ── Healthcare Exclusion ──────────────────────────────────────────────────

// CheckHealthcareExclusion screens against healthcare provider exclusion lists.
func (c *Client) CheckHealthcareExclusion(ctx context.Context, req HealthcareExclusionRequest) (*HealthcareExclusionResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/healthcare/validate", req)
	if err != nil {
		return nil, err
	}
	var resp HealthcareExclusionResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// LookupHealthcareExclusion looks up healthcare provider exclusion details.
func (c *Client) LookupHealthcareExclusion(ctx context.Context, req HealthcareExclusionRequest) (*HealthcareExclusionResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/healthcare/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp HealthcareExclusionResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// ── Risk & Financial ──────────────────────────────────────────────────────

// CheckBankruptcyRisk checks if a company has filed for or is in bankruptcy proceedings.
func (c *Client) CheckBankruptcyRisk(ctx context.Context, req BankruptcyRequest) (*BankruptcyResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/risk/bankruptcy", req)
	if err != nil {
		return nil, err
	}
	var resp BankruptcyResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// LookupCreditScore looks up commercial credit score and financial stability.
func (c *Client) LookupCreditScore(ctx context.Context, req CreditScoreRequest) (*CreditScoreResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/risk/creditscore", req)
	if err != nil {
		return nil, err
	}
	var resp CreditScoreResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// LookupFailRate looks up payment failure rate and risk classification.
func (c *Client) LookupFailRate(ctx context.Context, req FailRateRequest) (*FailRateResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/risk/failrate", req)
	if err != nil {
		return nil, err
	}
	var resp FailRateResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// AssessEntityRisk assesses entity fraud risk and adverse media.
func (c *Client) AssessEntityRisk(ctx context.Context, req EntityRiskRequest) (*EntityRiskResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/entityrisk/validate", req)
	if err != nil {
		return nil, err
	}
	var resp EntityRiskResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// LookupCreditAnalysis performs comprehensive credit analysis on a business entity.
func (c *Client) LookupCreditAnalysis(ctx context.Context, req CreditAnalysisRequest) (*CreditAnalysisResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/creditanalysis/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp CreditAnalysisResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// ── ESG & Cybersecurity ───────────────────────────────────────────────────

// LookupESGScore looks up ESG (Environmental, Social, Governance) scores.
func (c *Client) LookupESGScore(ctx context.Context, req ESGRequest) (*ESGResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/esg/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp ESGResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// DomainSecurityReport assesses domain cybersecurity and threat intelligence.
func (c *Client) DomainSecurityReport(ctx context.Context, req DomainSecurityRequest) (*DomainSecurityResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/itsecurity/validate", req)
	if err != nil {
		return nil, err
	}
	var resp DomainSecurityResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// CheckIPQuality checks IP address quality and fraud risk.
func (c *Client) CheckIPQuality(ctx context.Context, req IPQualityRequest) (*IPQualityResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/ipquality/validate", req)
	if err != nil {
		return nil, err
	}
	var resp IPQualityResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// ── Corporate Structure ───────────────────────────────────────────────────

// LookupBeneficialOwnership looks up beneficial ownership for corporate transparency.
func (c *Client) LookupBeneficialOwnership(ctx context.Context, req BeneficialOwnershipRequest) (*BeneficialOwnershipResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/beneficialownership/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp BeneficialOwnershipResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// LookupCorporateHierarchy looks up corporate hierarchy and ownership structure (US only).
func (c *Client) LookupCorporateHierarchy(ctx context.Context, req CorporateHierarchyRequest) (*CorporateHierarchyResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/corporatehierarchy/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp CorporateHierarchyResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// LookupDUNS looks up a DUNS number for company identification.
func (c *Client) LookupDUNS(ctx context.Context, req DUNSRequest) (*DUNSResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/duns-number-lookup", req)
	if err != nil {
		return nil, err
	}
	var resp DUNSResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// LookupHierarchy looks up company parent-child hierarchy.
func (c *Client) LookupHierarchy(ctx context.Context, req HierarchyRequest) (*HierarchyResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/parentchild/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp HierarchyResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// ── Industry Specific ─────────────────────────────────────────────────────

// ValidateNPI validates a US National Provider Identifier.
func (c *Client) ValidateNPI(ctx context.Context, req NPIRequest) (*NPIResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/npi/validate", req)
	if err != nil {
		return nil, err
	}
	var resp NPIResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// ValidateMedpass validates a healthcare supplier via Medpass.
func (c *Client) ValidateMedpass(ctx context.Context, req MedpassRequest) (*MedpassResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/medpass/validate", req)
	if err != nil {
		return nil, err
	}
	var resp MedpassResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// LookupDOTCarrier looks up USDOT/FMCSA motor carrier safety data.
func (c *Client) LookupDOTCarrier(ctx context.Context, req DOTCarrierRequest) (*DOTCarrierResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/dot/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp DOTCarrierResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// ValidateIndiaIdentity validates Indian identity documents (Driver License, Voter Registration).
func (c *Client) ValidateIndiaIdentity(ctx context.Context, req IndiaIdentityRequest) (*IndiaIdentityResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/identity/validate", req)
	if err != nil {
		return nil, err
	}
	var resp IndiaIdentityResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// ── Certification ─────────────────────────────────────────────────────────

// ValidateCertification validates a business certification (MBE, WBE, DBE, etc.).
func (c *Client) ValidateCertification(ctx context.Context, req CertificationRequest) (*CertificationResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/certification/validate", req)
	if err != nil {
		return nil, err
	}
	var resp CertificationResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// LookupCertification looks up business certifications (diversity, small business).
func (c *Client) LookupCertification(ctx context.Context, req CertificationRequest) (*CertificationResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/certification/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp CertificationResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// ── Business Classification ───────────────────────────────────────────────

// LookupBusinessClassification looks up NAICS/SIC business classification codes.
func (c *Client) LookupBusinessClassification(ctx context.Context, req BusinessClassificationRequest) (*BusinessClassificationResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/businessclassification/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp BusinessClassificationResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// ── Financial Operations ──────────────────────────────────────────────────

// AnalyzePaymentTerms analyzes payment terms for optimization and early-pay discounts.
func (c *Client) AnalyzePaymentTerms(ctx context.Context, req PaymentTermsRequest) (*PaymentTermsResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/paymentterms/validate", req)
	if err != nil {
		return nil, err
	}
	var resp PaymentTermsResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// LookupExchangeRates looks up currency exchange rates for specific dates.
func (c *Client) LookupExchangeRates(ctx context.Context, req ExchangeRateRequest) (*ExchangeRateResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/currency/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp ExchangeRateResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// ── Supplier Profile (SAP Ariba) ──────────────────────────────────────────

// LookupAribaSupplier looks up a SAP Ariba supplier profile by ANID.
func (c *Client) LookupAribaSupplier(ctx context.Context, req AribaSupplierRequest) (*AribaSupplierResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/ariba/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp AribaSupplierResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// ValidateAribaSupplier validates a SAP Ariba supplier profile by ANID.
func (c *Client) ValidateAribaSupplier(ctx context.Context, req AribaSupplierRequest) (*AribaSupplierResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/ariba/validate", req)
	if err != nil {
		return nil, err
	}
	var resp AribaSupplierResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// ── Gender Identification ─────────────────────────────────────────────────

// IdentifyGender predicts gender from a person's name.
func (c *Client) IdentifyGender(ctx context.Context, req GenderRequest) (*GenderResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/genderize/validate", req)
	if err != nil {
		return nil, err
	}
	var resp GenderResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// ── Reference Endpoints ───────────────────────────────────────────────────

// GetSupportedTaxFormats lists all supported country + tax type combinations.
// Returns 193 countries and 242 tax types with format descriptions.
func (c *Client) GetSupportedTaxFormats(ctx context.Context) (*TaxFormatsResponse, error) {
	data, err := c.doGet(ctx, "/api/tax/format-validate/countries")
	if err != nil {
		return nil, err
	}
	var resp TaxFormatsResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// GetPeppolSchemes lists all supported Peppol ICD schemes.
func (c *Client) GetPeppolSchemes(ctx context.Context) (*PeppolSchemesResponse, error) {
	data, err := c.doGet(ctx, "/api/peppol/schemes")
	if err != nil {
		return nil, err
	}
	var resp PeppolSchemesResponse
	remarshal(data, &resp)
	resp.Raw = data
	return &resp, nil
}

// ── Helpers ──────────────────────────────────────────────────────────────

func remarshal(data map[string]interface{}, v interface{}) {
	b, _ := json.Marshal(data)
	json.Unmarshal(b, v) //nolint:errcheck
}

func backoff(attempt int) time.Duration {
	return time.Duration(math.Pow(2, float64(attempt))) * time.Second
}

func sleep(ctx context.Context, d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
	case <-t.C:
	}
}
