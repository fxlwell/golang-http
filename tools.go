package http

func GetServerOuterIpAddr() (string, error) {
	s, _, err := DefaultClient.Get("http://ifconfig.me").String()
	return s, err
}
