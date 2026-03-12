package main

import "os"

var (
	PORT, DB_HOST, DB_PORT, DB_NAME, DB_USER, DB_PASS, DB_SSL string
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
	DB_HOST = GetEnv("DB_HOST", "")
	DB_PORT = GetEnv("DB_PORT", "")
	DB_NAME = GetEnv("DB_NAME", "")
	DB_USER = GetEnv("DB_USER", "")
	DB_PASS = GetEnv("DB_PASS", "")
	DB_SSL = GetEnv("DB_SSL", "enable")
}
