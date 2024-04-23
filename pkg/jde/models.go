package jde

type AuthResponse struct {
	Username       string   `json:"username"`
	Environment    string   `json:"environment"`
	Role           string   `json:"role"`
	Jasserver      string   `json:"jasserver"`
	UserInfo       UserInfo `json:"userInfo"`
	UserAuthorized bool     `json:"userAuthorized"`
}

type UserInfo struct {
	Token       string `json:"token"`
	AlphaName   string `json:"alphaName"`
	AppsRelease string `json:"appsRelease"`
	Username    string `json:"username"`
}

type Resource struct {
	Title string `json:"title"`
	Data  Data   `json:"data"`
}

type Data struct {
	GridData GridData `json:"gridData"`
}

type GridData struct {
	Columns Columns   `json:"columns"`
	Rowset  []Columns `json:"rowset"`
	Summary Summary   `json:"summary"`
}

type Columns struct {
	F0092User      string `json:"F0092_USER,omitempty"`
	F0092Ugrp      string `json:"F0092_UGRP"`
	F95921FrRole   string `json:"F95921_FRROLE,omitempty"`
	F95921ToRole   string `json:"F95921_TOROLE,omitempty"`
	F00926User     string `json:"F00926_USER,omitempty"`
	F00926RoleDesc string `json:"F00926_ROLEDESC,omitempty"`
}

type Summary struct {
	Records     int  `json:"records"`
	MoreRecords bool `json:"moreRecords"`
}

type UsersResponse struct {
	Resource Resource `json:"fs_DATABROWSE_F0092"`
	Links    []Link   `json:"links,omitempty"`
}

type RolesResponse struct {
	Resource Resource `json:"fs_DATABROWSE_F00926"`
	Links    []Link   `json:"links,omitempty"`
}

type RoleUsersResponse struct {
	Resource Resource `json:"fs_DATABROWSE_F95921"`
	Links    []Link   `json:"links,omitempty"`
}

type ValidateTokenResponse struct {
	IsValidSession bool   `json:"isValidSession,omitempty"`
	Message        string `json:"message,omitempty"`
	Exception      string `json:"exception,omitempty"`
	TimeStamp      string `json:"timeStamp,omitempty"`
}

type ConfigResponse struct {
	JasHost                   string           `json:"jasHost,omitempty"`
	JasPort                   string           `json:"jasPort,omitempty"`
	JasProtocol               string           `json:"jasProtocol,omitempty"`
	DefaultEnvironment        string           `json:"defaultEnvironment,omitempty"`
	DefaultRole               string           `json:"defaultRole,omitempty"`
	AisVersion                string           `json:"aisVersion,omitempty"`
	CapabilityList            []CapabilityList `json:"capabilityList,omitempty"`
	RequiredCapabilityMissing bool             `json:"requiredCapabilityMissing,omitempty"`
}

type CapabilityList struct {
	Name             string `json:"name,omitempty"`
	ShortDescription string `json:"shortDescription,omitempty"`
	LongDescription  string `json:"longDescription,omitempty"`
	AsOfRelease      string `json:"asOfRelease,omitempty"`
	SinceVersion     string `json:"sinceVersion,omitempty"`
}

type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}
