package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseArgs() (map[uint64]string, error) {
	config, err := getConfig()
	if err != nil {
		return nil, err
	}

	res := make(map[uint64]string)
	for _, cfg := range strings.Split(config, ",") {
		port, host, err := parseSingleConfig(cfg)
		if err != nil {
			return nil, err
		}

		res[port] = host
	}

	return res, nil
}

func getConfig() (string, error) {
	config, ok := os.LookupEnv("CONFIG")
	if ok {
		return config, nil
	}

	if len(os.Args) == 2 {
		return os.Args[1], nil
	}

	return "", fmt.Errorf("no config provided")
}

func parseSingleConfig(config string) (uint64, string, error) {
	parts := strings.Split(config, ":")
	if len(parts) != 2 {
		return 0, "", fmt.Errorf("invalid config: %s", config)
	}

	port, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		return 0, "", fmt.Errorf("invalid port %s in config %s", parts[0], config)
	}

	return port, parts[1], nil
}
