// Copyright 2020 StrongDM Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

/*
PLEASE NOTE: this is sample code intended to demonstrate strongDM SDK functionality.
You should review and test thoroughly before deploying to production.
This code is provided AS-IS and may or may not be updated in the future at our discretion.
Our Support team will be happy to assist with general SDK questions or issues.
*/

/* This script is an update to our standard example located here:
		 https://github.com/strongdm/strongdm-sdk-go-examples/tree/master/contrib/okta-sync

It is written against Okta golang API version 1.x.

It partially implements user/Group sync from Okta > SDM, and switches access grants from the user to the Role (Group) level.

NB: you must set the following environment variables: SDM_API_ACCESS_KEY, SDM_API_SECRET_KEY, OKTA_CLIENT_TOKEN, and OKTA_CLIENT_ORGURL.

This script reads a separate JSON file, matchers.yml, which maps Okta groups to resources in SDM by type or name.
For each Group defined in the YML, an SDM Role will be created, and access to the defined resources will be granted to that Role.
Any user that matches the Okta search filter will be created in SDM (see the "oktaQueryString" definition just below).
If they belong to an Okta Group that is defined in the YML, they will be assigned to the corresponding SDM Role.

The script won't remove any Roles or Users in SDM.
However, it will remove any grants for Groups/Roles that are not defined in the YML.
It will also add/remove grants for Groups/Roles if you change the mapping in the YML.

An important consideration is that Okta supports multiple group assignment, but strongDM does not.
This means that a user with multiple Group memberships will be assigned to the first Group/Role provided by Okta.
We recommend that you consider creating SDM-specific Groups in Okta, e.g. sdm-qa, sdm-dev, and assign users in Okta accordingly.
Then define only these groups in the YML, with appropriate resource mapping.
You may wish to modify the oktaQueryString to match only users who belong to the SDM-specific Groups.
*/

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
	"github.com/pkg/errors"
	sdm "github.com/strongdm/strongdm-sdk-go"
	"gopkg.in/yaml.v2"
)

// modify this Okta query filter with any valid Okta API parameters to control user creation in SDM
var oktaQueryString = "profile.login eq \"david+test@strongdm.com\" and (status eq \"ACTIVE\")"

var verbose = flag.Bool("json", false, "dump a JSON report for debugging")
var plan = flag.Bool("plan", false, "do not apply changes just plan and output the result")

// global Groups list, used to store Okta groups, and create missing Roles
var oktaGroups []string

var debug = false

func init() {
	flag.Parse()
}

type syncReport struct {
	Start          time.Time  `json:"start"`
	Complete       time.Time  `json:"complete"`
	OktaUserCount  int        `json:"oktaUsersCount"`
	OktaUsers      []oktaUser `json:"oktaUsers"`
	OktaGroupCount int        `json:"oktaGroupCount"`
	OktaGroups     []string   `json:"oktaGroups"`

	SDMUserCount int       `json:"sdmUsersCount"`
	SDMUsers     []userRow `json:"sdmUsers"`
	SDMRoleCount int       `json:"sdmRoleCount"`
	SDMRoles     []string  `json:"sdmRoles"`

	BothUserCount int `json:"bothUsersCount"`

	SDMResourcesCount int           `json:"sdmResourcesCount"`
	SDMResources      []entitlement `json:"sdmResources"`

	PermissionsGranted int                  `json:"permissionsGranted"`
	PermissionsRevoked int                  `json:"permissionsRevoked"`
	Grants             []entitlement        `json:"grants"`
	Revocations        []rolePermissionsRow `json:"revocations"`

	Matchers *MatcherConfig `json:"matchers"`
}

func (rpt *syncReport) String() string {
	if !*verbose {
		return rpt.short()
	}

	out, err := json.MarshalIndent(rpt, "", "\t")
	if err != nil {
		return fmt.Sprintf("error building JSON report: %s\n\n%s", err, rpt.short())
	}
	return string(out)
}

func (rpt *syncReport) short() string {
	return fmt.Sprintf("%d Okta users, %d strongDM users, %d matching users, %d grants, %d revocations, %d Groups, %d Roles\n",
		rpt.OktaUserCount, rpt.SDMUserCount, rpt.BothUserCount,
		rpt.PermissionsGranted, rpt.PermissionsRevoked, rpt.OktaGroupCount, rpt.SDMRoleCount)
}

func main() {
	ctx := context.Background()

	if os.Getenv("SDM_API_ACCESS_KEY") == "" ||
		os.Getenv("SDM_API_SECRET_KEY") == "" ||
		os.Getenv("OKTA_CLIENT_TOKEN") == "" ||
		os.Getenv("OKTA_CLIENT_ORGURL") == "" {
		fmt.Println("SDM_API_ACCESS_KEY, SDM_API_SECRET_KEY, OKTA_CLIENT_TOKEN, and OKTA_CLIENT_ORGURL must be set")
		os.Exit(1)
		return
	}

	client, err := sdm.New(os.Getenv("SDM_API_ACCESS_KEY"), os.Getenv("SDM_API_SECRET_KEY"))
	if err != nil {
		fmt.Println("failed to initialize strongDM client: ", err)
		os.Exit(1)
		return
	}

	var rpt syncReport
	rpt.Start = time.Now()

	matchers, err := loadMatchers()
	if err != nil {
		fmt.Printf("error loading Matchers users: %v\n", err)
		os.Exit(1)
		return
	}
	rpt.Matchers = matchers

	// get all Okta users that match filter defined in loadOktaUsers()
	// this also populates oktaGroups
	oktaUsers, err := loadOktaUsers(ctx)
	if err != nil {
		fmt.Printf("error loading Okta users: %v\n", err)
		os.Exit(1)
		return
	}
	rpt.OktaUsers = oktaUsers
	rpt.OktaUserCount = len(oktaUsers)
	rpt.OktaGroupCount = len(oktaGroups)

	// determine set of datasources and servers they should have access to by group
	entitlements, err := matchEntitlements(ctx, client, matchers)
	if err != nil {
		fmt.Printf("error matching entitlements: %v\n", err)
		os.Exit(1)
		return
	}

	// for each defined entitlement, use the Okta Group name
	// to find/create a corresponding Role in SDM
	for a := range entitlements {
		if debug {
			println("range entitlements: ", a)
		}
		rpt.SDMRoleCount++
		myRole := &sdm.Role{Name: a}
		_, err = client.Roles().Create(ctx, myRole)
		if err != nil && !(strings.Contains(err.Error(), "item already exists")) {
			fmt.Printf("error creating Role: %v\n", err)
		}
	}

	rolePermissions, err := loadRoleGrants(ctx, client)
	if err != nil {
		fmt.Printf("error loading permissions: %v\n", err)
		os.Exit(1)
		return
	}

	users, err := loadAccounts(ctx, client)
	if err != nil {
		fmt.Printf("error loading users: %v\n", err)
		os.Exit(1)
		return
	}
	rpt.SDMUsers = users
	rpt.SDMUserCount = len(users)

	resources, err := loadResources(ctx, client)
	if err != nil {
		fmt.Printf("error loading datasources: %v\n", err)
		os.Exit(1)
		return
	}
	rpt.SDMResources = resources
	rpt.SDMResourcesCount = len(resources)

	// use list of Okta users & compare to SDM, create missing users in SDM
	// NB: Okta allows blank First/Last names, SDM does not
	// If desired, this could be modified to only create SDM users if their Okta Group matches one defined in matchers.yml
	for _, oktaUser := range oktaUsers {
		if findUser(users, oktaUser) {
			break
		} else {
			println("Okta user not found in SDM. Attempting user creation ...")
			_, err := client.Accounts().Create(ctx, &sdm.User{
				FirstName: oktaUser.FirstName,
				LastName:  oktaUser.LastName,
				Email:     oktaUser.Login,
			})
			if err != nil {
				log.Fatal("Error while creating user: ", err)
			} else {
				fmt.Println("User creation successful!!!")
				// reload list of SDM Accounts, for group assignment later
				users, _ = loadAccounts(ctx, client)
			}
		}
	}

	bothCount := make(map[string]bool)
	for _, oktaUser := range oktaUsers {
		if debug {
			fmt.Println("oktauser: ", oktaUser)
		}
		for _, sdmUser := range users {
			if debug {
				fmt.Println("sdmUser: ", sdmUser)
			}
			if strings.ToLower(sdmUser.Email) == strings.ToLower(oktaUser.Login) {
				if debug {
					println("User matches!")
				}
				bothCount[sdmUser.Email] = true
				for _, g := range oktaUser.Groups {
					if debug {
						fmt.Println("group: ", g)
					}

					resp, err := client.Roles().List(ctx, "name:\""+g+"\"")
					if err != nil {
						fmt.Println("error finding user Role: ", err)
					}
					for resp.Next() {
						role := resp.Value()
						if debug {
							fmt.Println("found role: ", role)
						}
						attachment := &sdm.AccountAttachment{
							AccountID: sdmUser.ID,
							RoleID:    role.ID,
						}
						attachmentResponse, err := client.AccountAttachments().Create(ctx, attachment)
						if err == nil {
							attachmentID := attachmentResponse.AccountAttachment.ID
							log.Printf("Successfully created account attachment: ID: %v\n", attachmentID)
						} else if !strings.Contains(err.Error(), "item already exists") {
							fmt.Println("error finding user Role: ", err)
						}
					}
				}
			}
		}
	}
	rpt.BothUserCount = len(bothCount)

	matchingByRole := make(map[roleRow]entitlementList)
	if debug {
		fmt.Println(matchingByRole)
	}

	for _, group := range oktaGroups {
		uniq := make(map[entitlement]bool)
		for _, e := range entitlements[group] {
			uniq[e] = true
		}

		for e := range uniq {
			if debug {
				fmt.Println("DEBUG: name:" + group)
			}
			resp, err := client.Roles().List(ctx, "name:\""+group+"\"")
			if err != nil {
				fmt.Println("error finding Role", err)
			}
			for resp.Next() {
				role := resp.Value()
				newRoleRow := roleRow{
					ID:   role.ID,
					Name: role.Name,
				}

				matchingByRole[newRoleRow] = append(matchingByRole[newRoleRow], e)
				if debug {
					fmt.Println(matchingByRole)
				}
				break
			}
		}
	}

	toGrant := []rolePermissionsRow{}
	toRevoke := []rolePermissionsRow{}
	for r, entitlements := range matchingByRole {
		// are there any entitlements not permitted? grant.
		for _, e := range entitlements {
			found := false
			for _, p := range rolePermissions {
				if p.RoleID == r.ID && p.DatasourceID == e.DatasourceID {
					found = true
				}
			}
			if !found {
				if !*plan {
					toGrant = append(toGrant, rolePermissionsRow{RoleID: r.ID, DatasourceID: e.DatasourceID})
					rpt.PermissionsGranted++
				} else {
					fmt.Printf("Plan: grant %v to user %v\n", e.DatasourceID, r.ID)
				}
			}
		}
	}

	rpt.Grants = []entitlement{}
	for _, g := range toGrant {
		rpt.Grants = append(rpt.Grants, entitlement{DatasourceID: g.DatasourceID})
	}

	// are there any permissions not entitled? revoke.
	for _, p := range rolePermissions {
		found := false
		for r, entitlements := range matchingByRole {
			if p.RoleID == r.ID {
				for _, e := range entitlements {
					if p.RoleID == r.ID && e.DatasourceID == p.DatasourceID {
						found = true
					}
				}
			}
		}
		if !found {
			if !*plan {
				toRevoke = append(toRevoke, p)
				rpt.PermissionsRevoked++
			} else {
				fmt.Printf("Plan: revoke %s from user %s\n", p.DatasourceID, p.RoleID)
			}
		}
	}

	rpt.Revocations = toRevoke

	if !*plan {
		for _, grant := range toGrant {
			_, err := client.RoleGrants().Create(ctx, &sdm.RoleGrant{
				RoleID:     grant.RoleID,
				ResourceID: grant.DatasourceID,
			})
			var alreadyExistsErr *sdm.AlreadyExistsError
			if err != nil && !errors.As(err, &alreadyExistsErr) {
				fmt.Println("error granting: ", err)
			}
		}
		for _, grant := range toRevoke {
			_, err := client.RoleGrants().Delete(ctx, grant.ID)
			var notFoundError *sdm.NotFoundError
			if err != nil && !errors.As(err, &notFoundError) {
				fmt.Println("error revoking", err)
			}
		}
	}

	rpt.Complete = time.Now()
	fmt.Println(rpt.String())
}

type oktaUser struct {
	Login     string   `json:"login"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Groups    []string `json:"groups"`
}

type userList []userRow

type userRow struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Role      string `json:"roleName"`
}

type roleList []roleRow

type roleRow struct {
	ID   string `json:"id"`
	Name string `json:"Name"`
}

type permissionsList []permissionsRow
type rolePermissionsList []rolePermissionsRow

type permissionsRow struct {
	ID           string `json:"-"`
	UserID       string `json:"userID"`
	DatasourceID string `json:"datasourceID"`
}

type rolePermissionsRow struct {
	ID           string `json:"-"`
	RoleID       string `json:"userID"`
	DatasourceID string `json:"datasourceID"`
}

type entitlementList []entitlement

type entitlement struct {
	DatasourceID string `json:"id"`
	Name         string `json:"name"`
}

func loadOktaUsers(ctx context.Context) ([]oktaUser, error) {
	client, err := okta.NewClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "invalid Okta configuration")
	}
	search := query.NewQueryParams(query.WithSearch(oktaQueryString))

	apiUsers, _, err := client.User.ListUsers(search)
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve okta users")
	}

	var users []oktaUser
	for _, u := range apiUsers {
		login := (*u.Profile)["login"].(string)
		firstName := (*u.Profile)["firstName"].(string)
		lastName := (*u.Profile)["lastName"].(string)

		groups, _, err := client.User.ListUserGroups(u.Id, search)

		if err != nil {
			return nil, errors.Wrap(err, "unable to retrieve okta user groups")
		}

		var groupNames []string
		for _, g := range groups {
			if debug {
				println("loadOktausers: ", login, g.Profile.Name)
			}
			groupNames = append(groupNames, g.Profile.Name)
			oktaGroups = AppendIfMissing(oktaGroups, g.Profile.Name)
		}

		var u oktaUser
		u.Login = login
		u.FirstName = firstName
		u.LastName = lastName
		u.Groups = groupNames
		users = append(users, u)
	}
	return users, nil
}

func loadAccountGrants(ctx context.Context, client *sdm.Client) ([]permissionsRow, error) {
	grants, err := client.AccountGrants().List(ctx, "")
	if err != nil {
		return nil, err
	}
	var result permissionsList
	for grants.Next() {
		grant := grants.Value()
		result = append(result, permissionsRow{
			ID:           grant.ID,
			UserID:       grant.AccountID,
			DatasourceID: grant.ResourceID,
		})
	}
	if grants.Err() != nil {
		return nil, grants.Err()
	}
	return result, nil
}

func loadRoleGrants(ctx context.Context, client *sdm.Client) ([]rolePermissionsRow, error) {
	roleGrants, err := client.RoleGrants().List(ctx, "")
	if err != nil {
		return nil, err
	}
	var result rolePermissionsList
	for roleGrants.Next() {
		grant := roleGrants.Value()
		result = append(result, rolePermissionsRow{
			ID:           grant.ID,
			RoleID:       grant.RoleID,
			DatasourceID: grant.ResourceID,
		})
	}
	if roleGrants.Err() != nil {
		return nil, roleGrants.Err()
	}
	return result, nil
}

func loadAccounts(ctx context.Context, client *sdm.Client) ([]userRow, error) {
	accountAttachments, err := client.AccountAttachments().List(ctx, "")
	if err != nil {
		return nil, err
	}
	roles := map[string]string{}
	for accountAttachments.Next() {
		attachment := accountAttachments.Value()
		roles[attachment.AccountID] = attachment.RoleID
	}
	if accountAttachments.Err() != nil {
		return nil, accountAttachments.Err()
	}

	accounts, err := client.Accounts().List(ctx, "type:user")
	if err != nil {
		return nil, err
	}
	var result userList
	for accounts.Next() {
		user := accounts.Value().(*sdm.User)
		result = append(result, userRow{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Role:      roles[user.ID],
		})
	}
	if accounts.Err() != nil {
		return nil, accounts.Err()
	}
	return result, nil
}

func loadResources(ctx context.Context, client *sdm.Client) ([]entitlement, error) {
	// limit grant/revoke to datasources and servers only, allowing websites
	// to be granted manually for the time being
	var resources entitlementList
	resp, err := client.Resources().List(ctx, "category:datasource")
	if err != nil {
		return nil, err
	}
	for resp.Next() {
		resource := resp.Value()
		resources = append(resources, entitlement{
			DatasourceID: resource.GetID(),
			Name:         resource.GetName(),
		})
	}
	if resp.Err() != nil {
		return nil, resp.Err()
	}
	resp, err = client.Resources().List(ctx, "category:server")
	if err != nil {
		return nil, err
	}
	for resp.Next() {
		resource := resp.Value()
		resources = append(resources, entitlement{
			DatasourceID: resource.GetID(),
			Name:         resource.GetName(),
		})
	}
	if resp.Err() != nil {
		return nil, resp.Err()
	}
	return resources, nil
}

// MatcherConfig stores mapping data from matchers.yml
type MatcherConfig struct {
	Groups []struct {
		Name      string   `yaml:"name"`
		Resources []string `yaml:"resources"`
	} `yaml:"groups"`
}

func loadMatchers() (*MatcherConfig, error) {
	body, err := ioutil.ReadFile("matchers.yml")
	if err != nil {
		return nil, errors.Wrap(err, "unable to read from matchers configuration file")
	}

	var m MatcherConfig
	err = yaml.UnmarshalStrict(body, &m)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling matcher configuration")
	}

	return &m, err
}

// matchEntitlements creates lists of concrete datasources and servers by group name
func matchEntitlements(ctx context.Context, client *sdm.Client, matchers *MatcherConfig) (map[string]entitlementList, error) {
	result := make(map[string]entitlementList)
	for _, matcher := range matchers.Groups {
		if debug {
			println("inside matchEntitlements: ", matcher.Name)
		}
		uniq := make(map[entitlement]bool)
		for _, expression := range matcher.Resources {
			resources, err := client.Resources().List(ctx, expression)
			if err != nil {
				return nil, err
			}
			for resources.Next() {
				rs := resources.Value()
				uniq[entitlement{DatasourceID: rs.GetID()}] = true
			}
			if resources.Err() != nil {
				return nil, err
			}
		}
		for u := range uniq {
			result[matcher.Name] = append(result[matcher.Name], u)
		}
	}
	return result, nil
}

func findUser(a []userRow, x oktaUser) bool {
	for _, n := range a {
		if strings.ToLower(n.Email) == strings.ToLower(x.Login) {
			return true
		}
	}
	return false
}

// AppendIfMissing adds Group/Role name to existing list, if not already present
func AppendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}
