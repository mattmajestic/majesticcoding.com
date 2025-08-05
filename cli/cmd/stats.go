package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats [provider]",
	Short: "Fetch stats from /api/stats/:provider",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		provider := args[0]
		valid := map[string]bool{
			"youtube": true,
			"github":  true,
			"twitch":  true,
		}
		if !valid[provider] {
			fmt.Println("Invalid provider. Use one of: youtube, github, twitch")
			return
		}

		resp, err := http.Get(fmt.Sprintf("https://majesticcoding.com/api/stats/%s", provider))
		if err != nil {
			fmt.Println("Request failed:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Failed to fetch stats:", resp.Status)
			return
		}

		body, _ := ioutil.ReadAll(resp.Body)
		var data map[string]interface{}
		json.Unmarshal(body, &data)

		fmt.Println(renderTable(data))
	},
}

func renderTable(data map[string]interface{}) string {
	keyStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("204"))
	valStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))

	var out string
	for k, v := range data {
		out += fmt.Sprintf("%s: %s\n",
			keyStyle.Render(k),
			valStyle.Render(fmt.Sprintf("%v", v)),
		)
	}
	return out
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
