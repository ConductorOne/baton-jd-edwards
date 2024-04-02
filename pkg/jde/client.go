package jde

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

const (
	apiPath      = "/jderest/v2"
	dataservice  = "dataservice"
	tokenrequest = "tokenrequest"

	applicationJson = "application/json"
)

type Client struct {
	httpClient *http.Client
	baseUrl    string
	token      string
}

func NewClient(httpClient *http.Client, aisUrl, token string) *Client {
	baseUrl, _ := url.JoinPath(aisUrl, apiPath)
	return &Client{
		httpClient: httpClient,
		baseUrl:    baseUrl,
		token:      token,
	}
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

// Authenticate authenticates the user with the JD Edwards EnterpriseOne AIS server and returns the token.
func Authenticate(ctx context.Context, ais, username, password string) (string, error) {
	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	if err != nil {
		return "", err
	}

	url, err := url.JoinPath(ais, apiPath, tokenrequest)
	if err != nil {
		return "", err
	}

	authBody := AuthRequestBody{
		Username: username,
		Password: password,
	}

	requestBody, err := json.Marshal(authBody)
	if err != nil {
		return "", fmt.Errorf("error marshalling auth request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Add("Accept", applicationJson)
	req.Header.Add("Content-Type", applicationJson)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var res AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	if res.UserInfo.Token == "" {
		return "", fmt.Errorf("token not obtained, check if provided credentials and environment are correct")
	}

	return res.UserInfo.Token, nil
}

// ListUsers returns all the users from the JD Edwards EnterpriseOne AIS server.
func (c *Client) ListUsers(ctx context.Context) ([]Columns, error) {
	url, _ := url.JoinPath(c.baseUrl, dataservice)

	dataRequest := DataRequestBody{
		TargetName:       "F0092",
		TargetType:       "table",
		DataServiceType:  "BROWSE",
		FindOnEntry:      "true",
		ReturnControlIDs: "F0092.USER|F0092.UGRP",
		Query: &Query{
			AutoFind: true,
			Condition: []Condition{
				{ControlId: "F0092.UGRP", Operator: "EQUAL", Value: []Value{
					{Content: "", SpecialValueID: "LITERAL"},
				}},
			}},
	}

	var res UsersResponse
	if err := c.doRequest(ctx, http.MethodPost, url, &res, dataRequest); err != nil {
		return nil, err
	}

	return res.Resource.Data.GridData.Rowset, nil
}

// ListRoles returns all the roles from the JD Edwards EnterpriseOne AIS server.
func (c *Client) ListRoles(ctx context.Context) ([]Columns, error) {
	url, _ := url.JoinPath(c.baseUrl, dataservice)
	dataRequest := DataRequestBody{
		TargetName:       "F00926",
		TargetType:       "table",
		DataServiceType:  "BROWSE",
		FindOnEntry:      "true",
		ReturnControlIDs: "F00926.USER|F00926.ROLEDESC",
	}

	var res RolesResponse
	if err := c.doRequest(ctx, http.MethodPost, url, &res, dataRequest); err != nil {
		return nil, err
	}

	return res.Resource.Data.GridData.Rowset, nil
}

// ListRoleUsers returns all the users that are assigned to a role from the JD Edwards EnterpriseOne AIS server.
func (c *Client) ListRoleUsers(ctx context.Context, roleID string) ([]Columns, error) {
	url, _ := url.JoinPath(c.baseUrl, dataservice)
	dataRequest := DataRequestBody{
		TargetName:       "F95921",
		TargetType:       "table",
		DataServiceType:  "BROWSE",
		FindOnEntry:      "true",
		ReturnControlIDs: "F95921.FRROLE|F95921.TOROLE",
		Query: &Query{
			AutoFind: true,
			Condition: []Condition{
				{ControlId: "F95921.FRROLE", Operator: "EQUAL", Value: []Value{
					{Content: roleID, SpecialValueID: "LITERAL"},
				}},
			}},
	}

	var res RoleUsersResponse
	if err := c.doRequest(ctx, http.MethodPost, url, &res, dataRequest); err != nil {
		return nil, err
	}

	return res.Resource.Data.GridData.Rowset, nil
}

// ValidateToken validates the current session token.
func (c *Client) ValidateToken(ctx context.Context) (ValidateTokenResponse, error) {
	url, _ := url.JoinPath(c.baseUrl, tokenrequest, "validate")

	var res ValidateTokenResponse
	if err := c.doRequest(ctx, http.MethodPost, url, &res, nil); err != nil {
		return ValidateTokenResponse{}, err
	}

	return res, nil
}

func (c *Client) doRequest(ctx context.Context, method string, url string, res interface{}, requestBody interface{}) error {
	var payload []byte
	var err error

	if requestBody != nil {
		payload, err = json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("error marshalling request body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Accept", applicationJson)
	req.Header.Add("Content-Type", applicationJson)
	req.Header.Add("jde-AIS-Auth", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	return nil
}
