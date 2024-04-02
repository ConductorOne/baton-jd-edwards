package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-jd-edwards/pkg/jde"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	grant "github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type roleBuilder struct {
	resourceType *v2.ResourceType
	client       *jde.Client
}

const roleMembership = "member"

func (r *roleBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return r.resourceType
}

// Create a new connector resource for a JD Edwards role.
func roleResource(role, roleDescription string) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"role_id":          role,
		"role_description": roleDescription,
	}

	roleTraitOptions := []rs.RoleTraitOption{
		rs.WithRoleProfile(profile),
	}

	ret, err := rs.NewRoleResource(
		role,
		roleResourceType,
		role,
		roleTraitOptions,
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *roleBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	roles, err := r.client.ListRoles(ctx)
	if err != nil {
		return nil, "", nil, fmt.Errorf("error fetching roles: %w", err)
	}

	var rv []*v2.Resource
	for _, role := range roles {
		rr, err := roleResource(role.F00926User, role.F00926RoleDesc)
		if err != nil {
			return nil, "", nil, fmt.Errorf("error creating role resource: %w", err)
		}
		rv = append(rv, rr)
	}

	return rv, "", nil, nil
}

func (r *roleBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var rv []*v2.Entitlement

	assignmentOptions := []ent.EntitlementOption{
		ent.WithGrantableTo(userResourceType),
		ent.WithDisplayName(fmt.Sprintf("%s Role %s", resource.DisplayName, roleMembership)),
		ent.WithDescription(fmt.Sprintf("Member of %s role in JD Edwards EnterpriseOne", resource.DisplayName)),
	}

	rv = append(rv, ent.NewAssignmentEntitlement(
		resource,
		roleMembership,
		assignmentOptions...,
	))

	return rv, "", nil, nil
}

func (r *roleBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	users, err := r.client.ListRoleUsers(ctx, resource.Id.Resource)
	if err != nil {
		return nil, "", nil, fmt.Errorf("error fetching role users: %w", err)
	}

	var rv []*v2.Grant
	for _, user := range users {
		ur, err := userResource(user.F95921ToRole)
		if err != nil {
			return nil, "", nil, fmt.Errorf("error creating user resource for role %s: %w", resource.Id.Resource, err)
		}

		rv = append(rv, grant.NewGrant(
			resource,
			roleMembership,
			ur.Id,
		))
	}
	return rv, "", nil, nil
}

func newRoleBuilder(client *jde.Client) *roleBuilder {
	return &roleBuilder{
		resourceType: roleResourceType,
		client:       client,
	}
}
