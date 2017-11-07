package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ddliu/go-httpclient"
	"github.com/tidwall/gjson"
	"os"
)

var (
	server     string
	hostname   string
	user       string
	auth_token string
	login      string
	pass       string
)

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil

}

func init() {
	// Get hostname of server automatically
	name, err := os.Hostname()
	if err != nil {
		fmt.Print("Cannot get hostname of server. Please use the flag")
	}

	// Define inputs
	flag.StringVar(&server, "server", "http://127.0.0.1:5000", "URL to PrivacyIDEA server.")
	flag.StringVar(&hostname, "hostname", name, "Hostname of server to validate")
	flag.StringVar(&user, "user", "", "Username to validate")
	flag.StringVar(&login, "login", "admin", "Login username to PrivacyIDEA")
	flag.StringVar(&pass, "pass", "test", "Login password to PrivacyIDEA")
	flag.Parse()
}

func main() {
	httpclient.Defaults(httpclient.Map{
		httpclient.OPT_USERAGENT: "SSH-Auth",
	})

	res, err := httpclient.
		Begin().
		WithOption(httpclient.OPT_TIMEOUT, 10).
		Post(server+"/auth", map[string]string{
			"username": login,
			"password": pass,
		})

	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	result, err := res.ToString()

	if isJSON(result) != true {
		fmt.Print("Did not get a JSON response")
		os.Exit(2)
	}

	auth_token := gjson.Get(result, "result.value.token")

	if auth_token.Exists() {
		res, err := httpclient.
			Begin().
			WithHeader("Authorization", auth_token.String()).
			Get(server+"/machine/authitem/ssh", map[string]string{
				"hostname": hostname,
				"user":     user,
			})

		if err != nil {
			fmt.Print(err)
			os.Exit(3)
		}

		result, err := res.ToString()

		if isJSON(result) != true {
			fmt.Print("Did not get a JSON response")
			os.Exit(4)
		}

		keys := gjson.Get(result, "result.value.ssh")
		keys.ForEach(func(key, value gjson.Result) bool {
			ssh_key := gjson.Get(value.String(), "sshkey")
			fmt.Print(ssh_key.String())
			return true
		})
	}
}
