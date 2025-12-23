package main

type config struct {
	Token      []string `fig:"token" validate:"required"`
	BalconPath string   `fig:"balconpath" validate:"required"`
	Address    string   `fig:"address" validate:"required"`
	LogLevel   string   `fig:"loglevel" validate:"required"`
}
