package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

// see https://developers.google.com/admin-sdk/directory/v1/limits
const USERS_LIMIT = 500

var googleHTTPClient *http.Client

func ValidateGoogleEnv() error {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		return errors.Errorf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, admin.AdminDirectoryUserReadonlyScope)
	if err != nil {
		return errors.Errorf("Unable to parse client secret file to config: %v", err)
	}
	googleHTTPClient = getClient(config)
	return nil
}

func LoadGoogleUsers(ctx context.Context, matchers *MatcherConfig) (idpUserList, error) {
	srv, err := admin.NewService(ctx, option.WithHTTPClient(googleHTTPClient))
	if err != nil {
		return nil, err
	}

	r, err := srv.Users.List().Customer("my_customer").MaxResults(USERS_LIMIT).Do()
	if err != nil {
		return nil, err
	}

	var users []idpUser
	rootRoleName := getRootRoleName(matchers)
	for _, googleUser := range r.Users {
		users = append(users, idpUser{
			Login:     googleUser.PrimaryEmail,
			FirstName: googleUser.Name.GivenName,
			LastName:  googleUser.Name.FamilyName,
			Groups:    getGroups(googleUser.OrgUnitPath, rootRoleName, matchers),
		})
	}
	return users, nil
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func getGroups(orgUnitPath string, rootRoleName string, matchers *MatcherConfig) []string {
	if orgUnitPath == "/" {
		if rootRoleName == "" {
			return []string{}
		} else {
			return []string{rootRoleName}
		}
	}
	orgUnits := strings.Split(orgUnitPath, "/")
	return []string{orgUnits[len(orgUnits)-1]}
}

func getRootRoleName(matchers *MatcherConfig) string {
	for _, group := range matchers.Groups {
		if group.Root {
			return group.Name
		}
	}
	return ""
}
