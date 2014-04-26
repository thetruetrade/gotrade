package common

func RoundToTwoDecimals(untruncated float64) float64 {
	tmp := int(untruncated * 100)
	last := int(untruncated*1000) - tmp*100
	if last >= 50 {
		tmp += 1
	}
	return float64(tmp) / 100.0
}
