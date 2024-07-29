package global

import (
	"time"

	"gorm.io/gorm"
)

const (
	//beehive modules
	IMODULE_WEBSERVER = "webserver"
	IMODULE_CORE      = "core"
	IMODULE_TRANSPORT = "transport"

	//ithings mqtt
	DefaultMqttRxQueueSize = int(4096)
	DefaultMqttTxQueueSize = int(4096)

	//Core
	DefaultMaxGoRoutines = int(1024)

	//Life Status.
	DeviceStatusInactive = "inactive"
	DeviceStatusActive   = "active"
	DeviceStatusOnline   = "online"
	DeviceStatusOffline  = "offline"

	//life control
	DeviceCreate = int(0)
	DeviceStart  = int(1)
	DeviceStop   = int(2)
	DeviceDelete = int(3)
	DeviceUpdate = int(4)

	DefaultEdgeMaxResponseTime    = 5 * time.Second
	DefaultLifeTimeOfDesiredValue = 30 * 1000 * time.Millisecond
	//device states
	DeviceStateStarted = "started"
	DeviceStateStopped = "stopped"

	//Ithings response error code
	IRespCodeOk            = "200"
	IRespOkString          = "sucess"
	IRespCodeInvalidMsg    = "201"
	IRespInvalidMsgString  = "invalid message format"
	IRespCodeInternalError = "202"
	IRespInternalErrString = "server internal error"
	IRespCodeError         = "205"
	IRespCodeNoSuchDevice  = "400"
	IRespNoSuchDeviceStr   = "No such device"
	IRespCodeOPNotFound    = "405"
	IRespOPNotFoundStr     = "operation not found"

	//system config type
	//email
	EmailType = "email"
	//alarm email
	NotificationEmailType = "notificationEmail"

	//system config type
	INFLUXDB = "influxdb"

	EmailAddress  = "smtp.qq.com"
	EmailPort     = 465
	EmailUsername = "james.dai@qq.com"
	EmailPassword = "waasgigffvjabjif"
)

var (
	// gorm DB
	DBAccess *gorm.DB
)
