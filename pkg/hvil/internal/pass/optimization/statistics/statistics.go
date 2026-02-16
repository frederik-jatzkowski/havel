package statistics

type Statistics[T any] struct {
	StatisticsPass T `parser:"" json:",omitempty"`
}
