package store

// Soar 性能分析
type Soar struct {
}

type SoarResult struct {
	Score float32
	Tips  string
}

func (s Soar) Analyst(sql string) (*SoarResult, error) {
	return &SoarResult{
		Score: 90.0,
		Tips:  "",
	}, nil
}
