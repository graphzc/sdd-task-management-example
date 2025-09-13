package timeutil

import "time"

func GetBangkokLocation() (*time.Location, error) {
	return time.LoadLocation("Asia/Bangkok")
}

func BangkokNow() time.Time {
	location, _ := GetBangkokLocation()

	return time.Now().In(location)
}
