package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/enescakir/emoji"
	"github.com/spf13/cobra"
)

var (
	myClient       = &http.Client{Timeout: 10 * time.Second}
	fahrenheitflag bool
	rootCmd        = &cobra.Command{
		Use:   "root",
		Short: "a root cobra program",
		Long:  "long a root cobra program",
	}
	getWeatherCmd = &cobra.Command{
		Use:   "getWeather",
		Short: "This command will get the weather today",
		Long:  `This get command will call OpenWeatherApi`,
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var URL string

			if fahrenheitflag {
				URL = "https://api.openweathermap.org/data/2.5/forecast?q=" + args[0] + "&units=imperial&appid=6bd5268b894241275a887ab1263f1c4f"
			} else {
				URL = "https://api.openweathermap.org/data/2.5/forecast?q=" + args[0] + "&units=metric&appid=6bd5268b894241275a887ab1263f1c4f"
			}

			fmt.Println("Try to get '" + args[0] + "' Weather")

			response, err := myClient.Get(URL)
			if err != nil {
				fmt.Println(err)
			}
			defer response.Body.Close()

			if response.StatusCode == 200 {
				body, err := ioutil.ReadAll(response.Body)
				if err != nil {
					panic(err)
				}

				var result Response
				if err := json.Unmarshal(body, &result); err != nil {
					panic(err)
				}

				for _, v := range result.List {
					clouds := ""
					for _, w := range v.Weather {
						clouds = w.Description
					}

					currentTime := time.Now()
					sl := strings.Split(v.DtTxt, " ")

					if currentTime.Format("2006-01-02") == sl[0] {
						fmt.Printf("	%v  %v\n", emoji.Cityscape, result.City.Name)
						fmt.Printf("	%v %v\n", emoji.Calendar, v.DtTxt)

						if fahrenheitflag {
							if v.Main.Temp < 35 {
								fmt.Printf("	%v  Temp: %vF\n", emoji.Snowflake, v.Main.Temp)
							} else {
								fmt.Printf("	%v Temp: %vF\n", emoji.SunWithFace, v.Main.Temp)
							}

							fmt.Printf("	%v Feels like: %vF\n", emoji.Man, v.Main.FeelsLike)
						} else {
							if v.Main.Temp < 0 {
								fmt.Printf("	%v  Temp: %v°\n", emoji.Snowflake, v.Main.Temp)
							} else {
								fmt.Printf("	%v Temp: %v°\n", emoji.SunWithFace, v.Main.Temp)
							}
							fmt.Printf("	%v Feels like: %v°\n", emoji.Man, v.Main.FeelsLike)
						}

						fmt.Printf("	%v  Clouds: %v\n", emoji.Cloud, clouds)
						fmt.Printf("	----------------------------\n")
					}
				}
			} else {
				fmt.Println("Error")
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(getWeatherCmd)
	getWeatherCmd.Flags().BoolVarP(&fahrenheitflag, "fahrenheitFlag", "F", false, "Fahrenheit")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
