package models

type ShareAccessLogEntry struct {
	ShareId    int
	UserEmail  string
	UserIp     string
	UserAgent  string
	AccessDate string
}
