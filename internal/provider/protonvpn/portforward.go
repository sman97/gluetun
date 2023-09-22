package protonvpn

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/netip"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/natpmp"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

var (
	ErrGatewayIPNotValid = errors.New("gateway IP address is not valid")
)

// PortForward obtains a VPN server side port forwarded from ProtonVPN gateway.
func (p *Provider) PortForward(ctx context.Context, _ *http.Client,
	logger utils.Logger, gateway netip.Addr, _ string) (
	port uint16, err error) {
	if !gateway.IsValid() {
		return 0, fmt.Errorf("%w", ErrGatewayIPNotValid)
	}

	client := natpmp.New()
	_, externalIPv4Address, err := client.ExternalAddress(ctx,
		gateway)
	if err != nil {
		return 0, fmt.Errorf("getting external IPv4 address: %w", err)
	}

	logger.Info("gateway external IPv4 address is " + externalIPv4Address.String())
	const internalPort, externalPort = 0, 0
	const lifetime = 60 * time.Second

	_, _, assignedUDPExternalPort, assignedLifetime, err :=
		client.AddPortMapping(ctx, gateway, "udp",
			internalPort, externalPort, lifetime)
	if err != nil {
		return 0, fmt.Errorf("adding UDP port mapping: %w", err)
	}
	checkLifetime(logger, "UDP", lifetime, assignedLifetime)

	_, _, assignedTCPExternalPort, assignedLifetime, err :=
		client.AddPortMapping(ctx, gateway, "tcp",
			internalPort, externalPort, lifetime)
	if err != nil {
		return 0, fmt.Errorf("adding TCP port mapping: %w", err)
	}
	checkLifetime(logger, "TCP", lifetime, assignedLifetime)

	checkExternalPorts(logger, assignedUDPExternalPort, assignedTCPExternalPort)
	port = assignedTCPExternalPort

	return port, nil
}

func checkLifetime(logger utils.Logger, protocol string,
	requested, actual time.Duration) {
	if requested != actual {
		logger.Warn(fmt.Sprintf("assigned %s port lifetime %s differs"+
			" from requested lifetime %s", strings.ToUpper(protocol),
			actual, requested))
	}
}

func checkExternalPorts(logger utils.Logger, udpPort, tcpPort uint16) {
	if udpPort != tcpPort {
		logger.Warn(fmt.Sprintf("UDP external port %d differs from TCP external port %d",
			udpPort, tcpPort))
	}
}

func (p *Provider) KeepPortForward(ctx context.Context, port uint16,
	gateway netip.Addr, _ string, logger utils.Logger) (err error) {
	client := natpmp.New()
	const refreshTimeout = 45 * time.Second
	timer := time.NewTimer(refreshTimeout)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}

		networkProtocols := []string{"udp", "tcp"}
		const internalPort = 0
		const lifetime = 60 * time.Second

		for _, networkProtocol := range networkProtocols {
			_, _, assignedExternalPort, assignedLiftetime, err :=
				client.AddPortMapping(ctx, gateway, networkProtocol,
					internalPort, port, lifetime)
			if err != nil {
				return fmt.Errorf("adding port mapping: %w", err)
			}

			if assignedLiftetime != lifetime {
				logger.Warn(fmt.Sprintf("assigned lifetime %s differs"+
					" from requested lifetime %s",
					assignedLiftetime, lifetime))
			}

			if port != assignedExternalPort {
				logger.Warn(fmt.Sprintf("external port assigned %d changed to %d",
					port, assignedExternalPort))
				port = assignedExternalPort
			}
		}

		timer.Reset(refreshTimeout)
	}
}
