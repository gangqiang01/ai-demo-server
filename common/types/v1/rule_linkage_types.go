package v1

const (
	//alert Info level
	AlertInfoLevel  int64 = 0
	AlertErrorLevel int64 = 1

	//event relation
	EventMoreRelation     int32 = 0
	EventLessRelation     int32 = 1
	EventEqualRelation    int32 = 2
	EventNotEqualRelation int32 = 3

	//ruleLinkage status
	RuleLinkageStatusEnable  string = "enable"
	RuleLinkageStatusDisable string = "disable"

	//ruleLinkage trigger type
	EVENTTRIGGER           string = "event"
	PARTICULAREVENTTRIGGER string = "particularEvent"

	//ruleLinkage filter type
	FilterEffectiveTime string = "effectiveTime"

	//ruleLinkage action type
	REPORTACTION             string = "alert"
	PROPERTYACTION           string = "property"
	PARTICULARPROPERTYACTION string = "particularProperty"
	CONTROLACTION            string = "control"
	PARTICULARCONTROLACTION  string = "particularControl"
	REBOOT                   string = "reboot"
	SHUTDOWN                 string = "shutdown"

	//ruleLinkage log status
	//running: 0, error: 1, success: 2
	RULELOGSTATUSRUNNING int32 = 0
	RULELOGSTATUSERROR   int32 = 1
	RULELOGSTATUSSUCCESS int32 = 2

	//rulelinkage handle action type
	RuleLinkageHandleTypeEvent   string = "event"
	RuleLinkageHandleTypeRecover string = "recover"
)

// ruleLinkage
type RuleLinkageEvent struct {
	ServiceName string `form:"serviceName" json:"serviceName,omitempty"`
	EventName   string `form:"eventName" json:"eventName,omitempty"`
}

type RuleLinkageParticularEvent struct {
	DeviceId              string `form:"deviceId" json:"deviceId,omitempty"`
	ParticularServiceName string `form:"particularServiceName" json:"particularServiceName,omitempty"`
	ParticularEventName   string `form:"particularEventName" json:"particularEventName,omitempty"`
}

type RuleLinkageProperty struct {
	ServiceName   string `form:"serviceName" json:"serviceName,omitempty"`
	PropertyName  string `form:"propertyName" json:"propertyName,omitempty"`
	PropertyValue string `form:"propertyValue" json:"propertyValue,omitempty"`
}

type RuleLinkageParticularProperty struct {
	DeviceId                string `form:"deviceId" json:"deviceId,omitempty"`
	ParticularServiceName   string `form:"particularServiceName" json:"particularServiceName,omitempty"`
	ParticularPropertyName  string `form:"particularPropertyName" json:"particularPropertyName,omitempty"`
	ParticularPropertyValue string `form:"particularPropertyValue" json:"particularPropertyValue,omitempty"`
}

type RuleLinkageAlert struct {
	AlertId *int64 `form:"alertId" json:"alertId,omitempty"`
}

type RuleLinkageControl struct {
	ControlType string `form:"controlType" json:"controlType,omitempty"`
}

type RuleLinkageParticularControl struct {
	EdgeId                string `form:"edgeId" json:"edgeId,omitempty"`
	ParticularControlType string `form:"particularControlType" json:"particularControlType,omitempty"`
}

type RuleLinkageTrigger struct {
	Type              string `form:"type" json:"type"`
	DeviceModel       string `form:"deviceModel" json:"deviceModel,omitempty"`
	*RuleLinkageEvent `json:",inline"`
	//particular property
	*RuleLinkageParticularEvent `json:",inline"`
}

type RuleLinkageFilter struct {
	Type string `form:"type" json:"type"`
	//effective time
	*EffectiveTime `json:",inline,omitempty"`
}

type EffectiveTime struct {
	StartTime string `form:"startTime" json:"startTime"`
	EndTime   string `form:"endTime" json:"endTime"`
}

type RuleLinkageAction struct {
	Type        string `form:"type" json:"type"`
	DeviceModel string `form:"deviceModel" json:"deviceModel,omitempty"`
	//alert
	*RuleLinkageAlert `json:",inline,omitempty"`
	//property
	*RuleLinkageProperty `json:",inline,omitempty"`
	//particular property
	*RuleLinkageParticularProperty `json:",inline,omitempty"`
	//control
	*RuleLinkageControl `json:",inline,omitempty"`
	//particular control
	*RuleLinkageParticularControl `json:",inline,omitempty"`
}
