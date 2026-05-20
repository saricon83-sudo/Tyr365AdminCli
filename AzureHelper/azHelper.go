package azurehelper

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/saricon83-sudo/Tyr365AdminCli/internal/config"
)

type MetricResponse struct {
	Cost           int      `json:"cost"`
	Interval       string   `json:"interval"`
	Namespace      string   `json:"namespace"`
	ResourceRegion string   `json:"resourceregion"`
	Timespan       string   `json:"timespan"`
	Values         []Metric `json:"value"`
}

type Metric struct {
	DisplayDescription string       `json:"displayDescription"`
	ErrorCode          string       `json:"errorCode"`
	ID                 string       `json:"id"`
	Name               MetricName   `json:"name"`
	ResourceGroup      string       `json:"resourceGroup"`
	Timeseries         []TimeSeries `json:"timeseries"`
	Type               string       `json:"type"`
	Unit               string       `json:"unit"`
}

type MetricName struct {
	LocalizedValue string `json:"localizedValue"`
	Value          string `json:"value"`
}

type TimeSeries struct {
	Data []MetricData `json:"data"`
}

type MetricData struct {
	TimeStamp string  `json:"timeStamp"`
	Total     float64 `json:"total"`
}

type MetricsResult struct {
	Http5xxCount        float64 `json:"http_5_xx_count,omitempty"`
	TotalRequests       float64 `json:"total_requests,omitempty"`
	AverageResponseTime float64 `json:"average_response_time,omitempty"`
}

func GetMetrics() (*MetricsResult, error) {
	cfg := config.Get()
	resourcePath := cfg.GetString("azure.monitorResource")
	if resourcePath == "" {
		// Safe fallback to original hardcoded path for backwards compatibility
		resourcePath = "/subscriptions/e61fd8d2-77dc-42a8-a356-5537e74d8a87/resourceGroups/teamsprovisioning/providers/Microsoft.Web/sites/TyrensTeamsGovWebApi"
	}

	azureCLICommand := []string{
		"az",
		"monitor",
		"metrics",
		"list",
		"--resource",
		resourcePath,
		"--metric",
		"Http5xx,Requests,AverageResponseTime",
		"--aggregation",
		"Total",
		"--output",
		"json",
	}
	execCmd := exec.Command(azureCLICommand[0], azureCLICommand[1:]...)
	output, err := execCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error executing Azure CLI command: %w (output: %s)", err, string(output))
	}

	results, err := calculateMetrics(output)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate metrics from output: %w", err)
	}
	return &results, nil
}

func calculateMetrics(jsonBlob []byte) (MetricsResult, error) {
	var resp MetricResponse
	if err := json.Unmarshal(jsonBlob, &resp); err != nil {
		return MetricsResult{}, fmt.Errorf("error unmarshalling: %v", err)
	}

	var sumHttp5xx, sumRequests, sumResponseTime, countResponseTimeEntries float64

	for _, metric := range resp.Values {
		switch metric.Name.Value {
		case "Http5xx":
			for _, series := range metric.Timeseries {
				for _, data := range series.Data {
					if data.Total >= 1.0 {
						sumHttp5xx += data.Total
					}
				}
			}
		case "Requests":
			for _, series := range metric.Timeseries {
				for _, data := range series.Data {
					sumRequests += data.Total
				}
			}
		case "AverageResponseTime":
			for _, series := range metric.Timeseries {
				for _, data := range series.Data {
					sumResponseTime += data.Total
					countResponseTimeEntries++
				}
			}
		}
	}

	if countResponseTimeEntries == 0 { // Avoid division by zero
		return MetricsResult{}, fmt.Errorf("no entries for average response time")
	}

	averageResponseTime := sumResponseTime / countResponseTimeEntries

	return MetricsResult{
		Http5xxCount:        sumHttp5xx,
		TotalRequests:       sumRequests,
		AverageResponseTime: averageResponseTime,
	}, nil
}
