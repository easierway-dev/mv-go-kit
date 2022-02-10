package weight_cal

func GetZoneWeight(zoneAdjuster *WeightAdjuster, localZone string, serviceZone string) float64 {
	ratio := zoneAdjuster.GetWeight(localZone)
	culRatio := GetRatioByStep(ratio, 0.04)
	//cul zone weight
	if localZone == serviceZone {
		return culRatio
	} else {
		return 1 - culRatio
	}
}

func GetServiceWeight(serviceWeight *WeightAdjuster, serviceNode string) float64 {
	ratio := serviceWeight.GetWeight(serviceNode)
	return GetRatioByStep(ratio, 0.02)
}

func GetRatioByStep(ratio float64, step float64) float64 {
	beginRatio := 1.0
	//cul ratio by step
	for {
		if ratio > beginRatio-step {
			break
		}
		beginRatio -= step
		if beginRatio <= 0 {
			beginRatio = step
			break
		}
	}
	return beginRatio
}
