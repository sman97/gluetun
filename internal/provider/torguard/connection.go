package torguard

import (
	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (t *Torguard) GetOpenVPNConnection(selection configuration.ServerSelection) (
	connection models.OpenVPNConnection, err error) {
	protocol := constants.UDP
	if selection.OpenVPN.TCP {
		protocol = constants.TCP
	}

	var port uint16 = 1912
	if selection.OpenVPN.CustomPort > 0 {
		port = selection.OpenVPN.CustomPort
	}

	servers, err := t.filterServers(selection)
	if err != nil {
		return connection, err
	}

	var connections []models.OpenVPNConnection
	for _, server := range servers {
		for _, IP := range server.IPs {
			connection := models.OpenVPNConnection{
				IP:       IP,
				Port:     port,
				Protocol: protocol,
			}
			connections = append(connections, connection)
		}
	}

	if selection.TargetIP != nil {
		return utils.GetTargetIPOpenVPNConnection(connections, selection.TargetIP)
	}

	return utils.PickRandomOpenVPNConnection(connections, t.randSource), nil
}
