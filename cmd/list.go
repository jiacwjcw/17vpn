package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/shawnpeng17/17vpn/internal/pritunl"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List profiles",
	Run: func(cmd *cobra.Command, args []string) {
		p := pritunl.New()

		profiles := p.Profiles()
		if len(profiles) == 0 {
			color.Yellow("No profiles found in Pritunl!")
			os.Exit(1)
		}

		conns := p.Connections()

		var rows [][]string
		for i, profile := range profiles {
			row := []string{strconv.Itoa(i + 1), profile.Server, profile.User, "disconnected", "", "", ""}
			if conn, ok := conns[profile.ID]; ok {
				row[3] = conn.Status
				if conn.Timestamp > 0 {
					row[4] = formatDuration(time.Since(time.Unix(conn.Timestamp, 0)))
				}
				row[5] = conn.ClientAddr
				row[6] = conn.ServerAddr
			}
			row[3] = formatStatus(row[3])
			rows = append(rows, row)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Server", "User", "Status", "Connected", "Client IP", "Server IP"})
		table.AppendBulk(rows)
		table.SetColWidth(1000)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetBorder(false)
		table.SetCenterSeparator(" ")
		table.SetColumnSeparator(" ")
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetHeaderLine(false)
		table.Render()
	},
}

func formatStatus(status string) string {
	status = strings.ToUpper(status)
	switch status {
	case "CONNECTED":
		return color.New(color.FgGreen, color.Bold).SprintfFunc()(status)
	case "CONNECTING":
		return color.New(color.FgYellow, color.Bold).SprintfFunc()(status + "...")
	case "DISCONNECTING":
		return color.New(color.FgBlack, color.Bold).SprintfFunc()(status + "...")
	case "DISCONNECTED":
		return color.New(color.FgBlack, color.Bold).SprintfFunc()(status)
	default:
		return color.New(color.FgBlack, color.Bold).SprintfFunc()("UNKNOWN")
	}
}

func formatDuration(sec time.Duration) string {
	var ret string
	day := 24 * time.Hour
	d := sec / day
	sec = sec % day
	h := sec / time.Hour
	sec = sec % time.Hour
	m := sec / time.Minute
	sec = sec % time.Minute
	s := sec / time.Second
	if d > 0 {
		ret += fmt.Sprintf("%dd", d)
	}
	if h > 0 {
		ret += fmt.Sprintf("%dh", h)
	}
	if m > 0 {
		ret += fmt.Sprintf("%dm", m)
	}
	ret += fmt.Sprintf("%ds", s)
	return ret
}

func init() {
	rootCmd.AddCommand(listCmd)
}