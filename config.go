package main

import "os"

var (
	PORT string
)

// creating function to allow default val of undef env vars
func GetEnv(var_name, def_val string) string {
	value := os.Getenv(var_name)
	if len(value) < 1 {
		value = def_val
	}
	return value
}

// use init to pull in env vars
func init() {
	PORT = GetEnv("PORT", "3000")
}
