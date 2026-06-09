package gopro

func isGoPro(name string) bool {
	if len(name) < 6 {
		return false
	}
	return name[:6] == "GoPro "
}
