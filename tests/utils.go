package tests

import (
	"GoHealthChecker/internal/model"
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
)

func ParseLinesForURL(content string, url string) [][]string {
	// Split the content into lines
	lines := strings.Split(content, "\n")

	// Iterate through each line and check if it contains the URL
	result := make([][]string, 0)
	for _, line := range lines {
		if strings.Contains(line, url) {
			tmp := make([]string, 0)
			parts := strings.Split(line, "|")
			// remove leading and trailing spaces
			for i := range parts {
				parts[i] = strings.TrimSpace(parts[i])
				if len(parts[i]) > 0 {
					tmp = append(tmp, parts[i])
				}
			}
			result = append(result, tmp)
		}
	}
	return result
}

func CreateConfiguration(timeout int, polling int) (*bytes.Buffer, context.Context, context.CancelFunc, model.AppSettings) {
	outputBuffer := new(bytes.Buffer)
	ctx, cancel := context.WithCancel(context.Background())
	settings := model.AppSettings{
		Timeout:         timeout,
		PollingInterval: polling,
		Context:         ctx,
		OutputStream:    outputBuffer,
	}
	return outputBuffer, ctx, cancel, settings
}

func VerifyRerenders(t *testing.T, results []string, e_status string, e_result string, e_latency float64, e_size int) {
	STATUS_INDEX := 1
	RESULT_INDEX := 2
	LATENCY_INDEX := 3
	SIZE_INDEX := 4

	status := results[STATUS_INDEX]
	result := results[RESULT_INDEX]
	latency, _ := strconv.ParseFloat(strings.TrimSpace(strings.ReplaceAll(results[LATENCY_INDEX], "ms", "")), 64)
	size, _ := strconv.Atoi(strings.TrimSpace(strings.ReplaceAll(results[SIZE_INDEX], "B", "")))

	assert.Equal(t, e_status, status)
	assert.Equal(t, e_result, result)
	// assert that latency is within 10% of expected
	assert.LessOrEqual(t, latency, e_latency*1.10)
	assert.GreaterOrEqual(t, latency, e_latency*0.90)
	assert.Equal(t, e_size, size)
}

func VerifyMetricsTable(t *testing.T, results []string, eHitRatio string, eHitPercentage string, eLatencyAvg float64, eSizeAvg int, eLatencyMin float64, eSizeMin int, eLatencyMax float64, eSizeMax int) {
	HIT_RATIO_INDEX := 1
	HIT_PERCENTAGE_INDEX := 2
	LATENCY_AVG_INDEX := 3
	SIZE_AVG_INDEX := 4
	LATENCY_MIN_INDEX := 5
	SIZE_MIN_INDEX := 6
	LATENCY_MAX_INDEX := 7
	SIZE_MAX_INDEX := 8

	hitRatio := results[HIT_RATIO_INDEX]
	hitPercentage := results[HIT_PERCENTAGE_INDEX]
	latencyAvg, _ := strconv.ParseFloat(strings.TrimSpace(strings.ReplaceAll(results[LATENCY_AVG_INDEX], "ms", "")), 64)
	sizeAvg, _ := strconv.Atoi(strings.TrimSpace(strings.ReplaceAll(results[SIZE_AVG_INDEX], "B", "")))
	latencyMin, _ := strconv.ParseFloat(strings.TrimSpace(strings.ReplaceAll(results[LATENCY_MIN_INDEX], "ms", "")), 64)
	sizeMin, _ := strconv.Atoi(strings.TrimSpace(strings.ReplaceAll(results[SIZE_MIN_INDEX], "B", "")))
	latencyMax, _ := strconv.ParseFloat(strings.TrimSpace(strings.ReplaceAll(results[LATENCY_MAX_INDEX], "ms", "")), 64)
	sizeMax, _ := strconv.Atoi(strings.TrimSpace(strings.ReplaceAll(results[SIZE_MAX_INDEX], "B", "")))

	assert.Equal(t, eHitRatio, hitRatio)
	assert.Equal(t, eHitPercentage, hitPercentage)
	// assert that latency is within 10% of expected
	assert.LessOrEqual(t, latencyAvg, eLatencyAvg*1.10)
	assert.GreaterOrEqual(t, latencyAvg, eLatencyAvg*0.90)
	assert.Equal(t, eSizeAvg, sizeAvg)
	assert.LessOrEqual(t, latencyMin, eLatencyMin*1.10)
	assert.GreaterOrEqual(t, latencyMin, eLatencyMin*0.90)
	assert.Equal(t, eSizeMin, sizeMin)
	assert.LessOrEqual(t, latencyMax, eLatencyMax*1.10)
	assert.GreaterOrEqual(t, latencyMax, eLatencyMax*0.90)
	assert.Equal(t, eSizeMax, sizeMax)
}
