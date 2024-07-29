package constants

const (
	// DataBaseDriverName is sqlite3
	DataBaseDriverName = "sqlite3"
	// DataBaseAliasName is default
	DataBaseAliasName = "default"
	//client version
	ClientVersion = "1.0.0"
	//client name
	ClientName = "AppHub-Agent"
)

const (
	//resource.
	/* For ithings resource*/
	ITHINGS_RSC_MMI     = "mmi"
	ITHINGS_RSC_COMM    = "comm"
	ITHINGS_RSC_ITSMGR  = "itsmgr"
	ITHINGS_RSC_ITS     = "ithings"
	ITHINGS_RSC_DTWIN   = "devicetwin"
	ITHINGS_RSC_DROUTER = "itsrouter"
	/* For m2mcore resource*/
	ITHINGS_OBJECT_MAPPER           = "mapper"
	ITHINGS_OBJECT_DEVICE_DATA      = "device_data"
	ITHINGS_OBJECT_DEVICE_EVENT     = "device_event"
	ITHINGS_OBJECT_EVENT_RECOVER    = "event_recover"
	ITHINGS_OBJECT_DEVICE_STATUS    = "device_status"
	ITHINGS_OBJECT_DEVICE_SPEC_META = "device_spec_meta"
	ITHINGS_OBJECT_DEVICE_BRIEF     = "device_brief"
	ITHINGS_OBJECT_DESIRED_TWINS    = "desired_twins"

	//operation.
	ITHINGS_OPERATION_REGISTER     = "register"
	ITHINGS_OPERATION_REPORT       = "report"
	ITHINGS_OPERATION_RESUME       = "resume"
	ITHINGS_OPERATION_REPLY        = "reply"
	ITHINGS_OPERATION_FETCH        = "fetch"
	ITHINGS_OPERATION_LIFE_CONTROL = "life_control"
	ITHINGS_OPERATION_SET_PROPERTY = "set_property"

	//Response Code;
	ITHINGS_RSP_MAPPER_NOT_FOUND       = "4.04"
	ITHINGS_RSP_OPERATION_NOT_FOUND    = "4.05"
	ITHINGS_RSP_MAPPER_NOT_REGISTER    = "4.06"
	ITHINGS_RSP_MAPPER_REGISTER_FAILED = "4.07"
	ITHINGS_RSP_INVALID_MQTT_MESSAGE   = "4.08"
	ITHINGS_RSP_SUCCEED                = "200"
	IRespCodeOk                        = "200"
	IRespOkString                      = "sucess"
	IRespCodeInvalidMsg                = "201"
	IRespInvalidMsgString              = "invalid message format"
	IRespCodeInternalError             = "202"
	IRespInternalErrString             = "server internal error"
	IRespCodeError                     = "205"
	IRespCodeNoSuchDevice              = "400"
	IRespNoSuchDeviceStr               = "No such device"
	IRespCodeOPNotFound                = "405"
	IRespOPNotFoundStr                 = "operation not found"
)

const (
	DEVICE_STATUS_ONLINE  = "online"
	DEVICE_STATUS_OFFLINE = "offline"

	//state.
	DEVICE_STATE_STARTED = "started"
	DEVICE_STATE_STOPPED = "stopped"
)

const (
	DEVICE_LIFE_CREATE = "create"
	DEVICE_LIFE_START  = "start"
	DEVICE_LIFE_STOP   = "stop"
	DEVICE_LIFE_UPDATE = "update"
	DEVICE_LIFE_DELETE = "delete"
)

var (
	EndpointID string
)
