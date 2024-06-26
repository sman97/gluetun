package env

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
	"github.com/qdm12/govalid/binary"
)

func (s *Source) readHTTPProxy() (httpProxy settings.HTTPProxy, err error) {
	httpProxy.User = s.env.Get("HTTPPROXY_USER",
		env.RetroKeys("PROXY_USER", "TINYPROXY_USER"),
		env.ForceLowercase(false))

	httpProxy.Password = s.env.Get("HTTPPROXY_PASSWORD",
		env.RetroKeys("PROXY_PASSWORD", "TINYPROXY_PASSWORD"),
		env.ForceLowercase(false))

	httpProxy.ListeningAddress, err = s.readHTTProxyListeningAddress()
	if err != nil {
		return httpProxy, err
	}

	httpProxy.Enabled, err = s.env.BoolPtr("HTTPPROXY", env.RetroKeys("PROXY", "TINYPROXY"))
	if err != nil {
		return httpProxy, err
	}

	httpProxy.Stealth, err = s.env.BoolPtr("HTTPPROXY_STEALTH")
	if err != nil {
		return httpProxy, err
	}

	httpProxy.Log, err = s.readHTTProxyLog()
	if err != nil {
		return httpProxy, err
	}

	return httpProxy, nil
}

func (s *Source) readHTTProxyListeningAddress() (listeningAddress string, err error) {
	const currentKey = "HTTPPROXY_LISTENING_ADDRESS"
	key := firstKeySet(s.env, "HTTPPROXY_PORT", "TINYPROXY_PORT", "PROXY_PORT",
		currentKey)
	switch key {
	case "":
		return "", nil
	case currentKey:
		return s.env.String(key), nil
	}

	// Retro-compatible keys using a port only
	s.handleDeprecatedKey(key, currentKey)
	port, err := s.env.Uint16Ptr(key)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(":%d", *port), nil
}

func (s *Source) readHTTProxyLog() (enabled *bool, err error) {
	const currentKey = "HTTPPROXY_LOG"
	key := firstKeySet(s.env, "PROXY_LOG", "TINYPROXY_LOG", "HTTPPROXY_LOG")
	switch key {
	case "":
		return nil, nil //nolint:nilnil
	case currentKey:
		return s.env.BoolPtr(key)
	}

	// Retro-compatible keys using different boolean verbs
	s.handleDeprecatedKey(key, currentKey)
	value := s.env.String(key)
	retroOption := binary.OptionEnabled("on", "info", "connect", "notice")

	enabled, err = binary.Validate(value, retroOption)
	if err != nil {
		return nil, fmt.Errorf("environment variable %s: %w", key, err)
	}

	return enabled, nil
}
