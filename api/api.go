package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	apiTypes "hammerpost/api-types"
	"hammerpost/model"

	"github.com/dghubble/sling"
)

var DbController string

var cli = http.Client{Timeout: 900 * time.Second}

func GetNodeInfo(s *apiTypes.DefaultSuccessResponse, e *apiTypes.DefaultErrorResponse) (string, error) {

	_, err := sling.New().Client(&cli).Get(DbController+"/info").Receive(s, e)
	if err != nil {
		return "", err
	}

	return s.Message, nil
}

func RestartDB(s *apiTypes.DefaultSuccessResponse, e *apiTypes.DefaultErrorResponse) error {
	err := stopDB(s, e)
	if err != nil {
		return err
	}
	return startDB(s, e)
}

func stopDB(s *apiTypes.DefaultSuccessResponse, e *apiTypes.DefaultErrorResponse) error {
	_, err := sling.New().Client(&cli).Get(DbController+"/stop").Receive(s, e)
	if err != nil {
		return err
	}

	return nil
}

func startDB(s *apiTypes.DefaultSuccessResponse, e *apiTypes.DefaultErrorResponse) error {
	_, err := sling.New().Client(&cli).Get(DbController+"/start").Receive(s, e)
	if err != nil {
		return err
	}

	return nil
}

func SetParams(params []model.Param, s *apiTypes.DefaultSuccessResponse, e *apiTypes.DefaultErrorResponse) error {
	_, err := sling.New().Client(&cli).Post(DbController+"/set-param").BodyJSON(params).Receive(s, e)
	if err != nil {
		return err
	}

	if e.Error != "" {
		return errors.New(e.Error)
	}

	return nil
}

func GetSystemLoad(s *apiTypes.DefaultSuccessResponse, e *apiTypes.DefaultErrorResponse) (string, error) {
	_, err := sling.New().Client(&cli).Get(DbController+"/load").Receive(s, e)
	if err != nil {
		return "", err
	}

	return s.Message, nil
}

func GetMetric(s *apiTypes.DefaultSuccessResponse, e *apiTypes.DefaultErrorResponse) (metric *model.Metric, err error) {
	_, err = sling.New().Client(&cli).Get(DbController+"/metrics").Receive(s, e)
	if err != nil {
		return &model.Metric{
			CpuUsage:    0,
			MemoryUsage: 0,
			DiskTps:     0,
			TestId:      0,
		}, err
	}

	metric = new(model.Metric)
	if err := json.Unmarshal([]byte(s.Message), &metric); err != nil {
		return &model.Metric{
			CpuUsage:    0,
			MemoryUsage: 0,
			DiskTps:     0,
			TestId:      0,
		}, err
	}
	return metric, nil
}

func GetLoad(s *apiTypes.DefaultSuccessResponse, e *apiTypes.DefaultErrorResponse) (load float64, err error) {
	_, err = sling.New().Client(&cli).Get(DbController+"/load").Receive(s, e)
	if err != nil {
		return 0, err
	}

	load, err = strconv.ParseFloat(s.Message, 64)
	if err != nil {
		return 0, err
	}
	return load, nil
}
