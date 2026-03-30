package qubiton

// ── Address Validation ────────────────────────────────────────────────────

// AddressRequest is the input for ValidateAddress.
type AddressRequest struct {
	Country      string `json:"country"`
	AddressLine1 string `json:"addressLine1,omitempty"`
	AddressLine2 string `json:"addressLine2,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	PostalCode   string `json:"postalCode,omitempty"`
	CompanyName  string `json:"companyName,omitempty"`
}

// AddressResponse is the output from ValidateAddress.
type AddressResponse struct {
	IsValid             bool              `json:"isValid"`
	ConfidenceScore     float64           `json:"confidenceScore"`
	StandardizedAddress map[string]string `json:"standardizedAddress"`
	Raw                 map[string]interface{}
}

// ── Tax ID Validation ─────────────────────────────────────────────────────

// TaxRequest is the input for ValidateTax.
type TaxRequest struct {
	TaxNumber          string `json:"taxNumber"`
	TaxType            string `json:"taxType"`
	Country            string `json:"country"`
	CompanyName        string `json:"companyName"`
	BusinessEntityType string `json:"businessEntityType,omitempty"`
}

// TaxResponse is the output from ValidateTax.
type TaxResponse struct {
	IsValid        bool   `json:"isValid"`
	TaxIDType      string `json:"taxIdType"`
	Country        string `json:"country"`
	RegisteredName string `json:"registeredName"`
	Raw            map[string]interface{}
}

// ── Tax Format Validation ─────────────────────────────────────────────────

// TaxFormatRequest is the input for ValidateTaxFormat.
type TaxFormatRequest struct {
	TaxNumber string `json:"taxNumber"`
	TaxType   string `json:"taxType"`
	Country   string `json:"country"`
}

// TaxFormatResponse is the output from ValidateTaxFormat.
type TaxFormatResponse struct {
	IsValid      bool   `json:"isValid"`
	FormatMatch  bool   `json:"formatMatch"`
	ChecksumPass bool   `json:"checksumPass"`
	TaxType      string `json:"taxType"`
	Country      string `json:"country"`
	Raw          map[string]interface{}
}

// ── Bank Account Validation ───────────────────────────────────────────────

// BankAccountRequest is the input for ValidateBankAccount.
type BankAccountRequest struct {
	BusinessEntityType string `json:"businessEntityType"`
	Country            string `json:"country"`
	BankAccountHolder  string `json:"bankAccountHolder"`
	AccountNumber      string `json:"accountNumber,omitempty"`
	BusinessName       string `json:"businessName,omitempty"`
	TaxIdNumber        string `json:"taxIdNumber,omitempty"`
	TaxType            string `json:"taxType,omitempty"`
	BankCode           string `json:"bankCode,omitempty"`
	IBAN               string `json:"iban,omitempty"`
	SwiftCode          string `json:"swiftCode,omitempty"`
}

// BankAccountResponse is the output from ValidateBankAccount.
type BankAccountResponse struct {
	IsValid     bool   `json:"isValid"`
	BankName    string `json:"bankName"`
	AccountType string `json:"accountType"`
	Raw         map[string]interface{}
}

// ── BankPro Validation ────────────────────────────────────────────────────

// BankProRequest is the input for ValidateBankPro.
type BankProRequest struct {
	BusinessEntityType string `json:"businessEntityType"`
	Country            string `json:"country"`
	BankAccountHolder  string `json:"bankAccountHolder"`
	AccountNumber      string `json:"accountNumber,omitempty"`
	BankCode           string `json:"bankCode,omitempty"`
	IBAN               string `json:"iban,omitempty"`
	SwiftCode          string `json:"swiftCode,omitempty"`
}

// BankProResponse is the output from ValidateBankPro.
type BankProResponse struct {
	IsValid         bool    `json:"isValid"`
	ConfidenceScore float64 `json:"confidenceScore"`
	OwnershipMatch  bool    `json:"ownershipMatch"`
	BankName        string  `json:"bankName"`
	Raw             map[string]interface{}
}

// ── Email Validation ──────────────────────────────────────────────────────

// EmailRequest is the input for ValidateEmail.
type EmailRequest struct {
	EmailAddress string `json:"emailAddress"`
}

// EmailResponse is the output from ValidateEmail.
type EmailResponse struct {
	IsValid      bool   `json:"isValid"`
	IsDeliverable bool  `json:"isDeliverable"`
	IsDisposable bool   `json:"isDisposable"`
	Domain       string `json:"domain"`
	Raw          map[string]interface{}
}

// ── Phone Validation ──────────────────────────────────────────────────────

// PhoneRequest is the input for ValidatePhone.
type PhoneRequest struct {
	PhoneNumber    string `json:"phoneNumber"`
	Country        string `json:"country"`
	PhoneExtension string `json:"phoneExtension,omitempty"`
}

// PhoneResponse is the output from ValidatePhone.
type PhoneResponse struct {
	IsValid     bool   `json:"isValid"`
	PhoneType   string `json:"phoneType"`
	Carrier     string `json:"carrier"`
	Country     string `json:"country"`
	Raw         map[string]interface{}
}

// ── Business Registration Lookup ──────────────────────────────────────────

// BusinessRegistrationRequest is the input for LookupBusinessRegistration.
type BusinessRegistrationRequest struct {
	CompanyName string `json:"companyName"`
	Country     string `json:"country"`
	State       string `json:"state,omitempty"`
	City        string `json:"city,omitempty"`
}

// BusinessRegistrationResponse is the output from LookupBusinessRegistration.
type BusinessRegistrationResponse struct {
	Found              bool              `json:"found"`
	CompanyName        string            `json:"companyName"`
	RegistrationNumber string            `json:"registrationNumber"`
	Status             string            `json:"status"`
	Address            map[string]string `json:"address"`
	Raw                map[string]interface{}
}

// ── Peppol Validation ─────────────────────────────────────────────────────

// PeppolRequest is the input for ValidatePeppol.
type PeppolRequest struct {
	ParticipantId   string `json:"participantId"`
	DirectoryLookup *bool  `json:"directoryLookup,omitempty"`
}

// PeppolResponse is the output from ValidatePeppol.
type PeppolResponse struct {
	IsValid       bool   `json:"isValid"`
	FormatValid   bool   `json:"formatValid"`
	DirectoryHit  bool   `json:"directoryHit"`
	SchemeId      string `json:"schemeId"`
	Raw           map[string]interface{}
}

// ── Sanctions Screening ──────────────────────────────────────────────────

// SanctionsRequest is the input for CheckSanctions.
type SanctionsRequest struct {
	CompanyName  string `json:"companyName"`
	Country      string `json:"country"`
	AddressLine1 string `json:"addressLine1,omitempty"`
	AddressLine2 string `json:"addressLine2,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	PostalCode   string `json:"postalCode,omitempty"`
}

// SanctionsResponse is the output from CheckSanctions.
type SanctionsResponse struct {
	HasMatches    bool                     `json:"hasMatches"`
	Matches       []map[string]interface{} `json:"matches"`
	ScreenedLists []string                 `json:"screenedLists"`
	Raw           map[string]interface{}
}

// ── PEP Screening ─────────────────────────────────────────────────────────

// PEPRequest is the input for ScreenPEP.
type PEPRequest struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

// PEPResponse is the output from ScreenPEP.
type PEPResponse struct {
	HasMatches bool                     `json:"hasMatches"`
	Matches    []map[string]interface{} `json:"matches"`
	Raw        map[string]interface{}
}

// ── Directors Check ───────────────────────────────────────────────────────

// DirectorsRequest is the input for CheckDirectors.
type DirectorsRequest struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Country    string `json:"country"`
	MiddleName string `json:"middlename,omitempty"`
}

// DirectorsResponse is the output from CheckDirectors.
type DirectorsResponse struct {
	HasDisqualified bool                     `json:"hasDisqualified"`
	Directors       []map[string]interface{} `json:"directors"`
	Raw             map[string]interface{}
}

// ── EPA Prosecution ───────────────────────────────────────────────────────

// EPARequest is the input for CheckEPAProsecution and LookupEPAProsecution.
type EPARequest struct {
	Name       string `json:"name,omitempty"`
	State      string `json:"state,omitempty"`
	FiscalYear string `json:"fiscalYear,omitempty"`
}

// EPAResponse is the output from CheckEPAProsecution and LookupEPAProsecution.
type EPAResponse struct {
	HasMatches bool                     `json:"hasMatches"`
	Records    []map[string]interface{} `json:"records"`
	Raw        map[string]interface{}
}

// ── Healthcare Exclusion ──────────────────────────────────────────────────

// HealthcareExclusionRequest is the input for CheckHealthcareExclusion and LookupHealthcareExclusion.
type HealthcareExclusionRequest struct {
	HealthCareType string `json:"healthCareType"` // HCO or HCP
	EntityName     string `json:"entityName,omitempty"`
	LastName       string `json:"lastName,omitempty"`
	FirstName      string `json:"firstName,omitempty"`
	Address        string `json:"address,omitempty"`
	City           string `json:"city,omitempty"`
	State          string `json:"state,omitempty"`
	ZipCode        string `json:"zipCode,omitempty"`
}

// HealthcareExclusionResponse is the output from CheckHealthcareExclusion and LookupHealthcareExclusion.
type HealthcareExclusionResponse struct {
	HasExclusions bool                     `json:"hasExclusions"`
	Exclusions    []map[string]interface{} `json:"exclusions"`
	Raw           map[string]interface{}
}

// ── Bankruptcy Risk ───────────────────────────────────────────────────────

// BankruptcyRequest is the input for CheckBankruptcyRisk.
type BankruptcyRequest struct {
	CompanyName string `json:"companyName"`
	Country     string `json:"country"`
}

// BankruptcyResponse is the output from CheckBankruptcyRisk.
type BankruptcyResponse struct {
	HasBankruptcy bool    `json:"hasBankruptcy"`
	RiskScore     float64 `json:"riskScore"`
	Raw           map[string]interface{}
}

// ── Credit Score ──────────────────────────────────────────────────────────

// CreditScoreRequest is the input for LookupCreditScore.
type CreditScoreRequest struct {
	CompanyName string `json:"companyName"`
	Country     string `json:"country"`
}

// CreditScoreResponse is the output from LookupCreditScore.
type CreditScoreResponse struct {
	Score           int    `json:"score"`
	Rating          string `json:"rating"`
	FinancialHealth string `json:"financialHealth"`
	Raw             map[string]interface{}
}

// ── Fail Rate ─────────────────────────────────────────────────────────────

// FailRateRequest is the input for LookupFailRate.
type FailRateRequest struct {
	CompanyName string `json:"companyName"`
	Country     string `json:"country"`
}

// FailRateResponse is the output from LookupFailRate.
type FailRateResponse struct {
	FailRate       float64 `json:"failRate"`
	Classification string  `json:"classification"`
	Raw            map[string]interface{}
}

// ── Entity Risk Assessment ────────────────────────────────────────────────

// EntityRiskRequest is the input for AssessEntityRisk.
type EntityRiskRequest struct {
	CompanyName        string `json:"companyName"`
	Country            string `json:"country,omitempty"`
	Category           string `json:"category,omitempty"` // Financial, Operational, Geographic, Reputational, Regulatory
	URL                string `json:"url,omitempty"`
	BusinessEntityType string `json:"businessEntityType,omitempty"`
}

// EntityRiskResponse is the output from AssessEntityRisk.
type EntityRiskResponse struct {
	RiskScore  float64 `json:"riskScore"`
	RiskLevel  string  `json:"riskLevel"`
	Indicators []map[string]interface{} `json:"indicators"`
	Raw        map[string]interface{}
}

// ── Credit Analysis ───────────────────────────────────────────────────────

// CreditAnalysisRequest is the input for LookupCreditAnalysis.
type CreditAnalysisRequest struct {
	CompanyName  string `json:"companyName"`
	AddressLine1 string `json:"addressLine1"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
	DunsNumber   string `json:"dunsNumber,omitempty"`
	PostalCode   string `json:"postalCode,omitempty"`
	AddressLine2 string `json:"addressLine2,omitempty"`
}

// CreditAnalysisResponse is the output from LookupCreditAnalysis.
type CreditAnalysisResponse struct {
	CreditScore    int     `json:"creditScore"`
	CreditLimit    float64 `json:"creditLimit"`
	PaymentBehavior string `json:"paymentBehavior"`
	Raw            map[string]interface{}
}

// ── ESG Score ─────────────────────────────────────────────────────────────

// ESGRequest is the input for LookupESGScore.
type ESGRequest struct {
	CompanyName string `json:"companyName"`
	Country     string `json:"country"`
	Domain      string `json:"domain,omitempty"`
}

// ESGResponse is the output from LookupESGScore.
type ESGResponse struct {
	OverallScore      float64 `json:"overallScore"`
	EnvironmentScore  float64 `json:"environmentScore"`
	SocialScore       float64 `json:"socialScore"`
	GovernanceScore   float64 `json:"governanceScore"`
	Raw               map[string]interface{}
}

// ── Domain Security ───────────────────────────────────────────────────────

// DomainSecurityRequest is the input for DomainSecurityReport.
type DomainSecurityRequest struct {
	DomainName string `json:"domainName"`
}

// DomainSecurityResponse is the output from DomainSecurityReport.
type DomainSecurityResponse struct {
	RiskScore    float64 `json:"riskScore"`
	ThreatLevel  string  `json:"threatLevel"`
	Raw          map[string]interface{}
}

// ── IP Quality ────────────────────────────────────────────────────────────

// IPQualityRequest is the input for CheckIPQuality.
type IPQualityRequest struct {
	IPAddress string `json:"ipAddress"`
	UserAgent string `json:"userAgent,omitempty"`
}

// IPQualityResponse is the output from CheckIPQuality.
type IPQualityResponse struct {
	FraudScore float64 `json:"fraudScore"`
	IsProxy    bool    `json:"isProxy"`
	IsVPN      bool    `json:"isVPN"`
	IsBot      bool    `json:"isBot"`
	Raw        map[string]interface{}
}

// ── Beneficial Ownership ──────────────────────────────────────────────────

// BeneficialOwnershipRequest is the input for LookupBeneficialOwnership.
type BeneficialOwnershipRequest struct {
	CompanyName  string `json:"companyName"`
	CountryISO2  string `json:"countryIso2"`
	UBOThreshold string `json:"uboThreshold,omitempty"`
	MaxLayers    string `json:"maxLayers,omitempty"`
}

// BeneficialOwnershipResponse is the output from LookupBeneficialOwnership.
type BeneficialOwnershipResponse struct {
	Owners []map[string]interface{} `json:"owners"`
	Raw    map[string]interface{}
}

// ── Corporate Hierarchy ───────────────────────────────────────────────────

// CorporateHierarchyRequest is the input for LookupCorporateHierarchy.
type CorporateHierarchyRequest struct {
	CompanyName  string `json:"companyName"`
	AddressLine1 string `json:"addressLine1"`
	City         string `json:"city"`
	State        string `json:"state"`
	ZipCode      string `json:"zipCode"`
}

// CorporateHierarchyResponse is the output from LookupCorporateHierarchy.
type CorporateHierarchyResponse struct {
	ParentCompany string                   `json:"parentCompany"`
	Subsidiaries  []map[string]interface{} `json:"subsidiaries"`
	Raw           map[string]interface{}
}

// ── DUNS Lookup ───────────────────────────────────────────────────────────

// DUNSRequest is the input for LookupDUNS.
type DUNSRequest struct {
	DunsNumber string `json:"dunsNumber"`
}

// DUNSResponse is the output from LookupDUNS.
type DUNSResponse struct {
	Found       bool   `json:"found"`
	CompanyName string `json:"companyName"`
	DunsNumber  string `json:"dunsNumber"`
	Raw         map[string]interface{}
}

// ── Hierarchy Lookup ──────────────────────────────────────────────────────

// HierarchyRequest is the input for LookupHierarchy.
type HierarchyRequest struct {
	Identifier     string `json:"identifier"`
	IdentifierType string `json:"identifierType"`
	Country        string `json:"country,omitempty"`
	Options        string `json:"options,omitempty"`
}

// HierarchyResponse is the output from LookupHierarchy.
type HierarchyResponse struct {
	Parent   map[string]interface{}   `json:"parent"`
	Children []map[string]interface{} `json:"children"`
	Raw      map[string]interface{}
}

// ── NPI Validation ────────────────────────────────────────────────────────

// NPIRequest is the input for ValidateNPI.
type NPIRequest struct {
	NPI              string `json:"npi"`
	OrganizationName string `json:"organizationName,omitempty"`
	LastName         string `json:"lastName,omitempty"`
	FirstName        string `json:"firstName,omitempty"`
	MiddleName       string `json:"middleName,omitempty"`
}

// NPIResponse is the output from ValidateNPI.
type NPIResponse struct {
	IsValid          bool   `json:"isValid"`
	OrganizationName string `json:"organizationName"`
	ProviderType     string `json:"providerType"`
	Raw              map[string]interface{}
}

// ── Medpass Validation ────────────────────────────────────────────────────

// MedpassRequest is the input for ValidateMedpass.
type MedpassRequest struct {
	ID                 string `json:"id"`
	BusinessEntityType string `json:"businessEntityType"`
	CompanyName        string `json:"companyName,omitempty"`
	TaxId              string `json:"taxId,omitempty"`
	Country            string `json:"country,omitempty"`
	State              string `json:"state,omitempty"`
	City               string `json:"city,omitempty"`
	PostalCode         string `json:"postalCode,omitempty"`
	AddressLine1       string `json:"addressLine1,omitempty"`
	AddressLine2       string `json:"addressLine2,omitempty"`
}

// MedpassResponse is the output from ValidateMedpass.
type MedpassResponse struct {
	IsValid bool `json:"isValid"`
	Raw     map[string]interface{}
}

// ── DOT Carrier Lookup ────────────────────────────────────────────────────

// DOTCarrierRequest is the input for LookupDOTCarrier.
type DOTCarrierRequest struct {
	DotNumber  string `json:"dotNumber"`
	EntityName string `json:"entityName,omitempty"`
}

// DOTCarrierResponse is the output from LookupDOTCarrier.
type DOTCarrierResponse struct {
	Found      bool   `json:"found"`
	CarrierName string `json:"carrierName"`
	DotNumber  string `json:"dotNumber"`
	SafetyRating string `json:"safetyRating"`
	Raw        map[string]interface{}
}

// ── India Identity Validation ─────────────────────────────────────────────

// IndiaIdentityRequest is the input for ValidateIndiaIdentity.
type IndiaIdentityRequest struct {
	IdentityNumber     string `json:"identityNumber"`
	IdentityNumberType string `json:"identityNumberType"` // "Driver License" or "Voter Registration"
	EntityName         string `json:"entityName,omitempty"`
	DOB                string `json:"dob,omitempty"` // YYYY-MM-DD
}

// IndiaIdentityResponse is the output from ValidateIndiaIdentity.
type IndiaIdentityResponse struct {
	IsValid bool   `json:"isValid"`
	Name    string `json:"name"`
	Raw     map[string]interface{}
}

// ── Certification ─────────────────────────────────────────────────────────

// CertificationRequest is the input for ValidateCertification and LookupCertification.
type CertificationRequest struct {
	CompanyName         string `json:"companyName"`
	Country             string `json:"country"`
	City                string `json:"city,omitempty"`
	State               string `json:"state,omitempty"`
	ZipCode             string `json:"zipCode,omitempty"`
	AddressLine1        string `json:"addressLine1,omitempty"`
	AddressLine2        string `json:"addressLine2,omitempty"`
	IdentityType        string `json:"identityType,omitempty"`
	CertificationType   string `json:"certificationType,omitempty"` // MBE, WBE, DBE, etc.
	CertificationGroup  string `json:"certificationGroup,omitempty"`
	CertificationNumber string `json:"certificationNumber,omitempty"`
}

// CertificationResponse is the output from ValidateCertification and LookupCertification.
type CertificationResponse struct {
	IsCertified    bool                     `json:"isCertified"`
	Certifications []map[string]interface{} `json:"certifications"`
	Raw            map[string]interface{}
}

// ── Business Classification ───────────────────────────────────────────────

// BusinessClassificationRequest is the input for LookupBusinessClassification.
type BusinessClassificationRequest struct {
	CompanyName string `json:"companyName"`
	City        string `json:"city"`
	State       string `json:"state"`
	Country     string `json:"country"`
	Address1    string `json:"address1,omitempty"`
	Address2    string `json:"address2,omitempty"`
	Phone       string `json:"phone,omitempty"`
	PostalCode  string `json:"postalCode,omitempty"`
}

// BusinessClassificationResponse is the output from LookupBusinessClassification.
type BusinessClassificationResponse struct {
	NAICSCode string `json:"naicsCode"`
	SICCode   string `json:"sicCode"`
	Industry  string `json:"industry"`
	Raw       map[string]interface{}
}

// ── Payment Terms ─────────────────────────────────────────────────────────

// PaymentTermsRequest is the input for AnalyzePaymentTerms.
type PaymentTermsRequest struct {
	CurrentPayTerm float64 `json:"currentPayTerm"`
	AnnualSpend    float64 `json:"annualSpend"`
	AvgDaysPay     float64 `json:"avgDaysPay"`
	SavingsRate    float64 `json:"savingsRate"`
	Threshold      float64 `json:"threshold"`
	VendorName     string  `json:"vendorName,omitempty"`
	Country        string  `json:"country,omitempty"`
}

// PaymentTermsResponse is the output from AnalyzePaymentTerms.
type PaymentTermsResponse struct {
	RecommendedTerm float64 `json:"recommendedTerm"`
	PotentialSavings float64 `json:"potentialSavings"`
	Raw             map[string]interface{}
}

// ── Exchange Rates ────────────────────────────────────────────────────────

// ExchangeRateRequest is the input for LookupExchangeRates.
type ExchangeRateRequest struct {
	BaseCurrency string `json:"baseCurrency"`
	Dates        string `json:"dates"` // Comma-separated ISO dates
}

// ExchangeRateResponse is the output from LookupExchangeRates.
type ExchangeRateResponse struct {
	BaseCurrency string                 `json:"baseCurrency"`
	Rates        map[string]interface{} `json:"rates"`
	Raw          map[string]interface{}
}

// ── SAP Ariba Supplier ────────────────────────────────────────────────────

// AribaSupplierRequest is the input for LookupAribaSupplier and ValidateAribaSupplier.
type AribaSupplierRequest struct {
	ANID string `json:"anid"`
}

// AribaSupplierResponse is the output from LookupAribaSupplier and ValidateAribaSupplier.
type AribaSupplierResponse struct {
	Found       bool   `json:"found"`
	CompanyName string `json:"companyName"`
	ANID        string `json:"anid"`
	Raw         map[string]interface{}
}

// ── Gender Identification ─────────────────────────────────────────────────

// GenderRequest is the input for IdentifyGender.
type GenderRequest struct {
	Name    string `json:"name"`
	Country string `json:"country,omitempty"`
}

// GenderResponse is the output from IdentifyGender.
type GenderResponse struct {
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
	Raw         map[string]interface{}
}

// ── Tax Format Reference ──────────────────────────────────────────────────

// TaxFormatsResponse is the output from GetSupportedTaxFormats.
type TaxFormatsResponse struct {
	Countries []map[string]interface{} `json:"countries"`
	Raw       map[string]interface{}
}

// ── Peppol Schemes Reference ──────────────────────────────────────────────

// PeppolSchemesResponse is the output from GetPeppolSchemes.
type PeppolSchemesResponse struct {
	Schemes []map[string]interface{} `json:"schemes"`
	Raw     map[string]interface{}
}
