package main

type OrdsItem struct {
	StatementId   int    `json:"statementId"`
	StatementType string `json:"statementType"`
	StatementPos  struct {
		StartLine int `json:"startLine"`
		EndLine   int `json:"endLine"`
	} `json:"statementPos"`
	StatementText string `json:"statementText"`
	ErrorCode     int    `json:"errorCode"`
	ErrorLine     int    `json:"errorLine"`
	ErrorColumn   int    `json:"errorColumn"`
	ErrorDetails  string `json:"errorDetails"`
}

type OrdsResponse struct {
	Env struct {
		DefaultTimeZone string `json:"defaultTimeZone"`
	} `json:"env"`
	Items []OrdsItem `json:"items"`
}
