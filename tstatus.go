package main

// go get golang.org/x/text/language

import (
    "os"
    "os/exec"
    "log"
    "net/http"
    "encoding/json"
    "strconv"
    "golang.org/x/text/language"
    "golang.org/x/text/message"

)

func main() {

	// accepts one optional argument the path, defaults to /
	searchPath := "/"
	if len(os.Args) > 1 { searchPath = os.Args[1] }
	// setup a pretty printer
	printer := message.NewPrinter(language.English)

	// run the `df -H <searchPath>` command
	output, err := exec.Command("df", "-H", searchPath).Output()
	if err != nil { log.Fatal(err) }
	// Display the results
	printer.Print(string(output))
	printer.Printf("\n")

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
	result   := jResp["result"].(map[string]interface{})
	syncInfo := result["sync_info"].(map[string]interface{})
	latest_block_height, err := strconv.ParseInt(syncInfo["latest_block_height"].(string), 0, 64)
	catching_up := syncInfo["catching_up"]
	voting_power, err := strconv.ParseInt(result["validator_info"].(map[string]interface{})["voting_power"].(string), 0, 64)
	if err != nil { log.Fatalln(err) }

	// Display intreting server status
	printer.Printf("%12s %12s %15s\n", "block height", "catching up", "voting power")
	printer.Printf("%12d %12t %15d\n", latest_block_height, catching_up, voting_power)

}
