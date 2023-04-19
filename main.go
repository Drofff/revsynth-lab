package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/Drofff/revsynth-lab/lab"
	"github.com/Drofff/revsynth/aco"
	"github.com/Drofff/revsynth/circuit"
	"github.com/Drofff/revsynth/logging"
)

var desiredVector = circuit.TruthVector{Inputs: [][]int{
	{0, 0, 0},
	{0, 0, 1},
	{0, 1, 0},
	{0, 1, 1},
	{1, 0, 0},
	{1, 0, 1},
	{1, 1, 0},
	{1, 1, 1},
}, Vector: []int{7, 5, 2, 4, 6, 1, 0, 3}}

func formatResult(config aco.Config, res lab.SynthesisResult) string {
	return fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n",
		config.NumOfAnts,
		config.NumOfIterations,
		config.Alpha,
		config.Beta,
		config.EvaporationRate,
		config.DepositStrength,
		config.LocalLoops,
		config.SearchDepth,
		res.NumOfGates,
		res.Complexity,
		res.Duration)
}

func main() {
	l := lab.Lab{
		GateFactory:   circuit.NewToffoliGateFactory(),
		Logger:        logging.NewLogger(logging.LevelInfo),
		TargetCircuit: desiredVector,
	}

	for ei, experiment := range lab.Experiments {
		fmt.Printf("starting experiment #%v (%v)\n", ei, experiment.Filename)

		experimentConfig := experiment.Config
		experimentResult := ""
		for mod := 0; mod <= experiment.ModificationsCount; mod++ {
			fmt.Printf("modification %v\n", mod)
			modRes := l.RunSynthesis(experimentConfig)
			experimentResult += formatResult(experimentConfig, modRes)

			experimentConfig = experiment.ConfigModifier(experimentConfig)
		}

		err := os.WriteFile("./experiments/"+experiment.Filename+".csv", []byte(experimentResult), fs.ModePerm)
		if err != nil {
			log.Fatalln("error writing results to file:", err)
		}
	}
}
