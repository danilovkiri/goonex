package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/danilovkiri/goonex/internal/client"
	"github.com/danilovkiri/goonex/internal/logger"
	"github.com/olekukonko/tablewriter"
)

func main() {
	mainLogger := logger.InitLog()
	trackingID := flag.String("id", "0", "Valid tracking ID")
	flag.Parse()
	mainLogger.Debug().Str("trackingID", *trackingID).Msg("tracking ID from CLI was received")
	mainLogger.Debug().Msg("initializing client")
	mainClient := client.NewClient(mainLogger)
	mainLogger.Debug().Msg("trying to get onex parcel identifiers")
	parcelid, idbox, err := mainClient.NewPostRequestToTracker(*trackingID)
	if err != nil {
		mainLogger.Fatal().Err(err).Msg(fmt.Sprintf("could not find parcelID and IDbox for %s", *trackingID))
	}
	mainLogger.Debug().Str("parcelID", parcelid).Str("IDbox", idbox).Msg("successfully resolved parcelID and IDbox")
	if parcelid != "" && idbox != "" {
		mainLogger.Debug().Msg("trying to get onex parcel location history")
		hubs, err := mainClient.NewPostRequestToHub(parcelid, idbox)
		if err != nil {
			mainLogger.Fatal().Err(err).Msg(fmt.Sprintf("could not find hubs for %s", *trackingID))
		}
		sort.Slice(hubs, func(i, j int) bool {
			return hubs[i].Date.Before(hubs[j].Date)
		})

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Hub", "In/Out", "Timestamp"})
		for _, hub := range hubs {
			table.Append([]string{hub.Hub, hub.Type, hub.Date.Format(time.RFC1123)})
		}
		table.Render()
	}
}
