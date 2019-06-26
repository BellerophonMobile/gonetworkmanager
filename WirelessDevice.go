package gonetworkmanager

import (
	"encoding/json"

	"github.com/godbus/dbus"
)

const (
	WirelessDeviceInterface = DeviceInterface + ".Wireless"

	WirelessDeviceGetAccessPoints = WirelessDeviceInterface + ".GetAccessPoints"
	WirelessDeviceRequestScan     = WirelessDeviceInterface + ".RequestScan"
)

type WirelessDevice interface {
	Device

	// GetAccessPoints gets the list of access points visible to this device.
	// Note that this list does not include access points which hide their SSID.
	// To retrieve a list of all access points (including hidden ones) use the
	// GetAllAccessPoints() method.
	GetAccessPoints() ([]AccessPoint, error)

	RequestScan() error
}

func NewWirelessDevice(objectPath dbus.ObjectPath) (WirelessDevice, error) {
	var d wirelessDevice
	return &d, d.init(NetworkManagerInterface, objectPath)
}

type wirelessDevice struct {
	device
}

func (d *wirelessDevice) GetAccessPoints() ([]AccessPoint, error) {
	var apPaths []dbus.ObjectPath

	err := d.call(&apPaths, WirelessDeviceGetAccessPoints)
	if err != nil {
		return nil, err
	}
	aps := make([]AccessPoint, len(apPaths))

	for i, path := range apPaths {
		aps[i], err = NewAccessPoint(path)
		if err != nil {
			return nil, err
		}
	}

	return aps, nil
}

func (d *wirelessDevice) RequestScan() error {
	var options map[string]interface{}
	return d.obj.Call(WirelessDeviceRequestScan, 0, options).Store()
}

func (d *wirelessDevice) MarshalJSON() ([]byte, error) {
	m, err := d.device.marshalMap()
	if err != nil {
		return nil, err
	}
	aps, err := d.GetAccessPoints()
	if err != nil {
		return nil, err
	}
	m["AccessPoints"] = aps
	return json.Marshal(m)
}
