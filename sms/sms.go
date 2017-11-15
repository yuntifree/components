package sms

//SMS for sms send
type SMS interface {
	//Send return 0 for success, others for failure
	Send(phone string, code int) int
}
