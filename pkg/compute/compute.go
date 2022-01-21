package compute

import (
	"github.com/prometheus/common/model"

	"github.com/sherine-k/kube-carbon-footprint/pkg/dataset"
)

func ComputeCarbonFootprint(matrix model.Matrix, instancetype *dataset.Instance, region *dataset.Region) model.Matrix {
	var carbonFP model.Matrix
	for _, samplestream := range matrix {
		aMetric := model.Metric{}
		for label, value := range samplestream.Metric {
			aMetric[label] = value
		}
		aSamplePairs := []model.SamplePair{}
		for _, samplePair := range samplestream.Values {
			aSamplePair := model.SamplePair{}
			aSamplePair.Timestamp = samplePair.Timestamp
			aSamplePair.Value = carbonFootprintFromLoad(samplePair.Value, *instancetype, *region)
			aSamplePairs = append(aSamplePairs, aSamplePair)
		}

		aSampleStream := model.SampleStream{
			Metric: aMetric,
			Values: aSamplePairs,
		}
		carbonFP = append(carbonFP, &aSampleStream)
	}
	return carbonFP

}

func carbonFootprintFromLoad(value model.SampleValue, instancetype dataset.Instance, region dataset.Region) model.SampleValue {
	//gCO₂eq = PUE * Power * ZoneCO2e / 1000
	var cfp model.SampleValue
	power := model.SampleValue(0)
	if value <= 10 {
		power = interpolate(value, 0, 10, model.SampleValue(instancetype.LoadIdle), model.SampleValue(instancetype.Load10))
	} else if value <= 50 {
		power = interpolate(value, 10, 50, model.SampleValue(instancetype.Load10), model.SampleValue(instancetype.Load50))
	} else {
		power = interpolate(value, 50, 100, model.SampleValue(instancetype.Load50), model.SampleValue(instancetype.Load100))
	}
	cfp = power * model.SampleValue(region.PUE) * model.SampleValue(region.CO2e) / 1000
	//region.PUE * region.CO2e * insta
	return cfp
}

func interpolate(value, usageLow, usageHigh, powerLow, powerHigh model.SampleValue) model.SampleValue {
	return powerLow + value*(powerHigh-powerLow)/(usageHigh-usageLow)
}
