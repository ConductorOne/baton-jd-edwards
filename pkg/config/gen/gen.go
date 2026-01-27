package main

import (
	cfg "github.com/conductorone/baton-jd-edwards/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("jd-edwards", cfg.Config)
}
