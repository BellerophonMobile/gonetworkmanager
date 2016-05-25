package gonetworkmanager

//go:generate stringer -type=NmConnectivity
type NmConnectivity uint32

const (
	NmConnectivityUnknown NmConnectivity = 0
	NmConnectivityNone    NmConnectivity = 1
	NmConnectivityPortal  NmConnectivity = 2
	NmConnectivityLimited NmConnectivity = 3
	NmConnectivityFull    NmConnectivity = 4
)

//go:generate stringer -type=NmState
type NmState uint32

const (
	NmStateUnknown         NmState = 0
	NmStateAsleep          NmState = 10
	NmStateDisconnected    NmState = 20
	NmStateDisconnecting   NmState = 30
	NmStateConnecting      NmState = 40
	NmStateConnectedLocal  NmState = 50
	NmStateConnectedSite   NmState = 60
	NmStateConnectedGlobal NmState = 70
)

//go:generate stringer -type=NmDeviceState
type NmDeviceState uint32

const (
	NmDeviceStateUnknown      NmDeviceState = 0
	NmDeviceStateUnmanaged    NmDeviceState = 10
	NmDeviceStateUnavailable  NmDeviceState = 20
	NmDeviceStateDisconnected NmDeviceState = 30
	NmDeviceStatePrepare      NmDeviceState = 40
	NmDeviceStateConfig       NmDeviceState = 50
	NmDeviceStateNeed_auth    NmDeviceState = 60
	NmDeviceStateIp_config    NmDeviceState = 70
	NmDeviceStateIp_check     NmDeviceState = 80
	NmDeviceStateSecondaries  NmDeviceState = 90
	NmDeviceStateActivated    NmDeviceState = 100
	NmDeviceStateDeactivating NmDeviceState = 110
	NmDeviceStateFailed       NmDeviceState = 120
)

//go:generate stringer -type=NmDeviceType
type NmDeviceType uint32

const (
	NmDeviceTypeUnknown    NmDeviceType = 0
	NmDeviceTypeEthernet   NmDeviceType = 1
	NmDeviceTypeWifi       NmDeviceType = 2
	NmDeviceTypeUnused1    NmDeviceType = 3
	NmDeviceTypeUnused2    NmDeviceType = 4
	NmDeviceTypeBt         NmDeviceType = 5
	NmDeviceTypeOlpcMesh   NmDeviceType = 6
	NmDeviceTypeWimax      NmDeviceType = 7
	NmDeviceTypeModem      NmDeviceType = 8
	NmDeviceTypeInfiniband NmDeviceType = 9
	NmDeviceTypeBond       NmDeviceType = 10
	NmDeviceTypeVlan       NmDeviceType = 11
	NmDeviceTypeAdsl       NmDeviceType = 12
	NmDeviceTypeBridge     NmDeviceType = 13
	NmDeviceTypeGeneric    NmDeviceType = 14
	NmDeviceTypeTeam       NmDeviceType = 15
)
