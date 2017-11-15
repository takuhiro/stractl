// Copyright © 2017 Takuhiro Yoshida
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/strava/go.strava"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get strava resources",
	Long:  `Get strava resources`,
	Run: func(cmd *cobra.Command, args []string) {
		resource := args[0]
		if args[0] == "summary" {
			getSummary(cmd)
		} else {
			fmt.Printf("Unknown resource: %s\n", resource)
		}
	},
}

func init() {
	RootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getSummary(cmd *cobra.Command) {
	accessToken := viper.GetString("access_token")
	athleteId := viper.GetInt64("athlete_id")

	client := strava.NewClient(accessToken)
	service := strava.NewAthletesService(client)

	now := time.Now()
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	activities, err := service.ListActivities(athleteId).After(firstDay.Unix()).Before(now.Unix()).Do()
	if err != nil {
		fmt.Errorf("Failed to fetch activities: %v\n", err)
		return
	}

	sum := Summary{Start:firstDay, End:now}
	for _, a := range activities {
		sum.Count++
		sum.Distance += a.Distance
		sum.Time += a.MovingTime
		sum.Elevation += a.TotalElevationGain
	}
	printSummary(&sum)
}

// Summary is a container to save values for a specific period of time
type Summary struct {
	Start time.Time
	End time.Time
	Count int
	Distance float64
	Time int
	Elevation float64
}

func printSummary(summary *Summary) {
	fmt.Printf("%d activities (%s-%s)\n", summary.Count, summary.Start.Format("2006.01.02"), summary.End.Format("2006.01.02"))
	fmt.Printf("Distance(km):\t%f\n", summary.Distance/1000)
	fmt.Printf("Time(hh:mm:ss):\t%s\n", formatTime(summary.Time))
	fmt.Printf("Elevation(m):\t%f\n", summary.Elevation)
}

// formatTime returns time in "hh:mm:ss"
func formatTime(second int) string {
	hh := second / 3600
	mm := second % 3600 / 60
	ss := second % 60
	return fmt.Sprintf("%02d:%02d:%02d", hh, mm, ss)
}

