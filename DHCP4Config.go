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
	GetOptions() (DHCP4Options, error)

	MarshalJSON() ([]byte, error)
}

func NewDHCP4Config(objectPath dbus.ObjectPath) (DHCP4Config, error) {
	var c dhcp4Config
	return &c, c.init(NetworkManagerInterface, objectPath)
}

type dhcp4Config struct {
	dbusBase
}

func (c *dhcp4Config) GetOptions() (DHCP4Options, error) {
	options, err := c.getMapStringVariantProperty(DHCP4ConfigPropertyOptions)
	if err != nil {
		return nil, err
	}
	rv := make(DHCP4Options)

	for k, v := range options {
		rv[k] = v.Value()
	}

	return rv, nil
}

func (c *dhcp4Config) MarshalJSON() ([]byte, error) {
	Options, err := c.GetOptions()
	if err != nil {
		return nil, err
	}
	return json.Marshal(map[string]interface{}{
		"Options": Options,
	})
}
