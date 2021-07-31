package stats

type Stats struct {
	Usage Usage
}

type Usage struct {
	SpaceUsed   float64 `json:"space_used"`
	SpaceTotal  float64 `json:"space_total"`
	NodesOnline int     `json:"nodes_online"`
}
