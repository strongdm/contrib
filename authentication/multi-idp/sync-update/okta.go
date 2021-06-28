package main

import (
	"context"
	"net/url"
	"os"
	"strings"

	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
	"github.com/pkg/errors"
	"github.com/tomnomnom/linkheader"
)

const QUERY_STRING = "(status eq \"ACTIVE\")"
const API_LIMIT = 200

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

	var users []idpUser
	after := ""
	for {
		search := query.NewQueryParams(query.WithSearch(QUERY_STRING), query.WithLimit(API_LIMIT), query.WithAfter(after))
		apiUsers, resp, err := client.User.ListUsers(search)
		if err != nil {
			return nil, errors.Wrap(err, "unable to retrieve okta users")
		}

		for _, u := range apiUsers {
			groups, _, err := client.User.ListUserGroups(u.Id, nil)
			if err != nil {
				return nil, errors.Wrap(err, "unable to retrieve okta user groups")
			}

			var groupNames []string
			for _, g := range groups {
				groupNames = append(groupNames, g.Profile.Name)
			}

			profile := (*u.Profile)
			users = append(users, idpUser{
				Login:     profile["login"].(string),
				FirstName: profile["firstName"].(string),
				LastName:  profile["lastName"].(string),
				Groups:    groupNames,
			})
		}

		after, err = getQueryAfter(resp)
		if err != nil {
			return nil, errors.Wrap(err, "unable to parse after value from query")
		}
		if after == "" {
			break
		}
	}
	return users, nil
}

func getQueryAfter(resp *okta.Response) (string, error) {
	links := linkheader.Parse(strings.Join(resp.Header["Link"], ","))
	for _, link := range links {
		if link.Rel == "next" {
			u, err := url.Parse(link.URL)
			if err != nil {
				return "", err
			}
			m, _ := url.ParseQuery(u.RawQuery)
			return m["after"][0], nil
		}
	}
	return "", nil
}
