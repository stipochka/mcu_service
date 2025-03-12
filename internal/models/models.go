package models

type Record struct {
	ID            int        `json:"id" yaml:"id"`
	Status        string     `json:"status" yaml:"status"` // ON, OFF
	PollingPeriod int        `json:"polling_period" yaml:"polling_period"`
	Temperature   int        `json:"temperature"` // int type to avoid issues with float values
	LampStatus    string     `json:"lamp_status"` // ON, OFF
	Voltage       int        `json:"voltage"`     // same reason as Temperature
	Thresholds    Thresholds `json:"thresholds"`
}

type Thresholds struct {
	Temperature []SensorData `json:"temperature"`
	Humidity    []SensorData `json:"humidity"`
	Voltage     []SensorData `json:"voltage"`
}

type SensorData struct {
	Value         int `json:"value"`
	PollingPeriod int `json:"polling_period"`
}
