package gonetworkmanager

import (
	"github.com/godbus/dbus"
)

const (
	NetworkManagerInterface  = "org.freedesktop.NetworkManager"
	NetworkManagerObjectPath = "/org/freedesktop/NetworkManager"

	NetworkManagerGetDevices    = NetworkManagerInterface + ".GetDevices"
	NetworkManagerPropertyState = NetworkManagerInterface + ".state"
)

type NetworkManager interface {

	// GetDevices gets the list of network devices.
	GetDevices() []Device

	// GetState returns the overall networking state as determined by the
	// NetworkManager daemon, based on the state of network devices under it's
	// management.
	GetState() NmState
}

func NewNetworkManager() (NetworkManager, error) {
	var nm networkManager
	return &nm, nm.init(NetworkManagerInterface, NetworkManagerObjectPath)
}

type networkManager struct {
	dbusBase
}

func (n *networkManager) GetDevices() []Device {
	var devicePaths []dbus.ObjectPath

	n.call(&devicePaths, NetworkManagerGetDevices)
	devices := make([]Device, len(devicePaths))

	var err error
	for i, path := range devicePaths {
		devices[i], err = DeviceFactory(path)
		if err != nil {
			panic(err)
		}
	}

	return devices
}

func (n *networkManager) GetState() NmState {
	return NmState(n.getUint32Property(NetworkManagerPropertyState))
}
