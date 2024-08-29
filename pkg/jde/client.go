package jde

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

const (
	apiPathv1    = "/jderest"
	apiPathv2    = "/jderest/v2"
	dataservice  = "dataservice"
	tokenrequest = "tokenrequest"
	config       = "defaultconfig"
	validate     = "validate"

	// pageSizeNoMax is used to return all records from v1 since it does not support pagination.
	noMax = "No max"
	// default in v2, but needs to be set in v1.
	outputType = "GRID_DATA"
)

type Client struct {
	httpClient *uhttp.BaseHttpClient
	baseUrl    string
	token      string
	version    string
}

func NewClient(httpClient *http.Client, aisUrl, token, env, version string) (*Client, error) {
	path := getApiPath(version)
	baseUrl, _ := url.JoinPath(aisUrl, path)

	return &Client{
		httpClient: uhttp.NewBaseHttpClient(httpClient),
		baseUrl:    baseUrl,
		token:      token,
		version:    version,
	}, nil
}

type AuthRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DataRequestBody struct {
	Token                    string `json:"token,omitempty"`
	TargetName               string `json:"targetName,omitempty"`
	TargetType               string `json:"targetType,omitempty"`
	DataServiceType          string `json:"dataServiceType,omitempty"`
	MaxPageSize              string `json:"maxPageSize,omitempty"`
	ReturnControlIDs         string `json:"returnControlIDs,omitempty"`
	EnableNextPageProcessing string `json:"enableNextPageProcessing,omitempty"`
	FindOnEntry              string `json:"findOnEntry,omitempty"`
	Query                    *Query `json:"query,omitempty"`
	OutputType               string `json:"outputType,omitempty"`
}

type Query struct {
	AutoFind  bool        `json:"autoFind"`
	Condition []Condition `json:"condition"`
}

type Condition struct {
	ControlId string  `json:"controlId"`
	Operator  string  `json:"operator"`
	Value     []Value `json:"value"`
}

type Value struct {
	Content        string `json:"content"`
	SpecialValueID string `json:"specialValueId"`
}

type RequestParams struct {
	Url         string
	Method      string
	Res         interface{}
	Payload     interface{}
	QueryParams url.Values
}

type ValidateTokenBody struct {
	Token string `json:"token"`
}

// Authenticate authenticates the user with the JD Edwards EnterpriseOne AIS server and returns the token.
func Authenticate(ctx context.Context, ais, username, password, env, version string) (string, error) {
	authBody := AuthRequestBody{
		Username: username,
		Password: password,
	}

	path := getApiPath(version)
	url, err := url.JoinPath(ais, path, tokenrequest)
	if err != nil {
		return "", err
	}

	var res AuthResponse
	requestParams := RequestParams{
		Url:         url,
		Method:      http.MethodPost,
		Res:         &res,
		Payload:     authBody,
		QueryParams: nil,
	}

	err = doRequestDefaultClient(ctx, requestParams)
	if err != nil {
		return "", err
	}

	if res.UserInfo.Token == "" {
		return "", fmt.Errorf("token not obtained, check if provided credentials and environment are correct")
	}

	return res.UserInfo.Token, nil
}

// ListUsers returns all the users from the JD Edwards EnterpriseOne AIS server.
func (c *Client) ListUsers(ctx context.Context, pageSize string, enableNextPage bool) ([]Columns, string, error) {
	if c.version == "v1" && enableNextPage {
		pageSize = noMax
	}

	dataRequest := DataRequestBody{
		TargetName:               "F0092",
		TargetType:               "table",
		DataServiceType:          "BROWSE",
		FindOnEntry:              "true",
		ReturnControlIDs:         "F0092.USER|F0092.UGRP",
		MaxPageSize:              pageSize,
		EnableNextPageProcessing: fmt.Sprint(enableNextPage),
		OutputType:               outputType,
		Query: &Query{
			AutoFind: true,
			Condition: []Condition{
				{ControlId: "F0092.UGRP", Operator: "EQUAL", Value: []Value{
					{Content: "", SpecialValueID: "LITERAL"},
				}},
			}},
	}

	url, _ := url.JoinPath(c.baseUrl, dataservice)
	var res UsersResponse
	err := c.doRequest(ctx, http.MethodPost, url, dataRequest, &res)
	if err != nil {
		return nil, "", err
	}

	if res.Resource.Data.GridData.Summary.MoreRecords && c.version == "v2" {
		nextUrl := res.Links[0].Href
		return res.Resource.Data.GridData.Rowset, nextUrl, nil
	}

	return res.Resource.Data.GridData.Rowset, "", nil
}

func (c *Client) FetchMoreUsers(ctx context.Context, nextUrl string) ([]Columns, string, error) {
	var res UsersResponse
	err := c.doRequest(ctx, http.MethodPost, nextUrl, nil, &res)
	if err != nil {
		return nil, "", err
	}

	if res.Resource.Data.GridData.Summary.MoreRecords {
		nextUrl := res.Links[0].Href
		return res.Resource.Data.GridData.Rowset, nextUrl, nil
	}

	return res.Resource.Data.GridData.Rowset, "", nil
}

// ListRoles returns all the roles from the JD Edwards EnterpriseOne AIS server.
func (c *Client) ListRoles(ctx context.Context, pageSize string) ([]Columns, string, error) {
	if c.version == "v1" {
		pageSize = noMax
	}

	dataRequest := DataRequestBody{
		TargetName:               "F00926",
		TargetType:               "table",
		DataServiceType:          "BROWSE",
		FindOnEntry:              "true",
		ReturnControlIDs:         "F00926.USER|F00926.ROLEDESC",
		MaxPageSize:              pageSize,
		EnableNextPageProcessing: "true",
		OutputType:               outputType,
	}

	url, _ := url.JoinPath(c.baseUrl, dataservice)
	var res RolesResponse
	err := c.doRequest(ctx, http.MethodPost, url, dataRequest, &res)
	if err != nil {
		return nil, "", err
	}

	if res.Resource.Data.GridData.Summary.MoreRecords && c.version == "v2" {
		nextUrl := res.Links[0].Href
		return res.Resource.Data.GridData.Rowset, nextUrl, nil
	}

	return res.Resource.Data.GridData.Rowset, "", nil
}

func (c *Client) FetchMoreRoles(ctx context.Context, nextUrl string) ([]Columns, string, error) {
	var res RolesResponse
	err := c.doRequest(ctx, http.MethodPost, nextUrl, nil, &res)
	if err != nil {
		return nil, "", err
	}

	if res.Resource.Data.GridData.Summary.MoreRecords {
		nextUrl := res.Links[0].Href
		return res.Resource.Data.GridData.Rowset, nextUrl, nil
	}

	return res.Resource.Data.GridData.Rowset, "", nil
}

// ListRoleUsers returns all the users that are assigned to a role from the JD Edwards EnterpriseOne AIS server.
func (c *Client) ListRoleUsers(ctx context.Context, roleID string, pageSize string) ([]Columns, string, error) {
	if c.version == "v1" {
		pageSize = noMax
	}

	dataRequest := DataRequestBody{
		TargetName:               "F95921",
		TargetType:               "table",
		DataServiceType:          "BROWSE",
		FindOnEntry:              "true",
		ReturnControlIDs:         "F95921.FRROLE|F95921.TOROLE",
		MaxPageSize:              pageSize,
		EnableNextPageProcessing: "true",
		OutputType:               outputType,
		Query: &Query{
			AutoFind: true,
			Condition: []Condition{
				{ControlId: "F95921.FRROLE", Operator: "EQUAL", Value: []Value{
					{Content: roleID, SpecialValueID: "LITERAL"},
				}},
			}},
	}

	url, _ := url.JoinPath(c.baseUrl, dataservice)
	var res RoleUsersResponse
	err := c.doRequest(ctx, http.MethodPost, url, dataRequest, &res)
	if err != nil {
		return nil, "", err
	}

	if res.Resource.Data.GridData.Summary.MoreRecords && c.version == "v2" {
		nextUrl := res.Links[0].Href
		return res.Resource.Data.GridData.Rowset, nextUrl, nil
	}

	return res.Resource.Data.GridData.Rowset, "", nil
}

func (c *Client) FetchMoreRoleUsers(ctx context.Context, nextUrl string) ([]Columns, string, error) {
	var res RoleUsersResponse
	err := c.doRequest(ctx, http.MethodPost, nextUrl, nil, &res)
	if err != nil {
		return nil, "", err
	}

	if res.Resource.Data.GridData.Summary.MoreRecords {
		nextUrl := res.Links[0].Href
		return res.Resource.Data.GridData.Rowset, nextUrl, nil
	}

	return res.Resource.Data.GridData.Rowset, "", nil
}

// ValidateToken validates the current session token.
func (c *Client) ValidateToken(ctx context.Context) (ValidateTokenResponse, error) {
	url, _ := url.JoinPath(c.baseUrl, tokenrequest, validate)
	body := ValidateTokenBody{
		Token: c.token,
	}

	var res ValidateTokenResponse
	err := c.doRequest(ctx, http.MethodPost, url, body, &res)
	if err != nil {
		return ValidateTokenResponse{}, err
	}

	return res, nil
}

func GetConfigv1(ctx context.Context, ais string) (ConfigResponse, bool, error) {
	params := url.Values{}
	params.Add("requiredCapabilities", "dataservice,outputType")
	var res ConfigResponse

	url, err := url.JoinPath(ais, apiPathv1, config)
	if err != nil {
		return ConfigResponse{}, false, err
	}

	requestParams := RequestParams{
		Url:         url,
		Method:      http.MethodPost,
		Res:         &res,
		Payload:     nil,
		QueryParams: params,
	}
	if err := doRequestDefaultClient(ctx, requestParams); err != nil {
		return ConfigResponse{}, false, err
	}

	if res.RequiredCapabilityMissing {
		return ConfigResponse{}, true, nil
	}

	return res, false, nil
}

func GetConfigv2(ctx context.Context, ais string) (ConfigResponse, bool, error) {
	params := url.Values{}
	params.Add("all", "true")
	params.Add("requiredCapabilities", "dataservice,outputType")

	url, err := url.JoinPath(ais, apiPathv2, config)
	if err != nil {
		return ConfigResponse{}, false, err
	}

	var res ConfigResponse
	requestParams := RequestParams{
		Url:         url,
		Method:      http.MethodPost,
		Res:         &res,
		Payload:     nil,
		QueryParams: params,
	}
	if err := doRequestDefaultClient(ctx, requestParams); err != nil {
		return ConfigResponse{}, false, err
	}

	if res.RequiredCapabilityMissing {
		return ConfigResponse{}, true, nil
	}

	return res, false, nil
}

func GetConfig(ctx context.Context, ais string) (ConfigResponse, bool, string, error) {
	var config ConfigResponse
	var capabilityMissing bool
	var err error
	version := "v2"

	// try to fetch config from v2, if it fails, try v1
	config, capabilityMissing, err = GetConfigv2(ctx, ais)
	if err != nil {
		config, capabilityMissing, err = GetConfigv1(ctx, ais)
		if err != nil {
			return ConfigResponse{}, capabilityMissing, "", err
		}
		version = "v1"
	}

	return config, capabilityMissing, version, nil
}

// ValidateTokenV1 validates the current session token.
func (c *Client) ValidateTokenV1(ctx context.Context) error {
	_, _, err := c.ListUsers(ctx, "1", false)
	if err != nil {
		return fmt.Errorf("error validating token while fetching user: %w", err)
	}

	return nil
}

func (c *Client) doRequest(ctx context.Context, method string, reqUrl string, payload interface{}, res interface{}) error {
	u, err := url.Parse(reqUrl)
	if err != nil {
		return err
	}

	req, err := c.httpClient.NewRequest(
		ctx,
		method,
		u,
		uhttp.WithJSONBody(payload),
		uhttp.WithAcceptJSONHeader(),
		uhttp.WithContentTypeJSONHeader(),
		uhttp.WithHeader("jde-AIS-Auth", c.token),
	)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req, uhttp.WithJSONResponse(res))
	if err != nil {
		return err
	}

	resp.Body.Close()

	return nil
}

func doRequestDefaultClient(ctx context.Context, requestParams RequestParams) error {
	httpClient, err := (&uhttp.NoAuth{}).GetClient(ctx)
	if err != nil {
		return err
	}

	client := uhttp.NewBaseHttpClient(httpClient)
	u, err := url.Parse(requestParams.Url)
	if err != nil {
		return err
	}

	req, err := client.NewRequest(
		ctx,
		requestParams.Method,
		u,
		uhttp.WithJSONBody(requestParams.Payload),
		uhttp.WithAcceptJSONHeader(),
		uhttp.WithContentTypeJSONHeader(),
	)
	if err != nil {
		return err
	}

	if requestParams.QueryParams != nil {
		req.URL.RawQuery = requestParams.QueryParams.Encode()
	}

	resp, err := client.Do(req, uhttp.WithJSONResponse(requestParams.Res))
	if err != nil {
		return err
	}

	resp.Body.Close()

	return nil
}

func getApiPath(version string) string {
	path := apiPathv2
	if version == "v1" {
		path = apiPathv1
	}
	return path
}

func CanValidateToken(ctx context.Context, config ConfigResponse) bool {
	for _, capability := range config.CapabilityList {
		if capability.Name == validate {
			return true
		}
	}
	return false
}
