package gonetworkmanager

import (
	"encoding/binary"
	"net"

	"github.com/godbus/dbus"
)

type dbusBase struct {
	conn *dbus.Conn
	obj  dbus.BusObject
}

func (d *dbusBase) init(iface string, objectPath dbus.ObjectPath) error {
	var err error

	d.conn, err = dbus.SystemBus()
	if err != nil {
		return err
	}

	d.obj = d.conn.Object(iface, objectPath)

	return nil
}

func (d *dbusBase) call(value interface{}, method string, args ...interface{}) {
	err := d.callError(value, method, args...)
	if err != nil {
		panic(err)
	}
}

func (d *dbusBase) callError(value interface{}, method string, args ...interface{}) error {
	return d.obj.Call(method, 0, args...).Store(value)
}

func (d *dbusBase) getProperty(iface string) interface{} {
	variant, err := d.obj.GetProperty(iface)
	if err != nil {
		panic(err)
	}
	return variant.Value()
}

func (d *dbusBase) getObjectProperty(iface string) dbus.ObjectPath {
	value, ok := d.getProperty(iface).(dbus.ObjectPath)
	if !ok {
		panic(ErrVariantType)
	}
	return value
}

func (d *dbusBase) getSliceObjectProperty(iface string) []dbus.ObjectPath {
	value, ok := d.getProperty(iface).([]dbus.ObjectPath)
	if !ok {
		panic(ErrVariantType)
	}
	return value
}

func (d *dbusBase) getStringProperty(iface string) string {
	value, ok := d.getProperty(iface).(string)
	if !ok {
		panic(ErrVariantType)
	}
	return value
}

func (d *dbusBase) getSliceStringProperty(iface string) []string {
	value, ok := d.getProperty(iface).([]string)
	if !ok {
		panic(ErrVariantType)
	}
	return value
}

func (d *dbusBase) getUint32Property(iface string) uint32 {
	value, ok := d.getProperty(iface).(uint32)
	if !ok {
		panic(ErrVariantType)
	}
	return value
}

func (d *dbusBase) getSliceUint32Property(iface string) []uint32 {
	value, ok := d.getProperty(iface).([]uint32)
	if !ok {
		panic(ErrVariantType)
	}
	return value
}

func (d *dbusBase) getSliceSliceUint32Property(iface string) [][]uint32 {
	value, ok := d.getProperty(iface).([][]uint32)
	if !ok {
		panic(ErrVariantType)
	}
	return value
}

func (d *dbusBase) getSliceByteProperty(iface string) []byte {
	value, ok := d.getProperty(iface).([]byte)
	if !ok {
		panic(ErrVariantType)
	}
	return value
}

func ip4ToString(ip uint32) string {
	bs := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(bs, ip)
	return net.IP(bs).String()
}
