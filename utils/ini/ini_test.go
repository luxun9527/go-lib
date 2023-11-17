package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"testing"
)

func TestIni(t *testing.T) {
	cfg, err := ini.Load("./test.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	// Classic read of values, default section can be represented as empty string
	fmt.Println("App Mode:", cfg.Section("").Key("app_mode").String())
	fmt.Println("Data Path:", cfg.Section("paths").Key("data1").String())

	// Let's do some candidate value limitation
	fmt.Println("Server Protocol:",
		cfg.Section("server").Key("protocol").In("http", []string{"http", "https"}))
	// Value read that is not in candidates will be discarded and fall back to given default value
	fmt.Println("Email Protocol:", cfg.Section("server").Key("protocol").In("smtp", []string{"imap", "smtp"}))

	c := &Config{}
	if err := cfg.MapTo(&c); err != nil {
		t.Error(err)
		return
	}
	t.Log(c)
}

type Config struct {
	AppMode string `ini:"app_mode"`
	Paths   struct {
		Data string `ini:"data"`
	} `ini:"paths"`
	Server struct {
		Protocol      string `ini:"protocol"`
		HttpPort      string `ini:"http_port"`
		EnforceDomain string `ini:"enforce_domain"`
	} `ini:"server"`
}
