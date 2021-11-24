package main

type GetBankGroupsResponse struct {
	JSON            string      `json:"JSON"`
	AccountNickname interface{} `json:"accountNickname"`
	AccountNumber   interface{} `json:"accountNumber"`
	AccountType     interface{} `json:"accountType"`
	BankGroups      []BankGroup `json:"bankGroups"`
	BankToID        interface{} `json:"bankToId"`
	Banks           interface{} `json:"banks"`
	DocumentNumber  interface{} `json:"documentNumber"`
	DocumentType    interface{} `json:"documentType"`
	DocumentTypes   interface{} `json:"documentTypes"`
	FiduciariaTypes interface{} `json:"fiduciariaTypes"`
	ProductName     interface{} `json:"productName"`
	ProductTypes1   interface{} `json:"productTypes1"`
	ProductTypes2   interface{} `json:"productTypes2"`
}

type BankGroup struct {
	AccType     string      `json:"accType"`
	AccTypeCode string      `json:"accTypeCode"`
	AcctID      interface{} `json:"acctId"`
	Balance     interface{} `json:"balance"`
	BankID      string      `json:"bankId"`
	Category    interface{} `json:"category"`
	Description string      `json:"description"`
	EXTData     struct{}    `json:"extData"`
	ID          string      `json:"id"`
	NickName    interface{} `json:"nickName"`
}

type GetProductTypesResponse struct {
	JSON            string         `json:"JSON"`
	AccountNickname interface{}    `json:"accountNickname"`
	AccountNumber   interface{}    `json:"accountNumber"`
	AccountType     interface{}    `json:"accountType"`
	BankGroups      interface{}    `json:"bankGroups"`
	BankToID        interface{}    `json:"bankToId"`
	Banks           interface{}    `json:"banks"`
	DocumentNumber  interface{}    `json:"documentNumber"`
	DocumentType    interface{}    `json:"documentType"`
	DocumentTypes   interface{}    `json:"documentTypes"`
	FiduciariaTypes interface{}    `json:"fiduciariaTypes"`
	ProductName     interface{}    `json:"productName"`
	ProductTypes1   []ProductTypes `json:"productTypes1"`
	ProductTypes2   []ProductTypes `json:"productTypes2"`
}

type ProductTypes struct {
	Code           string `json:"code"`
	CreationDate   string `json:"creationDate"`
	Description    string `json:"description"`
	ExpirationDate string `json:"expirationDate"`
	LanguageCode   string `json:"languageCode"`
}

type GetDocumentTypesResponse struct {
	JSON            string         `json:"JSON"`
	AccountNickname interface{}    `json:"accountNickname"`
	AccountNumber   interface{}    `json:"accountNumber"`
	AccountType     interface{}    `json:"accountType"`
	BankGroups      interface{}    `json:"bankGroups"`
	BankToID        interface{}    `json:"bankToId"`
	Banks           interface{}    `json:"banks"`
	DocumentNumber  interface{}    `json:"documentNumber"`
	DocumentType    interface{}    `json:"documentType"`
	DocumentTypes   []DocumentType `json:"documentTypes"`
	FiduciariaTypes interface{}    `json:"fiduciariaTypes"`
	ProductName     interface{}    `json:"productName"`
	ProductTypes1   interface{}    `json:"productTypes1"`
	ProductTypes2   interface{}    `json:"productTypes2"`
}

type DocumentType struct {
	AllowedSegments string      `json:"allowedSegments"`
	BankID          interface{} `json:"bankId"`
	Code            string      `json:"code"`
	CodeBank        string      `json:"codeBank"`
	CreationDate    string      `json:"creationDate"`
	Description     string      `json:"description"`
	ExpirationDate  string      `json:"expirationDate"`
	LanguageCode    string      `json:"languageCode"`
}

type GetBanksResponse struct {
	JSON            string      `json:"JSON"`
	AccountNickname interface{} `json:"accountNickname"`
	AccountNumber   interface{} `json:"accountNumber"`
	AccountType     interface{} `json:"accountType"`
	BankGroups      interface{} `json:"bankGroups"`
	BankToID        interface{} `json:"bankToId"`
	Banks           []Bank      `json:"banks"`
	DocumentNumber  interface{} `json:"documentNumber"`
	DocumentType    interface{} `json:"documentType"`
	DocumentTypes   interface{} `json:"documentTypes"`
	FiduciariaTypes interface{} `json:"fiduciariaTypes"`
	ProductName     interface{} `json:"productName"`
	ProductTypes1   interface{} `json:"productTypes1"`
	ProductTypes2   interface{} `json:"productTypes2"`
}

type Bank struct {
	Address         interface{} `json:"address"`
	BankIdentifier  interface{} `json:"bankIdentifier"`
	BankName        interface{} `json:"bankName"`
	BranchCode      interface{} `json:"branchCode"`
	BranchName      interface{} `json:"branchName"`
	CityID          interface{} `json:"cityId"`
	CityName        interface{} `json:"cityName"`
	ClgFlag         int64       `json:"clgFlag"`
	Code            string      `json:"code"`
	CodeType        interface{} `json:"codeType"`
	CountryID       interface{} `json:"countryId"`
	CountryName     interface{} `json:"countryName"`
	CreationDate    string      `json:"creationDate"`
	Credit          int64       `json:"credit"`
	Debit           int64       `json:"debit"`
	Description     string      `json:"description"`
	ExpirationDate  string      `json:"expirationDate"`
	IsHomeBank      int64       `json:"isHomeBank"`
	IsInternational int64       `json:"isInternational"`
	IsSubsidiary    int64       `json:"isSubsidiary"`
	LanguageCode    string      `json:"languageCode"`
	MICRBankCode    interface{} `json:"micrBankCode"`
	MICRBranchCode  interface{} `json:"micrBranchCode"`
	PaymentSystemID interface{} `json:"paymentSystemId"`
	SortName        interface{} `json:"sortName"`
	State           interface{} `json:"state"`
}
