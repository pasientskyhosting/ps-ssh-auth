package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ddliu/go-httpclient"
	"github.com/tidwall/gjson"
	"os"
	"os/user"
)

var (
	server     string
	hostname   string
	username   string
	auth_token string
	login      string
	pass       string
	unsafetls  bool
)

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func init() {
	// Get hostname of server automatically
	name, err := os.Hostname()
	if err != nil {
		fmt.Println("Cannot get hostname of server. Please use the commandline to hardcode it")
	}

	// Define inputs
	flag.StringVar(&server, "server", "https://127.0.0.1:5000", "URL to PrivacyIDEA server.")
	flag.StringVar(&hostname, "hostname", name, "Hostname of server to validate")
	flag.StringVar(&username, "user", "", "Username to validate")
	flag.StringVar(&login, "login", "admin", "Login username to PrivacyIDEA")
	flag.StringVar(&pass, "pass", "test", "Login password to PrivacyIDEA")
	flag.BoolVar(&unsafetls, "unsafe", false, "Do not do SSL/TLS certificate check")
	flag.Parse()
}

func main() {
	httpclient.Defaults(httpclient.Map{
		httpclient.OPT_USERAGENT:      "SSH-Auth/1.1",
		httpclient.OPT_TIMEOUT:        5,
		httpclient.OPT_CONNECTTIMEOUT: 5,
	})

	res, err := httpclient.
		Begin().
		WithOption(httpclient.OPT_UNSAFE_TLS, unsafetls).
		Post(server+"/auth", map[string]string{
			"username": login,
			"password": pass,
		})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	result, err := res.ToString()

	if isJSON(result) != true {
		// did not get a json response from the server
		os.Exit(1)
	}

	// Get auth token for privacyIDEA
	auth_token := gjson.Get(result, "result.value.token")

	if auth_token.Exists() {

		// Get SSH keys for the machine and user
		res, err := httpclient.
			Begin().
			WithOption(httpclient.OPT_UNSAFE_TLS, unsafetls).
			WithHeader("Authorization", auth_token.String()).
			Get(server+"/machine/authitem/ssh", map[string]string{
				"hostname": hostname,
				"user":     username,
			})

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		result, err := res.ToString()

		if isJSON(result) != true {
			// did not get a json response from the server
			os.Exit(1)
		}

		keys := gjson.Get(result, "result.value.ssh")

		if keys.Exists() {
			if unix_user_exists(username) != true {
				if create_unix_user(username) != true {
					fmt.Println("Cannot create local user as it does not exist. User wont be able to log in")
					os.Exit(1)
				}
			}
		}

		keys.ForEach(func(key, value gjson.Result) bool {
			ssh_key := gjson.Get(value.String(), "sshkey")
			fmt.Println(ssh_key.String())

			return true
		})
	}
}

func unix_user_exists(username string) bool {
	user, err := user.Lookup(username)

	if err != nil {
		return false
	}

	if user != nil {
		return true
	}
	return true
}

func create_unix_user(username string) bool {
	return true
}
