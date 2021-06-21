package main

import (
	"context"
	"os"

	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
	"github.com/pkg/errors"
)

const QUERY_STRING = "(status eq \"ACTIVE\")"

func ValidateOktaEnv() error {
	if os.Getenv("OKTA_CLIENT_TOKEN") == "" || os.Getenv("OKTA_CLIENT_ORGURL") == "" {
		return errors.Errorf("OKTA_CLIENT_TOKEN and OKTA_CLIENT_ORGURL must be set when using Okta")
	}
	return nil
}

func LoadOktaUsers(ctx context.Context) (idpUserList, error) {
	client, err := okta.NewClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "invalid Okta configuration")
	}
	search := query.NewQueryParams(query.WithSearch(QUERY_STRING), query.WithLimit(USERS_LIMIT))

	apiUsers, _, err := client.User.ListUsers(search)
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve okta users")
	}

	var users []idpUser
	for _, u := range apiUsers {
		profile := (*u.Profile)

		groups, _, err := client.User.ListUserGroups(u.Id, nil)
		if err != nil {
			return nil, errors.Wrap(err, "unable to retrieve okta user groups")
		}

		var groupNames []string
		for _, g := range groups {
			groupNames = append(groupNames, g.Profile.Name)
		}

		var u idpUser
		u.Login = profile["login"].(string)
		u.FirstName = profile["firstName"].(string)
		u.LastName = profile["lastName"].(string)
		u.Groups = groupNames
		users = append(users, u)
	}
	return users, nil
}
