package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-jd-edwards/pkg/jde"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type userBuilder struct {
	resourceType *v2.ResourceType
	client       *jde.Client
}

func (u *userBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return u.resourceType
}

// Create a new connector resource for a JD Edwards user.
func userResource(user string) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"user_id": user,
	}

	userTraitOptions := []rs.UserTraitOption{
		rs.WithUserProfile(profile),
		rs.WithStatus(v2.UserTrait_Status_STATUS_ENABLED),
	}

	ret, err := rs.NewUserResource(
		user,
		userResourceType,
		user,
		userTraitOptions,
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (u *userBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	bag, page, isInitial, err := parsePageToken(pToken.Token, &v2.ResourceId{ResourceType: userResourceType.Id})
	if err != nil {
		return nil, "", nil, err
	}

	var allUsers []jde.Columns
	var nextToken string
	if page == "" && isInitial {
		users, nextUrl, err := u.client.ListUsers(ctx, "100", true)
		if err != nil {
			return nil, "", nil, fmt.Errorf("error fetching users: %w", err)
		}
		allUsers = append(allUsers, users...)

		nextToken, err = bag.NextToken(nextUrl)
		if err != nil {
			return nil, "", nil, err
		}
	} else {
		users, nextUrl, err := u.client.FetchMoreUsers(ctx, page)
		if err != nil {
			return nil, "", nil, fmt.Errorf("error fetching users: %w", err)
		}
		allUsers = append(allUsers, users...)

		nextToken, err = bag.NextToken(nextUrl)
		if err != nil {
			return nil, "", nil, err
		}
	}

	var rv []*v2.Resource
	for _, user := range allUsers {
		ur, err := userResource(user.F0092User)
		if err != nil {
			return nil, "", nil, fmt.Errorf("error creating user resource: %w", err)
		}
		rv = append(rv, ur)
	}

	return rv, nextToken, nil, nil
}

// Entitlements always returns an empty slice for users.
func (u *userBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (u *userBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newUserBuilder(client *jde.Client) *userBuilder {
	return &userBuilder{
		resourceType: userResourceType,
		client:       client,
	}
}
