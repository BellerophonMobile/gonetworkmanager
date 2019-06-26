package gonetworkmanager

import (
	"github.com/godbus/dbus"
)

const (
	ActiveConnectionInterface             = NetworkManagerInterface + ".Connection.Active"
	ActiveConnectionProperyConnection     = ActiveConnectionInterface + ".Connection"
	ActiveConnectionProperySpecificObject = ActiveConnectionInterface + ".SpecificObject"
	ActiveConnectionProperyID             = ActiveConnectionInterface + ".Id"
	ActiveConnectionProperyUUID           = ActiveConnectionInterface + ".Uuid"
	ActiveConnectionProperyType           = ActiveConnectionInterface + ".Type"
	ActiveConnectionProperyDevices        = ActiveConnectionInterface + ".Devices"
	ActiveConnectionProperyState          = ActiveConnectionInterface + ".State"
	ActiveConnectionProperyStateFlags     = ActiveConnectionInterface + ".StateFlags"
	ActiveConnectionProperyDefault        = ActiveConnectionInterface + ".Default"
	ActiveConnectionProperyIP4Config      = ActiveConnectionInterface + ".Ip4Config"
	ActiveConnectionProperyDHCP4Config    = ActiveConnectionInterface + ".Dhcp4Config"
	ActiveConnectionProperyDefault6       = ActiveConnectionInterface + ".Default6"
	ActiveConnectionProperyVPN            = ActiveConnectionInterface + ".Vpn"
	ActiveConnectionProperyMaster         = ActiveConnectionInterface + ".Master"
)

type ActiveConnection interface {
	// GetConnection gets connection object of the connection.
	GetConnection() (Connection, error)

	// GetSpecificObject gets a specific object associated with the active connection.
	GetSpecificObject() (AccessPoint, error)

	// GetID gets the ID of the connection.
	GetID() (string, error)

	// GetUUID gets the UUID of the connection.
	GetUUID() (string, error)

	// GetType gets the type of the connection.
	GetType() (string, error)

	// GetDevices gets array of device objects which are part of this active connection.
	GetDevices() ([]Device, error)

	// GetState gets the state of the connection.
	GetState() (uint32, error)

	// GetStateFlags gets the state flags of the connection.
	GetStateFlags() (uint32, error)

	// GetDefault gets the default IPv4 flag of the connection.
	GetDefault() (bool, error)

	// GetIP4Config gets the IP4Config of the connection.
	GetIP4Config() (IP4Config, error)

	// GetDHCP4Config gets the DHCP4Config of the connection.
	GetDHCP4Config() (DHCP4Config, error)

	// GetVPN gets the VPN flag of the connection.
	GetVPN() (bool, error)

	// GetMaster gets the master device of the connection.
	GetMaster() (Device, error)
}

func NewActiveConnection(objectPath dbus.ObjectPath) (ActiveConnection, error) {
	var a activeConnection
	return &a, a.init(NetworkManagerInterface, objectPath)
}

type activeConnection struct {
	dbusBase
}

func (a *activeConnection) GetConnection() (Connection, error) {
	path, err := a.getObjectProperty(ActiveConnectionProperyConnection)
	if err != nil {
		return nil, err
	}
	con, err := NewConnection(path)
	if err != nil {
		return nil, err
	}
	return con, nil
}

func (a *activeConnection) GetSpecificObject() (AccessPoint, error) {
	path, err := a.getObjectProperty(ActiveConnectionProperySpecificObject)
	if err != nil {
		return nil, err
	}
	ap, err := NewAccessPoint(path)
	if err != nil {
		return nil, err
	}
	return ap, nil
}

func (a *activeConnection) GetID() (string, error) {
	return a.getStringProperty(ActiveConnectionProperyID)
}

func (a *activeConnection) GetUUID() (string, error) {
	return a.getStringProperty(ActiveConnectionProperyUUID)
}

func (a *activeConnection) GetType() (string, error) {
	return a.getStringProperty(ActiveConnectionProperyType)
}

func (a *activeConnection) GetDevices() ([]Device, error) {
	paths, err := a.getSliceObjectProperty(ActiveConnectionProperyDevices)
	if err != nil {
		return nil, err
	}
	devices := make([]Device, len(paths))
	for i, path := range paths {
		devices[i], err = DeviceFactory(path)
		if err != nil {
			return nil, err
		}
	}
	return devices, nil
}

func (a *activeConnection) GetState() (uint32, error) {
	return a.getUint32Property(ActiveConnectionProperyState)
}

func (a *activeConnection) GetStateFlags() (uint32, error) {
	return a.getUint32Property(ActiveConnectionProperyStateFlags)
}

func (a *activeConnection) GetDefault() (bool, error) {
	b, err := a.getProperty(ActiveConnectionProperyDefault)
	if err != nil {
		return false, err
	}
	return b.(bool), nil
}

func (a *activeConnection) GetIP4Config() (IP4Config, error) {
	path, err := a.getObjectProperty(ActiveConnectionProperyIP4Config)
	if err != nil {
		return nil, err
	}
	r, err := NewIP4Config(path)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (a *activeConnection) GetDHCP4Config() (DHCP4Config, error) {
	path, err := a.getObjectProperty(ActiveConnectionProperyDHCP4Config)
	if err != nil {
		return nil, err
	}
	r, err := NewDHCP4Config(path)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (a *activeConnection) GetVPN() (bool, error) {
	ret, err := a.getProperty(ActiveConnectionProperyVPN)
	if err != nil {
		return false, err
	}
	return ret.(bool), nil
}

func (a *activeConnection) GetMaster() (Device, error) {
	path, err := a.getObjectProperty(ActiveConnectionProperyMaster)
	if err != nil {
		return nil, err
	}
	r, err := DeviceFactory(path)
	if err != nil {
		return nil, err
	}
	return r, nil
}
