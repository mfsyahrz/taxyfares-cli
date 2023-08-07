package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mfsyahrz/taxyfare-cli/pkg/taxyfare"
	"github.com/mfsyahrz/taxyfare-cli/pkg/utils/logger"
)

func Execute() {
	tm := taxyfare.New()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Input action [add [record] | sum]: ")
		scanner.Scan()

		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			logger.Error("empty input")
			os.Exit(1)
		}

		args := strings.Split(input, " ")

		switch args[0] {
		case "add":
			if len(args) != 3 {
				fmt.Println("Invalid \"add\" input")
				continue
			}

			record := fmt.Sprintf("%s %s", args[1], args[2])
			if err := tm.AddRecord(record); err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

		case "sum":
			if len(args) != 1 {
				fmt.Println("Invalid \"sum\" input")
				continue
			}

			currentFare, err := tm.GetCurrentFare()
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			fmt.Println("Current Fare: ", currentFare)
			fmt.Println("Record History: \n", tm.GetRecordHistory())

			logger.Info("summary from the trip - current fare: ", currentFare)
		default:
			fmt.Println("Invalid action")
		}

	}

}
