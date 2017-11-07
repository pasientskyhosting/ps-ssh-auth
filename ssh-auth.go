package main

import (
	"encoding/json"
	"flag"
	"github.com/ddliu/go-httpclient"
	"github.com/tidwall/gjson"
	"os"
)

var (
	server     string
	hostname   string
	user       string
	auth_token string
)

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil

}

func init() {
	// Get hostname of server automatically
	name, err := os.Hostname()
	if err != nil {
		println("Cannot get hostname of server. Please use the flag")
	}

	// Define inputs
	flag.StringVar(&server, "server", "http://127.0.0.1:5000", "URL to PrivacyIDEA server.")
	flag.StringVar(&hostname, "hostname", name, "Hostname of server to validate")
	flag.StringVar(&user, "user", "", "Username to validate")
	flag.Parse()
}

func main() {
	httpclient.Defaults(httpclient.Map{
		httpclient.OPT_USERAGENT: "SSH-Auth",
		"Accept-Language":        "en-us",
	})

	res, err := httpclient.
		Begin().
		WithOption(httpclient.OPT_TIMEOUT, 10).
		Post(server+"/auth", map[string]string{
			"username": "admin",
			"password": "test",
		})

	if err != nil {
		println(err)
    os.Exit(0)
	}

	result, err := res.ToString()

	if isJSON(result) != true {
		println("Did not get a JSON response")
    os.Exit(0)
	}

	auth_token := gjson.Get(result, "result.value.token")

	if auth_token.Exists() {
		res, err := httpclient.
			Begin().
			WithHeader("Authorization", auth_token.String()).
			Get(server+"/machine/authitem/ssh", map[string]string{
				"hostname": hostname,
				"user": user,
			})

		if err != nil {
			println(err)
      os.Exit(0)
		}

		result, err := res.ToString()

		if isJSON(result) != true {
			println("Did not get a JSON response")
      os.Exit(0)
		}

		keys := gjson.Get(result, "result.value.ssh")
		keys.ForEach(func(key, value gjson.Result) bool {
			ssh_key := gjson.Get(value.String(), "sshkey")
			println(ssh_key.String())
			return true // keep iterating
		})
	}
}
