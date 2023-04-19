package lab

import (
	"math"
	"time"

	"github.com/Drofff/revsynth/aco"
	"github.com/Drofff/revsynth/circuit"
	"github.com/Drofff/revsynth/logging"
)

type Lab struct {
	GateFactory   circuit.GateFactory
	Logger        logging.Logger
	TargetCircuit circuit.TruthVector
}

type ExperimentSettings struct {
	Filename           string
	Config             aco.Config
	ConfigModifier     func(aco.Config) aco.Config
	ModificationsCount int
}

type SynthesisResult struct {
	NumOfGates int
	Complexity int
	Duration   int64
}

const iterationsPerConfig = 3

var Experiments = []ExperimentSettings{
	{
		Filename: "num_of_ants",
		Config: aco.Config{
			NumOfAnts:       10,
			NumOfIterations: 15,
			Alpha:           2.0,
			Beta:            1.0,
			EvaporationRate: 0.3,
			DepositStrength: 100,

			LocalLoops:  4,
			SearchDepth: 6,
		},
		ConfigModifier: func(config aco.Config) aco.Config {
			config.NumOfAnts += 10
			return config
		},
		ModificationsCount: 9,
	},
	{
		Filename: "num_of_iterations",
		Config: aco.Config{
			NumOfAnts:       40,
			NumOfIterations: 5,
			Alpha:           2.0,
			Beta:            1.0,
			EvaporationRate: 0.3,
			DepositStrength: 100,

			LocalLoops:  4,
			SearchDepth: 6,
		},
		ConfigModifier: func(config aco.Config) aco.Config {
			config.NumOfIterations += 5
			return config
		},
		ModificationsCount: 5,
	},
	{
		Filename: "evaporation_rate",
		Config: aco.Config{
			NumOfAnts:       40,
			NumOfIterations: 15,
			Alpha:           2.0,
			Beta:            1.0,
			EvaporationRate: 0.1,
			DepositStrength: 100,

			LocalLoops:  4,
			SearchDepth: 6,
		},
		ConfigModifier: func(config aco.Config) aco.Config {
			config.EvaporationRate += 0.2
			return config
		},
		ModificationsCount: 4,
	},
	{
		Filename: "local_loops",
		Config: aco.Config{
			NumOfAnts:       40,
			NumOfIterations: 15,
			Alpha:           2.0,
			Beta:            1.0,
			EvaporationRate: 0.3,
			DepositStrength: 100,

			LocalLoops:  2,
			SearchDepth: 6,
		},
		ConfigModifier: func(config aco.Config) aco.Config {
			config.LocalLoops += 4
			return config
		},
		ModificationsCount: 9,
	},
	{
		Filename: "search_depth",
		Config: aco.Config{
			NumOfAnts:       40,
			NumOfIterations: 15,
			Alpha:           2.0,
			Beta:            1.0,
			EvaporationRate: 0.3,
			DepositStrength: 100,

			LocalLoops:  4,
			SearchDepth: 2,
		},
		ConfigModifier: func(config aco.Config) aco.Config {
			config.SearchDepth += 3
			return config
		},
		ModificationsCount: 5,
	},
}

func summarizeResults(results []SynthesisResult) SynthesisResult {
	summary := SynthesisResult{}
	resultsLen := float64(len(results))
	for _, result := range results {
		summary.NumOfGates += int(math.Round(float64(result.NumOfGates) / resultsLen))
		summary.Complexity += int(math.Round(float64(result.Complexity) / resultsLen))
		summary.Duration += int64(math.Round(float64(result.Duration) / resultsLen))
	}
	return summary
}

func (l Lab) RunSynthesis(config aco.Config) SynthesisResult {
	s := aco.NewSynthesizer(config, l.GateFactory, l.Logger)

	results := make([]SynthesisResult, 0)
	for i := 0; i < iterationsPerConfig; i++ {
		startedAt := time.Now().UnixMilli()
		res := s.Synthesise(l.TargetCircuit)

		duration := time.Now().UnixMilli() - startedAt
		results = append(results, SynthesisResult{
			NumOfGates: len(res.Gates),
			Complexity: res.Complexity,
			Duration:   duration,
		})
	}

	return summarizeResults(results)
}
