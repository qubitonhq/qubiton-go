package qubiton

import "time"

// ── Base Request ──────────────────────────────────────────────────────────

// BaseRequest mirrors the canonical .NET BaseRequest (smartvm.BusinessEntities
// .Client.API.Bases.BaseRequest). The server marks RequestedByClient as
// [Required] [StringLength(350)] — every authenticated request must populate
// it, otherwise the server returns 400 with a validation error.
//
// All concrete *Request types embed BaseRequest so callers can populate these
// fields once on a builder or set them per-request. The SDK's User-Agent is
// substituted for RequestedByClient automatically by the client when the
// caller leaves it empty.
type BaseRequest struct {
	// RequestedByClient identifies the calling application. Required by the
	// server (StringLength=350). The SDK populates this with its User-Agent
	// when callers leave it empty.
	RequestedByClient string `json:"requestedByClient,omitempty"`
	// RequestedByApplication is an optional friendlier app name (StringLength=350).
	RequestedByApplication string `json:"requestedByApplication,omitempty"`
	// RequestedByApplicationUrl is an optional URL identifying the app (StringLength=3000).
	RequestedByApplicationUrl string `json:"requestedByApplicationUrl,omitempty"`
	// RequestedByIPAddress is an optional caller IP (StringLength=100).
	// Server-side property name is RequestedByIPAddress; with .NET CamelCase
	// JsonNamingPolicy that serialises to "requestedByIPAddress".
	RequestedByIPAddress string `json:"requestedByIPAddress,omitempty"`
	// SourceUniqueId is an optional caller-supplied id echoed back in
	// responses for correlation (StringLength=128).
	SourceUniqueId string `json:"sourceUniqueId,omitempty"`
	// QubitOnUniqueId is an optional caller-supplied id echoed back in
	// responses for correlation (StringLength=128).
	QubitOnUniqueId string `json:"qubitOnUniqueId,omitempty"`
}

// ── Address Validation ────────────────────────────────────────────────────

// AddressRequest is the input for ValidateAddress. Mirrors the canonical
// AddressRequest : BaseRequest, IBaseEntityAddress.
type AddressRequest struct {
	BaseRequest
	Country         string `json:"country"`
	AddressLine1    string `json:"addressLine1,omitempty"`
	AddressLine2    string `json:"addressLine2,omitempty"`
	AddressLine3    string `json:"addressLine3,omitempty"`
	AddressLine4    string `json:"addressLine4,omitempty"`
	AddressLine5    string `json:"addressLine5,omitempty"`
	AddressLine6    string `json:"addressLine6,omitempty"`
	AddressLine7    string `json:"addressLine7,omitempty"`
	AddressLine8    string `json:"addressLine8,omitempty"`
	NameFull        string `json:"nameFull,omitempty"`
	NameLast        string `json:"nameLast,omitempty"`
	City            string `json:"city,omitempty"`
	State           string `json:"state,omitempty"`
	PostalCode      string `json:"postalCode,omitempty"`
	CompanyName     string `json:"companyName,omitempty"`
	EmailAddress    string `json:"emailAddress,omitempty"`
	PhoneNumber     string `json:"phoneNumber,omitempty"`
	ReferenceNumber string `json:"referenceNumber,omitempty"`
	AddressType     string `json:"addressType,omitempty"`
	// OutputInLatin requests that the standardized response be returned in
	// Latin script. Default false (server-side default).
	OutputInLatin bool `json:"outputInLatin,omitempty"`
	// CallbackUrl is an optional callback URL for bulk webhook notifications.
	CallbackUrl string `json:"callbackUrl,omitempty"`
}

// AddressResponse is the output from ValidateAddress.
type AddressResponse struct {
	// Standardized address components
	AddressType  string   `json:"addressType,omitempty"`
	Address1     string   `json:"address1,omitempty"`
	Address2     string   `json:"address2,omitempty"`
	Address3     string   `json:"address3,omitempty"`
	Address4     string   `json:"address4,omitempty"`
	SuiteNumber  string   `json:"suiteNumber,omitempty"`
	StreetNumber string   `json:"streetNumber,omitempty"`
	City         string   `json:"city,omitempty"`
	State        string   `json:"state,omitempty"`
	StateName    string   `json:"stateName,omitempty"`
	PostalCode   string   `json:"postalCode,omitempty"`
	Country      *Country `json:"country,omitempty"`
	PremiseType  string   `json:"premiseType,omitempty"`
	Province     string   `json:"province,omitempty"`
	CareOf       string   `json:"careOf,omitempty"`

	// Local-language address variants
	LocalAddress1     string `json:"localAddress1,omitempty"`
	LocalAddress2     string `json:"localAddress2,omitempty"`
	LocalAddress3     string `json:"localAddress3,omitempty"`
	LocalAddress4     string `json:"localAddress4,omitempty"`
	LocalSuiteNumber  string `json:"localSuiteNumber,omitempty"`
	LocalStreetNumber string `json:"localStreetNumber,omitempty"`
	LocalCity         string `json:"localCity,omitempty"`
	LocalState        string `json:"localState,omitempty"`
	LocalPostalCode   string `json:"localPostalCode,omitempty"`
	LocalProvince     string `json:"localProvince,omitempty"`
	LocalPremiseType  string `json:"localPremiseType,omitempty"`
	LocalCareOf       string `json:"localCareOf,omitempty"`
	LocalCountry      string `json:"localCountry,omitempty"`

	// PO Box
	POBoxNumber     string `json:"poBoxNumber,omitempty"`
	POBoxCity       string `json:"poBoxCity,omitempty"`
	POBoxState      string `json:"poBoxState,omitempty"`
	POBoxPostalCode string `json:"poBoxPostalCode,omitempty"`
	POBoxCountry    string `json:"poBoxCountry,omitempty"`

	// Geocoding and classification
	GeoCode       string `json:"geoCode,omitempty"`
	InCityLimit   string `json:"inCityLimit,omitempty"`
	IsResidential *bool  `json:"isResidential,omitempty"`

	// Original request echo
	RequestCompanyName  string `json:"requestCompanyName,omitempty"`
	RequestFullName     string `json:"requestFullName,omitempty"`
	RequestAddressLine1 string `json:"requestAddressLine1,omitempty"`
	RequestAddressLine2 string `json:"requestAddressLine2,omitempty"`
	RequestCity         string `json:"requestCity,omitempty"`
	RequestState        string `json:"requestState,omitempty"`
	RequestPostalCode   string `json:"requestPostalCode,omitempty"`
	RequestEmailAddress string `json:"requestEmailAddress,omitempty"`
	RequestPhoneNumber  string `json:"requestPhoneNumber,omitempty"`

	PreferredLanguage string                   `json:"preferredLanguage,omitempty"`
	AdditionalInfo    map[string]interface{}   `json:"additionalInfo,omitempty"`
	AddressCodesInfo  []map[string]interface{} `json:"addressCodesInfo,omitempty"`

	// Validation
	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
	SourceUniqueId    string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId   string             `json:"qubitOnUniqueId,omitempty"`
}

// Country represents a country with ISO codes. Mirrors the canonical
// Country DTO referenced from many response shapes (AddressResponse,
// PhoneResponse, IPQualityResponse, etc.).
//
// Wire field naming: CountryISO2 and CountryISO3 wire as "countryISO2" /
// "countryISO3" because the .NET CamelCase JsonNamingPolicy lowercases only
// the leading character run; the trailing acronym keeps its caps. This is
// distinct from the IBAN/CLABE/etc. case earlier in this file where the
// acronym is at the start of the property name and is fully lower-cased.
type Country struct {
	CountryName      string `json:"countryName,omitempty"`
	CountryISO2      string `json:"countryISO2,omitempty"`
	CountryISO3      string `json:"countryISO3,omitempty"`
	Continent        string `json:"continent,omitempty"`
	AlternateCountry string `json:"alternateCountry,omitempty"`
}

// ValidationResult represents a single validation check result.
type ValidationResult struct {
	Key         string `json:"key,omitempty"`
	Value       string `json:"value,omitempty"`
	Description string `json:"description,omitempty"`
}

// ── Tax ID Validation ─────────────────────────────────────────────────────

// TaxRequest is the input for ValidateTax.
type TaxRequest struct {
	BaseRequest
	// IdentityNumber is the tax ID or identity number. Required.
	IdentityNumber string `json:"identityNumber"`
	// IdentityNumberType is the type of tax/identity number (e.g., EIN, VAT, GST, TIN, ABN, PAN, SSN). Required.
	IdentityNumberType string `json:"identityNumberType"`
	// EntityName is the business/vendor/company name. Required by the canonical BaseTaxRequest.
	EntityName string `json:"entityName"`
	// Country is ISO 3166-1 alpha-2, alpha-3, or full country name. Required.
	Country string `json:"country"`
	// BusinessEntityType is optional (e.g., "Corporation", "LLC", "Individual").
	BusinessEntityType string `json:"businessEntityType,omitempty"`
	// IdentityState is an optional state for state-issued identity numbers.
	IdentityState string `json:"identityState,omitempty"`
	// AddressRequest carries optional address validation info for tax requests that support it.
	AddressRequest *TaxAddressRequest `json:"addressRequest,omitempty"`
	// CallbackUrl is an optional callback URL for bulk webhook notifications.
	CallbackUrl string `json:"callbackUrl,omitempty"`
}

// TaxAddressRequest carries optional address fields for tax requests.
type TaxAddressRequest struct {
	AddressLine1 string `json:"addressLine1,omitempty"`
	AddressLine2 string `json:"addressLine2,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	PostalCode   string `json:"postalCode,omitempty"`
	Country      string `json:"country,omitempty"`
}

// TaxAdditionalInfo mirrors the canonical AdditionalInfo nested in TaxResponse.
type TaxAdditionalInfo struct {
	ConfirmingTaxAuthority              string `json:"confirmingTaxAuthority,omitempty"`
	EntityNameMatchesTaxAuthorityRecord bool   `json:"entityNameMatchesTaxAuthorityRecord,omitempty"`
	EntityTypeMatchesTaxAuthorityRecord bool   `json:"entityTypeMatchesTaxAuthorityRecord,omitempty"`
	EntityNameFromTaxAuthorityRecord    string `json:"entityNameFromTaxAuthorityRecord,omitempty"`
}

// TaxResponse is the output from ValidateTax. Mirrors the canonical
// BaseTaxResponse + TaxResponse hierarchy in the .NET server.
type TaxResponse struct {
	// From TaxValidationResponse (legacy fields, also returned).
	TaxValid          bool  `json:"taxValid"`
	IsEntityNameMatch *bool `json:"isEntityNameMatch,omitempty"`

	// From BaseTaxResponse
	IdentityNumberValidationID int64              `json:"identityNumberValidationID,omitempty"`
	IdentityNumberType         string             `json:"identityNumberType,omitempty"`
	IdentityNumber             string             `json:"identityNumber,omitempty"`
	EntityName                 string             `json:"entityName,omitempty"`
	RequestedByIPAddress       string             `json:"requestedByIPAddress,omitempty"`
	TaxAddress                 *AddressResponse   `json:"taxAddress,omitempty"`
	BusinessEntityType         string             `json:"businessEntityType,omitempty"`
	SmartvmBusinessEntityType  string             `json:"smartvmBusinessEntityType,omitempty"`
	TaxAdditionalInfo          *TaxAdditionalInfo `json:"taxAdditionalInfo,omitempty"`
	SubSection                 *int               `json:"subSection,omitempty"`
	SubSectionCode             *string            `json:"subSectionCode,omitempty"`
	SubSectionDescription      *string            `json:"subSectionDescription,omitempty"`

	// From TaxResponse (Common)
	FileNumber        string `json:"fileNumber,omitempty"`
	DayOfRegistration string `json:"dayOfRegistration,omitempty"`

	// From BaseValidationResponse
	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
	SourceUniqueId    string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId   string             `json:"qubitOnUniqueId,omitempty"`
}

// ── Tax Format Validation ─────────────────────────────────────────────────

// TaxFormatRequest is the input for ValidateTaxFormat.
// TaxFormatRequest is the input to ValidateTaxFormat.
//
// Field names mirror the canonical TaxFormatValidationRequest
// (identityNumber, identityNumberType, countryIso2). Older field
// names (taxNumber, taxType, country) were renamed during the
// upstream AutoMapper -> Facet migration; the wire JSON keys are now:
//
//	{"identityNumber":"…","identityNumberType":"…","countryIso2":"US"}
type TaxFormatRequest struct {
	BaseRequest
	// IdentityNumber is the tax ID number to validate. Required.
	IdentityNumber string `json:"identityNumber"`
	// IdentityNumberType is the tax/identity number type (e.g. "TIN", "VAT", "ABN", "EIN"). Required.
	IdentityNumberType string `json:"identityNumberType"`
	// CountryIso2 is ISO 3166-1 alpha-2, alpha-3, or full country name. Required.
	CountryIso2 string `json:"countryIso2"`
}

// TaxFormatResponse is the output from ValidateTaxFormat.
type TaxFormatResponse struct {
	IsValid      bool   `json:"isValid"`
	FormatMatch  bool   `json:"formatMatch"`
	ChecksumPass bool   `json:"checksumPass"`
	TaxType      string `json:"taxType,omitempty"`
	Country      string `json:"country,omitempty"`
	Message      string `json:"message,omitempty"`
}

// SupportedTaxFormat is one entry in the GetSupportedTaxFormats array.
type SupportedTaxFormat struct {
	Country     string `json:"country,omitempty"`
	CountryCode string `json:"countryCode,omitempty"`
	TaxType     string `json:"taxType,omitempty"`
	Pattern     string `json:"pattern,omitempty"`
	Description string `json:"description,omitempty"`
	HasChecksum bool   `json:"hasChecksum,omitempty"`
}

// ── Bank Account Validation (Ownership via BankPro) ─────────────────────

// BankAccountRequest is the input for ValidateBankAccount.
//
// Acronym JSON tags follow .NET CamelCase JsonNamingPolicy semantics —
// empirically verified against System.Text.Json with .NET 10 (the runtime
// fully lower-cases consecutive leading capitals): IBAN→iban, CLABE→clabe,
// CBU→cbu, BACS→bacs, CHAPS→chaps. SwiftCode is preserved separately because
// the canonical request distinguishes the BIC8/BIC11 SwiftCode field from
// the SWIFT discriminator value.
//
// BankNumberType is the discriminator. Valid values:
// IBAN, SWIFT, ROUTING, BankAccount, SortCode, CLABE, CBU, GIRO, QRIBAN, AccountNumber, IFSC.
type BankAccountRequest struct {
	BaseRequest

	BankNumberType string `json:"bankNumberType"`
	// Country is ISO 3166-1 alpha-2, alpha-3, or full country name. Required.
	Country string `json:"country"`

	// BankAccountValidationId is an optional server-side correlation id.
	BankAccountValidationId *int64 `json:"bankAccountValidationId,omitempty"`

	// Account holder identification (use FirstName+LastName for individuals, BusinessName otherwise).
	FirstName          string `json:"firstName,omitempty"`
	LastName           string `json:"lastName,omitempty"`
	BusinessName       string `json:"businessName,omitempty"`
	BankAccountHolder  string `json:"bankAccountHolder,omitempty"`
	BusinessEntityType string `json:"businessEntityType,omitempty"`
	// DateOfBirth is optional for individual account holders. Marshalled as an
	// RFC3339 timestamp; the .NET server's default DateTimeConverter accepts
	// either a date-only or full timestamp. Using *time.Time keeps the type
	// consistent with ExchangeRateRequest.Dates.
	DateOfBirth *time.Time `json:"dateOfBirth,omitempty"`

	// Account / routing identifiers — populate the field(s) appropriate for BankNumberType.
	AccountNumber string `json:"accountNumber,omitempty"`
	BankCode      string `json:"bankCode,omitempty"`
	IBAN          string `json:"iban,omitempty"`
	// SwiftCode is the BIC8 / BIC11 (server property name "SwiftCode").
	SwiftCode string `json:"swiftCode,omitempty"`
	CBU       string `json:"cbu,omitempty"`
	CLABE     string `json:"clabe,omitempty"`
	BankGiro  string `json:"bankGiro,omitempty"`
	AccountType string `json:"accountType,omitempty"`

	// Bank metadata (optional).
	BankName         string `json:"bankName,omitempty"`
	LocalBankName    string `json:"localBankName,omitempty"`
	BankCurrencyCode string `json:"bankCurrencyCode,omitempty"`
	EntityCountry    string `json:"entityCountry,omitempty"`

	// Tax identifiers (optional, used for cross-validation).
	TaxIdNumber string `json:"taxIdNumber,omitempty"`
	TaxType     string `json:"taxType,omitempty"`

	// Account holder address (optional).
	AddressLine1 string `json:"addressLine1,omitempty"`
	AddressLine2 string `json:"addressLine2,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	PostalCode   string `json:"postalCode,omitempty"`

	// IfscDetails toggles inclusion of IFSC details in response (default true).
	IfscDetails *bool `json:"ifscDetails,omitempty"`
	// ContinueBankValidationOnTaxFailure controls whether bank pro continues if tax fails.
	ContinueBankValidationOnTaxFailure *bool `json:"continueBankValidationOnTaxFailure,omitempty"`
	// CallbackUrl is an optional callback URL for bulk webhook notifications.
	CallbackUrl string `json:"callbackUrl,omitempty"`
}

// BankAccountResponse is the output from ValidateBankAccount.
//
// Server hierarchy: BankAccountNumberResponse → CanadaBankAccountResponse →
// CbuResponse → ... → BaseValidationResponse. Notably this is a NON-Pro
// response, so it carries `BankAccountValidations` (a list-style detail set
// returned via BankValidationDetail in JSON) and `validationResults`. The Pro
// endpoint uses a different parent (RoutingResponse) — see BankProResponse.
//
// Acronym JSON tags follow .NET CamelCase JsonNamingPolicy (empirically
// verified against System.Text.Json with .NET 10): IBAN→iban,
// CLABENumber→clabeNumber.
type BankAccountResponse struct {
	// Account details
	BankAccountNumber string `json:"bankAccountNumber,omitempty"`
	BankCurrencyCode  string `json:"bankCurrencyCode,omitempty"`
	AccountType       string `json:"accountType,omitempty"`
	AccountHolder     string `json:"accountHolder,omitempty"`
	IBAN              string `json:"iban,omitempty"`
	SwiftCode         string `json:"swiftCode,omitempty"`
	CLABENumber       string `json:"clabeNumber,omitempty"`
	BankGiro          string `json:"bankGiro,omitempty"`

	// Tax/identity
	TaxIdNumber string `json:"taxIdNumber,omitempty"`
	TaxType     string `json:"taxType,omitempty"`

	// Bank information
	BankBranchCode string           `json:"bankBranchCode,omitempty"`
	BankName       string           `json:"bankName,omitempty"`
	BranchName     string           `json:"branchName,omitempty"`
	BankKey        string           `json:"bankKey,omitempty"`
	LocalBankName  string           `json:"localBankName,omitempty"`
	BankAddress    *AddressResponse `json:"bankAddress,omitempty"`
	BankCode       string           `json:"bankCode,omitempty"`
	BankNumberType string           `json:"bankNumberType,omitempty"`
	Country        *Country         `json:"country,omitempty"`

	// Validation details (non-Pro)
	ValidationStatus     string                 `json:"validationStatus,omitempty"`
	BankValidationDetail []BankValidationDetail `json:"bankValidationDetail,omitempty"`

	// AdditionalInfo carries country-specific bank metadata (e.g. UK BACS/CHAPS
	// flags, SE company number, PL request id). Mirrors the canonical
	// AdditionalBankInfo on BankAccountNumberResponse.
	AdditionalInfo *AdditionalBankInfo `json:"additionalInfo,omitempty"`

	// IframeDetails is the optional Inverite iframe handshake (CA bank account
	// pre-validation). Mirrors the inherited InveriteCreateServiceResponse
	// from CanadaBankAccountResponse.
	IframeDetails *InveriteCreateServiceResponse `json:"iframeDetails,omitempty"`

	// BankProValidations is the BankPro analytics block returned by the
	// /api/bank/validate route (BankResponse extends BankAccountNumberResponse
	// with this field). Will be nil if the legacy BankAccountValidationController
	// shape is in use.
	BankProValidations *BankProDetails `json:"bankProValidations,omitempty"`

	// Request/response holder info
	RequestedInfo *BankAccountHolderInfo `json:"requestedInfo,omitempty"`
	ResponseInfo  *BankAccountHolderInfo `json:"responseInfo,omitempty"`

	// Metadata
	SourceUniqueId  string `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId string `json:"qubitOnUniqueId,omitempty"`

	// Validation results (BaseValidationResponse).
	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
}

// AdditionalBankInfo mirrors the canonical AdditionalBankInfo
// (smartvm.BusinessEntities.Client.API.DataLookup.EntityBank).
type AdditionalBankInfo struct {
	DirectDebits  string `json:"directDebits,omitempty"`
	PfsPayments   string `json:"pfsPayments,omitempty"`
	CHAPS         string `json:"chaps,omitempty"`
	BACS          string `json:"bacs,omitempty"`
	CccPayments   string `json:"cccPayments,omitempty"`
	ChapsBIC      string `json:"chapsBIC,omitempty"`
	RequestID     string `json:"requestID,omitempty"`
	CompanyNumber string `json:"companyNumber,omitempty"`
}

// InveriteCreateServiceResponse mirrors the canonical Inverite handshake DTO
// (smartvm.BusinessEntities.Internal.Bases.Bank.BankAccount.InveriteCreateServiceResponse).
// All field names use explicit snake_case JSON tags as defined on the .NET model
// (these survive the CamelCase JsonNamingPolicy because they have explicit
// [JsonPropertyName] attributes).
type InveriteCreateServiceResponse struct {
	Password    string   `json:"password,omitempty"`
	RequestGuid string   `json:"request_guid,omitempty"`
	UserName    string   `json:"username,omitempty"`
	IframeUrl   string   `json:"iframeurl,omitempty"`
	Reused      int      `json:"reused,omitempty"`
	DownloadUrl string   `json:"downloadurl,omitempty"`
	ApiKey      string   `json:"api_key,omitempty"`
	Error       []string `json:"errors,omitempty"`
}

// BankValidationDetail is one validation source/level result.
type BankValidationDetail struct {
	ValidationSource string `json:"validationSource,omitempty"`
	ValidationLevel  string `json:"validationLevel,omitempty"`
	ValidationStatus string `json:"validationStatus,omitempty"`
	ValidationCode   string `json:"validationCode,omitempty"`
	ValidationDate   string `json:"validationDate,omitempty"`
}

// BankAccountValidationDetail contains account and customer validation scores.
type BankAccountValidationDetail struct {
	ValidationSource       string             `json:"validationSource,omitempty"`
	AccountScore           int                `json:"accountScore,omitempty"`
	CustomerScore          int                `json:"customerScore,omitempty"`
	AccountResponse        string             `json:"accountResponse,omitempty"`
	AccountResponseCode    string             `json:"accountResponseCode,omitempty"`
	AccountAddedDate       string             `json:"accountAddedDate,omitempty"`
	AccountLastUpdatedDate string             `json:"accountLastUpdatedDate,omitempty"`
	AccountClosedDate      string             `json:"accountClosedDate,omitempty"`
	CustomerResponseCode   string             `json:"customerResponseCode,omitempty"`
	CustomerResponse       string             `json:"customerResponse,omitempty"`
	IsNameMatchOverridden  bool               `json:"isNameMatchOverridden,omitempty"`
	ValidationResults      []ValidationResult `json:"validationResults,omitempty"`
}

// BankProDetails contains ownership verification and fraud detection analytics.
type BankProDetails struct {
	// Match scores
	VendorMatch        string `json:"vendorMatch,omitempty"`
	AccountHolderMatch string `json:"accountHolderMatch,omitempty"`
	BankNameMatch      string `json:"bankNameMatch,omitempty"`
	BankAccountMatch   string `json:"bankAccountMatch,omitempty"`
	BankCurrencyMatch  string `json:"bankCurrencyMatch,omitempty"`
	BankCountryMatch   string `json:"bankCountryMatch,omitempty"`
	VendorCountryMatch string `json:"vendorCountryMatch,omitempty"`
	AccountTypeMatch   string `json:"accountTypeMatch,omitempty"`
	TaxIDMatch         string `json:"taxIDMatch,omitempty"`
	TaxIDTypeMatch     string `json:"taxIDTypeMatch,omitempty"`

	// Composite scores
	OverallMatchScore       *float64 `json:"overallMatchScore,omitempty"`
	ValidationHistoryScore  *float64 `json:"validationHistoryScore,omitempty"`
	VendorNameMatchScore    *float64 `json:"vendorNameMatchScore,omitempty"`
	AccountHolderMatchScore *float64 `json:"accountHolderMatchScore,omitempty"`
	BankNameMatchScore      *float64 `json:"bankNameMatchScore,omitempty"`

	// History
	PreviousValidationMatch string `json:"previousValidationMatch,omitempty"`
	ApexHistoryMatch        string `json:"apexHistoryMatch,omitempty"`
	DaysSinceLastValidation *int   `json:"daysSinceLastValidation,omitempty"`
	RecentlyValidated       string `json:"recentlyValidated,omitempty"`
	RecentlySeen            string `json:"recentlySeen,omitempty"`
	ValidationAttempts      *int   `json:"validationAttempts,omitempty"`
	ValidationPASSCount     *int   `json:"validationPASSCount,omitempty"`
	ValidationFAILCount     *int   `json:"validationFAILCount,omitempty"`

	// Risk flags
	IsRedFlagsToConsider             *bool    `json:"isRedFlagsToConsider,omitempty"`
	RedFlagReasons                   []string `json:"redFlagReasons,omitempty"`
	FoundBankAccountWithSameVendor   string   `json:"foundBankAccountWithSameVendor,omitempty"`
	FoundBankAccountWithSameTaxID    string   `json:"foundBankAccountWithSameTaxID,omitempty"`
	FoundBankAccountDifferentVendor  string   `json:"foundBankAccountDifferentVendor,omitempty"`
	FoundSameVendorDifferentAccounts string   `json:"foundSameVendorDifferentAccounts,omitempty"`
	FoundSameTaxIDDifferentAccounts  string   `json:"foundSameTaxIDDifferentAccounts,omitempty"`

	// Result descriptions
	ResultDescription     string `json:"resultDescription,omitempty"`
	MatchCodeDescription  string `json:"matchCodeDescription,omitempty"`
	ResultCodeDescription string `json:"resultCodeDescription,omitempty"`

	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
	Errors            string             `json:"errors,omitempty"`
}

// BankAccountHolderInfo contains account holder details for request/response comparison.
type BankAccountHolderInfo struct {
	BankAccountHolder  string `json:"bankAccountHolder,omitempty"`
	BusinessEntityType string `json:"businessEntityType,omitempty"`
	FirstName          string `json:"firstName,omitempty"`
	LastName           string `json:"lastName,omitempty"`
	BusinessName       string `json:"businessName,omitempty"`
	AddressLine1       string `json:"addressLine1,omitempty"`
	AddressLine2       string `json:"addressLine2,omitempty"`
	City               string `json:"city,omitempty"`
	State              string `json:"state,omitempty"`
	PostalCode         string `json:"postalCode,omitempty"`
	TaxIdNumber        string `json:"taxIdNumber,omitempty"`
}

// ── BankPro Validation ────────────────────────────────────────────────────

// BankProRequest is the input for ValidateBankPro. The wire shape matches
// BankAccountRequest; this is a named type (not an alias) so each SDK request
// model corresponds to a single canonical .NET request DTO.
type BankProRequest BankAccountRequest

// BankProResponse is the output from ValidateBankPro.
//
// Server hierarchy: BankAccountNumberProResponse → CanadaProBankAccountResponse
// → RoutingResponse → ... — note this differs from BankAccountResponse, which
// inherits via CbuResponse. The Pro response carries a single
// `BankAccountValidations` detail object plus Pro analytics (BankProValidations,
// Threshold, MatchCode, Score, ResultCode) and explicitly suppresses
// validationResults via [JsonIgnore].
//
// Acronym JSON tags (empirically verified against System.Text.Json with
// .NET 10): IBAN→iban, CLABENumber→clabeNumber, Clabe→clabe.
type BankProResponse struct {
	// Account details
	BankAccountNumber string `json:"bankAccountNumber,omitempty"`
	BankCurrencyCode  string `json:"bankCurrencyCode,omitempty"`
	AccountType       string `json:"accountType,omitempty"`
	AccountHolder     string `json:"accountHolder,omitempty"`
	IBAN              string `json:"iban,omitempty"`
	SwiftCode         string `json:"swiftCode,omitempty"`
	CLABENumber       string `json:"clabeNumber,omitempty"`
	BankGiro          string `json:"bankGiro,omitempty"`
	// Clabe is the Pro-only secondary CLABE field (server property "Clabe").
	Clabe string `json:"clabe,omitempty"`

	// Tax/identity
	TaxIdNumber string `json:"taxIdNumber,omitempty"`
	TaxType     string `json:"taxType,omitempty"`

	// Pro scoring (server returns double; nil when not present).
	Threshold  *float64 `json:"threshold,omitempty"`
	MatchCode  *float64 `json:"matchCode,omitempty"`
	Score      *float64 `json:"score,omitempty"`
	ResultCode *float64 `json:"resultCode,omitempty"`

	// Bank information
	BankBranchCode string           `json:"bankBranchCode,omitempty"`
	BankName       string           `json:"bankName,omitempty"`
	BranchName     string           `json:"branchName,omitempty"`
	BankKey        string           `json:"bankKey,omitempty"`
	LocalBankName  string           `json:"localBankName,omitempty"`
	BankAddress    *AddressResponse `json:"bankAddress,omitempty"`
	BankCode       string           `json:"bankCode,omitempty"`
	BankNumberType string           `json:"bankNumberType,omitempty"`
	Country        *Country         `json:"country,omitempty"`

	// Validation details (Pro: single object + analytics)
	ValidationStatus       string                       `json:"validationStatus,omitempty"`
	BankAccountValidations *BankAccountValidationDetail `json:"bankAccountValidations,omitempty"`
	BankProValidations     *BankProDetails              `json:"bankProValidations,omitempty"`

	// IframeDetails is the optional Inverite iframe handshake (CA Pro). Mirrors
	// the inherited field on CanadaProBankAccountResponse.
	IframeDetails *InveriteCreateServiceResponse `json:"iframeDetails,omitempty"`

	// Request/response holder info
	RequestedInfo *BankAccountHolderInfo `json:"requestedInfo,omitempty"`
	ResponseInfo  *BankAccountHolderInfo `json:"responseInfo,omitempty"`

	// Metadata
	Id              string `json:"id,omitempty"`
	VendorName      string `json:"vendorName,omitempty"`
	VendorCountry   string `json:"vendorCountry,omitempty"`
	SourceUniqueId  string `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId string `json:"qubitOnUniqueId,omitempty"`
	// Note: ValidationResults is intentionally absent — the canonical Pro response
	// uses [JsonIgnore] on its inherited ValidationResults property.
}

// ── Email Validation ──────────────────────────────────────────────────────

// EmailRequest is the input for ValidateEmail.
type EmailRequest struct {
	BaseRequest
	EmailAddress string `json:"emailAddress"`
}

// EmailResponse is the output from ValidateEmail.
// Mirrors the canonical EmailResponse : BaseValidationResponse exactly.
type EmailResponse struct {
	EmailType    string `json:"emailType,omitempty"`
	EmailAddress string `json:"emailAddress,omitempty"`
	FraudScore   *int   `json:"fraudScore,omitempty"`

	// From BaseValidationResponse
	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
	SourceUniqueId    string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId   string             `json:"qubitOnUniqueId,omitempty"`
}

// ── Phone Validation ──────────────────────────────────────────────────────

// PhoneRequest is the input for ValidatePhone. Mirrors the canonical
// PhoneRequest : BaseRequest from smartvm.BusinessEntities.Client.API.Phone.
//
// Configuration toggles (e.g. whether to search for alternate countries) live
// inside Options — the .NET model exposes them as PhoneRequestOptions and
// reads them via SearchForAlternateCountries() defaulting to true. Passing a
// non-nil Options always overrides the server default.
type PhoneRequest struct {
	BaseRequest
	PhoneNumber    string               `json:"phoneNumber"`
	Country        string               `json:"country"`
	PhoneExtension string               `json:"phoneExtension,omitempty"`
	Options        *PhoneRequestOptions `json:"options,omitempty"`
}

// PhoneRequestOptions mirrors the canonical PhoneRequestOptions nested type.
type PhoneRequestOptions struct {
	// SearchForAlternateCountries enables the alternate-country search. The
	// server defaults this to true when Options is omitted.
	SearchForAlternateCountries bool `json:"searchForAlternateCountries,omitempty"`
}

// PhoneResponse is the output from ValidatePhone.
type PhoneResponse struct {
	PhoneNumber             string                  `json:"phoneNumber,omitempty"`
	PhoneExtension          string                  `json:"phoneExtension,omitempty"`
	PhoneCountryCode        string                  `json:"phoneCountryCode,omitempty"`
	FullPhoneNumber         string                  `json:"fullPhoneNumber,omitempty"`
	FraudScore              *int                    `json:"fraudScore,omitempty"`
	Country                 *Country                `json:"country,omitempty"`
	AlternatePhoneCountries []AlternatePhoneCountry `json:"alternatePhoneCountries,omitempty"`

	// Validation (from BaseValidationResponse)
	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
	SourceUniqueId    string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId   string             `json:"qubitOnUniqueId,omitempty"`
}

// AlternatePhoneCountry represents a possible alternate country for a phone number.
type AlternatePhoneCountry struct {
	FullPhoneNumber string   `json:"fullPhoneNumber,omitempty"`
	Country         *Country `json:"country,omitempty"`
}

// ── Business Registration Lookup ──────────────────────────────────────────

// BusinessRegistrationRequest is the input for LookupBusinessRegistration.
type BusinessRegistrationRequest struct {
	BaseRequest
	// EntityName is the business/vendor name. Required.
	EntityName string `json:"entityName"`
	// Country is ISO 3166-1 alpha-2, alpha-3, or full country name. Required.
	Country string `json:"country"`
	// State is optional.
	State string `json:"state,omitempty"`
	// City is optional.
	City string `json:"city,omitempty"`
}

// BusinessRegistrationResponse mirrors the canonical
// BusinessRegistrationResponse : BaseLookupResponse from
// smartvm.BusinessEntities.Client.API.API.BusinessRegistration. The server
// returns an OUTER WRAPPER that carries a list of matching registrations
// (BusinessRegistrations) plus validation metadata; the per-registration
// detail lives on the nested BusinessRegistration type.
type BusinessRegistrationResponse struct {
	// BusinessRegistrations is the list of matching registry records.
	BusinessRegistrations []BusinessRegistration `json:"businessRegistrations,omitempty"`
	// ValidationDescription is a human-readable explanation of the validation outcome.
	ValidationDescription string `json:"validationDescription,omitempty"`
	// ValidationPass is nil when not asserted by the server; true/false otherwise.
	ValidationPass *bool `json:"validationPass,omitempty"`

	// Lookup metadata (BaseLookupResponse)
	Score                float64            `json:"score,omitempty"`
	SourceResultCode     string             `json:"sourceResultCode,omitempty"`
	ValidationResultCode string             `json:"validationResultCode,omitempty"`
	ValidationDate       string             `json:"validationDate,omitempty"`
	SourceUniqueId       string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId      string             `json:"qubitOnUniqueId,omitempty"`
	ValidationResults    []ValidationResult `json:"validationResults,omitempty"`
}

// BusinessRegistration mirrors the canonical BusinessRegistration nested
// inside BusinessRegistrationResponse — one record per matching government
// registry hit, with optional nested addresses/persons/phones.
type BusinessRegistration struct {
	RegistrationId                string                        `json:"registrationId,omitempty"`
	EntityName                    string                        `json:"entityName,omitempty"`
	Status                        string                        `json:"status,omitempty"`
	StatusReason                  string                        `json:"statusReason,omitempty"`
	BusinessEntityType            string                        `json:"businessEntityType,omitempty"`
	BusinessEntityTypeDescription string                        `json:"businessEntityTypeDescription,omitempty"`
	Jurisdiction                  string                        `json:"jurisdiction,omitempty"`
	RegistrationDate              string                        `json:"registrationDate,omitempty"`
	FormationDate                 string                        `json:"formationDate,omitempty"`
	ExpirationDate                string                        `json:"expirationDate,omitempty"`
	Duration                      string                        `json:"duration,omitempty"`
	TaxNumber                     string                        `json:"taxNumber,omitempty"`
	PreviousEntityNames           []string                      `json:"previousEntityNames,omitempty"`
	Addresses                     []BusinessRegistrationAddress `json:"addresses,omitempty"`
	Persons                       []BusinessRegistrationPerson  `json:"persons,omitempty"`
	Phones                        []BusinessRegistrationPhone   `json:"phones,omitempty"`
}

// BusinessRegistrationAddress is an address associated with a registered business.
type BusinessRegistrationAddress struct {
	AddressType  string   `json:"addressType,omitempty"`
	AddressLine1 string   `json:"addressLine1,omitempty"`
	AddressLine2 string   `json:"addressLine2,omitempty"`
	City         string   `json:"city,omitempty"`
	State        string   `json:"state,omitempty"`
	Zip          string   `json:"zip,omitempty"`
	Country      *Country `json:"country,omitempty"`
}

// BusinessRegistrationPerson is an officer, director, or agent associated with a business.
type BusinessRegistrationPerson struct {
	Name        string                        `json:"name,omitempty"`
	Designation string                        `json:"designation,omitempty"`
	Addresses   []BusinessRegistrationAddress `json:"addresses,omitempty"`
}

// BusinessRegistrationPhone is a phone number associated with a registered business.
type BusinessRegistrationPhone struct {
	PhoneType   string `json:"phoneType,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

// ── Peppol Validation ─────────────────────────────────────────────────────

// PeppolRequest is the input for ValidatePeppol.
type PeppolRequest struct {
	BaseRequest
	ParticipantId   string `json:"participantId"`
	DirectoryLookup *bool  `json:"directoryLookup,omitempty"`
}

// PeppolResponse is the output from ValidatePeppol.
type PeppolResponse struct {
	IsValid               bool     `json:"isValid"`
	ValidationType        string   `json:"validationType,omitempty"`
	ParticipantId         string   `json:"participantId,omitempty"`
	IcdCode               string   `json:"icdCode,omitempty"`
	IcdScheme             string   `json:"icdScheme,omitempty"`
	Identifier            string   `json:"identifier,omitempty"`
	CountryIso2           string   `json:"countryIso2,omitempty"`
	RegisteredInDirectory *bool    `json:"registeredInDirectory,omitempty"`
	ParticipantName       string   `json:"participantName,omitempty"`
	DocumentTypes         []string `json:"documentTypes,omitempty"`
	Errors                []string `json:"errors,omitempty"`
}

// SupportedPeppolScheme is one entry in the GetPeppolSchemes array.
type SupportedPeppolScheme struct {
	IcdCode     string `json:"icdCode,omitempty"`
	IcdScheme   string `json:"icdScheme,omitempty"`
	CountryIso2 string `json:"countryIso2,omitempty"`
	Description string `json:"description,omitempty"`
}

// ── Sanctions Screening ──────────────────────────────────────────────────

// SanctionsRequest is the input for CheckSanctions. Maps to the canonical
// ProhibitedListRequest (POST /api/prohibited/lookup).
type SanctionsRequest struct {
	BaseRequest
	// CompanyName is the entity name to screen. Required.
	CompanyName string `json:"companyName"`
	// CompanyNameDBA contains "doing business as" name aliases (optional).
	CompanyNameDBA []string `json:"companyNameDBA,omitempty"`
	// FirstName, MiddleName, LastName for individual screening (optional).
	FirstName  string `json:"firstName,omitempty"`
	MiddleName string `json:"middleName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
	// IdentityNumber for additional matching (optional).
	IdentityNumber string `json:"identityNumber,omitempty"`
	// BusinessEntityType is optional.
	BusinessEntityType string `json:"businessEntityType,omitempty"`
	// RequestType: "Initial", "Ongoing", or "Re-check".
	RequestType string `json:"requestType,omitempty"`
	// Threshold: confidence threshold (0.0–1.0). Defaults to 0.0.
	Threshold *float64 `json:"threshold,omitempty"`
	// IdentityId is a list of identity type/number pairs for screening.
	IdentityId []SanctionsIdentity `json:"identityId,omitempty"`
	// Addresses is a list of addresses associated with the entity for screening.
	Addresses []SanctionsAddress `json:"addresses,omitempty"`
	// ProhibitedListSearchList narrows screening to specific lists.
	ProhibitedListSearchList []string `json:"prohibitedListSearchList,omitempty"`
	// ProhibitedLookupType: default "Entity"; use "Individual" for person screening.
	ProhibitedLookupType string `json:"prohibitedLookupType,omitempty"`
	// BlockedCountryUNNumber is an optional UN country code (3 chars). Server field is
	// "BlockedCountryUN_Number" — case-sensitive snake-style. Verified against the
	// canonical ProhibitedListRequest in the .NET project.
	BlockedCountryUNNumber string `json:"blockedCountryUN_Number,omitempty"`
	// AddressLine1/2/3, City, State, PostalCode, Country are legacy/optional address fields.
	AddressLine1 string `json:"addressLine1,omitempty"`
	AddressLine2 string `json:"addressLine2,omitempty"`
	AddressLine3 string `json:"addressLine3,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	PostalCode   string `json:"postalCode,omitempty"`
	Country      string `json:"country,omitempty"`
}

// SanctionsIdentity is an identity type/number pair for prohibited list screening.
type SanctionsIdentity struct {
	IdentityType   string `json:"identityType,omitempty"`
	IdentityNumber string `json:"identityNumber,omitempty"`
}

// SanctionsAddress is an address associated with an entity for prohibited list screening.
type SanctionsAddress struct {
	AddressLine1 string `json:"addressLine1,omitempty"`
	AddressLine2 string `json:"addressLine2,omitempty"`
	AddressLine3 string `json:"addressLine3,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	PostalCode   string `json:"postalCode,omitempty"`
	Country      string `json:"country,omitempty"`
}

// SanctionsResponse is one entry in the array returned by CheckSanctions.
// One entry per matching list/source.
type SanctionsResponse struct {
	IsMatch              bool                     `json:"isMatch"`
	Description          string                   `json:"description,omitempty"`
	AdditionalInfo       *SanctionsAdditionalInfo `json:"additionalInfo,omitempty"`
	SourceId             int                      `json:"sourceId,omitempty"`
	Score                float64                  `json:"score,omitempty"`
	SourceResultCode     string                   `json:"sourceResultCode,omitempty"`
	ValidationResultCode string                   `json:"validationResultCode,omitempty"`
	ValidationDate       string                   `json:"validationDate,omitempty"`
	SourceUniqueId       string                   `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId      string                   `json:"qubitOnUniqueId,omitempty"`
}

// SanctionsAdditionalInfo contains the detailed match results per sanctions list.
type SanctionsAdditionalInfo struct {
	ResultCode string                `json:"resultCode,omitempty"`
	Lists      []SanctionsListDetail `json:"lists,omitempty"`
}

// SanctionsListDetail represents matches from a specific sanctions list.
type SanctionsListDetail struct {
	Code                  string            `json:"code,omitempty"`
	Name                  string            `json:"name,omitempty"`
	Description           string            `json:"description,omitempty"`
	TotalMatches          *int              `json:"totalMatches,omitempty"`
	EntityMatches         *int              `json:"entityMatches,omitempty"`
	BlockedCountryMatches *int              `json:"blockedCountryMatches,omitempty"`
	DateAdded             string            `json:"dateAdded,omitempty"`
	Entities              []SanctionsEntity `json:"entities,omitempty"`
}

// SanctionsEntity represents a matched entity on a sanctions list.
type SanctionsEntity struct {
	Id              string                       `json:"id,omitempty"`
	Number          string                       `json:"number,omitempty"`
	MaxScore        *float64                     `json:"maxScore,omitempty"`
	FirstName       string                       `json:"firstName,omitempty"`
	LastName        string                       `json:"lastName,omitempty"`
	MiddleName      string                       `json:"middleName,omitempty"`
	OtherName       string                       `json:"otherName,omitempty"`
	WholeName       string                       `json:"wholeName,omitempty"`
	HasAliasMatch   *bool                        `json:"hasAliasMatch,omitempty"`
	HasAddressMatch *bool                        `json:"hasAddressMatch,omitempty"`
	Score           *float64                     `json:"score,omitempty"`
	Type            string                       `json:"type,omitempty"`
	Programs        string                       `json:"programs,omitempty"`
	Remarks         string                       `json:"remarks,omitempty"`
	Aliases         []SanctionsEntityAlias       `json:"aliases,omitempty"`
	Addresses       []SanctionsEntityAddress     `json:"addresses,omitempty"`
	Identifiers     []SanctionsEntityIdentifier  `json:"identifiers,omitempty"`
	Relationships   []SanctionsEntityRelationship `json:"relationships,omitempty"`
	Countries       []string                     `json:"countries,omitempty"`
}

// SanctionsEntityIdentifier is an identifier (passport, national ID, etc.)
// for a sanctions-listed entity.
type SanctionsEntityIdentifier struct {
	Type      string   `json:"type,omitempty"`
	Number    string   `json:"number,omitempty"`
	Country   *Country `json:"country,omitempty"`
	IssueDate string   `json:"issueDate,omitempty"`
	ExpiresOn string   `json:"expiresOn,omitempty"`
	Note      string   `json:"note,omitempty"`
}

// SanctionsEntityRelationship is a related entity (e.g., owner, subsidiary).
type SanctionsEntityRelationship struct {
	RelatedEntityId   string  `json:"relatedEntityId,omitempty"`
	RelatedEntityName string  `json:"relatedEntityName,omitempty"`
	Type              string  `json:"type,omitempty"`
	Description       string  `json:"description,omitempty"`
	Score             *float64 `json:"score,omitempty"`
}

// SanctionsEntityAlias is an alias for a sanctions entity.
type SanctionsEntityAlias struct {
	Name  string   `json:"name,omitempty"`
	Type  string   `json:"type,omitempty"`
	Score *float64 `json:"score,omitempty"`
}

// SanctionsEntityAddress is an address for a sanctions entity.
type SanctionsEntityAddress struct {
	Address1   string   `json:"address1,omitempty"`
	Address2   string   `json:"address2,omitempty"`
	Address3   string   `json:"address3,omitempty"`
	City       string   `json:"city,omitempty"`
	State      string   `json:"state,omitempty"`
	StateName  string   `json:"stateName,omitempty"`
	PostalCode string   `json:"postalCode,omitempty"`
	Country    *Country `json:"country,omitempty"`
	Score      *float64 `json:"score,omitempty"`
	Remark     string   `json:"remark,omitempty"`
}

// ── PEP Screening ─────────────────────────────────────────────────────────

// PEPRequest is the input for ScreenPEP.
type PEPRequest struct {
	BaseRequest
	Name    string `json:"name"`
	Country string `json:"country"`
}

// PEPResponse is one entry in the array returned by ScreenPEP. One entry per
// matching source/dataset.
type PEPResponse struct {
	Persons       []PEPPerson       `json:"persons,omitempty"`
	Organizations []PEPOrganization `json:"organizations,omitempty"`
	Members       []PEPMember       `json:"members,omitempty"`
	Areas         []PEPArea         `json:"areas,omitempty"`

	// Lookup metadata (from BaseLookupResponse)
	Score                float64 `json:"score,omitempty"`
	SourceResultCode     string  `json:"sourceResultCode,omitempty"`
	ValidationResultCode string  `json:"validationResultCode,omitempty"`
	ValidationDate       string  `json:"validationDate,omitempty"`
	SourceUniqueId       string  `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId      string  `json:"qubitOnUniqueId,omitempty"`
}

// PEPPerson represents a politically exposed person.
type PEPPerson struct {
	Id                string             `json:"id,omitempty"`
	Name              string             `json:"name,omitempty"`
	SortName          string             `json:"sortName,omitempty"`
	GivenName         string             `json:"givenName,omitempty"`
	FamilyName        string             `json:"familyName,omitempty"`
	HonorificPrefix   string             `json:"honorificPrefix,omitempty"`
	Gender            string             `json:"gender,omitempty"`
	BirthDate         string             `json:"birthDate,omitempty"`
	DeathDate         string             `json:"deathDate,omitempty"`
	Email             string             `json:"email,omitempty"`
	Image             string             `json:"image,omitempty"`
	ContactDetails    []PEPContactDetail `json:"contactDetails,omitempty"`
	PersonIdentifiers []PEPIdentifier    `json:"personIdentifiers,omitempty"`
	PersonLinks       []PEPLink          `json:"personLinks,omitempty"`
	PersonOtherNames  []PEPOtherName     `json:"personOtherNames,omitempty"`
	PersonImages      []PEPImage         `json:"personImages,omitempty"`
}

// PEPOrganization represents an organization associated with a PEP.
type PEPOrganization struct {
	Id                      string          `json:"id,omitempty"`
	Name                    string          `json:"name,omitempty"`
	Classification          string          `json:"classification,omitempty"`
	OrganizationType        string          `json:"organizationType,omitempty"`
	Seats                   *int            `json:"seats,omitempty"`
	Image                   string          `json:"image,omitempty"`
	OrganizationIdentifiers []PEPIdentifier `json:"organizationIdentifiers,omitempty"`
	OrganizationLinks       []PEPLink       `json:"organizationLinks,omitempty"`
	OrganizationOtherNames  []PEPOtherName  `json:"organizationOtherNames,omitempty"`
}

// PEPMember represents a membership record linking persons to organizations.
type PEPMember struct {
	PersonId            string    `json:"personId,omitempty"`
	OrganizationId      string    `json:"organizationId,omitempty"`
	AreaId              string    `json:"areaId,omitempty"`
	Role                string    `json:"role,omitempty"`
	StartDate           string    `json:"startDate,omitempty"`
	EndDate             string    `json:"endDate,omitempty"`
	LegislativePeriodId string    `json:"legislativePeriodId,omitempty"`
	OnBehalfOfId        string    `json:"onBehalfOfId,omitempty"`
	Sources             []PEPLink `json:"sources,omitempty"`
}

// PEPArea represents a geographic area associated with PEP data.
type PEPArea struct {
	Id              string          `json:"id,omitempty"`
	Name            string          `json:"name,omitempty"`
	AreaType        string          `json:"areaType,omitempty"`
	AreaIdentifiers []PEPIdentifier `json:"areaIdentifiers,omitempty"`
	AreaOtherNames  []PEPOtherName  `json:"areaOtherNames,omitempty"`
}

// PEPContactDetail is a contact entry for a PEP person.
type PEPContactDetail struct {
	ContactType string `json:"contactType,omitempty"`
	Value       string `json:"value,omitempty"`
}

// PEPIdentifier is an identifier (scheme + value) for PEP entities.
type PEPIdentifier struct {
	Identifier string `json:"identifier,omitempty"`
	Scheme     string `json:"scheme,omitempty"`
}

// PEPLink is a URL link with optional note. Acronym JSON tag URL→url per
// .NET CamelCase JsonNamingPolicy (empirically verified against
// System.Text.Json with .NET 10 — consecutive leading capitals are fully
// lower-cased).
type PEPLink struct {
	URL  string `json:"url,omitempty"`
	Note string `json:"note,omitempty"`
}

// PEPOtherName is an alternate name for a PEP entity.
type PEPOtherName struct {
	Name     string `json:"name,omitempty"`
	Language string `json:"language,omitempty"`
	Note     string `json:"note,omitempty"`
	Source   string `json:"source,omitempty"`
}

// PEPImage is an image URL for a PEP person. Acronym JSON tag URL→url.
type PEPImage struct {
	URL string `json:"url,omitempty"`
}

// ── Directors Check ───────────────────────────────────────────────────────

// DirectorsRequest is the input for CheckDirectors.
type DirectorsRequest struct {
	BaseRequest
	FirstName  string `json:"firstName,omitempty"`
	MiddleName string `json:"middleName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
	// Country is ISO 3166-1 alpha-2 (exactly 2 characters). Required.
	Country string `json:"country"`
}

// DirectorsResponse is the output from CheckDirectors. Mirrors the canonical
// DisqualifiedDirectorResponse : BaseValidationResponse from
// smartvm.BusinessEntities.Client.API.DisqualifiedDirectors.
type DirectorsResponse struct {
	// TotalResults is the total number of disqualified-director records found.
	TotalResults int `json:"totalResults,omitempty"`
	// Roles lists each disqualification record. Mirrors DisqualifiedRoles.
	Roles []DisqualifiedRole `json:"roles,omitempty"`

	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
	SourceUniqueId    string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId   string             `json:"qubitOnUniqueId,omitempty"`
}

// DisqualifiedRole mirrors the canonical DisqualifiedRoles entry.
type DisqualifiedRole struct {
	FirstName                string                    `json:"firstName,omitempty"`
	MiddleName               string                    `json:"middleName,omitempty"`
	LastName                 string                    `json:"lastName,omitempty"`
	DisqualifiedDirectorId   int                       `json:"disqualifiedDirectorId,omitempty"`
	Addresses                []AddressResponse         `json:"addresses,omitempty"`
	Associations             *DisqualifiedAssociations `json:"associations,omitempty"`
	DisqualificationCriteria *DisqualificationCriteria `json:"disqualificationCriteria,omitempty"`
	Aliases                  *DisqualifiedAliases      `json:"aliases,omitempty"`
}

// DisqualifiedAssociations mirrors the canonical DisqualifiedAssociations.
// Field set is loose because the .NET model nests further per-jurisdiction
// data; consumers can also marshal the raw JSON.
type DisqualifiedAssociations struct {
	Associations []map[string]interface{} `json:"associations,omitempty"`
}

// DisqualificationCriteria mirrors the canonical NZDisqualificationCriteria
// from smartvm.BusinessEntities.Client.API.DisqualifiedDirectors. The model
// is an outer wrapper with a list of Criterion entries (the .NET property
// is lowercase "criteria" — preserved verbatim under CamelCase since the
// leading character is already lowercase).
type DisqualificationCriteria struct {
	Criteria []Criterion `json:"criteria,omitempty"`
}

// Criterion mirrors the canonical Criterion entry inside
// DisqualificationCriteria.
type Criterion struct {
	StartDate string `json:"startDate,omitempty"`
	Criteria  string `json:"criteria,omitempty"`
	EndDate   string `json:"endDate,omitempty"`
	Comments  string `json:"comments,omitempty"`
}

// DisqualifiedAliases mirrors the canonical DisqualifiedAliases.
type DisqualifiedAliases struct {
	Aliases []string `json:"aliases,omitempty"`
}

// ── EPA Prosecution ───────────────────────────────────────────────────────

// EPARequest is the input for CheckEPAProsecution and LookupEPAProsecution.
type EPARequest struct {
	BaseRequest
	Name       string `json:"name,omitempty"`
	State      string `json:"state,omitempty"`
	FiscalYear string `json:"fiscalYear,omitempty"`
}

// EPAResponse is one entry in the array returned by EPA endpoints.
type EPAResponse struct {
	ProsecutionSummaryId int64  `json:"prosecutionSummaryId,omitempty"`
	Name                 string `json:"name,omitempty"`
	State                string `json:"state,omitempty"`
	Year                 string `json:"year,omitempty"`
	Action               string `json:"action,omitempty"`

	// Validation (from BaseValidationResponse)
	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
	SourceUniqueId    string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId   string             `json:"qubitOnUniqueId,omitempty"`
}

// ── Healthcare Exclusion ──────────────────────────────────────────────────

// HealthcareExclusionRequest is the input for CheckHealthcareExclusion and LookupHealthcareExclusion.
type HealthcareExclusionRequest struct {
	BaseRequest
	HealthCareType string `json:"healthCareType"` // HCO or HCP
	EntityName     string `json:"entityName,omitempty"`
	LastName       string `json:"lastName,omitempty"`
	FirstName      string `json:"firstName,omitempty"`
	Address        string `json:"address,omitempty"`
	City           string `json:"city,omitempty"`
	State          string `json:"state,omitempty"`
	ZipCode        string `json:"zipCode,omitempty"`
}

// HealthcareExclusionResponse is one entry in the array returned by
// CheckHealthcareExclusion (POST /api/providerexclusion/validate). Mirrors the
// canonical HealthCareResponse : BaseValidationResponse exactly. Acronym JSON
// tags follow .NET CamelCase JsonNamingPolicy (empirically verified against
// System.Text.Json with .NET 10 — consecutive leading capitals are fully
// lower-cased): UPIN→upin, NPI→npi.
type HealthcareExclusionResponse struct {
	LastName           string           `json:"lastName,omitempty"`
	FirstName          string           `json:"firstName,omitempty"`
	MiddleName         string           `json:"middleName,omitempty"`
	BusinessName       string           `json:"businessName,omitempty"`
	General            string           `json:"general,omitempty"`
	Speciality         string           `json:"speciality,omitempty"`
	UPIN               string           `json:"upin,omitempty"`
	NPI                string           `json:"npi,omitempty"`
	DateOfBirth        string           `json:"dateOfBirth,omitempty"`
	HealthCareAddress  *AddressResponse `json:"healthCareAddress,omitempty"`
	Address            string           `json:"address,omitempty"`
	City               string           `json:"city,omitempty"`
	State              string           `json:"state,omitempty"`
	ZipCode            string           `json:"zipCode,omitempty"`
	ExclusionType      string           `json:"exclusionType,omitempty"`
	ExclusionDate      string           `json:"exclusionDate,omitempty"`
	StatusUpdateTime   string           `json:"statusUpdateTime,omitempty"`
	ReinstatedDate     string           `json:"reinstatedDate,omitempty"`
	WaiverDate         string           `json:"waiverDate,omitempty"`
	WaiverState        string           `json:"waiverState,omitempty"`
	IsProviderExcluded bool             `json:"isProviderExcluded,omitempty"`
	NetworkEntityId    int64            `json:"networkEntityId,omitempty"`

	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
	SourceUniqueId    string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId   string             `json:"qubitOnUniqueId,omitempty"`
}

// HealthcareExclusionLookupResponse is the single-object response for
// LookupHealthcareExclusion (POST /api/providerexclusion/lookup). Mirrors the
// canonical HealthCareLookupResponse : BaseLookupResponse exactly.
type HealthcareExclusionLookupResponse struct {
	HealthCareLookupEntity []HealthCareEntity `json:"healthCareLookupEntity,omitempty"`

	Score                float64            `json:"score,omitempty"`
	SourceResultCode     string             `json:"sourceResultCode,omitempty"`
	ValidationResultCode string             `json:"validationResultCode,omitempty"`
	ValidationDate       string             `json:"validationDate,omitempty"`
	ValidationResults    []ValidationResult `json:"validationResults,omitempty"`
	SourceUniqueId       string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId      string             `json:"qubitOnUniqueId,omitempty"`
}

// HealthCareEntity mirrors the canonical HealthCareLookupResponse.HealthCareEntity
// nested type. Acronym JSON tags follow .NET CamelCase JsonNamingPolicy.
type HealthCareEntity struct {
	LastName           string           `json:"lastName,omitempty"`
	FirstName          string           `json:"firstName,omitempty"`
	MiddleName         string           `json:"middleName,omitempty"`
	BusinessName       string           `json:"businessName,omitempty"`
	General            string           `json:"general,omitempty"`
	Speciality         string           `json:"speciality,omitempty"`
	UPIN               string           `json:"upin,omitempty"`
	NPI                string           `json:"npi,omitempty"`
	DateOfBirth        string           `json:"dateOfBirth,omitempty"`
	HealthCareAddress  *AddressResponse `json:"healthCareAddress,omitempty"`
	Address            string           `json:"address,omitempty"`
	City               string           `json:"city,omitempty"`
	State              string           `json:"state,omitempty"`
	ZipCode            string           `json:"zipCode,omitempty"`
	ExclusionType      string           `json:"exclusionType,omitempty"`
	ExclusionDate      string           `json:"exclusionDate,omitempty"`
	ReinstatedDate     string           `json:"reinstatedDate,omitempty"`
	WaiverDate         string           `json:"waiverDate,omitempty"`
	WaiverState        string           `json:"waiverState,omitempty"`
	IsProviderExcluded bool             `json:"isProviderExcluded,omitempty"`
}

// ── Risk Lookup (adverse media) ───────────────────────────────────────────

// RiskLookupRequest is the input for LookupRisk (/api/risk/lookup).
// Server only accepts categories Social, Governance, Environmental.
type RiskLookupRequest struct {
	BaseRequest
	EntityName          string `json:"entityName"`
	Category            string `json:"category"`
	Country             string `json:"country,omitempty"`
	MaximumArticleCount *int   `json:"maximumArticleCount,omitempty"`
	AddressLine1        string `json:"addressLine1,omitempty"`
	AddressLine2        string `json:"addressLine2,omitempty"`
	City                string `json:"city,omitempty"`
	State               string `json:"state,omitempty"`
	PostalCode          string `json:"postalCode,omitempty"`
}

// RiskResponse is one entry in the array returned by /api/risk/lookup.
// Mirrors the canonical RiskResponse : BaseLookupResponse.
type RiskResponse struct {
	TotalArticles       float64            `json:"totalArticles,omitempty"`
	TotalSentimentCount map[string]float64 `json:"totalSentimentCount,omitempty"`
	Categories          []RiskCategory     `json:"categories,omitempty"`

	// Lookup metadata (BaseLookupResponse)
	Score                float64 `json:"score,omitempty"`
	SourceResultCode     string  `json:"sourceResultCode,omitempty"`
	ValidationResultCode string  `json:"validationResultCode,omitempty"`
	ValidationDate       string  `json:"validationDate,omitempty"`
	SourceUniqueId       string  `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId      string  `json:"qubitOnUniqueId,omitempty"`
}

// RiskCategory mirrors the canonical Category type for risk responses.
type RiskCategory struct {
	Name           string             `json:"name,omitempty"`
	SentimentCount map[string]float64 `json:"sentimentCount,omitempty"`
	TotalArticles  float64            `json:"totalArticles,omitempty"`
	Articles       []RiskArticle      `json:"articles,omitempty"`
}

// RiskArticle is one adverse-media article result. Acronym JSON tag URL→url.
type RiskArticle struct {
	Title           string  `json:"title,omitempty"`
	URL             string  `json:"url,omitempty"`
	Source          string  `json:"source,omitempty"`
	PublicationDate string  `json:"publicationDate,omitempty"`
	Snippet         string  `json:"snippet,omitempty"`
	Score           float64 `json:"score,omitempty"`
	Sentiment       string  `json:"sentiment,omitempty"`
}

// ── Risk Control (Bankruptcy / Credit Score / Fail Rate) ──────────────────

// BankruptcyRequest is the input for CheckBankruptcyRisk.
type BankruptcyRequest struct {
	BaseRequest
	CompanyName string `json:"companyName"`
	Country     string `json:"country"`
}

// BankruptcyResponse mirrors the canonical RiskControlResponse for
// category=Bankruptcy. The server returns a single object with cases and a
// recommendation; RiskScore is a string on the wire.
type BankruptcyResponse struct {
	CompanyName       string             `json:"companyname,omitempty"`
	Country           string             `json:"country,omitempty"`
	Category          string             `json:"category,omitempty"`
	Cases             []RiskControlCase  `json:"cases,omitempty"`
	Recommendation    string             `json:"recommendation,omitempty"`
	RiskScore         string             `json:"riskScore,omitempty"`
	CategoryValue     string             `json:"categoryValue,omitempty"`
	QubitOnNumber     *int64             `json:"qubitOnNumber,omitempty"`
	Score             float64            `json:"score,omitempty"`
	SourceUniqueId    string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId   string             `json:"qubitOnUniqueId,omitempty"`
	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
}

// CreditScoreRequest is the input for LookupCreditScore.
type CreditScoreRequest struct {
	BaseRequest
	CompanyName string `json:"companyName"`
	Country     string `json:"country"`
}

// CreditScoreResponse mirrors RiskControlResponse for category=Credit Score.
type CreditScoreResponse struct {
	CompanyName       string             `json:"companyname,omitempty"`
	Country           string             `json:"country,omitempty"`
	Category          string             `json:"category,omitempty"`
	Cases             []RiskControlCase  `json:"cases,omitempty"`
	Recommendation    string             `json:"recommendation,omitempty"`
	RiskScore         string             `json:"riskScore,omitempty"`
	CategoryValue     string             `json:"categoryValue,omitempty"`
	QubitOnNumber     *int64             `json:"qubitOnNumber,omitempty"`
	Score             float64            `json:"score,omitempty"`
	SourceUniqueId    string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId   string             `json:"qubitOnUniqueId,omitempty"`
	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
}

// FailRateRequest is the input for LookupFailRate.
type FailRateRequest struct {
	BaseRequest
	CompanyName string `json:"companyName"`
	Country     string `json:"country"`
}

// FailRateResponse mirrors RiskControlResponse for category=Fail Rate.
type FailRateResponse struct {
	CompanyName       string             `json:"companyname,omitempty"`
	Country           string             `json:"country,omitempty"`
	Category          string             `json:"category,omitempty"`
	Cases             []RiskControlCase  `json:"cases,omitempty"`
	Recommendation    string             `json:"recommendation,omitempty"`
	RiskScore         string             `json:"riskScore,omitempty"`
	CategoryValue     string             `json:"categoryValue,omitempty"`
	QubitOnNumber     *int64             `json:"qubitOnNumber,omitempty"`
	Score             float64            `json:"score,omitempty"`
	SourceUniqueId    string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId   string             `json:"qubitOnUniqueId,omitempty"`
	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
}

// RiskControlCase is one case row inside a RiskControlResponse.
//
// JSON tags use the wire-format names (with embedded spaces) from the
// canonical Cases DTO. The .NET model uses
// `[JsonPropertyName("case number")]`, `[JsonPropertyName("case name")]`,
// `[JsonPropertyName("case type")]` overrides — encoding/json supports
// embedded spaces in struct tags, so the wire compatibility is exact.
//
// Pitfall: linters and code generators frequently warn about or rewrite
// JSON tags containing whitespace. Do NOT "fix" these tags by removing
// the embedded spaces or replacing them with camelCase — the wire
// contract from System.Text.Json explicitly transmits "case number" /
// "case name" / "case type" with literal spaces, and any rewrite will
// silently drop these fields on decode. If your linter complains, suppress
// it on this struct rather than touching the tag values.
type RiskControlCase struct {
	CaseNumber  string `json:"case number,omitempty"`
	CaseName    string `json:"case name,omitempty"`
	CaseType    string `json:"case type,omitempty"`
	CaseDate    string `json:"date,omitempty"`
	Status      string `json:"status,omitempty"`
	Location    string `json:"location,omitempty"`
	Source      string `json:"source,omitempty"`
	Description string `json:"description,omitempty"`
}

// ── Entity Risk Assessment ────────────────────────────────────────────────

// EntityRiskRequest is the input for AssessEntityRisk. Mirrors the canonical
// EntityRiskRequest : BaseRequest, IEntityHeader. Acronym JSON tag URL→url
// per .NET CamelCase JsonNamingPolicy. CountryOfIncorporation has an explicit
// [JsonPropertyName("CountryOfIncorporation")] override on the .NET model so
// the wire field preserves the leading capital.
type EntityRiskRequest struct {
	BaseRequest
	CompanyName            string                       `json:"companyName"`
	CompanyNameDBA         []string                     `json:"companyNameDBA,omitempty"`
	QubitOnEntityId        string                       `json:"qubitOnEntityId,omitempty"`
	BusinessEntityType     string                       `json:"businessEntityType,omitempty"`
	BusinessEntitySubType  string                       `json:"businessEntitySubType,omitempty"`
	MaximumArticleCount    *int                         `json:"maximumArticleCount,omitempty"`
	Category               string                       `json:"category,omitempty"`
	URL                    []string                     `json:"url,omitempty"`
	CountryOfIncorporation string                       `json:"CountryOfIncorporation,omitempty"`
	StateOfIncorporation   string                       `json:"stateOfIncorporation,omitempty"`
	YearStarted            string                       `json:"yearStarted,omitempty"`
	SalesVolume            *float64                     `json:"salesVolume,omitempty"`
	EmployeesOnSite        *float64                     `json:"employeesOnSite,omitempty"`
	EmployeesTotal         *float64                     `json:"employeesTotal,omitempty"`
	LocationType           string                       `json:"locationType,omitempty"`
	LineOfBusiness         string                       `json:"lineOfBusiness,omitempty"`
	EntityAddresses        []EntityAddressRiskRequest   `json:"entityAddresses,omitempty"`
	EntityPhones           []EntityPhoneRiskRequest     `json:"entityPhones,omitempty"`
	EntityIdentities       []EntityIdentityRiskRequest  `json:"entityIdentities,omitempty"`
	EntityPersons          []EntityPersonRiskRequest    `json:"entityPersons,omitempty"`
}

// EntityAddressRiskRequest is one address attached to an EntityRiskRequest.
type EntityAddressRiskRequest struct {
	AddressLine1 string `json:"addressLine1,omitempty"`
	AddressLine2 string `json:"addressLine2,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	PostalCode   string `json:"postalCode,omitempty"`
	Country      string `json:"country,omitempty"`
}

// EntityPhoneRiskRequest is one phone attached to an EntityRiskRequest.
type EntityPhoneRiskRequest struct {
	PhoneNumber string `json:"phoneNumber,omitempty"`
	PhoneType   string `json:"phoneType,omitempty"`
}

// EntityIdentityRiskRequest is one identity-document entry attached to an EntityRiskRequest.
type EntityIdentityRiskRequest struct {
	IdentityType   string `json:"identityType,omitempty"`
	IdentityNumber string `json:"identityNumber,omitempty"`
}

// EntityPersonRiskRequest is one person attached to an EntityRiskRequest.
type EntityPersonRiskRequest struct {
	FirstName  string `json:"firstName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
	MiddleName string `json:"middleName,omitempty"`
}

// EntityRiskResponse is the output from AssessEntityRisk.
// Mirrors the canonical EntityRiskResponse : BaseEntityRiskResponse.
type EntityRiskResponse struct {
	// Company identity
	CompanyName           string   `json:"companyName,omitempty"`
	Name                  string   `json:"name,omitempty"`
	QubitOnCompanyName    string   `json:"qubitOnCompanyName,omitempty"`
	CompanyNameDBA        []string `json:"companyNameDBA,omitempty"`
	QubitOnCompanyNameDBA []string `json:"qubitOnCompanyNameDBA,omitempty"`
	QubitOnEntityId       string   `json:"qubitOnEntityId,omitempty"`

	// Entity classification
	BusinessEntityType        string `json:"businessEntityType,omitempty"`
	QubitOnBusinessEntityType string `json:"qubitOnBusinessEntityType,omitempty"`
	BusinessEntitySubType     string `json:"businessEntitySubType,omitempty"`
	CountryOfIncorporation    string `json:"CountryOfIncorporation,omitempty"`
	StateOfIncorporation      string `json:"StateOfIncorporation,omitempty"`

	// Operations. Acronym JSON tag URL→url.
	URL             []string `json:"url,omitempty"`
	YearStarted     string   `json:"yearStarted,omitempty"`
	SalesVolume     *float64 `json:"salesVolume,omitempty"`
	EmployeesOnSite *float64 `json:"employeesOnSite,omitempty"`
	EmployeesTotal  *float64 `json:"employeesTotal,omitempty"`
	LocationType    string   `json:"locationType,omitempty"`
	LineOfBusiness  string   `json:"lineOfBusiness,omitempty"`

	// Risk scoring (real types from BaseEntityRiskResponse + EntityRiskResponse)
	HighRiskScore  float64            `json:"highRiskScore,omitempty"`
	IsHighRisk     bool               `json:"isHighRisk,omitempty"`
	SentimentCount map[string]float64 `json:"sentimentCount,omitempty"`
	TotalArticles  float64            `json:"totalArticles,omitempty"`
	Categories     []RiskCategory     `json:"categories,omitempty"`
	RiskScores     []RiskScoreEntry   `json:"riskScores,omitempty"`
	FraudScores    []FraudScoreEntry  `json:"fraudScores,omitempty"`

	// Nested entities
	EntityAddresses  []EntityAddressRiskResponse  `json:"entityAddresses,omitempty"`
	EntityPhones     []EntityPhoneRiskResponse    `json:"entityPhones,omitempty"`
	EntityIdentities []EntityIdentityRiskResponse `json:"entityIdentities,omitempty"`
	EntityPersons    []EntityPersonRiskResponse   `json:"entityPersons,omitempty"`

	// Validation
	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
	SourceUniqueId    string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId   string             `json:"qubitOnUniqueId,omitempty"`
}

// RiskScoreEntry mirrors the canonical RiskScore type.
type RiskScoreEntry struct {
	ScoreType string  `json:"scoreType,omitempty"`
	Score     float64 `json:"score,omitempty"`
}

// FraudScoreEntry mirrors the canonical FraudScore type.
type FraudScoreEntry struct {
	FraudType string  `json:"fraudType,omitempty"`
	Score     float64 `json:"score,omitempty"`
}

// EntityAddressRiskResponse is a nested address risk entry.
type EntityAddressRiskResponse struct {
	AddressLine1  string            `json:"addressLine1,omitempty"`
	AddressLine2  string            `json:"addressLine2,omitempty"`
	City          string            `json:"city,omitempty"`
	State         string            `json:"state,omitempty"`
	PostalCode    string            `json:"postalCode,omitempty"`
	Country       string            `json:"country,omitempty"`
	HighRiskScore float64           `json:"highRiskScore,omitempty"`
	IsHighRisk    bool              `json:"isHighRisk,omitempty"`
	FraudScores   []FraudScoreEntry `json:"fraudScores,omitempty"`
}

// EntityPhoneRiskResponse is a nested phone risk entry.
type EntityPhoneRiskResponse struct {
	PhoneNumber   string            `json:"phoneNumber,omitempty"`
	PhoneType     string            `json:"phoneType,omitempty"`
	HighRiskScore float64           `json:"highRiskScore,omitempty"`
	IsHighRisk    bool              `json:"isHighRisk,omitempty"`
	FraudScores   []FraudScoreEntry `json:"fraudScores,omitempty"`
}

// EntityIdentityRiskResponse is a nested identity-document risk entry.
type EntityIdentityRiskResponse struct {
	IdentityType   string            `json:"identityType,omitempty"`
	IdentityNumber string            `json:"identityNumber,omitempty"`
	HighRiskScore  float64           `json:"highRiskScore,omitempty"`
	IsHighRisk     bool              `json:"isHighRisk,omitempty"`
	FraudScores    []FraudScoreEntry `json:"fraudScores,omitempty"`
}

// EntityPersonRiskResponse is a nested person risk entry.
type EntityPersonRiskResponse struct {
	FirstName     string            `json:"firstName,omitempty"`
	LastName      string            `json:"lastName,omitempty"`
	MiddleName    string            `json:"middleName,omitempty"`
	HighRiskScore float64           `json:"highRiskScore,omitempty"`
	IsHighRisk    bool              `json:"isHighRisk,omitempty"`
	FraudScores   []FraudScoreEntry `json:"fraudScores,omitempty"`
}

// ── Credit Analysis ───────────────────────────────────────────────────────

// CreditAnalysisRequest is the input for LookupCreditAnalysis.
type CreditAnalysisRequest struct {
	BaseRequest
	CompanyName  string `json:"companyName"`
	AddressLine1 string `json:"addressLine1"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
	DunsNumber   string `json:"dunsNumber,omitempty"`
	PostalCode   string `json:"postalCode,omitempty"`
	AddressLine2 string `json:"addressLine2,omitempty"`
}

// CreditAnalysisResponse is one entry in the array returned by LookupCreditAnalysis.
// Mirrors the canonical CreditAnalysisResponse : BaseLookupResponse from
// smartvm.BusinessEntities.Internal.Bases.CreditAnalysis.
type CreditAnalysisResponse struct {
	ApplicationCreated   bool                       `json:"applicationCreated,omitempty"`
	ApplicationEcf       *CreditApplicationEcf      `json:"applicationEcf,omitempty"`
	ApplicationId        string                     `json:"applicationId,omitempty"`
	BureauCompanyList    []CreditBureauCompany      `json:"bureauCompanyList,omitempty"`
	BureauPullErrorFlag  bool                       `json:"bureauPullErrorFlag,omitempty"`
	CompanyInfoDto       *CreditCompanyInfo         `json:"companyInfoDto,omitempty"`
	CorporateLinkageDto  *CreditCorporateLinkage    `json:"corporateLinkageDto,omitempty"`
	DateOfDecisioning    string                     `json:"dateOfDecisioning,omitempty"`
	ExactMatchFound      bool                       `json:"exactMatchFound,omitempty"`
	ListOfSimilarsFound  bool                       `json:"listOfSimilarsFound,omitempty"`
	NoMatchFound         bool                       `json:"noMatchFound,omitempty"`

	// BaseLookupResponse
	Score                float64            `json:"score,omitempty"`
	SourceResultCode     string             `json:"sourceResultCode,omitempty"`
	ValidationResultCode string             `json:"validationResultCode,omitempty"`
	ValidationDate       string             `json:"validationDate,omitempty"`
	SourceUniqueId       string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId      string             `json:"qubitOnUniqueId,omitempty"`
	ValidationResults    []ValidationResult `json:"validationResults,omitempty"`
}

// CreditApplicationEcf mirrors the canonical CreditApplicationEcf nested type.
type CreditApplicationEcf struct {
	Application      *CreditApplication       `json:"application,omitempty"`
	ApplicationValues []CreditApplicationValue `json:"applicationValues,omitempty"`
	CreditTerms      *CreditTerms             `json:"creditTerms,omitempty"`
	DecisionOutcome  *CreditDecisionOutcome   `json:"decisionOutcome,omitempty"`
	ScoreList        []CreditScoreListEntry   `json:"scoreList,omitempty"`
}

// CreditApplication mirrors the canonical CreditApplication.
type CreditApplication struct {
	ApplicationId string                  `json:"applicationId,omitempty"`
	BureauIdList  []CreditBureauIdEntry   `json:"bureauIdList,omitempty"`
	Business      *CreditBusiness         `json:"business,omitempty"`
	Status        string                  `json:"status,omitempty"`
}

// CreditApplicationValue mirrors the canonical CreditApplicationValue.
type CreditApplicationValue struct {
	FieldType          int    `json:"fieldType,omitempty"`
	FieldTypeSpecified bool   `json:"fieldTypeSpecified,omitempty"`
	Index              int    `json:"index,omitempty"`
	IndexSpecified     bool   `json:"indexSpecified,omitempty"`
	Label              string `json:"label,omitempty"`
	Name               string `json:"name,omitempty"`
	Value              string `json:"value,omitempty"`
}

// CreditBureauIdEntry mirrors the canonical BureauIdList.
type CreditBureauIdEntry struct {
	BusinessBureauName     string `json:"businessBureauName,omitempty"`
	BureauNameSpecified    bool   `json:"bureauNameSpecified,omitempty"`
	BusinessBureauId       string `json:"businessBureauId,omitempty"`
}

// CreditBusiness mirrors the canonical Business.
type CreditBusiness struct {
	Address      *CreditAddress `json:"address,omitempty"`
	BusinessName string         `json:"businessName,omitempty"`
}

// CreditAddress mirrors the canonical Address (CreditAnalysis namespace).
type CreditAddress struct {
	City    string `json:"city,omitempty"`
	Country string `json:"country,omitempty"`
	State   string `json:"state,omitempty"`
	Street  string `json:"street,omitempty"`
	ZipCode string `json:"zipCode,omitempty"`
}

// CreditCompanyInfo mirrors the canonical CompanyInfoDto.
type CreditCompanyInfo struct {
	Address     string `json:"address,omitempty"`
	City        string `json:"city,omitempty"`
	CompanyName string `json:"companyName,omitempty"`
	Country     string `json:"country,omitempty"`
	Phone       string `json:"phone,omitempty"`
	State       string `json:"state,omitempty"`
	ZipCode     string `json:"zipCode,omitempty"`
}

// CreditCorporateLinkage mirrors the canonical CorporateLinkageDto.
type CreditCorporateLinkage struct {
	DunsNo            string `json:"dunsNo,omitempty"`
	GlbUltimateDunsNo string `json:"glbUltimateDunsNo,omitempty"`
	HqDunsNo          string `json:"hqDunsNo,omitempty"`
	ParentDunsNo      string `json:"parentDunsNo,omitempty"`
}

// CreditTerms mirrors the canonical CreditTerms.
type CreditTerms struct {
	CreditLimit          float64 `json:"creditLimit,omitempty"`
	CreditLimitSpecified bool    `json:"creditLimitSpecified,omitempty"`
	CreditTermsStatus    string  `json:"creditTermsStatus,omitempty"`
	EarlyPaymentDiscount string  `json:"earlyPaymentDiscount,omitempty"`
	PaymentTerms         string  `json:"paymentTerms,omitempty"`
}

// CreditDecisionOutcome mirrors the canonical DecisionOutcome.
type CreditDecisionOutcome struct {
	AnalystInstruction    string                       `json:"analystInstruction,omitempty"`
	Outcome               string                       `json:"outcome,omitempty"`
	RecmmendedCreditTerms *CreditRecommendedTerms      `json:"recmmendedCreditTerms,omitempty"`
	TriggeredBy           string                       `json:"triggeredBy,omitempty"`
	AndFields             []CreditDecisionAndField     `json:"andFields,omitempty"`
}

// CreditRecommendedTerms mirrors the canonical RecmmendedCreditTerms (server typo preserved).
type CreditRecommendedTerms struct {
	CreditLimit          float64 `json:"creditLimit,omitempty"`
	CreditLimitSpecified bool    `json:"creditLimitSpecified,omitempty"`
	CreditTermsStatus    string  `json:"creditTermsStatus,omitempty"`
	EarlyPaymentDiscount string  `json:"earlyPaymentDiscount,omitempty"`
	PaymentTerms         string  `json:"paymentTerms,omitempty"`
}

// CreditDecisionAndField mirrors the canonical AndField.
type CreditDecisionAndField struct {
	ApplicationValue string `json:"applicationValue,omitempty"`
	CriteriaType     string `json:"criteriaType,omitempty"`
	EntityType       string `json:"entityType,omitempty"`
	FieldName        string `json:"fieldName,omitempty"`
	OutCome          string `json:"outCome,omitempty"`
	Rule             string `json:"rule,omitempty"`
	Value            string `json:"value,omitempty"`
}

// CreditScoreListEntry mirrors the canonical ScoreList entry.
type CreditScoreListEntry struct {
	DateCreated          string  `json:"dateCreated,omitempty"`
	DateCreatedSpecified bool    `json:"dateCreatedSpecified,omitempty"`
	ScoreName            string  `json:"scoreName,omitempty"`
	Value                float64 `json:"value,omitempty"`
	ValueSpecified       bool    `json:"valueSpecified,omitempty"`
}

// CreditBureauCompany mirrors the canonical BureauCompany.
type CreditBureauCompany struct {
	BureauCompanyAddress         *CreditAddress `json:"bureauCompanyAddress,omitempty"`
	BureauIdentifierNumber       string         `json:"bureauIdentifierNumber,omitempty"`
	BureauName                   string         `json:"bureauName,omitempty"`
	BusinessName                 string         `json:"businessName,omitempty"`
	LocationIndicator            string         `json:"locationIndicator,omitempty"`
	LocationIndicatorExplanation string         `json:"locationIndicatorExplanation,omitempty"`
	MatchScore                   int            `json:"matchScore,omitempty"`
	Telephone                    string         `json:"telephone,omitempty"`
}

// ── ESG Score ─────────────────────────────────────────────────────────────

// ESGRequest is the input for LookupESGScore. Mirrors the canonical
// ESGScoringRequest from smartvm.BusinessEntities.Client.API.API.ESG —
// the body consumes CompanyName and ESGId.
//
// Country and Domain are bound on the server as [FromQuery] (not body)
// and are serialised by LookupESGScore into the URL query string with
// percent-encoding. They use json:"-" so json.Marshal omits them from
// the request body.
//
// Acronym JSON tag ESGId→esgId per the upstream CamelCase naming policy
// (empirically verified against System.Text.Json — consecutive leading
// capitals are fully lower-cased).
type ESGRequest struct {
	BaseRequest
	CompanyName string `json:"companyName,omitempty"`
	ESGId       int    `json:"esgId,omitempty"`
	// Country is sent as a query-string parameter, not body. ISO 3166-1 alpha-2/3 or full name. Optional.
	Country string `json:"-"`
	// Domain is sent as a query-string parameter, not body. e.g. "example.com". Optional.
	Domain string `json:"-"`
}

// ESGResponse is one entry in the array returned by /api/esg/Scores.
// The canonical model has [JsonPropertyName("ESGId")] which overrides the
// CamelCase policy — wire literal is "ESGId" (PascalCase preserved).
type ESGResponse struct {
	ESGId            int            `json:"ESGId,omitempty"`
	Name             string         `json:"name,omitempty"`
	CompanyDesc      string         `json:"company_desc,omitempty"`
	Grade            string         `json:"grade,omitempty"`
	ExchangeSymbol   string         `json:"exchange_symbol,omitempty"`
	StockSymbol      string         `json:"stock_symbol,omitempty"`
	// Industry and CountryCode use pointers so callers can distinguish
	// "absent" (server returned no value) from "0" — both are valid
	// industry/country code dictionary values on some providers.
	Industry         *int           `json:"industry,omitempty"`
	CountryCode      *int           `json:"country,omitempty"`
	EnvironmentScore int            `json:"e"`
	SocialScore      int            `json:"s"`
	GovernanceScore  int            `json:"g"`
	Total            int            `json:"total"`
	Progress         []ESGProgress  `json:"progress,omitempty"`
	TopRecords       *ESGTopRecords `json:"toprecords,omitempty"`

	// Lookup metadata
	Score                float64 `json:"score,omitempty"`
	SourceResultCode     string  `json:"sourceResultCode,omitempty"`
	ValidationResultCode string  `json:"validationResultCode,omitempty"`
	ValidationDate       string  `json:"validationDate,omitempty"`
	SourceUniqueId       string  `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId      string  `json:"qubitOnUniqueId,omitempty"`
}

// ESGProgress represents a historical ESG score data point.
type ESGProgress struct {
	Date             string `json:"date"`
	EnvironmentScore int    `json:"e"`
	SocialScore      int    `json:"s"`
	GovernanceScore  int    `json:"g"`
	Total            int    `json:"total"`
}

// ESGTopRecords contains competitor benchmark data.
type ESGTopRecords struct {
	Top []ESGCompetitor `json:"top,omitempty"`
}

// ESGCompetitor represents a competitor's ESG scores for benchmarking.
type ESGCompetitor struct {
	Id               int     `json:"id"`
	Name             string  `json:"name"`
	Grade            string  `json:"grade,omitempty"`
	ExchangeSymbol   string  `json:"exchange_symbol,omitempty"`
	StockSymbol      string  `json:"stock_symbol,omitempty"`
	Industry         int     `json:"industry,omitempty"`
	CountryCode      int     `json:"country,omitempty"`
	CompanyDesc      string  `json:"company_desc,omitempty"`
	EnvironmentScore float32 `json:"e"`
	SocialScore      float32 `json:"s"`
	GovernanceScore  float32 `json:"g"`
	Total            float32 `json:"total"`
}

// ── Domain Security ───────────────────────────────────────────────────────

// DomainSecurityRequest is the input for DomainSecurityReport.
type DomainSecurityRequest struct {
	BaseRequest
	DomainName string `json:"domain"`
}

// DomainSecurityResponse is the output from DomainSecurityReport.
// Acronym JSON tag IPAddress→ipAddress per .NET CamelCase JsonNamingPolicy
// (empirically verified against System.Text.Json with .NET 10).
type DomainSecurityResponse struct {
	DomainName         string                   `json:"domain,omitempty"`
	Score              *float64                 `json:"score,omitempty"`
	Grade              string                   `json:"grade,omitempty"`
	IPAddress          string                   `json:"ipAddress,omitempty"`
	Issues             []map[string]interface{} `json:"issues,omitempty"`
	BreachedEmails     []map[string]interface{} `json:"breachedEmails,omitempty"`
	SubdomainSummary   map[string]interface{}   `json:"subdomainSummary,omitempty"`
	NetworkSummary     map[string]interface{}   `json:"networkSummary,omitempty"`
	WebSecuritySummary map[string]interface{}   `json:"webSecuritySummary,omitempty"`

	// Validation
	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
	SourceUniqueId    string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId   string             `json:"qubitOnUniqueId,omitempty"`
}

// ── IP Quality ────────────────────────────────────────────────────────────

// IPQualityRequest is the input for CheckIPQuality.
// Acronym JSON tag IPAddress→ipAddress per .NET CamelCase JsonNamingPolicy
// (empirically verified against System.Text.Json with .NET 10).
type IPQualityRequest struct {
	BaseRequest
	IPAddress string `json:"ipAddress"`
	UserAgent string `json:"userAgent,omitempty"`
}

// IPQualityResponse is the output from CheckIPQuality.
//
// All boolean flags from the server are pointers so that "absent" can be
// distinguished from "present and false" — server-omitted booleans should not
// silently become false.
//
// Acronym JSON tags follow .NET CamelCase JsonNamingPolicy (empirically
// verified against System.Text.Json with .NET 10 — consecutive leading
// capitals are fully lower-cased): IPAddress→ipAddress, IsVPN→isVPN,
// IsTOR→isTOR, ISP→isp.
//
// Note: AbuseVelocity, BotStatus, DeviceModel, DeviceBrand, Mobile,
// FraudScore, OperatingSystem, Browser are declared as PUBLIC FIELDS (not
// properties) on the canonical .NET response; System.Text.Json does not
// serialise public fields by default, so those values never appear on the
// wire. They are intentionally omitted from this struct.
type IPQualityResponse struct {
	IPAddress      string   `json:"ipAddress,omitempty"`
	UserAgent      string   `json:"userAgent,omitempty"`
	IsProxy        *bool    `json:"isProxy,omitempty"`
	IsVPN          *bool    `json:"isVPN,omitempty"`
	IsActiveVPN    *bool    `json:"isActiveVPN,omitempty"`
	ISP            string   `json:"isp,omitempty"`
	IsCrawler      *bool    `json:"isCrawler,omitempty"`
	IsTOR          *bool    `json:"isTOR,omitempty"`
	IsActiveTOR    *bool    `json:"isActiveTOR,omitempty"`
	RecentAbuse    *bool    `json:"recentAbuse,omitempty"`
	Organization   string   `json:"organization,omitempty"`
	Latitude       *float64 `json:"latitude,omitempty"`
	Longitude      *float64 `json:"longitude,omitempty"`
	Timezone       string   `json:"timezone,omitempty"`
	City           string   `json:"city,omitempty"`
	Region         string   `json:"region,omitempty"`
	Host           string   `json:"host,omitempty"`
	ASN            *int     `json:"asn,omitempty"`
	ConnectionType string   `json:"connectionType,omitempty"`
	Country        *Country `json:"country,omitempty"`

	// Validation
	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
	SourceUniqueId    string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId   string             `json:"qubitOnUniqueId,omitempty"`
}

// ── Beneficial Ownership ──────────────────────────────────────────────────

// BeneficialOwnershipRequest is the input for LookupBeneficialOwnership.
type BeneficialOwnershipRequest struct {
	BaseRequest
	CompanyName  string `json:"companyName"`
	CountryISO2  string `json:"countryIso2"`
	UBOThreshold string `json:"uboThreshold,omitempty"`
	MaxLayers    string `json:"maxLayers,omitempty"`
}

// BeneficialOwnershipResponse is one entry in the array returned by
// LookupBeneficialOwnership. Mirrors the canonical
// BeneficialOwnershipLookupResponse : BaseLookupResponse exactly.
type BeneficialOwnershipResponse struct {
	OwnershipTreeField   *RetrieveOwnershipTree `json:"ownershipTreeField,omitempty"`
	ResponseCodeField    string                 `json:"responseCodeField,omitempty"`
	ResponseDetailsField string                 `json:"responseDetailsField,omitempty"`
	Message              string                 `json:"message,omitempty"`
	// ErrorResponse uses a PascalCase JSON tag because the canonical model
	// declares an explicit [JsonPropertyName("ErrorResponse")] override.
	ErrorResponse *InternalErrorResponse `json:"ErrorResponse,omitempty"`

	// Lookup metadata (BaseLookupResponse)
	Score                float64            `json:"score,omitempty"`
	SourceResultCode     string             `json:"sourceResultCode,omitempty"`
	ValidationResultCode string             `json:"validationResultCode,omitempty"`
	ValidationDate       string             `json:"validationDate,omitempty"`
	SourceUniqueId       string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId      string             `json:"qubitOnUniqueId,omitempty"`
	ValidationResults    []ValidationResult `json:"validationResults,omitempty"`
}

// RetrieveOwnershipTree mirrors the canonical RetrieveOwnershipTree nested
// inside BeneficialOwnershipLookupResponse. Field names are the canonical
// "*Field" names with explicit camelCase JSON tags (the .NET model uses
// [JsonPropertyName] on every field).
type RetrieveOwnershipTree struct {
	IdField                   string                       `json:"idField,omitempty"`
	NameField                 string                       `json:"nameField,omitempty"`
	AddressField              string                       `json:"addressField,omitempty"`
	NodesField                []OwnershipNode              `json:"nodesField,omitempty"`
	UltimateBeneficialOwners  []UltimateBeneficialOwner    `json:"ultimateBeneficialOwnersField,omitempty"`
	StatusField               string                       `json:"statusField,omitempty"`
	ErrorStatusExtraInfoField string                       `json:"errorStatusExtraInfoField,omitempty"`
	UboThresholdField         float64                      `json:"uboThresholdField,omitempty"`
	MaxCreditCostField        int                          `json:"maxCreditCostField,omitempty"`
	MaxLayersField            int                          `json:"maxLayersField,omitempty"`
	DateCreatedField          string                       `json:"dateCreatedField,omitempty"`
	RegAuthField              string                       `json:"regAuthField,omitempty"`
	RetrievalLocationField    string                       `json:"retrievalLocationField,omitempty"`
	TotalCreditsSpentField    int                          `json:"totalCreditsSpentField,omitempty"`
	UnwrapFeeField            int                          `json:"unwrapFeeField,omitempty"`
	CountryISOField           string                       `json:"countryISOField,omitempty"`
	CompanyCodeField          string                       `json:"companyCodeField,omitempty"`
	ContinuationKeysUsedField []string                     `json:"continuationKeysUsedField,omitempty"`
	DataSource                string                       `json:"dataSource,omitempty"`
	RegistryDataDate          string                       `json:"registryDataDate,omitempty"`
}

// OwnershipNode mirrors the canonical Node nested in RetrieveOwnershipTree.
type OwnershipNode struct {
	LevelField            int            `json:"levelField,omitempty"`
	EntityField           *OwnershipEntity `json:"entityField,omitempty"`
	EdgesField            []OwnershipEdge `json:"edgesField,omitempty"`
	RollupPercentageField float64        `json:"rollupPercentageField,omitempty"`
}

// OwnershipEntity mirrors the canonical Entity nested in OwnershipNode.
type OwnershipEntity struct {
	IdField                        string                  `json:"idField,omitempty"`
	TypeField                      string                  `json:"typeField,omitempty"`
	NameField                      string                  `json:"nameField,omitempty"`
	NameInEnglishField             string                  `json:"nameInEnglishField,omitempty"`
	ProhibitedListResponse         map[string]interface{}  `json:"prohibitedListResponse,omitempty"`
	CountryISOField                string                  `json:"countryISOField,omitempty"`
	CompanyCodeField               string                  `json:"companyCodeField,omitempty"`
	AddressField                   string                  `json:"addressField,omitempty"`
	LinksField                     *OwnershipLinks         `json:"linksField,omitempty"`
	RegistrationAuthorityField     string                  `json:"registrationAuthorityField,omitempty"`
	RegistrationAuthorityCodeField string                  `json:"registrationAuthorityCodeField,omitempty"`
	DateOfBirth                    string                  `json:"dateOfBirth,omitempty"`
	Nationality                    string                  `json:"nationality,omitempty"`
	CountryOfResidence             string                  `json:"countryOfResidence,omitempty"`
	NatureOfControl                string                  `json:"natureOfControl,omitempty"`
	ReasonForNonContinuationField  *ReasonForNonContinuation `json:"reasonForNonContinuationField,omitempty"`
}

// OwnershipLinks mirrors the canonical Links nested in OwnershipEntity.
type OwnershipLinks struct {
	EnhancedProfileField string `json:"enhancedProfileField,omitempty"`
}

// OwnershipEdge mirrors the canonical Edge nested in OwnershipNode.
type OwnershipEdge struct {
	NodeIdField       string                  `json:"nodeIdField,omitempty"`
	TypeField         string                  `json:"typeField,omitempty"`
	PercentageField   *float64                `json:"percentageField,omitempty"`
	IsCircularField   bool                    `json:"isCircularField,omitempty"`
	RoleField         string                  `json:"roleField,omitempty"`
	ShareholdingsField []OwnershipShareholding `json:"shareholdingsField,omitempty"`
}

// OwnershipShareholding mirrors the canonical Shareholding nested in OwnershipEdge.
type OwnershipShareholding struct {
	IsJointlyHeldField        *bool    `json:"isJointlyHeldField,omitempty"`
	JointHoldingGroupIdField  string   `json:"jointHoldingGroupIdField,omitempty"`
	PercentageField           *float64 `json:"percentageField,omitempty"`
}

// UltimateBeneficialOwner mirrors the canonical UltimateBeneficialOwner nested
// in RetrieveOwnershipTree.
type UltimateBeneficialOwner struct {
	IdField                string                 `json:"idField,omitempty"`
	NameField              string                 `json:"nameField,omitempty"`
	NameInEnglishField     string                 `json:"nameInEnglishField,omitempty"`
	PercentageField        string                 `json:"percentageField,omitempty"`
	EntityTypeField        string                 `json:"entityTypeField,omitempty"`
	AddressField           string                 `json:"addressField,omitempty"`
	DateOfBirth            string                 `json:"dateOfBirth,omitempty"`
	Nationality            string                 `json:"nationality,omitempty"`
	CountryOfResidence     string                 `json:"countryOfResidence,omitempty"`
	NatureOfControl        string                 `json:"natureOfControl,omitempty"`
	NotifiedDate           string                 `json:"notifiedDate,omitempty"`
	ControlType            string                 `json:"controlType,omitempty"`
	Status                 string                 `json:"status,omitempty"`
	ProhibitedListResponse map[string]interface{} `json:"prohibitedListResponse,omitempty"`
}

// ReasonForNonContinuation mirrors the canonical ReasonForNonContinuation.
type ReasonForNonContinuation struct {
	DetailsField        string                  `json:"detailsField,omitempty"`
	TypeField           string                  `json:"typeField,omitempty"`
	CandidatesField     []OwnershipCandidate    `json:"candidatesField,omitempty"`
	ContinuationKeyField string                 `json:"continuationKeyField,omitempty"`
}

// OwnershipCandidate mirrors the canonical CandidateField.
type OwnershipCandidate struct {
	AddressField               string `json:"addressField,omitempty"`
	CompanyIdField             string `json:"companyIdField,omitempty"`
	CompanyNameField           string `json:"companyNameField,omitempty"`
	CountryIsoField            string `json:"countryIsoField,omitempty"`
	RegistrationAuthorityField string `json:"registrationAuthorityField,omitempty"`
	ContinuationKeyField       string `json:"continuationKeyField,omitempty"`
}

// InternalErrorResponse mirrors the canonical InternalErrorResponse used in
// BeneficialOwnership and other lookup envelopes. The .NET CamelCase
// JsonNamingPolicy lowercases the first character of each property name.
type InternalErrorResponse struct {
	StatusCode int    `json:"statusCode,omitempty"`
	Status     string `json:"status,omitempty"`
	Message    string `json:"message,omitempty"`
}

// ── Corporate Hierarchy ───────────────────────────────────────────────────

// CorporateHierarchyRequest is the input for LookupCorporateHierarchy.
type CorporateHierarchyRequest struct {
	BaseRequest
	CompanyName  string `json:"companyName"`
	AddressLine1 string `json:"addressLine1"`
	City         string `json:"city"`
	State        string `json:"state"`
	ZipCode      string `json:"zipCode"`
}

// CorporateHierarchyResponse is one entry in the array returned by
// LookupCorporateHierarchy. Mirrors the canonical CorporateHierarchyResponse :
// BaseLookupResponse exactly. The body is wrapped in an ELSGenericMessage
// payload — the shape is provider-specific (Equifax ELS) so its inner
// fields are exposed via the typed wrapper rather than as a free-form map.
type CorporateHierarchyResponse struct {
	ELSGenericMessage *ELSGenericMessage `json:"elsGenericMessage,omitempty"`

	// Lookup metadata (BaseLookupResponse)
	Score                float64            `json:"score,omitempty"`
	SourceResultCode     string             `json:"sourceResultCode,omitempty"`
	ValidationResultCode string             `json:"validationResultCode,omitempty"`
	ValidationDate       string             `json:"validationDate,omitempty"`
	SourceUniqueId       string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId      string             `json:"qubitOnUniqueId,omitempty"`
	ValidationResults    []ValidationResult `json:"validationResults,omitempty"`
}

// ELSGenericMessage mirrors the canonical ELSGenericMessage envelope returned
// by the Equifax ELS corporate-hierarchy provider. The nested members
// (Stage2Data, Message, StandardizedAddress, InputData) are intentionally
// kept as map[string]interface{} because they are deeply nested provider-
// specific shapes that are not declared as Client.API DTOs.
type ELSGenericMessage struct {
	Stage2Data          map[string]interface{} `json:"stage2Data,omitempty"`
	Message             map[string]interface{} `json:"message,omitempty"`
	StandardizedAddress map[string]interface{} `json:"standardizedAddress,omitempty"`
	// TransactionId uses an explicit [JsonPropertyName("TRANSACTION_ID")]
	// override on the .NET model.
	TransactionId int                    `json:"TRANSACTION_ID,omitempty"`
	InputData     map[string]interface{} `json:"inputData,omitempty"`
}

// ── DUNS Lookup ───────────────────────────────────────────────────────────

// DUNSRequest is the input for LookupDUNS.
type DUNSRequest struct {
	BaseRequest
	DunsNumber string `json:"dunsNumber"`
}

// DUNSResponse is one entry in the array returned by LookupDUNS.
type DUNSResponse struct {
	DunsNumber      string `json:"dunsNumber,omitempty"`
	CompanyName     string `json:"companyName,omitempty"`
	TradeStyle      string `json:"tradeStyle,omitempty"`
	AddressLine1    string `json:"addressLine1,omitempty"`
	AddressLine2    string `json:"addressLine2,omitempty"`
	City            string `json:"city,omitempty"`
	State           string `json:"state,omitempty"`
	PostalCode      string `json:"postalCode,omitempty"`
	Country         string `json:"country,omitempty"`
	PrimaryNAICS    string `json:"primaryNaics,omitempty"`
	PrimarySIC      string `json:"primarySic,omitempty"`
	OperatingStatus string `json:"operatingStatus,omitempty"`
	YearStarted     string `json:"yearStarted,omitempty"`
	EmployeesTotal  *int   `json:"employeesTotal,omitempty"`

	Score             float64            `json:"score,omitempty"`
	SourceUniqueId    string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId   string             `json:"qubitOnUniqueId,omitempty"`
	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
}

// ── Hierarchy Lookup ──────────────────────────────────────────────────────

// HierarchyRequest is the input for LookupHierarchy.
type HierarchyRequest struct {
	BaseRequest
	Identifier     string `json:"identifier"`
	IdentifierType string `json:"identifierType"`
	Country        string `json:"country,omitempty"`
	Options        string `json:"options,omitempty"`
}

// HierarchyResponse is the output from LookupHierarchy. Mirrors the canonical
// HierarchyLookupResponse : BaseValidationResponse from
// smartvm.BusinessEntities.Client.API.HierarchyLookup. The Parent property is
// a Company struct (whose Children are also Company), not a free-form map.
type HierarchyResponse struct {
	Parent *HierarchyCompany `json:"parent,omitempty"`

	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
	SourceUniqueId    string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId   string             `json:"qubitOnUniqueId,omitempty"`
}

// HierarchyCompany mirrors the canonical Company DTO used by HierarchyLookupResponse.
type HierarchyCompany struct {
	CompanyName        string                 `json:"companyName,omitempty"`
	CompanyNameDBA     string                 `json:"companyNameDBA,omitempty"`
	BusinessEntityType string                 `json:"businessEntityType,omitempty"`
	Identifiers        map[string][]string    `json:"identifiers,omitempty"`
	Address            []HierarchyAddress     `json:"address,omitempty"`
	Phones             []HierarchyPhone       `json:"phones,omitempty"`
	Emails             []HierarchyEmail       `json:"emails,omitempty"`
	Persons            []HierarchyPerson      `json:"persons,omitempty"`
	Country            *Country               `json:"country,omitempty"`
	Children           []HierarchyCompany     `json:"children,omitempty"`
}

// HierarchyAddress is one address attached to a HierarchyCompany. Mirrors the
// canonical CompanyListAddress.
type HierarchyAddress struct {
	AddressType  string   `json:"addressType,omitempty"`
	AddressLine1 string   `json:"addressLine1,omitempty"`
	AddressLine2 string   `json:"addressLine2,omitempty"`
	AddressLine3 string   `json:"addressLine3,omitempty"`
	AddressLine4 string   `json:"addressLine4,omitempty"`
	City         string   `json:"city,omitempty"`
	State        string   `json:"state,omitempty"`
	PostalCode   string   `json:"postalCode,omitempty"`
	Country      *Country `json:"country,omitempty"`
	Rank         *int     `json:"rank,omitempty"`
}

// HierarchyPhone is one phone attached to a HierarchyCompany. Mirrors the
// canonical CompanyListPhones.
type HierarchyPhone struct {
	PhoneNumber  string `json:"phoneNumber,omitempty"`
	PhoneType    string `json:"phoneType,omitempty"`
	PhoneSubType string `json:"phoneSubType,omitempty"`
	Rank         *int   `json:"rank,omitempty"`
}

// HierarchyEmail is one email attached to a HierarchyCompany. Mirrors the
// canonical CompanyListEmails.
type HierarchyEmail struct {
	EmailAddress string `json:"emailAddress,omitempty"`
	EmailType    string `json:"emailType,omitempty"`
	Rank         *int   `json:"rank,omitempty"`
}

// HierarchyPerson is one person attached to a HierarchyCompany. Mirrors the
// canonical CompanyListPersons. Acronym JSON tag LinkedInUrl→linkedInUrl per
// .NET CamelCase JsonNamingPolicy.
type HierarchyPerson struct {
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	PersonType  string `json:"personType,omitempty"`
	LinkedInUrl string `json:"linkedInUrl,omitempty"`
	Rank        *int   `json:"rank,omitempty"`
}

// ── NPI Validation ────────────────────────────────────────────────────────

// NPIRequest is the input for ValidateNPI.
// Acronym JSON tag NPI→npi per .NET CamelCase JsonNamingPolicy (empirically
// verified against System.Text.Json with .NET 10 — consecutive leading
// capitals are fully lower-cased).
type NPIRequest struct {
	BaseRequest
	NPI              string `json:"npi"`
	OrganizationName string `json:"organizationName,omitempty"`
	LastName         string `json:"lastName,omitempty"`
	FirstName        string `json:"firstName,omitempty"`
	MiddleName       string `json:"middleName,omitempty"`
}

// NPIResponse is the output from ValidateNPI. Mirrors the canonical
// NPIResponse : BaseValidationResponse exactly. NPI is int64 on the wire
// (server property is `public long NPI { get; set; }`).
//
// Acronym JSON tags follow .NET CamelCase JsonNamingPolicy (empirically
// verified against System.Text.Json with .NET 10 — consecutive leading
// capitals are fully lower-cased): NPI→npi,
// NPIDeactivationReasonCode→npiDeactivationReasonCode, etc.
//
// Property names with embedded underscores (HealthcareProviderTaxonomyCode_1
// through _15, ProviderLicenseNumber_*, HealthcareProviderTaxonomyGroup_*)
// are preserved verbatim because .NET CamelCase JsonNamingPolicy only
// lowercases the leading character run; the underscore + digit suffix is
// kept as-is.
type NPIResponse struct {
	NPI                                              int64  `json:"npi,omitempty"`
	ProviderOrganizationName                         string `json:"providerOrganizationName,omitempty"`
	ProviderLastName                                 string `json:"providerLastName,omitempty"`
	ProviderFirstName                                string `json:"providerFirstName,omitempty"`
	ProviderMiddleName                               string `json:"providerMiddleName,omitempty"`
	ProviderNameSuffixText                           string `json:"providerNameSuffixText,omitempty"`
	ProviderFirstLineBusinessMailingAddress          string `json:"providerFirstLineBusinessMailingAddress,omitempty"`
	ProviderSecondLineBusinessMailingAddress         string `json:"providerSecondLineBusinessMailingAddress,omitempty"`
	ProviderBusinessMailingAddressCityName           string `json:"providerBusinessMailingAddressCityName,omitempty"`
	ProviderBusinessMailingAddressStateName          string `json:"providerBusinessMailingAddressStateName,omitempty"`
	ProviderBusinessMailingAddressPostalCode         string `json:"providerBusinessMailingAddressPostalCode,omitempty"`
	ProviderBusinessMailingAddressCountryCode        string `json:"providerBusinessMailingAddressCountryCode,omitempty"`
	ProviderBusinessMailingAddressTelephoneNumber    string `json:"providerBusinessMailingAddressTelephoneNumber,omitempty"`
	LastUpdateDate                                   string `json:"lastUpdateDate,omitempty"`
	NPIDeactivationReasonCode                        string `json:"npiDeactivationReasonCode,omitempty"`
	NPIDeactivationDate                              string `json:"npiDeactivationDate,omitempty"`
	NPIReactivationDate                              string `json:"npiReactivationDate,omitempty"`
	AuthorizedOfficialLastName                       string `json:"authorizedOfficialLastName,omitempty"`
	AuthorizedOfficialFirstName                      string `json:"authorizedOfficialFirstName,omitempty"`
	AuthorizedOfficialMiddleName                     string `json:"authorizedOfficialMiddleName,omitempty"`
	AuthorizedOfficialTitleorPosition                string `json:"authorizedOfficialTitleorPosition,omitempty"`
	AuthorizedOfficialTelephoneNumber                string `json:"authorizedOfficialTelephoneNumber,omitempty"`
	HealthcareProviderTaxonomyCode_1                 string `json:"healthcareProviderTaxonomyCode_1,omitempty"`
	ProviderLicenseNumber_1                          string `json:"providerLicenseNumber_1,omitempty"`
	ProviderLicenseNumberStateCode_1                 string `json:"providerLicenseNumberStateCode_1,omitempty"`
	HealthcareProviderTaxonomyCode_2                 string `json:"healthcareProviderTaxonomyCode_2,omitempty"`
	ProviderLicenseNumber_2                          string `json:"providerLicenseNumber_2,omitempty"`
	ProviderLicenseNumberStateCode_2                 string `json:"providerLicenseNumberStateCode_2,omitempty"`
	HealthcareProviderTaxonomyCode_3                 string `json:"healthcareProviderTaxonomyCode_3,omitempty"`
	ProviderLicenseNumber_3                          string `json:"providerLicenseNumber_3,omitempty"`
	ProviderLicenseNumberStateCode_3                 string `json:"providerLicenseNumberStateCode_3,omitempty"`
	HealthcareProviderTaxonomyCode_4                 string `json:"healthcareProviderTaxonomyCode_4,omitempty"`
	ProviderLicenseNumber_4                          string `json:"providerLicenseNumber_4,omitempty"`
	ProviderLicenseNumberStateCode_4                 string `json:"providerLicenseNumberStateCode_4,omitempty"`
	HealthcareProviderTaxonomyCode_5                 string `json:"healthcareProviderTaxonomyCode_5,omitempty"`
	ProviderLicenseNumber_5                          string `json:"providerLicenseNumber_5,omitempty"`
	ProviderLicenseNumberStateCode_5                 string `json:"providerLicenseNumberStateCode_5,omitempty"`
	HealthcareProviderTaxonomyCode_6                 string `json:"healthcareProviderTaxonomyCode_6,omitempty"`
	ProviderLicenseNumber_6                          string `json:"providerLicenseNumber_6,omitempty"`
	ProviderLicenseNumberStateCode_6                 string `json:"providerLicenseNumberStateCode_6,omitempty"`
	HealthcareProviderTaxonomyCode_7                 string `json:"healthcareProviderTaxonomyCode_7,omitempty"`
	ProviderLicenseNumber_7                          string `json:"providerLicenseNumber_7,omitempty"`
	ProviderLicenseNumberStateCode_7                 string `json:"providerLicenseNumberStateCode_7,omitempty"`
	HealthcareProviderTaxonomyCode_8                 string `json:"healthcareProviderTaxonomyCode_8,omitempty"`
	ProviderLicenseNumber_8                          string `json:"providerLicenseNumber_8,omitempty"`
	ProviderLicenseNumberStateCode_8                 string `json:"providerLicenseNumberStateCode_8,omitempty"`
	HealthcareProviderTaxonomyCode_9                 string `json:"healthcareProviderTaxonomyCode_9,omitempty"`
	ProviderLicenseNumber_9                          string `json:"providerLicenseNumber_9,omitempty"`
	ProviderLicenseNumberStateCode_9                 string `json:"providerLicenseNumberStateCode_9,omitempty"`
	HealthcareProviderTaxonomyCode_10                string `json:"healthcareProviderTaxonomyCode_10,omitempty"`
	ProviderLicenseNumber_10                         string `json:"providerLicenseNumber_10,omitempty"`
	ProviderLicenseNumberStateCode_10                string `json:"providerLicenseNumberStateCode_10,omitempty"`
	HealthcareProviderTaxonomyCode_11                string `json:"healthcareProviderTaxonomyCode_11,omitempty"`
	ProviderLicenseNumber_11                         string `json:"providerLicenseNumber_11,omitempty"`
	ProviderLicenseNumberStateCode_11                string `json:"providerLicenseNumberStateCode_11,omitempty"`
	HealthcareProviderTaxonomyCode_12                string `json:"healthcareProviderTaxonomyCode_12,omitempty"`
	ProviderLicenseNumber_12                         string `json:"providerLicenseNumber_12,omitempty"`
	ProviderLicenseNumberStateCode_12                string `json:"providerLicenseNumberStateCode_12,omitempty"`
	HealthcareProviderTaxonomyCode_13                string `json:"healthcareProviderTaxonomyCode_13,omitempty"`
	ProviderLicenseNumber_13                         string `json:"providerLicenseNumber_13,omitempty"`
	ProviderLicenseNumberStateCode_13                string `json:"providerLicenseNumberStateCode_13,omitempty"`
	HealthcareProviderTaxonomyCode_14                string `json:"healthcareProviderTaxonomyCode_14,omitempty"`
	ProviderLicenseNumber_14                         string `json:"providerLicenseNumber_14,omitempty"`
	ProviderLicenseNumberStateCode_14                string `json:"providerLicenseNumberStateCode_14,omitempty"`
	HealthcareProviderTaxonomyCode_15                string `json:"healthcareProviderTaxonomyCode_15,omitempty"`
	ProviderLicenseNumber_15                         string `json:"providerLicenseNumber_15,omitempty"`
	ProviderLicenseNumberStateCode_15                string `json:"providerLicenseNumberStateCode_15,omitempty"`
	IsSoleProprietor                                 string `json:"isSoleProprietor,omitempty"`
	IsOrganizationSubpart                            string `json:"isOrganizationSubpart,omitempty"`
	ParentOrganizationLBN                            string `json:"parentOrganizationLBN,omitempty"`
	ParentOrganizationTIN                            string `json:"parentOrganizationTIN,omitempty"`
	HealthcareProviderTaxonomyGroup_1                string `json:"healthcareProviderTaxonomyGroup_1,omitempty"`
	HealthcareProviderTaxonomyGroup_2                string `json:"healthcareProviderTaxonomyGroup_2,omitempty"`
	HealthcareProviderTaxonomyGroup_3                string `json:"healthcareProviderTaxonomyGroup_3,omitempty"`
	HealthcareProviderTaxonomyGroup_4                string `json:"healthcareProviderTaxonomyGroup_4,omitempty"`
	HealthcareProviderTaxonomyGroup_5                string `json:"healthcareProviderTaxonomyGroup_5,omitempty"`
	HealthcareProviderTaxonomyGroup_6                string `json:"healthcareProviderTaxonomyGroup_6,omitempty"`
	HealthcareProviderTaxonomyGroup_7                string `json:"healthcareProviderTaxonomyGroup_7,omitempty"`
	HealthcareProviderTaxonomyGroup_8                string `json:"healthcareProviderTaxonomyGroup_8,omitempty"`
	HealthcareProviderTaxonomyGroup_9                string `json:"healthcareProviderTaxonomyGroup_9,omitempty"`
	HealthcareProviderTaxonomyGroup_10               string `json:"healthcareProviderTaxonomyGroup_10,omitempty"`
	HealthcareProviderTaxonomyGroup_11               string `json:"healthcareProviderTaxonomyGroup_11,omitempty"`
	HealthcareProviderTaxonomyGroup_12               string `json:"healthcareProviderTaxonomyGroup_12,omitempty"`
	HealthcareProviderTaxonomyGroup_13               string `json:"healthcareProviderTaxonomyGroup_13,omitempty"`
	HealthcareProviderTaxonomyGroup_14               string `json:"healthcareProviderTaxonomyGroup_14,omitempty"`
	HealthcareProviderTaxonomyGroup_15               string `json:"healthcareProviderTaxonomyGroup_15,omitempty"`

	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
	SourceUniqueId    string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId   string             `json:"qubitOnUniqueId,omitempty"`
}

// ── Medpass Validation ────────────────────────────────────────────────────

// MedpassRequest is the input for ValidateMedpass. Mirrors the canonical
// MedpassRequest : BaseRequest, IBaseEntityAddress.
//
// Acronym JSON tags follow .NET CamelCase JsonNamingPolicy (empirically
// verified against System.Text.Json with .NET 10 — consecutive leading
// capitals are fully lower-cased): ID→id, TaxID→taxID (internal acronym;
// only the leading run is lower-cased).
type MedpassRequest struct {
	BaseRequest
	ID                 string                `json:"id"`
	BusinessEntityType string                `json:"businessEntityType"`
	Identities         []CompanyListIdentity `json:"identities,omitempty"`
	CompanyName        string                `json:"companyName,omitempty"`
	CompanyNameDBA     string                `json:"companyNameDBA,omitempty"`
	LegalCompanyName   string                `json:"legalCompanyName,omitempty"`
	TaxID              string                `json:"taxID,omitempty"`
	Country            string                `json:"country,omitempty"`
	State              string                `json:"state,omitempty"`
	City               string                `json:"city,omitempty"`
	PostalCode         string                `json:"postalCode,omitempty"`
	AddressLine1       string                `json:"addressLine1,omitempty"`
	AddressLine2       string                `json:"addressLine2,omitempty"`
}

// CompanyListIdentity mirrors the canonical CompanyListIdentity DTO from
// smartvm.BusinessEntities.Client.API.Bases. Used by MedpassRequest and
// other endpoints that accept a typed list of identity numbers.
type CompanyListIdentity struct {
	Numbers []string `json:"numbers,omitempty"`
	Type    string   `json:"type,omitempty"`
	Pattern string   `json:"pattern,omitempty"`
}

// MedpassResponse is the output from ValidateMedpass. Mirrors the canonical
// MedpassResponse : BaseValidationResponse exactly — only Status and
// LastUpdatedDate are exposed; everything else lives on the base.
type MedpassResponse struct {
	Status          string `json:"status,omitempty"`
	LastUpdatedDate string `json:"lastUpdatedDate,omitempty"`

	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
	SourceUniqueId    string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId   string             `json:"qubitOnUniqueId,omitempty"`
}

// ── DOT Carrier Lookup ────────────────────────────────────────────────────

// DOTCarrierRequest is the input for LookupDOTCarrier.
type DOTCarrierRequest struct {
	BaseRequest
	DotNumber  string `json:"dotNumber"`
	EntityName string `json:"entityName,omitempty"`
}

// DOTCarrierResponse is one entry in the array returned by LookupDOTCarrier.
type DOTCarrierResponse struct {
	DotNumber           string   `json:"dotNumber,omitempty"`
	CarrierName         string   `json:"carrierName,omitempty"`
	LegalName           string   `json:"legalName,omitempty"`
	DBAName             string   `json:"dbaName,omitempty"`
	OperatingStatus     string   `json:"operatingStatus,omitempty"`
	SafetyRating        string   `json:"safetyRating,omitempty"`
	SafetyRatingDate    string   `json:"safetyRatingDate,omitempty"`
	MCS150Date          string   `json:"mcs150Date,omitempty"`
	MCS150Mileage       *int64   `json:"mcs150Mileage,omitempty"`
	PowerUnits          *int     `json:"powerUnits,omitempty"`
	Drivers             *int     `json:"drivers,omitempty"`
	OutOfServiceDate    string   `json:"outOfServiceDate,omitempty"`
	PhysicalAddress     string   `json:"physicalAddress,omitempty"`
	PhysicalCity        string   `json:"physicalCity,omitempty"`
	PhysicalState       string   `json:"physicalState,omitempty"`
	PhysicalZipCode     string   `json:"physicalZipCode,omitempty"`
	MailingAddress      string   `json:"mailingAddress,omitempty"`
	Telephone           string   `json:"telephone,omitempty"`
	CrashesTotal        *int     `json:"crashesTotal,omitempty"`
	CrashesFatal        *int     `json:"crashesFatal,omitempty"`
	CrashesInjury       *int     `json:"crashesInjury,omitempty"`
	InspectionsTotal    *int     `json:"inspectionsTotal,omitempty"`
	OutOfServicePercent *float64 `json:"outOfServicePercent,omitempty"`

	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
	SourceUniqueId    string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId   string             `json:"qubitOnUniqueId,omitempty"`
}

// ── India Identity Validation ─────────────────────────────────────────────

// IndiaIdentityRequest is the input for ValidateIndiaIdentity. Mirrors the
// canonical INIdentityRequest : BaseTaxRequest. The Dob property has an
// explicit [JsonPropertyName("dob")] override that bypasses CamelCase, so the
// wire field is plain "dob".
//
// BusinessEntityType / IdentityState / BusinessEntityTypeCheckRequired /
// EntityNameMatchCheckRequired are inherited from BaseTaxRequest. The two
// "*CheckRequired" booleans are marked [JsonIgnore] on the canonical model
// so they are not transmitted to the server today; the SDK still exposes
// them as pointer fields with omitempty so they remain absent from the wire
// unless the caller sets them, future-proofing against the [JsonIgnore]
// being lifted.
type IndiaIdentityRequest struct {
	BaseRequest
	IdentityNumber                   string `json:"identityNumber"`
	IdentityNumberType               string `json:"identityNumberType"`
	EntityName                       string `json:"entityName,omitempty"`
	BusinessEntityType               string `json:"businessEntityType,omitempty"`
	IdentityState                    string `json:"identityState,omitempty"`
	BusinessEntityTypeCheckRequired  *bool  `json:"businessEntityTypeCheckRequired,omitempty"`
	EntityNameMatchCheckRequired     *bool  `json:"entityNameMatchCheckRequired,omitempty"`
	// Dob is the date of birth in YYYY-MM-DD format. Required for Driver
	// License validation; optional for Voter Registration. Wire field "dob".
	Dob string `json:"dob,omitempty"`
}

// IndiaIdentityResponse is the output from ValidateIndiaIdentity. Mirrors the
// canonical INIdentityResponse : BaseTaxResponse — the response inherits all
// fields from BaseTaxResponse and adds nothing of its own. The fields exposed
// here are those inherited from BaseTaxResponse / BaseValidationResponse.
type IndiaIdentityResponse struct {
	// From BaseTaxResponse / BaseValidationResponse
	IdentityNumberValidationID int64              `json:"identityNumberValidationID,omitempty"`
	IdentityNumberType         string             `json:"identityNumberType,omitempty"`
	IdentityNumber             string             `json:"identityNumber,omitempty"`
	EntityName                 string             `json:"entityName,omitempty"`
	RequestedByIPAddress       string             `json:"requestedByIPAddress,omitempty"`
	TaxAddress                 *AddressResponse   `json:"taxAddress,omitempty"`
	BusinessEntityType         string             `json:"businessEntityType,omitempty"`
	SmartvmBusinessEntityType  string             `json:"smartvmBusinessEntityType,omitempty"`
	TaxAdditionalInfo          *TaxAdditionalInfo `json:"taxAdditionalInfo,omitempty"`
	SubSection                 *int               `json:"subSection,omitempty"`
	SubSectionCode             *string            `json:"subSectionCode,omitempty"`
	SubSectionDescription      *string            `json:"subSectionDescription,omitempty"`

	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
	SourceUniqueId    string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId   string             `json:"qubitOnUniqueId,omitempty"`
}

// ── Certification ─────────────────────────────────────────────────────────

// CertificationRequest is the input for ValidateCertification and LookupCertification.
type CertificationRequest struct {
	BaseRequest
	CompanyName         string `json:"companyName"`
	Country             string `json:"country"`
	City                string `json:"city,omitempty"`
	State               string `json:"state,omitempty"`
	ZipCode             string `json:"zipCode,omitempty"`
	AddressLine1        string `json:"addressLine1,omitempty"`
	AddressLine2        string `json:"addressLine2,omitempty"`
	IdentityType        string `json:"identityType,omitempty"`
	CertificationType   string `json:"certificationType,omitempty"`
	CertificationGroup  string `json:"certificationGroup,omitempty"`
	CertificationNumber string `json:"certificationNumber,omitempty"`
}

// CertificationResponse is the output from ValidateCertification (single object)
// and one entry in the array returned by LookupCertification.
type CertificationResponse struct {
	IsCertified         bool                     `json:"isCertified,omitempty"`
	CertificationType   string                   `json:"certificationType,omitempty"`
	CertificationGroup  string                   `json:"certificationGroup,omitempty"`
	CertificationNumber string                   `json:"certificationNumber,omitempty"`
	CertifyingBody      string                   `json:"certifyingBody,omitempty"`
	IssuedDate          string                   `json:"issuedDate,omitempty"`
	ExpirationDate      string                   `json:"expirationDate,omitempty"`
	Status              string                   `json:"status,omitempty"`
	CompanyName         string                   `json:"companyName,omitempty"`
	Certifications      []map[string]interface{} `json:"certifications,omitempty"`

	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
	SourceUniqueId    string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId   string             `json:"qubitOnUniqueId,omitempty"`
}

// ── Business Classification ───────────────────────────────────────────────

// BusinessClassificationRequest is the input for LookupBusinessClassification.
type BusinessClassificationRequest struct {
	BaseRequest
	CompanyName string `json:"companyName"`
	City        string `json:"city"`
	State       string `json:"state"`
	Country     string `json:"country"`
	Address1    string `json:"address1,omitempty"`
	Address2    string `json:"address2,omitempty"`
	Phone       string `json:"phone,omitempty"`
	PostalCode  string `json:"postalCode,omitempty"`
}

// BusinessClassificationResponse is one entry in the array returned by
// LookupBusinessClassification.
type BusinessClassificationResponse struct {
	NAICSCode        string   `json:"naicsCode,omitempty"`
	NAICSDescription string   `json:"naicsDescription,omitempty"`
	SICCode          string   `json:"sicCode,omitempty"`
	SICDescription   string   `json:"sicDescription,omitempty"`
	Industry         string   `json:"industry,omitempty"`
	Confidence       *float64 `json:"confidence,omitempty"`

	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
	SourceUniqueId    string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId   string             `json:"qubitOnUniqueId,omitempty"`
}

// ── Payment Terms ─────────────────────────────────────────────────────────

// PaymentTermsRequest is the input for AnalyzePaymentTerms.
type PaymentTermsRequest struct {
	BaseRequest
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
	RecommendedTerm  *float64 `json:"recommendedTerm,omitempty"`
	PotentialSavings *float64 `json:"potentialSavings,omitempty"`
	Recommendation   string   `json:"recommendation,omitempty"`

	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
	SourceUniqueId    string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId   string             `json:"qubitOnUniqueId,omitempty"`
}

// ── Exchange Rates ────────────────────────────────────────────────────────

// ExchangeRateRequest is the input for LookupExchangeRates.
//
// Dates is a slice of time.Time values; they are marshalled to RFC3339 UTC
// timestamps in the JSON body. The API expects a non-null array — callers
// passing a nil/empty slice will get an empty JSON array on the wire.
//
// Note: BaseCurrency is sent as a path parameter (URL-escaped) and Dates is
// the JSON body; this struct intentionally does not embed BaseRequest because
// the underlying server endpoint reads only the path + body.
type ExchangeRateRequest struct {
	BaseCurrency string      `json:"-"`
	Dates        []time.Time `json:"-"`
}

// ExchangeRateResponse is one entry in the array returned by
// /api/currency/exchange-rates/{baseCurrency}. Mirrors the canonical
// CurrencyDayExchangeRates exactly. ExchangeRates is a list (not a map) where
// each entry carries Currency / Rate / DataSource.
type ExchangeRateResponse struct {
	Date          string                 `json:"date,omitempty"`
	BaseCurrency  string                 `json:"baseCurrency,omitempty"`
	ExchangeRates []CurrencyExchangeRate `json:"exchangeRates,omitempty"`
}

// CurrencyExchangeRate mirrors the canonical CurrencyExchangeRate.
// Rate is decimal on the .NET side; we use float64 for portability.
type CurrencyExchangeRate struct {
	Currency   string  `json:"currency,omitempty"`
	Rate       float64 `json:"rate,omitempty"`
	DataSource string  `json:"dataSource,omitempty"`
}

// ── SAP Ariba Supplier ────────────────────────────────────────────────────

// AribaSupplierRequest is the input for LookupAribaSupplier and ValidateAribaSupplier.
// Acronym JSON tag ANID→anid per .NET CamelCase JsonNamingPolicy (empirically
// verified against System.Text.Json with .NET 10 — consecutive leading
// capitals are fully lower-cased).
type AribaSupplierRequest struct {
	BaseRequest
	ANID string `json:"anid"`
}

// AribaSupplierResponse is the output from LookupAribaSupplier (one entry per
// matching supplier, returned as a JSON array) and ValidateAribaSupplier
// (single object). Acronym JSON tag ANID→anid.
type AribaSupplierResponse struct {
	ANID         string                 `json:"anid,omitempty"`
	CompanyName  string                 `json:"companyName,omitempty"`
	DunsNumber   string                 `json:"dunsNumber,omitempty"`
	TaxId        string                 `json:"taxId,omitempty"`
	Country      string                 `json:"country,omitempty"`
	City         string                 `json:"city,omitempty"`
	State        string                 `json:"state,omitempty"`
	PostalCode   string                 `json:"postalCode,omitempty"`
	AddressLine1 string                 `json:"addressLine1,omitempty"`
	Status       string                 `json:"status,omitempty"`
	Profile      map[string]interface{} `json:"profile,omitempty"`

	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
	SourceUniqueId    string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId   string             `json:"qubitOnUniqueId,omitempty"`
}

// ── Gender Identification ─────────────────────────────────────────────────

// GenderRequest is the input for IdentifyGender.
type GenderRequest struct {
	BaseRequest
	Name    string `json:"name"`
	Country string `json:"country,omitempty"`
}

// GenderResponse is the output from IdentifyGender. Mirrors the canonical
// GenderizeResponse exactly — Confidence is a string on the wire (server
// property `public string Confidence`), not a float.
type GenderResponse struct {
	Name       string `json:"name,omitempty"`
	Gender     string `json:"gender,omitempty"`
	Confidence string `json:"confidence,omitempty"`
	Country    string `json:"country,omitempty"`
}

// ── Continuous Screening ──────────────────────────────────────────────────

// ContinuousScreeningRequest is the input for ScreenContinuous.
//
// NOTE: server-side the endpoint /api/continuous-screening/screen currently
// returns HTTP 501 Not Implemented. The fields below are best-guess from the
// planned schema and may change once the server implementation ships.
type ContinuousScreeningRequest struct {
	BaseRequest
	EntityName         string `json:"entityName"`
	Country            string `json:"country,omitempty"`
	IdentityNumber     string `json:"identityNumber,omitempty"`
	IdentityNumberType string `json:"identityNumberType,omitempty"`
	CallbackUrl        string `json:"callbackUrl,omitempty"`
}

// ContinuousScreeningResponse is the output from ScreenContinuous.
//
// NOTE: the server currently returns 501 Not Implemented for this endpoint;
// the response shape below is the planned schema and may change. Field names
// are best-guess and not yet stable.
type ContinuousScreeningResponse struct {
	MonitoringId    string                   `json:"monitoringId,omitempty"`
	Status          string                   `json:"status,omitempty"`
	EnrolledAt      string                   `json:"enrolledAt,omitempty"`
	NextScreeningAt string                   `json:"nextScreeningAt,omitempty"`
	Matches         []map[string]interface{} `json:"matches,omitempty"`
	SourceUniqueId  string                   `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId string                   `json:"qubitOnUniqueId,omitempty"`
}

// ── Check Status (bulk callback) ──────────────────────────────────────────

// CallBackRequest mirrors the canonical CallBackRequest : BaseRequest from
// smartvm.BusinessEntities.Client.API.Status. The server expects a single
// CallBackID Guid; the property name uses CamelCase JsonNamingPolicy which
// preserves the leading capital "C" but lowercases only the first character —
// the wire field is "callBackID".
type CallBackRequest struct {
	BaseRequest
	// CallBackID is the bulk request callback identifier (UUID/Guid).
	CallBackID string `json:"callBackID"`
}

// CheckStatusResponse mirrors the canonical CheckStatusResponse from
// smartvm.BusinessEntities.Client.API.API.CheckStatus.
type CheckStatusResponse struct {
	// CallBackID is the bulk request callback identifier (UUID/Guid).
	CallBackID string `json:"callBackID,omitempty"`
	// PickUpUrl is the URL to retrieve the completed bulk results.
	PickUpUrl string `json:"pickUpUrl,omitempty"`
	// Status is the bulk processing status code.
	Status string `json:"status,omitempty"`
	// StatusDescription is a human-readable status description.
	StatusDescription string `json:"statusDescription,omitempty"`
	// TotalRecords is the number of records submitted for this callback ID.
	TotalRecords int `json:"totalRecords,omitempty"`
	// ProcessedRecords is the number of records processed for this callback ID.
	ProcessedRecords int `json:"processedRecords,omitempty"`
}
