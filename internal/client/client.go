package client

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/danilovkiri/goonex/internal/dto"
	"github.com/rs/zerolog"
)

type Client struct {
	client *http.Client
	logger *zerolog.Logger
}

func (c *Client) decompressAndUnmarshallJsonObj(reader io.Reader, logger *zerolog.Logger, hubData *dto.HubData) error {
	gr, err := gzip.NewReader(reader)
	if err != nil {
		logger.Error().Err(err).Msg("could not create a gzip reader")
		return err
	}
	defer gr.Close()

	err = json.NewDecoder(gr).Decode(&hubData)
	if err != nil {
		logger.Error().Err(err).Msg("could not create decode compressed data")
		return err
	}
	return nil
}

func NewClient(logger *zerolog.Logger) *Client {
	return &Client{
		client: &http.Client{Timeout: 10 * time.Second},
		logger: logger,
	}
}

func (c *Client) NewPostRequestToHub(parcelid, idbox string) ([]dto.Hub, error) {
	const method = "POST"
	const url = "https://onex.am/parcel/hub"
	// create payload
	payload := strings.NewReader(fmt.Sprintf("parcel_id=%v&idbox=%v", parcelid, idbox))
	// create request
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to create a new hub request")
		return nil, err
	}
	// add headers
	req.Header.Add("Authority", "onex.am")
	req.Header.Add("Method", "POST")
	req.Header.Add("Path", "/parcel/hub")
	req.Header.Add("Scheme", "https")
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Add("Content-Length", "26")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Dnt", "1")
	req.Header.Add("Origin", "https://onex.am")
	req.Header.Add("Referer", "https://onex.am/ru/onextrack")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	// execute request
	res, err := c.client.Do(req)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to execute hub request")
		return nil, err
	}
	if res.StatusCode != 200 {
		c.logger.Error().Err(err).Msg(fmt.Sprintf("request returned code %v", res.StatusCode))
		return nil, err
	}
	defer res.Body.Close()

	// decompress and deserialize
	var post dto.HubData
	err = c.decompressAndUnmarshallJsonObj(res.Body, c.logger, &post)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to decompress and deserialize request response")
		return nil, err
	}

	return post.Data, nil
}

func (c *Client) NewPostRequestToTracker(trackingID string) (parcelid, idbox string, err error) {
	const method = "POST"
	const url = "https://onex.am/onextrack/findtrackingcodeimport"
	// create payload
	payload, payloadType, err := createMultipartFormDataPayload(trackingID)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to create multipart payload for tracking request")
		return "", "", err
	}
	// create request
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to create a new tracking request")
		return "", "", err
	}
	// set headers
	req.Header.Set("Content-Type", payloadType)
	// execute request
	res, err := c.client.Do(req)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to execute tracking request")
		return "", "", err
	}
	if res.StatusCode != 200 {
		c.logger.Error().Err(err).Msg(fmt.Sprintf("request returned code %v", res.StatusCode))
		return "", "", err
	}
	defer res.Body.Close()
	// read response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to read tracking request response")
		return "", "", err
	}

	// deserialize response into struct
	var post dto.TrackingCodeImportData
	err = json.Unmarshal(body, &post)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to unmarshal tracking request response")
		return "", "", err
	}
	return post.Data.Import.Parcelid, post.Data.Import.Idbox, nil

}

func createMultipartFormDataPayload(trackingID string) (io.Reader, string, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("tcode", trackingID)
	err := writer.Close()
	return payload, writer.FormDataContentType(), err
}
