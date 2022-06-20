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

func formatMeteringData(metricsSet map[string]string, resourceMetrics map[string]map[string]string, resourceNames map[string]string) {
	metricsList := []string{}
	for metric := range metricsSet {
		metricsList = append(metricsList, metric)
	}
	sort.Strings(metricsList)
	resources := make([]string, 0, len(resourceMetrics))
	for resource := range resourceMetrics {
		resources = append(resources, resource)
	}
	sort.Strings(resources)
	data := make(map[string][]interface{})
	for _, resource := range resources {
		resourceData := make(map[string]string)
		for _, metric := range metricsList {
			if _, ok := resourceData["machine_id"]; !ok {
				resourceData["machine_id"] = resource
				resourceData["name"] = resourceNames[resource]
			}
			resourceData[metric] = resourceMetrics[resource][metric]
		}
		if _, ok := data["data"]; !ok {
			data["data"] = make([]interface{}, 0)
		}
		data["data"] = append(data["data"], resourceData)
	}
	metricSums := make(map[string]float64)
	for _, resourceData := range data["data"] {
		for _, metric := range metricsList {
			valueString, ok := resourceData.((map[string]string))[metric]
			if !ok || valueString == "" {
				continue
			}
			value, err := strconv.ParseFloat(valueString, 64)
			if err != nil {
				fmt.Printf("metric: %s, value: %s\n", metric, valueString)
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
	if err := cli.Formatter.Format(data, &viper.Viper{}, cli.CLIOutputOptions{append([]string{"name"}, metricsList...), append([]string{"machine_id", "name"}, metricsList...), append([]string{"TOTAL"}, sums...), append([]string{"TOTAL", ""}, sums...), map[string]string{}}); err != nil {
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

func getResourceNamesIDMap(search string) map[string]string {
	resourceNames := make(map[string]string)
	paramsListResources := viper.New()
	paramsListResources.Set("search", search)
	_, decoded, _, err := MistApiV2ListMachines(paramsListResources)
	if err != nil {
		logger.Fatalf("Error calling operation: %s", err.Error())
	}
	for _, item := range decoded["data"].([]interface{}) {
		resourceNames[item.(map[string]interface{})["id"].(string)] = item.(map[string]interface{})["name"].(string)
	}
	_, decoded, _, err = MistApiV2ListVolumes(paramsListResources)
	if err != nil {
		logger.Fatalf("Error calling operation: %s", err.Error())
	}
	for _, item := range decoded["data"].([]interface{}) {
		resourceNames[item.(map[string]interface{})["id"].(string)] = item.(map[string]interface{})["name"].(string)
	}
	return resourceNames
}

func mapResourceNamesWithMetrics(response promqlResponse, search string) (map[string]string, map[string]map[string]string, map[string]string) {
	metricsNameSet := make(map[string]string)
	resourceIDToMetricMap := make(map[string]map[string]string)
	resourceIDToNameMap := getResourceNamesIDMap(search)

	for _, item := range response.Data.DataPromql.Result {
		resourceID, ok := item.Metric["machine_id"]
		if !ok {
			resourceID, ok = item.Metric["volume_id"]
			if !ok {
				continue
			}
		}
		if resourceIDToMetricMap[resourceID] == nil {
			resourceIDToMetricMap[resourceID] = make(map[string]string)
		}
		if item.Value != nil {
			resourceIDToMetricMap[resourceID][item.Metric["__name__"]] = item.Value[1].(string)
		}
		metricsNameSet[item.Metric["__name__"]] = item.Metric["value_type"]
	}

	return metricsNameSet, resourceIDToMetricMap, resourceIDToNameMap
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

	rawResponse, err := json.Marshal(decoded)
	if err != nil {
		fmt.Println("error:", err)
	}

	var response promqlResponse
	err = json.Unmarshal(rawResponse, &response)
	if err != nil {
		fmt.Println("error:", err)
	}

	return mapResourceNamesWithMetrics(response, search)
}

func calculateDiffs(resourceMetricsStart map[string]map[string]string, resourceMetricsEnd map[string]map[string]string, metricsSet map[string]string) map[string]map[string]string {
	for resourceID, metrics := range resourceMetricsEnd {
		for metric, valueEnd := range metrics {
			if metricsSet[metric] != "counter" {
				continue
			}
			if _, ok := resourceMetricsStart[resourceID]; !ok {
				continue
			}
			if _, ok := resourceMetricsStart[resourceID][metric]; !ok {
				continue
			}
			valueStart := resourceMetricsStart[resourceID][metric]
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
			resourceMetricsEnd[resourceID][metric] = fmt.Sprintf("%f", valueEndFloat-valueStartFloat)
		}
	}
	return resourceMetricsEnd
}
