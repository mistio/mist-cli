package main

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gitlab.ops.mist.io/mistio/openapi-cli-generator/cli"
)

func formatMeteringData(metricsSet map[string]string, machineMetrics map[string]map[string]string, machineNames map[string]string) {
	metricsList := []string{}
	for metric, _ := range metricsSet {
		metricsList = append(metricsList, metric)
	}
	sort.Strings(metricsList)
	machines := make([]string, 0, len(machineMetrics))
	for machine, _ := range machineMetrics {
		machines = append(machines, machine)
	}
	sort.Strings(machines)
	data := make(map[string][]interface{})
	for _, machine := range machines {
		machineData := make(map[string]string)
		for _, metric := range metricsList {
			if _, ok := machineData["machine_id"]; !ok {
				machineData["machine_id"] = machine
				machineData["name"] = machineNames[machine]
			}
			machineData[metric] = machineMetrics[machine][metric]
		}
		if _, ok := data["data"]; !ok {
			data["data"] = make([]interface{}, 0)
		}
		data["data"] = append(data["data"], machineData)
	}
	metricSums := make(map[string]float64)
	for _, machineData := range data["data"] {
		for _, metric := range metricsList {
			value, err := strconv.ParseFloat(machineData.((map[string]string))[metric], 64)
			if err != nil {
				fmt.Println(err)
			} else {
				metricSums[metric] += value
			}
		}
	}
	sums := make([]string, len(metricsList))
	for i, metric := range metricsList {
		sums[i] = fmt.Sprintf("%f", metricSums[metric])
	}
	if err := cli.Formatter.Format(data, cli.CLIOutputOptions{append([]string{"name"}, metricsList...), append([]string{"machine_id", "name"}, metricsList...),append([]string{"TOTAL",}, sums...),append([]string{"TOTAL", ""}, sums...)}); err != nil {
		logger.Fatalf("Formatting failed: %s", err.Error())
	}
}

func parseTime(s string) (time.Time, error) {
	if t, err := strconv.ParseFloat(s, 64); err == nil {
		s, ns := math.Modf(t)
		return time.Unix(int64(s), int64(ns*float64(time.Second))).UTC(), nil
	}
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t, nil
	}
	return time.Time{}, errors.Errorf("cannot parse %q to a valid timestamp", s)
}

func getMeteringData(dtStart, dtEnd, search, queryTemplate string) (map[string]string, map[string]map[string]string, map[string]string) {
	paramsGetDatapoints := viper.New()
	paramsGetDatapoints.Set("time", dtEnd)
	paramsGetDatapoints.Set("search", search)
	dtStartTime, _ := parseTime(dtStart)
	dtEndTime, _ := parseTime(dtEnd)
	timeRange := int((dtEndTime.Sub(dtStartTime)).Seconds())
	query := fmt.Sprintf(queryTemplate, timeRange)
	_, decoded, _, err := MistApiV2GetDatapoints(query, paramsGetDatapoints)
	if err != nil {
		logger.Fatalf("Error calling operation: %s", err.Error())
	}

	type resultItem struct {
		Metric map[string]string `json:"metric"`
		Value  []interface{}     `json:"value"`
	}

	type promqlResponse struct {
		Data struct {
			DataPromql struct {
				Result []resultItem `json:"result"`
			} `json:"data"`
		} `json:"data"`
		Metadata map[string]interface{}
	}

	rawResponse, err := json.Marshal(decoded)
	if err != nil {
		fmt.Println("error:", err)
	}

	var response promqlResponse
	err = json.Unmarshal(rawResponse, &response)
	if err != nil {
		fmt.Println("error:", err)
	}

	metricsSet := make(map[string]string)
	machineMetrics := make(map[string]map[string]string)
	machineNames := make(map[string]string)

	for _, item := range response.Data.DataPromql.Result {
		if machineMetrics[item.Metric["machine_id"]] == nil {
			machineMetrics[item.Metric["machine_id"]] = make(map[string]string)
			machineNames[item.Metric["machine_id"]] = item.Metric["name"]
		}
		if item.Value != nil {
			machineMetrics[item.Metric["machine_id"]][item.Metric["__name__"]] = item.Value[1].(string)
		}
		metricsSet[item.Metric["__name__"]] = item.Metric["value_type"]
	}

	return metricsSet, machineMetrics, machineNames
}

func calculateDiffs(machineMetricsStart map[string]map[string]string, machineMetricsEnd map[string]map[string]string, metricsSet map[string]string) map[string]map[string]string {
	for machineId, metrics := range machineMetricsEnd {
		for metric, valueEnd := range metrics {
			if metricsSet[metric] != "counter" {
				continue
			}
			if _, ok := machineMetricsStart[machineId]; !ok {
				continue
			}
			if _, ok := machineMetricsStart[machineId][metric]; !ok {
				continue
			}
			valueStart := machineMetricsStart[machineId][metric]
			valueStartFloat, err := strconv.ParseFloat(valueStart, 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			valueEndFloat, err := strconv.ParseFloat(valueEnd, 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			machineMetricsEnd[machineId][metric] = fmt.Sprintf("%f", valueEndFloat-valueStartFloat)
		}
	}
	return machineMetricsEnd
}