package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-jd-edwards/pkg/jde"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

type Connector struct {
	client  *jde.Client
	version string
	aisUrl  string
}

// ResourceSyncers returns a ResourceSyncer for each resource type that should be synced from the upstream service.
func (d *Connector) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		newUserBuilder(d.client),
		newRoleBuilder(d.client),
	}
}

// Metadata returns metadata about the connector.
func (d *Connector) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "JD Edwards Connector",
		Description: "Connector syncing users and roles from JD Edwards EnterpriseOne.",
	}, nil
}

// Validate is called to ensure that the connector is properly configured. It should exercise any API credentials
// to be sure that they are valid.
func (d *Connector) Validate(ctx context.Context) (annotations.Annotations, error) {
	// check if all capabilities are configured
	config, capabilitiesMissing, version, err := jde.GetConfig(ctx, d.aisUrl)
	if err != nil {
		return nil, fmt.Errorf("error fetching config: %w", err)
	}

	if capabilitiesMissing {
		return nil, fmt.Errorf("capabilities missing, make sure dataservice and tokenrequest capabilities are configured on the AIS server")
	}

	validateTokenConfigured := jde.CanValidateToken(ctx, config)

	// if we are using v1 or don't have validate token configured, we need to validate token differently.
	if version == "v1" || !validateTokenConfigured {
		err = d.client.ValidateTokenV1(ctx)
		if err != nil {
			return nil, fmt.Errorf("error validating token: %w", err)
		}
		return nil, nil
	}

	res, err := d.client.ValidateToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("error validating token: %w", err)
	}

	if res.Message != "" {
		return nil, fmt.Errorf("error validating token: %s", res.Message)
	}

	return nil, nil
}

// New returns a new instance of the connector.
func New(ctx context.Context, aisUrl, username, password, env string) (*Connector, error) {
	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	if err != nil {
		return nil, err
	}

	// call config to see which AIS version we are using
	_, _, version, err := jde.GetConfig(ctx, aisUrl)
	if err != nil {
		return nil, fmt.Errorf("error fetching config: %w", err)
	}

	token, err := jde.Authenticate(ctx, aisUrl, username, password, env, version)
	if err != nil {
		return nil, fmt.Errorf("error authenticating: %w", err)
	}

	client, err := jde.NewClient(httpClient, aisUrl, token, env, version)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	return &Connector{
		client:  client,
		version: version,
		aisUrl:  aisUrl,
	}, nil
}
