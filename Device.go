package gonetworkmanager

import (
	"encoding/json"
	"errors"

	"github.com/godbus/dbus"
)

const (
	DeviceInterface = NetworkManagerInterface + ".Device"

	DevicePropertyInterface            = DeviceInterface + ".Interface"
	DevicePropertyIpInterface          = DeviceInterface + ".IpInterface"
	DevicePropertyState                = DeviceInterface + ".State"
	DevicePropertyIP4Config            = DeviceInterface + ".Ip4Config"
	DevicePropertyDeviceType           = DeviceInterface + ".DeviceType"
	DevicePropertyAvailableConnections = DeviceInterface + ".AvailableConnections"
	DevicePropertyDhcp4Config          = DeviceInterface + ".Dhcp4Config"
)

func DeviceFactory(objectPath dbus.ObjectPath) (Device, error) {
	d, err := NewDevice(objectPath)
	if err != nil {
		return nil, err
	}

	dt, err := d.GetDeviceType()
	if err != nil {
		return nil, err
	}
	switch dt {
	case NmDeviceTypeWifi:
		return NewWirelessDevice(objectPath)
	}

	return d, nil
}

type Device interface {
	GetPath() dbus.ObjectPath

	// GetInterface gets the name of the device's control (and often data)
	// interface.
	GetInterface() (string, error)

	// GetIpInterface gets the IP interface name of the device.
	GetIpInterface() (string, error)

	// GetState gets the current state of the device.
	GetState() (NmDeviceState, error)

	// GetIP4Config gets the Ip4Config object describing the configuration of the
	// device. Only valid when the device is in the NM_DEVICE_STATE_ACTIVATED
	// state.
	GetIP4Config() (IP4Config, error)

	// GetDHCP4Config gets the Dhcp4Config object describing the configuration of the
	// device. Only valid when the device is in the NM_DEVICE_STATE_ACTIVATED
	// state.
	GetDHCP4Config() (DHCP4Config, error)

	// GetDeviceType gets the general type of the network device; ie Ethernet,
	// WiFi, etc.
	GetDeviceType() (NmDeviceType, error)

	// GetAvailableConnections gets an array of object paths of every configured
	// connection that is currently 'available' through this device.
	GetAvailableConnections() ([]Connection, error)

	MarshalJSON() ([]byte, error)
}

func NewDevice(objectPath dbus.ObjectPath) (Device, error) {
	var d device
	return &d, d.init(NetworkManagerInterface, objectPath)
}

type device struct {
	dbusBase
}

func (d *device) GetPath() dbus.ObjectPath {
	return d.obj.Path()
}

func (d *device) GetInterface() (string, error) {
	return d.getStringProperty(DevicePropertyInterface)
}

func (d *device) GetIpInterface() (string, error) {
	return d.getStringProperty(DevicePropertyIpInterface)
}

func (d *device) GetState() (NmDeviceState, error) {
	r, err := d.getUint32Property(DevicePropertyState)
	if err != nil {
		return NmDeviceStateFailed, err
	}
	return NmDeviceState(r), nil
}

func (d *device) GetIP4Config() (IP4Config, error) {
	path, err := d.getObjectProperty(DevicePropertyIP4Config)
	if err != nil {
		return nil, err
	}
	if path == "/" {
		return nil, errors.New("device path was empty")
	}

	cfg, err := NewIP4Config(path)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (d *device) GetDHCP4Config() (DHCP4Config, error) {
	path, err := d.getObjectProperty(DevicePropertyDhcp4Config)
	if err != nil {
		return nil, err
	}
	if path == "/" {
		return nil, errors.New("device path was empty")
	}

	cfg, err := NewDHCP4Config(path)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (d *device) GetDeviceType() (NmDeviceType, error) {
	r, err := d.getUint32Property(DevicePropertyDeviceType)
	if err != nil {
		return NmDeviceTypeUnknown, err
	}
	return NmDeviceType(r), nil
}

func (d *device) GetAvailableConnections() ([]Connection, error) {
	connPaths, err := d.getSliceObjectProperty(DevicePropertyAvailableConnections)
	if err != nil {
		return nil, err
	}
	conns := make([]Connection, len(connPaths))

	for i, path := range connPaths {
		conns[i], err = NewConnection(path)
		if err != nil {
			return nil, err
		}
	}

	return conns, nil
}

func (d *device) marshalMap() (map[string]interface{}, error) {
	Interface, err := d.GetInterface()
	if err != nil {
		return nil, err
	}
	IPinterface, err := d.GetIpInterface()
	if err != nil {
		return nil, err
	}
	State, err := d.GetState()
	if err != nil {
		return nil, err
	}
	IP4Config, err := d.GetIP4Config()
	if err != nil {
		return nil, err
	}
	DHCP4Config, err := d.GetDHCP4Config()
	if err != nil {
		return nil, err
	}
	DeviceType, err := d.GetDeviceType()
	if err != nil {
		return nil, err
	}
	AvailableConnections, err := d.GetAvailableConnections()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Interface":            Interface,
		"IP interface":         IPinterface,
		"State":                State.String(),
		"IP4Config":            IP4Config,
		"DHCP4Config":          DHCP4Config,
		"DeviceType":           DeviceType.String(),
		"AvailableConnections": AvailableConnections,
	}, nil
}

func (d *device) MarshalJSON() ([]byte, error) {
	m, err := d.marshalMap()
	if err != nil {
		return nil, err
	}
	return json.Marshal(m)
}
