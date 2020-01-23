package models

//go:generate obfuscate Secret
type Secret string

//go:generate obfuscate AnotherSecret
type AnotherSecret struct {
	Key1 string
	Key2 int
}
