package main

// go get golang.org/x/text/language

import (
    "os"
    "fmt"
    "os/exec"
    "log"
    "net/http"
    "encoding/json"
    "strconv"

)

func Format(n int64) string {
	in := strconv.FormatInt(n, 10)
	numOfDigits := len(in)
	if n < 0 {
		numOfDigits-- // First character is the - sign (not a digit)
	}
	numOfCommas := (numOfDigits - 1) / 3

	out := make([]byte, len(in)+numOfCommas)
	if n < 0 {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ','
		}
	}
}

func main() {

	// accepts one optional argument the path, defaults to /
	searchPath := "/"
	if len(os.Args) > 1 { searchPath = os.Args[1] }

	// run the `df -H <searchPath>` command
	output, err := exec.Command("df", "-H", searchPath).Output()
	if err != nil { log.Fatal(err) }
	// Display the results
	fmt.Print(string(output))
	fmt.Printf("\n")

	// Query the server status using the REST API 
	// curl http://127.0.0.1:26657/status
	resp, err := http.Get("http://localhost:26657/status")
	if err != nil { log.Fatalln(err) }
	defer resp.Body.Close()

	// Parse the JSON results for intresting data
	var jResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&jResp); err != nil {
		log.Fatal("ooopsss! an error occurred, please try again")
	}
	status    := jResp["result"].(map[string]interface{})
	sync_info := status["sync_info"].(map[string]interface{})
	latest_block_height, err := strconv.ParseInt(sync_info["latest_block_height"].(string), 0, 64)
	catching_up := sync_info["catching_up"]
	voting_power, err := strconv.ParseInt(status["validator_info"].(map[string]interface{})["voting_power"].(string), 0, 64)
	if err != nil { log.Fatalln(err) }

	// Display intreting server status
	fmt.Printf("%19s %11s %15s\n", "latest_block_height", "catching_up", "voting_power")
	fmt.Printf("%19s %11t %15s\n", Format(latest_block_height), catching_up, Format(voting_power))

}
