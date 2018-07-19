package gonetworkmanager

import (
	"encoding/json"

	"github.com/godbus/dbus"
)

const (
	DHCP4ConfigInterface = NetworkManagerInterface + ".DHCP4Config"

	DHCP4ConfigPropertyOptions = DHCP4ConfigInterface + ".Options"
)

type DHCP4Options map[string]interface{}

type DHCP4Config interface {
	// GetOptions gets options map of configuration returned by the IPv4 DHCP server.
	GetOptions() DHCP4Options

	MarshalJSON() ([]byte, error)
}

func NewDHCP4Config(objectPath dbus.ObjectPath) (DHCP4Config, error) {
	var c dhcp4Config
	return &c, c.init(NetworkManagerInterface, objectPath)
}

type dhcp4Config struct {
	dbusBase
}

func (c *dhcp4Config) GetOptions() DHCP4Options {
	options := c.getMapStringVariantProperty(DHCP4ConfigPropertyOptions)
	rv := make(DHCP4Options)

	for k, v := range options {
		rv[k] = v.Value()
	}

	return rv
}

func (c *dhcp4Config) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Options": c.GetOptions(),
	})
}
