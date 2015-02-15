package mst

import "log"

func MustNotErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func MustInt(i int, err error) int {
	if err != nil {
		log.Panic(err)
	}
	return i
}

func MustString(s string, err error) string {
	if err != nil {
		log.Panic(err)
	}
	return s
}

func MustStringArray(s []string, err error) []string {
	if err != nil {
		log.Panic(err)
	}
	return s
}
