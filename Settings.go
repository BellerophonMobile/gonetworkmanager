package gonetworkmanager

import (
	"github.com/godbus/dbus"
)

const (
	SettingsInterface  = NetworkManagerInterface + ".Settings"
	SettingsObjectPath = NetworkManagerObjectPath + "/Settings"

	SettingsListConnections = SettingsInterface + ".ListConnections"
	SettingsAddConnection   = SettingsInterface + ".AddConnection"
)

type Settings interface {

	// ListConnections gets list the saved network connections known to NetworkManager
	ListConnections() ([]Connection, error)

	// AddConnection call new connection and save it to disk.
	AddConnection(settings ConnectionSettings) (Connection, error)
}

func NewSettings() (Settings, error) {
	var s settings
	return &s, s.init(NetworkManagerInterface, SettingsObjectPath)
}

type settings struct {
	dbusBase
}

func (s *settings) ListConnections() ([]Connection, error) {
	var connectionPaths []dbus.ObjectPath

	err := s.call(&connectionPaths, SettingsListConnections)
	if err != nil {
		return nil, err
	}
	connections := make([]Connection, len(connectionPaths))

	for i, path := range connectionPaths {
		connections[i], err = NewConnection(path)
		if err != nil {
			return nil, err
		}
	}

	return connections, nil
}

func (s *settings) AddConnection(settings ConnectionSettings) (Connection, error) {
	var path dbus.ObjectPath
	err := s.call(&path, SettingsAddConnection, settings)
	if err != nil {
		return nil, err
	}
	con, err := NewConnection(path)
	if err != nil {
		return nil, err
	}
	return con, nil
}
