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
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
	"github.com/pkg/errors"
	sdm "github.com/strongdm/strongdm-sdk-go"
	"gopkg.in/yaml.v2"
)

const OKTA_USERS_LIMIT = 500

var oktaQueryString = "(status eq \"ACTIVE\")"

var jsonFlag = flag.Bool("json", false, "dump a JSON report for debugging")
var planFlag = flag.Bool("plan", false, "do not apply changes just plan and output the result")

// carefully use this flags
var deleteRolesNotInOktaFlag = flag.Bool("delete-roles-not-in-okta", false, "delete roles not present in okta")
var deleteUsersNotInOktaFlag = flag.Bool("delete-users-not-in-okta", false, "delete users not present in okta")

func init() {
	flag.Parse()
}

type syncReport struct {
	Start         time.Time    `json:"start"`
	Complete      time.Time    `json:"complete"`
	OktaUserCount int          `json:"oktaUsersCount"`
	OktaUsers     oktaUserList `json:"oktaUsers"`

	SDMUsersInOktaCount   int      `json:"sdmUsersInOktaCount"`
	SDMUsersInOkta        userList `json:"sdmUsersInOkta"`
	SDMUserNotInOktaCount int      `json:"sdmUsersNotInOktaCount"`
	SDMUsersNotInOkta     userList `json:"sdmUsersNotInOkta"`

	SDMRoleInOktaCount    int      `json:"sdmRolesInOktaCount"`
	SDMRolesInOkta        roleList `json:"sdmRolesInOkta"`
	SDMRoleNotInOktaCount int      `json:"sdmRolesNotInOktaCount"`
	SDMRolesNotInOkta     roleList `json:"sdmRolesNotInOkta"`

	Matchers *MatcherConfig `json:"matchers"`
}

func (rpt *syncReport) String() string {
	if !*jsonFlag {
		return rpt.short()
	}

	out, err := json.MarshalIndent(rpt, "", "\t")
	if err != nil {
		return fmt.Sprintf("error building JSON report: %s\n\n%s", err, rpt.short())
	}
	return string(out)
}

func (rpt *syncReport) short() string {
	return fmt.Sprintf("%d Okta users, %d strongDM users in okta, %d strongDM roles in okta\n",
		rpt.OktaUserCount, rpt.SDMUsersInOktaCount, rpt.SDMRoleInOktaCount)
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
		fmt.Fprintf(os.Stderr, "failed to initialize strongDM client: %v\n", err)
		os.Exit(1)
		return
	}

	var rpt syncReport
	rpt.Start = time.Now()

	matchers, err := loadMatchers()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading Matchers users: %v\n", err)
		os.Exit(1)
		return
	}
	rpt.Matchers = matchers

	oktaUsers, err := loadOktaUsers(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading Okta users: %v\n", err)
		os.Exit(1)
		return
	}
	rpt.OktaUsers = oktaUsers
	rpt.OktaUserCount = len(oktaUsers)

	initialRoles, err := loadRoles(ctx, client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading roles: %v\n", err)
		os.Exit(1)
		return
	}

	initialUsers, err := loadUsers(ctx, client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading users: %v\n", err)
		os.Exit(1)
		return
	}

	if !*planFlag {
		matchingRoles, unmatchingRoles, err := syncRoles(ctx, client, initialRoles, matchers)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error synchronizing roles: %v\n", err)
			os.Exit(1)
			return
		}
		rpt.SDMRolesInOkta = matchingRoles
		rpt.SDMRoleInOktaCount = len(matchingRoles)
		rpt.SDMRolesNotInOkta = unmatchingRoles
		rpt.SDMRoleNotInOktaCount = len(unmatchingRoles)

		matchingUsers, unmatchingUsers, err := syncUsers(ctx, client, initialUsers, matchingRoles, oktaUsers, matchers)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error synchronizing users: %v\n", err)
			os.Exit(1)
			return
		}
		rpt.SDMUsersInOkta = matchingUsers
		rpt.SDMUsersInOktaCount = len(matchingUsers)
		rpt.SDMUsersNotInOkta = unmatchingUsers
		rpt.SDMUserNotInOktaCount = len(unmatchingUsers)
	}

	rpt.Complete = time.Now()
	fmt.Println(rpt.String())
}

type oktaUserList []oktaUser

type oktaUser struct {
	Login     string   `json:"login"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Groups    []string `json:"groups"`
}

type roleList []roleRow

type roleRow struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type userList []userRow

type userRow struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Role      string `json:"roleName"`
}

type MatcherConfig struct {
	Groups []struct {
		Name      string   `yaml:"name"`
		Resources []string `yaml:"resources"`
	} `yaml:"groups"`
}

type entitlementList []entitlement

type entitlement struct {
	ResourceID string `json:"id"`
	Name       string `json:"name"`
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

func loadOktaUsers(ctx context.Context) (oktaUserList, error) {
	client, err := okta.NewClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "invalid Okta configuration")
	}
	search := query.NewQueryParams(query.WithSearch(oktaQueryString), query.WithLimit(OKTA_USERS_LIMIT))

	apiUsers, _, err := client.User.ListUsers(search)
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve okta users")
	}

	var users []oktaUser
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

		var u oktaUser
		u.Login = profile["login"].(string)
		u.FirstName = profile["firstName"].(string)
		u.LastName = profile["lastName"].(string)
		u.Groups = groupNames
		users = append(users, u)
	}
	return users, nil
}

func loadRoles(ctx context.Context, client *sdm.Client) (roleList, error) {
	roles, err := client.Roles().List(ctx, "")
	if err != nil {
		return nil, err
	}
	var result roleList
	for roles.Next() {
		role := roles.Value()
		result = append(result, roleRow{
			ID:   role.ID,
			Name: role.Name,
		})
	}
	if roles.Err() != nil {
		return nil, roles.Err()
	}
	return result, nil
}

func loadUsers(ctx context.Context, client *sdm.Client) (userList, error) {
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

func syncRoles(ctx context.Context, client *sdm.Client, initialRoles roleList, matchers *MatcherConfig) (roleList, roleList, error) {
	entitlementsByGroup, err := matchEntitlements(ctx, client, matchers)
	if err != nil {
		return nil, nil, err
	}
	matchingRoles, err := createMatchingRoles(ctx, client, entitlementsByGroup)
	if err != nil {
		return nil, nil, err
	}
	unmatchingRoles := calculateUnmatchingRoles(initialRoles, matchingRoles)
	if *deleteRolesNotInOktaFlag {
		err = deleteUnmatchingRoles(ctx, client, unmatchingRoles)
		if err != nil {
			return nil, nil, err
		}
	}
	return matchingRoles, unmatchingRoles, nil
}

// matchEntitlements creates lists of concrete datasources and servers by group name
func matchEntitlements(ctx context.Context, client *sdm.Client, matchers *MatcherConfig) (map[string]entitlementList, error) {
	result := make(map[string]entitlementList)
	for _, matcher := range matchers.Groups {
		uniq := make(map[entitlement]bool)
		for _, expression := range matcher.Resources {
			resources, err := client.Resources().List(ctx, expression)
			if err != nil {
				return nil, err
			}
			for resources.Next() {
				rs := resources.Value()
				uniq[entitlement{ResourceID: rs.GetID()}] = true
			}
			if resources.Err() != nil {
				return nil, err
			}
		}
		result[matcher.Name] = make(entitlementList, 0) // for creating groups without available resources
		for u := range uniq {
			result[matcher.Name] = append(result[matcher.Name], u)
		}
	}
	return result, nil
}

func createMatchingRoles(ctx context.Context, client *sdm.Client, entitlementsByGroup map[string]entitlementList) (roleList, error) {
	finalRoles := roleList{}
	for groupName, entitlements := range entitlementsByGroup {
		role, err := loadOrCreateRole(ctx, client, groupName, false)
		if err != nil {
			return nil, err
		}
		for _, e := range entitlements {
			err := createRoleGrant(ctx, client, role.ID, e.ResourceID)
			if err != nil {
				return nil, err
			}
		}
		finalRoles = append(finalRoles, roleRow{
			ID:   role.ID,
			Name: role.Name,
		})
	}
	return finalRoles, nil
}

func loadOrCreateRole(ctx context.Context, client *sdm.Client, roleName string, isComposite bool) (*sdm.Role, error) {
	roles, err := client.Roles().List(ctx, fmt.Sprintf("name:\"%s\"", roleName))
	if err != nil {
		return nil, err
	}
	if roles.Next() {
		return roles.Value(), nil
	}

	resp, err := client.Roles().Create(ctx, &sdm.Role{
		Name:      roleName,
		Composite: isComposite,
	})
	if err != nil {
		return nil, err
	}
	return resp.Role, nil
}

func createRoleGrant(ctx context.Context, client *sdm.Client, roleID string, resourceID string) error {
	_, err := client.RoleGrants().Create(ctx, &sdm.RoleGrant{
		RoleID:     roleID,
		ResourceID: resourceID,
	})
	var alreadyExistsErr *sdm.AlreadyExistsError
	if err != nil && !errors.As(err, &alreadyExistsErr) {
		return err
	}
	return nil
}

func calculateUnmatchingRoles(initialRoles roleList, matchingRoles roleList) roleList {
	unmatchingRoles := roleList{}
	for _, irole := range initialRoles {
		found := false
		for _, mrole := range matchingRoles {
			if irole.ID == mrole.ID {
				found = true
				break
			}
		}
		if !found {
			unmatchingRoles = append(unmatchingRoles, irole)
		}
	}
	return unmatchingRoles
}

func deleteUnmatchingRoles(ctx context.Context, client *sdm.Client, unmatchingRoles roleList) error {
	for _, role := range unmatchingRoles {
		_, err := client.Roles().Delete(ctx, role.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func syncUsers(ctx context.Context, client *sdm.Client, initialUsers userList, roles roleList, oktaUsers oktaUserList, matchers *MatcherConfig) (userList, userList, error) {
	matchingUsers, err := createMatchingUsers(ctx, client, roles, oktaUsers, matchers)
	if err != nil {
		return nil, nil, err
	}
	unmatchingUsers := calculateUnmatchingUsers(initialUsers, matchingUsers)
	if *deleteUsersNotInOktaFlag {
		err = deleteUnmatchingUsers(ctx, client, unmatchingUsers)
		if err != nil {
			return nil, nil, err
		}
	}
	return matchingUsers, unmatchingUsers, nil
}

func createMatchingUsers(ctx context.Context, client *sdm.Client, roles roleList, oktaUsers oktaUserList, matchers *MatcherConfig) (userList, error) {
	matchingUsers := userList{}
	for _, oktaUser := range oktaUsers {
		if !oktaUserHasMatchingGroup(oktaUser, matchers) {
			fmt.Fprintf(os.Stderr, "ignoring user %s - no group in matchers assigned to it\n", oktaUser.Login)
			continue
		}
		user, err := loadOrCreateUser(ctx, client, oktaUser)
		var alreadyExistsErr *sdm.AlreadyExistsError
		if errors.As(err, &alreadyExistsErr) {
			fmt.Fprintf(os.Stderr, "ignoring user %s - might be assigned to a different org\n", oktaUser.Login)
			continue
		}
		if err != nil {
			return nil, err
		}
		err = removePreviousAccountAttachments(ctx, client, user.ID)
		if err != nil {
			return nil, err
		}
		oktaGroups := matchingGroups(oktaUser.Groups, matchers)
		var roleName string
		if len(oktaGroups) == 1 {
			roleID, err := findRoleID(oktaGroups[0], roles)
			if err != nil {
				return nil, err
			}
			err = assignRole(ctx, client, user.ID, roleID)
			if err != nil {
				return nil, err
			}
			roleName = oktaGroups[0]
		} else if len(oktaGroups) > 1 {
			compositeRole, err := createCompositeRole(ctx, client, roles, oktaGroups)
			if err != nil {
				return nil, err
			}
			err = assignRole(ctx, client, user.ID, compositeRole.ID)
			if err != nil {
				return nil, err
			}
			roleName = compositeRole.Name
		}
		matchingUsers = append(matchingUsers, userRow{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Role:      roleName,
		})
	}
	return matchingUsers, nil
}

func oktaUserHasMatchingGroup(oktaUser oktaUser, matchers *MatcherConfig) bool {
	for _, oktaGroup := range oktaUser.Groups {
		for _, matcherGroup := range matchers.Groups {
			if oktaGroup == matcherGroup.Name {
				return true
			}
		}
	}
	return false
}

func loadOrCreateUser(ctx context.Context, client *sdm.Client, oktaUser oktaUser) (*sdm.User, error) {
	users, err := client.Accounts().List(ctx, fmt.Sprintf("email:\"%s\"", oktaUser.Login))
	if err != nil {
		return nil, err
	}
	if users.Next() {
		return users.Value().(*sdm.User), nil
	}

	resp, err := client.Accounts().Create(ctx, &sdm.User{
		Email:     oktaUser.Login,
		FirstName: oktaUser.FirstName,
		LastName:  oktaUser.LastName,
	})
	if err != nil {
		return nil, err
	}
	return resp.Account.(*sdm.User), nil
}

func removePreviousAccountAttachments(ctx context.Context, client *sdm.Client, userID string) error {
	attachments, err := client.AccountAttachments().List(ctx, fmt.Sprintf("accountId:\"%s\"", userID))
	if err != nil {
		return err
	}
	for attachments.Next() {
		attachmentID := attachments.Value().ID
		_, err := client.AccountAttachments().Delete(ctx, attachmentID)
		if err != nil {
			return err
		}
	}
	return nil
}

func matchingGroups(oktaGroups []string, matchers *MatcherConfig) []string {
	result := []string{}
	for _, oktaGroup := range oktaGroups {
		for _, matcherGroup := range matchers.Groups {
			if oktaGroup == matcherGroup.Name {
				result = append(result, oktaGroup)
			}
		}
	}
	return result
}

func findRoleID(groupName string, roles roleList) (string, error) {
	for _, r := range roles {
		if r.Name == groupName {
			return r.ID, nil
		}
	}
	return "", fmt.Errorf("cannot find roleID for roleName = %s", groupName)
}

func assignRole(ctx context.Context, client *sdm.Client, userID string, roleID string) error {
	_, err := client.AccountAttachments().Create(ctx, &sdm.AccountAttachment{
		AccountID: userID,
		RoleID:    roleID,
	})
	var alreadyExistsErr *sdm.AlreadyExistsError
	if err != nil && !errors.As(err, &alreadyExistsErr) {
		return err
	}
	return nil
}

func createCompositeRole(ctx context.Context, client *sdm.Client, roles roleList, oktaGroups []string) (*sdm.Role, error) {
	compositeRoleName := strings.Join(oktaGroups, "_")
	compositeRole, err := loadOrCreateRole(ctx, client, compositeRoleName, true)
	if err != nil {
		return nil, err
	}
	err = removePreviousCompositeRoleAttachments(ctx, client, compositeRole.ID)
	if err != nil {
		return nil, err
	}
	err = assignNewCompositeRoleAttachments(ctx, client, compositeRole.ID, roles, oktaGroups)
	if err != nil {
		return nil, err
	}
	return compositeRole, nil
}

func removePreviousCompositeRoleAttachments(ctx context.Context, client *sdm.Client, compositeRoleID string) error {
	attachments, err := client.RoleAttachments().List(ctx, fmt.Sprintf("compositeRoleId:\"%s\"", compositeRoleID))
	if err != nil {
		return err
	}
	for attachments.Next() {
		attachmentID := attachments.Value().ID
		_, err := client.RoleAttachments().Delete(ctx, attachmentID)
		if err != nil {
			return err
		}
	}
	return nil
}

func assignNewCompositeRoleAttachments(ctx context.Context, client *sdm.Client, compositeRoleID string, roles roleList, oktaGroups []string) error {
	for _, group := range oktaGroups {
		roleID, err := findRoleID(group, roles)
		if err != nil {
			return err
		}
		_, err = client.RoleAttachments().Create(ctx, &sdm.RoleAttachment{
			CompositeRoleID: compositeRoleID,
			AttachedRoleID:  roleID,
		})
		var alreadyExistsErr *sdm.AlreadyExistsError
		if err != nil && !errors.As(err, &alreadyExistsErr) {
			return err
		}
	}
	return nil
}

func calculateUnmatchingUsers(initialUsers userList, matchingUsers userList) userList {
	unmatchingUsers := userList{}
	for _, iuser := range initialUsers {
		found := false
		for _, muser := range matchingUsers {
			if iuser.ID == muser.ID {
				found = true
				break
			}
		}
		if !found {
			unmatchingUsers = append(unmatchingUsers, iuser)
		}
	}
	return unmatchingUsers
}

func deleteUnmatchingUsers(ctx context.Context, client *sdm.Client, unmatchingUsers userList) error {
	for _, user := range unmatchingUsers {
		_, err := client.Accounts().Delete(ctx, user.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
