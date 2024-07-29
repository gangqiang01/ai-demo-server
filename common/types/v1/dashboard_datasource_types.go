package v1

const (
	DashboardQueryTypeTimeseries string = "timeserie"
	DashboardQueryTypeTable      string = "table"
)

type Dashboard_search_request struct {
	DataType string `form:"dataType" json:"dataType"`
	Type     string `form:"type" json:"type"`
	Device   string `form:"device" json:"device"`
	Ithing   string `form:"ithing" json:"ithing"`
	Module   string `form:"module" json:"module"`
	Property string `form:"property" json:"property"`
}
type Range struct {
	From string `from:"from" json:"from"`
	To   string `from:"to" json:"to"`
}

type Dashboard_query_request struct {
	Range         *Range    `from:"range" json:"range"`
	Targets       []*Target `form:"targets" json:"targets"`
	MaxDataPoints int64     `form:"maxDataPoints" json:"maxDataPoints"`
}

type Target struct {
	DataType    string `form:"dataType" json:"dataType"`
	Scene       string `form:"scene" json:"scene"`
	DeviceId    string `form:"deviceId" json:"deviceId"`
	Ithing      string `form:"ithing" json:"ithing"`
	Module      string `form:"module" json:"module"`
	Property    string `form:"property" json:"property"`
	Target      string `form:"target" json:"target"`
	DisplayName string `form:"displayName" json:"displayName"`
	RefId       string `form:"refId" json:"refId"`
	Type        string `form:"type" json:"type"`
}

type Dashboard_search_responce struct {
	Text  string `form:"text" json:"text,omitempty"`
	Value string `form:"value" json:"value,omitempty"`
}
type Dashboard_query_timeserie_responce struct {
	Target     string          `form:"target" json:"target,omitempty"`
	Datapoints [][]interface{} `form:"datapoints" json:"datapoints,omitempty"`
}

type Column struct {
	Text string `form:"text" json:"text"`
	Type string `form:"type" json:"type"`
}

type Dashboard_query_table_responce struct {
	Columns []*Column       `form:"columns" json:"columns"`
	Rows    [][]interface{} `form:"rows" json:"rows"`
	Type    string          `form:"type" json:"type"`
}

//scene

var (
	//overview series
	Dashboard_online_device         string = "Online Device"
	Dashboard_offline_device        string = "Offline Device"
	Dashboard_android_device        string = "Android Device"
	Dashboard_linux_device          string = "Linux Device"
	Dashboard_windows_device        string = "Windows Device"
	Dashboard_online_subDevice      string = "Online Sub Device"
	Dashboard_offline_subDevice     string = "Offline Sub Device"
	Dashboard_normal_device         string = "Normal Device"
	Dashboard_error_device          string = "Error Device"
	Dashboard_warning_device        string = "Warning Device"
	Dashboard_normal_subDevice      string = "Normal Sub Device"
	Dashboard_error_subDevice       string = "Error Sub Device"
	Dashboard_warning_subDevice     string = "Warning Sub Device"
	Dashboard_error_alert           string = "Error Alert"
	Dashboard_error_alert_handled   string = "Error Handled Alarm"
	Dashboard_warning_alert         string = "Warning Alert"
	Dashboard_warning_alert_handled string = "Warning Handled Alarm"
	//overview table
	Dashboard_error_alert_unhandled   string = "Unhandled Error Alarm"
	Dashboard_warning_alert_unhandled string = "Unhandled Warning Alarm"

	//information series
	Dashboard_total_memory_device      string = "Device Total Memory"
	Dashboard_free_memory_device       string = "Device Free Memory"
	Dashboard_usage_cpu_device         string = "Device CPU Usage"
	Dashboard_temp_cpu_device          string = "Device CPU temperature"
	Dashboard_total_storage_device     string = "Device Total Storage"
	Dashboard_free_storage_device      string = "Device Free Storage"
	Dashboard_battery_available_device string = "Device Battery Available"

	//information table
	Dashboard_monitor_app_device     string = "Device Monitor App"
	Dashboard_usb_device             string = "Device USB"
	Dashboard_monitor_process_device string = "Device Monitor Process"
	Dashboard_monitor_docker_device  string = "Device Monitor Docker"
)

var Dashboard_overview_scene_timeserie_resp = []*Dashboard_search_responce{
	//timeseries
	{Value: Dashboard_online_device, Text: Dashboard_online_device},
	{Value: Dashboard_offline_device, Text: Dashboard_offline_device},
	{Value: Dashboard_android_device, Text: Dashboard_android_device},
	{Value: Dashboard_linux_device, Text: Dashboard_linux_device},
	{Value: Dashboard_windows_device, Text: Dashboard_windows_device},

	{Value: Dashboard_online_subDevice, Text: Dashboard_online_subDevice},
	{Value: Dashboard_offline_subDevice, Text: Dashboard_offline_subDevice},
	{Value: Dashboard_normal_device, Text: Dashboard_normal_device},
	{Value: Dashboard_error_device, Text: Dashboard_error_device},
	{Value: Dashboard_warning_device, Text: Dashboard_warning_device},
	{Value: Dashboard_normal_subDevice, Text: Dashboard_normal_subDevice},
	{Value: Dashboard_error_subDevice, Text: Dashboard_error_subDevice},
	{Value: Dashboard_warning_subDevice, Text: Dashboard_warning_subDevice},
	{Value: Dashboard_error_alert, Text: Dashboard_error_alert},
	{Value: Dashboard_error_alert_handled, Text: Dashboard_error_alert_handled},
	{Value: Dashboard_warning_alert, Text: Dashboard_warning_alert},
	{Value: Dashboard_warning_alert_handled, Text: Dashboard_warning_alert_handled},
}
var Dashboard_overview_scene_table_resp = []*Dashboard_search_responce{
	//table
	{Value: Dashboard_error_alert_unhandled, Text: Dashboard_error_alert_unhandled},
	{Value: Dashboard_warning_alert_unhandled, Text: Dashboard_warning_alert_unhandled},
}

var Dashboard_information_scene_timeserie_resp = []*Dashboard_search_responce{
	//timeseries
	{Value: Dashboard_total_memory_device, Text: Dashboard_total_memory_device},
	{Value: Dashboard_free_memory_device, Text: Dashboard_free_memory_device},
	{Value: Dashboard_usage_cpu_device, Text: Dashboard_usage_cpu_device},
	{Value: Dashboard_temp_cpu_device, Text: Dashboard_temp_cpu_device},
	{Value: Dashboard_total_storage_device, Text: Dashboard_total_storage_device},
	{Value: Dashboard_free_storage_device, Text: Dashboard_free_storage_device},
	{Value: Dashboard_battery_available_device, Text: Dashboard_battery_available_device},
}

var Dashboard_information_scene_table_resp = []*Dashboard_search_responce{
	//table
	{Value: Dashboard_monitor_app_device, Text: Dashboard_monitor_app_device},
	{Value: Dashboard_usb_device, Text: Dashboard_usb_device},
	{Value: Dashboard_monitor_process_device, Text: Dashboard_monitor_process_device},
	{Value: Dashboard_monitor_docker_device, Text: Dashboard_monitor_docker_device},
}
