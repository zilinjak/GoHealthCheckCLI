package model

type Metrics struct {
	TotalRequests   int `json:"total_requests"`
	FailedRequests  int `json:"failed_requests"`
	SuccessRequests int `json:"success_requests"`

	LatencyAverage float64 `json:"latency_average"`
	LatencyMin     float64 `json:"latency_min"`
	LatencyMax     float64 `json:"latency_max"`

	SizeAverage uint64 `json:"size_average"`
	SizeMin     uint64 `json:"size_min"`
	SizeMax     uint64 `json:"size_max"`
}

func NewMetrics(result HealthCheckResult) Metrics {
	// Convert values to float64
	latency := float64(result.Latency.Milliseconds())
	metrics := Metrics{
		TotalRequests:   1,
		FailedRequests:  0,
		SuccessRequests: 0,
		LatencyAverage:  latency,
		LatencyMin:      latency,
		LatencyMax:      latency,
		SizeAverage:     result.Size,
		SizeMin:         result.Size,
		SizeMax:         result.Size,
	}

	// Set success/failure count based on the result
	if result.IsOk {
		metrics.SuccessRequests = 1
	} else {
		metrics.FailedRequests = 1
	}

	return metrics
}

func (m *Metrics) Update(result HealthCheckResult) {
	// Update counters
	m.TotalRequests++
	if result.IsOk {
		m.SuccessRequests++
	} else {
		m.FailedRequests++
	}

	// Convert values to float64
	latency := float64(result.Latency.Milliseconds())
	size := result.Size

	// Update latency statistics
	m.LatencyAverage = ((m.LatencyAverage * float64(m.TotalRequests-1)) + latency) / float64(m.TotalRequests)
	if latency < m.LatencyMin {
		m.LatencyMin = latency
	}
	if latency > m.LatencyMax {
		m.LatencyMax = latency
	}

	// Update size statistics
	m.SizeAverage = uint64(((float64(m.SizeAverage) * float64(m.TotalRequests-1)) + float64(size)) / float64(m.TotalRequests))
	if size < m.SizeMin {
		m.SizeMin = size
	}
	if size > m.SizeMax {
		m.SizeMax = size
	}
}
