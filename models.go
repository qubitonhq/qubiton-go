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
	GeoCode      string `json:"geoCode,omitempty"`
	InCityLimit  string `json:"inCityLimit,omitempty"`
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

	PreferredLanguage string                  `json:"preferredLanguage,omitempty"`
	AdditionalInfo    map[string]interface{}   `json:"additionalInfo,omitempty"`
	AddressCodesInfo  []map[string]interface{} `json:"addressCodesInfo,omitempty"`

	// Validation
	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
	SourceUniqueId    string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId   string             `json:"qubitOnUniqueId,omitempty"`
}

// Country represents a country with ISO codes.
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
	TaxNumber          string `json:"taxNumber"`
	TaxType            string `json:"taxType"`
	Country            string `json:"country"`
	CompanyName        string `json:"companyName"`
	BusinessEntityType string `json:"businessEntityType,omitempty"`
}

// TaxResponse is the output from ValidateTax.
type TaxResponse struct {
	TaxValid          bool   `json:"taxValid"`
	IsEntityNameMatch *bool  `json:"isEntityNameMatch,omitempty"`
	EntityName        string `json:"entityName,omitempty"`
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
}

// ── Bank Account Validation (Ownership via BankPro) ─────────────────────

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
	// Account details
	BankAccountNumber string `json:"bankAccountNumber,omitempty"`
	BankCurrencyCode  string `json:"bankCurrencyCode,omitempty"`
	AccountType       string `json:"accountType,omitempty"`
	AccountHolder     string `json:"accountHolder,omitempty"`
	IBAN              string `json:"iban,omitempty"`
	SwiftCode         string `json:"swiftCode,omitempty"`
	CLABENumber       string `json:"clabeNumber,omitempty"`

	// Tax/identity
	TaxIdNumber string `json:"taxIdNumber,omitempty"`
	TaxType     string `json:"taxType,omitempty"`

	// Scoring
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

	// Validation details
	ValidationStatus      string                       `json:"validationStatus,omitempty"`
	BankAccountValidations *BankAccountValidationDetail `json:"bankAccountValidations,omitempty"`
	BankProValidations     *BankProDetails              `json:"bankProValidations,omitempty"`

	// Request/response holder info
	RequestedInfo *BankAccountHolderInfo `json:"requestedInfo,omitempty"`
	ResponseInfo  *BankAccountHolderInfo `json:"responseInfo,omitempty"`

	// Metadata
	VendorName      string `json:"vendorName,omitempty"`
	VendorCountry   string `json:"vendorCountry,omitempty"`
	SourceUniqueId  string `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId string `json:"qubitOnUniqueId,omitempty"`
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
	VendorMatch       string `json:"vendorMatch,omitempty"`
	AccountHolderMatch string `json:"accountHolderMatch,omitempty"`
	BankNameMatch     string `json:"bankNameMatch,omitempty"`
	BankAccountMatch  string `json:"bankAccountMatch,omitempty"`
	BankCurrencyMatch string `json:"bankCurrencyMatch,omitempty"`
	BankCountryMatch  string `json:"bankCountryMatch,omitempty"`
	VendorCountryMatch string `json:"vendorCountryMatch,omitempty"`
	AccountTypeMatch  string `json:"accountTypeMatch,omitempty"`
	TaxIDMatch        string `json:"taxIDMatch,omitempty"`
	TaxIDTypeMatch    string `json:"taxIDTypeMatch,omitempty"`

	// Composite scores
	OverallMatchScore        *float64 `json:"overallMatchScore,omitempty"`
	ValidationHistoryScore   *float64 `json:"validationHistoryScore,omitempty"`
	VendorNameMatchScore     *float64 `json:"vendorNameMatchScore,omitempty"`
	AccountHolderMatchScore  *float64 `json:"accountHolderMatchScore,omitempty"`
	BankNameMatchScore       *float64 `json:"bankNameMatchScore,omitempty"`

	// History
	PreviousValidationMatch    string `json:"previousValidationMatch,omitempty"`
	ApexHistoryMatch           string `json:"apexHistoryMatch,omitempty"`
	DaysSinceLastValidation    *int   `json:"daysSinceLastValidation,omitempty"`
	RecentlyValidated          string `json:"recentlyValidated,omitempty"`
	RecentlySeen               string `json:"recentlySeen,omitempty"`
	ValidationAttempts         *int   `json:"validationAttempts,omitempty"`
	ValidationPASSCount        *int   `json:"validationPASSCount,omitempty"`
	ValidationFAILCount        *int   `json:"validationFAILCount,omitempty"`

	// Risk flags
	IsRedFlagsToConsider                 *bool    `json:"isRedFlagsToConsider,omitempty"`
	RedFlagReasons                       []string `json:"redFlagReasons,omitempty"`
	FoundBankAccountWithSameVendor       string   `json:"foundBankAccountWithSameVendor,omitempty"`
	FoundBankAccountWithSameTaxID        string   `json:"foundBankAccountWithSameTaxID,omitempty"`
	FoundBankAccountDifferentVendor      string   `json:"foundBankAccountDifferentVendor,omitempty"`
	FoundSameVendorDifferentAccounts     string   `json:"foundSameVendorDifferentAccounts,omitempty"`
	FoundSameTaxIDDifferentAccounts      string   `json:"foundSameTaxIDDifferentAccounts,omitempty"`

	// Result descriptions
	ResultDescription    string `json:"resultDescription,omitempty"`
	MatchCodeDescription string `json:"matchCodeDescription,omitempty"`
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

// BankProResponse is the output from ValidateBankPro (same structure as BankAccountResponse).
type BankProResponse = BankAccountResponse

// ── Email Validation ──────────────────────────────────────────────────────

// EmailRequest is the input for ValidateEmail.
type EmailRequest struct {
	EmailAddress string `json:"emailAddress"`
}

// EmailResponse is the output from ValidateEmail.
type EmailResponse struct {
	EmailAddress    string   `json:"emailAddress,omitempty"`
	EmailType       string   `json:"emailType,omitempty"`
	FraudScore      *int     `json:"fraudScore,omitempty"`
	FormatCheck     string   `json:"formatCheck,omitempty"`
	SMTPCheck       string   `json:"smtpCheck,omitempty"`
	DNSCheck        string   `json:"dnsCheck,omitempty"`
	FreeCheck       string   `json:"freeCheck,omitempty"`
	DisposableCheck string   `json:"disposableCheck,omitempty"`
	CatchAllCheck   string   `json:"catchAllCheck,omitempty"`
	MXRecords       []string `json:"mxRecords,omitempty"`
	Audit           *EmailAudit `json:"audit,omitempty"`

	// Validation (from BaseValidationResponse)
	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
	SourceUniqueId    string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId   string             `json:"qubitOnUniqueId,omitempty"`
}

// EmailAudit contains audit timestamps for email validation.
type EmailAudit struct {
	AuditCreatedDate string `json:"auditCreatedDate,omitempty"`
	AuditUpdatedDate string `json:"auditUpdatedDate,omitempty"`
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
	PhoneNumber            string                  `json:"phoneNumber,omitempty"`
	PhoneExtension         string                  `json:"phoneExtension,omitempty"`
	PhoneCountryCode       string                  `json:"phoneCountryCode,omitempty"`
	FullPhoneNumber        string                  `json:"fullPhoneNumber,omitempty"`
	FraudScore             *int                    `json:"fraudScore,omitempty"`
	Country                *Country                `json:"country,omitempty"`
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
	CompanyName string `json:"companyName"`
	Country     string `json:"country"`
	State       string `json:"state,omitempty"`
	City        string `json:"city,omitempty"`
}

// BusinessRegistrationResponse is the output from LookupBusinessRegistration.
type BusinessRegistrationResponse struct {
	BusinessRegistrations []BusinessRegistration `json:"businessRegistrations"`
	ValidationDescription string                 `json:"validationDescription"`
	ValidationPass        *bool                  `json:"validationPass"`
	Score                 float64                `json:"score"`
	SourceResultCode      string                 `json:"sourceResultCode"`
	ValidationResultCode  string                 `json:"validationResultCode"`
	ValidationDate        string                 `json:"validationDate"`
	SourceUniqueId        string                 `json:"sourceUniqueId"`
	QubitOnUniqueId       string                 `json:"qubitOnUniqueId"`
}

// BusinessRegistration represents a single business entity from a government registry.
type BusinessRegistration struct {
	RegistrationId                string   `json:"registrationId"`
	EntityName                    string   `json:"entityName"`
	Status                        string   `json:"status"`
	StatusReason                  string   `json:"statusReason,omitempty"`
	BusinessEntityType            string   `json:"businessEntityType"`
	BusinessEntityTypeDescription string   `json:"businessEntityTypeDescription,omitempty"`
	Jurisdiction                  string   `json:"jurisdiction,omitempty"`
	RegistrationDate              string   `json:"registrationDate,omitempty"`
	FormationDate                 string   `json:"formationDate,omitempty"`
	ExpirationDate                string   `json:"expirationDate,omitempty"`
	Duration                      string   `json:"duration,omitempty"`
	TaxNumber                     string   `json:"taxNumber,omitempty"`
	PreviousEntityNames           []string `json:"previousEntityNames,omitempty"`
	Addresses                     []BusinessRegistrationAddress `json:"addresses,omitempty"`
	Persons                       []BusinessRegistrationPerson  `json:"persons,omitempty"`
	Phones                        []BusinessRegistrationPhone   `json:"phones,omitempty"`
}

// BusinessRegistrationAddress is an address associated with a registered business.
type BusinessRegistrationAddress struct {
	AddressType  string `json:"addressType,omitempty"`
	AddressLine1 string `json:"addressLine1,omitempty"`
	AddressLine2 string `json:"addressLine2,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	Zip          string `json:"zip,omitempty"`
}

// BusinessRegistrationPerson is an officer, director, or agent associated with a business.
type BusinessRegistrationPerson struct {
	Name        string                        `json:"name"`
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
	ParticipantId   string `json:"participantId"`
	DirectoryLookup *bool  `json:"directoryLookup,omitempty"`
}

// PeppolResponse is the output from ValidatePeppol.
type PeppolResponse struct {
	IsValid                bool     `json:"isValid"`
	ValidationType         string   `json:"validationType,omitempty"`
	ParticipantId          string   `json:"participantId,omitempty"`
	IcdCode                string   `json:"icdCode,omitempty"`
	IcdScheme              string   `json:"icdScheme,omitempty"`
	Identifier             string   `json:"identifier,omitempty"`
	CountryIso2            string   `json:"countryIso2,omitempty"`
	RegisteredInDirectory  *bool    `json:"registeredInDirectory,omitempty"`
	ParticipantName        string   `json:"participantName,omitempty"`
	DocumentTypes          []string `json:"documentTypes,omitempty"`
	Errors                 []string `json:"errors,omitempty"`
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
	IsMatch        bool                        `json:"isMatch"`
	Description    string                      `json:"description,omitempty"`
	AdditionalInfo *SanctionsAdditionalInfo     `json:"additionalInfo,omitempty"`
	SourceId       int                         `json:"sourceId,omitempty"`
	Score          float64                     `json:"score,omitempty"`
	SourceResultCode    string                 `json:"sourceResultCode,omitempty"`
	ValidationResultCode string                `json:"validationResultCode,omitempty"`
	ValidationDate      string                 `json:"validationDate,omitempty"`
	SourceUniqueId      string                 `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId     string                 `json:"qubitOnUniqueId,omitempty"`
}

// SanctionsAdditionalInfo contains the detailed match results per sanctions list.
type SanctionsAdditionalInfo struct {
	ResultCode string                `json:"resultCode,omitempty"`
	Lists      []SanctionsListDetail `json:"lists,omitempty"`
}

// SanctionsListDetail represents matches from a specific sanctions list.
type SanctionsListDetail struct {
	Code                  string             `json:"code,omitempty"`
	Name                  string             `json:"name,omitempty"`
	Description           string             `json:"description,omitempty"`
	TotalMatches          *int               `json:"totalMatches,omitempty"`
	EntityMatches         *int               `json:"entityMatches,omitempty"`
	BlockedCountryMatches *int               `json:"blockedCountryMatches,omitempty"`
	DateAdded             string             `json:"dateAdded,omitempty"`
	Entities              []SanctionsEntity  `json:"entities,omitempty"`
}

// SanctionsEntity represents a matched entity on a sanctions list.
type SanctionsEntity struct {
	Id              string                  `json:"id,omitempty"`
	Number          string                  `json:"number,omitempty"`
	MaxScore        *float64                `json:"maxScore,omitempty"`
	FirstName       string                  `json:"firstName,omitempty"`
	LastName        string                  `json:"lastName,omitempty"`
	MiddleName      string                  `json:"middleName,omitempty"`
	OtherName       string                  `json:"otherName,omitempty"`
	WholeName       string                  `json:"wholeName,omitempty"`
	HasAliasMatch   *bool                   `json:"hasAliasMatch,omitempty"`
	HasAddressMatch *bool                   `json:"hasAddressMatch,omitempty"`
	Score           *float64                `json:"score,omitempty"`
	Type            string                  `json:"type,omitempty"`
	Programs        string                  `json:"programs,omitempty"`
	Remarks         string                  `json:"remarks,omitempty"`
	Aliases         []SanctionsEntityAlias   `json:"aliases,omitempty"`
	Addresses       []SanctionsEntityAddress `json:"addresses,omitempty"`
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
	Name    string `json:"name"`
	Country string `json:"country"`
}

// PEPResponse is the output from ScreenPEP.
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
	Id               string              `json:"id,omitempty"`
	Name             string              `json:"name,omitempty"`
	SortName         string              `json:"sortName,omitempty"`
	GivenName        string              `json:"givenName,omitempty"`
	FamilyName       string              `json:"familyName,omitempty"`
	HonorificPrefix  string              `json:"honorificPrefix,omitempty"`
	Gender           string              `json:"gender,omitempty"`
	BirthDate        string              `json:"birthDate,omitempty"`
	DeathDate        string              `json:"deathDate,omitempty"`
	Email            string              `json:"email,omitempty"`
	Image            string              `json:"image,omitempty"`
	ContactDetails   []PEPContactDetail  `json:"contactDetails,omitempty"`
	PersonIdentifiers []PEPIdentifier    `json:"personIdentifiers,omitempty"`
	PersonLinks      []PEPLink           `json:"personLinks,omitempty"`
	PersonOtherNames []PEPOtherName      `json:"personOtherNames,omitempty"`
	PersonImages     []PEPImage          `json:"personImages,omitempty"`
}

// PEPOrganization represents an organization associated with a PEP.
type PEPOrganization struct {
	Id                      string         `json:"id,omitempty"`
	Name                    string         `json:"name,omitempty"`
	Classification          string         `json:"classification,omitempty"`
	OrganizationType        string         `json:"organizationType,omitempty"`
	Seats                   *int           `json:"seats,omitempty"`
	Image                   string         `json:"image,omitempty"`
	OrganizationIdentifiers []PEPIdentifier `json:"organizationIdentifiers,omitempty"`
	OrganizationLinks       []PEPLink      `json:"organizationLinks,omitempty"`
	OrganizationOtherNames  []PEPOtherName `json:"organizationOtherNames,omitempty"`
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
	Id              string         `json:"id,omitempty"`
	Name            string         `json:"name,omitempty"`
	AreaType        string         `json:"areaType,omitempty"`
	AreaIdentifiers []PEPIdentifier `json:"areaIdentifiers,omitempty"`
	AreaOtherNames  []PEPOtherName `json:"areaOtherNames,omitempty"`
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

// PEPLink is a URL link with optional note.
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

// PEPImage is an image URL for a PEP person.
type PEPImage struct {
	URL string `json:"url,omitempty"`
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
}

// ── Entity Risk Assessment ────────────────────────────────────────────────

// EntityRiskRequest is the input for AssessEntityRisk.
type EntityRiskRequest struct {
	CompanyName            string `json:"companyName"`
	CountryOfIncorporation string `json:"countryOfIncorporation,omitempty"`
	Category               string `json:"category,omitempty"`
	URL                    string `json:"url,omitempty"`
	BusinessEntityType     string `json:"businessEntityType,omitempty"`
}

// EntityRiskResponse is the output from AssessEntityRisk.
type EntityRiskResponse struct {
	// Company identity
	CompanyName              string   `json:"companyName,omitempty"`
	Name                     string   `json:"name,omitempty"`
	QubitOnCompanyName       string   `json:"qubitOnCompanyName,omitempty"`
	CompanyNameDBA           []string `json:"companyNameDBA,omitempty"`
	QubitOnCompanyNameDBA    []string `json:"qubitOnCompanyNameDBA,omitempty"`
	QubitOnEntityId          string   `json:"qubitOnEntityId,omitempty"`

	// Entity classification
	BusinessEntityType       string `json:"businessEntityType,omitempty"`
	QubitOnBusinessEntityType string `json:"qubitOnBusinessEntityType,omitempty"`
	BusinessEntitySubType    string `json:"businessEntitySubType,omitempty"`
	CountryOfIncorporation   string `json:"countryOfIncorporation,omitempty"`
	StateOfIncorporation     string `json:"stateOfIncorporation,omitempty"`

	// Operations
	URL            []string `json:"url,omitempty"`
	YearStarted    string   `json:"yearStarted,omitempty"`
	SalesVolume    *float64 `json:"salesVolume,omitempty"`
	EmployeesOnSite *float64 `json:"employeesOnSite,omitempty"`
	EmployeesTotal *float64 `json:"employeesTotal,omitempty"`
	LocationType   string   `json:"locationType,omitempty"`
	LineOfBusiness string   `json:"lineOfBusiness,omitempty"`

	// Risk scoring
	HighRiskScore  float64                `json:"highRiskScore,omitempty"`
	IsHighRisk     bool                   `json:"isHighRisk,omitempty"`
	SentimentCount map[string]float64     `json:"sentimentCount,omitempty"`
	TotalArticles  float64                `json:"totalArticles,omitempty"`
	Categories     []map[string]interface{} `json:"categories,omitempty"`
	RiskScores     []map[string]interface{} `json:"riskScores,omitempty"`
	FraudScores    []map[string]interface{} `json:"fraudScores,omitempty"`

	// Nested entities
	EntityAddresses  []map[string]interface{} `json:"entityAddresses,omitempty"`
	EntityPhones     []map[string]interface{} `json:"entityPhones,omitempty"`
	EntityIdentities []map[string]interface{} `json:"entityIdentities,omitempty"`
	EntityPersons    []map[string]interface{} `json:"entityPersons,omitempty"`

	// Validation
	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
	SourceUniqueId    string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId   string             `json:"qubitOnUniqueId,omitempty"`
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
	CreditScore     int     `json:"creditScore"`
	CreditLimit     float64 `json:"creditLimit"`
	PaymentBehavior string  `json:"paymentBehavior"`
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
	ESGId            int    `json:"ESGId,omitempty"`
	Name             string `json:"name,omitempty"`
	CompanyDesc      string `json:"company_desc,omitempty"`
	Grade            string `json:"grade,omitempty"`
	ExchangeSymbol   string `json:"exchange_symbol,omitempty"`
	StockSymbol      string `json:"stock_symbol,omitempty"`
	Industry         int    `json:"industry,omitempty"`
	CountryCode      int    `json:"country,omitempty"`
	EnvironmentScore int    `json:"e"`
	SocialScore      int    `json:"s"`
	GovernanceScore  int    `json:"g"`
	Total            int    `json:"total"`
	Progress         []ESGProgress       `json:"progress,omitempty"`
	TopRecords       *ESGTopRecords      `json:"toprecords,omitempty"`

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
	DomainName string `json:"domain"`
}

// DomainSecurityResponse is the output from DomainSecurityReport.
type DomainSecurityResponse struct {
	RiskScore   float64 `json:"riskScore"`
	ThreatLevel string  `json:"threatLevel"`
}

// ── IP Quality ────────────────────────────────────────────────────────────

// IPQualityRequest is the input for CheckIPQuality.
type IPQualityRequest struct {
	IPAddress string `json:"ipAddress"`
	UserAgent string `json:"userAgent,omitempty"`
}

// IPQualityResponse is the output from CheckIPQuality.
type IPQualityResponse struct {
	IPAddress      string   `json:"ipAddress,omitempty"`
	UserAgent      string   `json:"userAgent,omitempty"`
	FraudScore     int      `json:"fraudScore"`
	IsProxy        bool     `json:"isProxy"`
	IsVPN          bool     `json:"isVPN"`
	IsActiveVPN    bool     `json:"isActiveVPN,omitempty"`
	ISP            string   `json:"isp,omitempty"`
	IsCrawler      bool     `json:"isCrawler,omitempty"`
	IsTOR          bool     `json:"isTOR,omitempty"`
	IsActiveTOR    bool     `json:"isActiveTOR,omitempty"`
	RecentAbuse    bool     `json:"recentAbuse,omitempty"`
	Organization   string   `json:"organization,omitempty"`
	Latitude       float64  `json:"latitude,omitempty"`
	Longitude      float64  `json:"longitude,omitempty"`
	Timezone       string   `json:"timezone,omitempty"`
	City           string   `json:"city,omitempty"`
	Region         string   `json:"region,omitempty"`
	Host           string   `json:"host,omitempty"`
	ASN            int      `json:"asn,omitempty"`
	ConnectionType string   `json:"connectionType,omitempty"`
	AbuseVelocity  string   `json:"abuseVelocity,omitempty"`
	BotStatus      bool     `json:"botStatus,omitempty"`
	DeviceModel    string   `json:"deviceModel,omitempty"`
	DeviceBrand    string   `json:"deviceBrand,omitempty"`
	Mobile         bool     `json:"mobile,omitempty"`
	OperatingSystem string  `json:"operatingSystem,omitempty"`
	Browser        string   `json:"browser,omitempty"`
	Country        *Country `json:"country,omitempty"`

	// Validation
	ValidationResults []ValidationResult `json:"validationResults,omitempty"`
	SourceUniqueId    string             `json:"sourceUniqueId,omitempty"`
	QubitOnUniqueId   string             `json:"qubitOnUniqueId,omitempty"`
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
}

// ── DOT Carrier Lookup ────────────────────────────────────────────────────

// DOTCarrierRequest is the input for LookupDOTCarrier.
type DOTCarrierRequest struct {
	DotNumber  string `json:"dotNumber"`
	EntityName string `json:"entityName,omitempty"`
}

// DOTCarrierResponse is the output from LookupDOTCarrier.
type DOTCarrierResponse struct {
	Found        bool   `json:"found"`
	CarrierName  string `json:"carrierName"`
	DotNumber    string `json:"dotNumber"`
	SafetyRating string `json:"safetyRating"`
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
	RecommendedTerm  float64 `json:"recommendedTerm"`
	PotentialSavings float64 `json:"potentialSavings"`
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
}

// ── Tax Format Reference ──────────────────────────────────────────────────

// TaxFormatsResponse is the output from GetSupportedTaxFormats.
type TaxFormatsResponse struct {
	Countries []map[string]interface{} `json:"countries"`
}

// ── Peppol Schemes Reference ──────────────────────────────────────────────

// PeppolSchemesResponse is the output from GetPeppolSchemes.
type PeppolSchemesResponse struct {
	Schemes []map[string]interface{} `json:"schemes"`
}
