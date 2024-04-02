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
}

type RolesResponse struct {
	Resource Resource `json:"fs_DATABROWSE_F00926"`
}

type RoleUsersResponse struct {
	Resource Resource `json:"fs_DATABROWSE_F95921"`
}

type ValidateTokenResponse struct {
	IsValidSession bool  `json:"isValidSession,omitempty"`
	Error          Error `json:"error,omitempty"`
}

type Error struct {
	Message   string `json:"message"`
	Exception string `json:"exception"`
	TimeStamp string `json:"timeStamp"`
}
