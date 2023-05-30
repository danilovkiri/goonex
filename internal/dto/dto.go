package dto

import (
	"encoding/json"
	"strings"
	"time"
)

type HubDateTime time.Time

type Hub struct {
	Hub  string      `json:"hub"`
	Date HubDateTime `json:"date"`
	Type string      `json:"type"`
}

func (j *HubDateTime) Before(u HubDateTime) bool {
	a, _ := time.Parse(time.RFC1123, j.Format(time.RFC1123))
	b, _ := time.Parse(time.RFC1123, u.Format(time.RFC1123))
	return a.Before(b)
}

func (j *HubDateTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}
	*j = HubDateTime(t)
	return nil
}

func (j HubDateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(j))
}

func (j HubDateTime) Format(s string) string {
	t := time.Time(j)
	return t.Format(s)
}

type HubData struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Data   []Hub  `json:"data"`
}

type TrackingCodeImportData struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   struct {
		Import struct {
			Country       string `json:"country"`
			CountryTo     string `json:"countryTo"`
			Trackingcode  string `json:"trackingcode"`
			Weight        string `json:"weight"`
			Orderstatus   string `json:"orderstatus"`
			Parcelid      string `json:"parcelid"`
			Idbox         string `json:"idbox"`
			WoScanneddate string `json:"wo_scanneddate"`
			VWeight       string `json:"v_weight"`
			Inusadate     string `json:"inusadate"`
			Inmywaydate   string `json:"inmywaydate"`
			Inarmeniadate any    `json:"inarmeniadate"`
			Receiveddate  any    `json:"receiveddate"`
			Estimateddate string `json:"estimateddate"`
		} `json:"import"`
		Track struct {
			ID                      int    `json:"id"`
			TrackingNumber          string `json:"tracking_number"`
			TrackingNumberSecondary any    `json:"tracking_number_secondary"`
			TrackingNumberCurrent   any    `json:"tracking_number_current"`
			Courier                 struct {
				Slug        string `json:"slug"`
				Name        string `json:"name"`
				NameAlt     string `json:"name_alt"`
				CountryCode string `json:"country_code"`
				ReviewCount int    `json:"review_count"`
				ReviewScore string `json:"review_score"`
			} `json:"courier"`
			IsActive    bool   `json:"is_active"`
			IsDelivered bool   `json:"is_delivered"`
			LastCheck   string `json:"last_check"`
			Checkpoints []any  `json:"checkpoints"`
			Extra       []any  `json:"extra"`
		} `json:"track"`
		Iherb bool `json:"iherb"`
	} `json:"data"`
}

type TrackingCodeImportDataEmpty struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   struct {
		Import bool  `json:"import"`
		Track  bool  `json:"track"`
		Iherb  []any `json:"iherb"`
	} `json:"data"`
}
