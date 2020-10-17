package util

func NormalizeRgb(r, g, b uint8) (float64, float64, float64) {
	const uint8ValueRange = 0x100
	return float64(r) / uint8ValueRange, float64(g) / uint8ValueRange, float64(b) / uint8ValueRange
}
