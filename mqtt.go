package main

import (
	"encoding/json"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQTTDiscovery - Simple struct to aid the creation/publishing of Home-Assistant MQTT Discovery
type MQTTDiscovery struct {
	DiscoveryPrefix string
	Component       string
	NodeID          string
	ObjectID        string

	AuxCommandTopic          string `json:"aux_command_topic,omitempty"`
	AuxStateTemplate         string `json:"aux_state_template,omitempty"`
	AuxStateTopic            string `json:"aux_state_topic,omitempty"`
	AvailabilityTopic        string `json:"availability_topic,omitempty"`
	AwayModeCommandTopic     string `json:"away_mode_command_topic,omitempty"`
	AwayModeStateTemplate    string `json:"away_mode_state_template,omitempty"`
	AwayModeStateTopic       string `json:"away_mode_state_topic,omitempty"`
	BrightnessCommandTopic   string `json:"brightness_command_topic,omitempty"`
	BrightnessScale          string `json:"brightness_scale,omitempty"`
	BrightnessStateTopic     string `json:"brightness_state_topic,omitempty"`
	BrightnessValueTemplate  string `json:"brightness_value_template,omitempty"`
	ColorTempCommandTopic    string `json:"color_temp_command_topic,omitempty"`
	ColorTempStateTopic      string `json:"color_temp_state_topic,omitempty"`
	ColorTempValueTemplate   string `json:"color_temp_value_template,omitempty"`
	CommandTopic             string `json:"command_topic,omitempty"`
	CurrentTemperatureTopic  string `json:"current_temperature_topic,omitempty"`
	DeviceClass              string `json:"device_class,omitempty"`
	EffectCommandTopic       string `json:"effect_command_topic,omitempty"`
	EffectList               string `json:"effect_list,omitempty"`
	EffectStateTopic         string `json:"effect_state_topic,omitempty"`
	EffectValueTemplate      string `json:"effect_value_template,omitempty"`
	ExpireAfter              string `json:"expire_after,omitempty"`
	FanModeCommandTopic      string `json:"fan_mode_command_topic,omitempty"`
	FanModeStateTemplate     string `json:"fan_mode_state_template,omitempty"`
	FanModeStateTopic        string `json:"fan_mode_state_topic,omitempty"`
	ForceUpdate              string `json:"force_update,omitempty"`
	HoldCommandTopic         string `json:"hold_command_topic,omitempty"`
	HoldStateTemplate        string `json:"hold_state_template,omitempty"`
	HoldStateTopic           string `json:"hold_state_topic,omitempty"`
	Icon                     string `json:"icon,omitempty"`
	Initial                  string `json:"initial,omitempty"`
	JSONAttributes           string `json:"json_attributes,omitempty"`
	MaxTemp                  string `json:"max_temp,omitempty"`
	MinTemp                  string `json:"min_temp,omitempty"`
	ModeCommandTopic         string `json:"mode_command_topic,omitempty"`
	ModeStateTemplate        string `json:"mode_state_template,omitempty"`
	ModeStateTopic           string `json:"mode_state_topic,omitempty"`
	Name                     string `json:"name,omitempty"`
	OnCommandType            string `json:"on_command_type,omitempty"`
	Optimistic               string `json:"optimistic,omitempty"`
	OscillationCommandTopic  string `json:"oscillation_command_topic,omitempty"`
	OscillationStateTopic    string `json:"oscillation_state_topic,omitempty"`
	OscillationValueTemplate string `json:"oscillation_value_template,omitempty"`
	PayloadArmAway           string `json:"payload_arm_away,omitempty"`
	PayloadArmHome           string `json:"payload_arm_home,omitempty"`
	PayloadAvailable         string `json:"payload_available,omitempty"`
	PayloadClose             string `json:"payload_close,omitempty"`
	PayloadDisarm            string `json:"payload_disarm,omitempty"`
	PayloadHighSpeed         string `json:"payload_high_speed,omitempty"`
	PayloadLock              string `json:"payload_lock,omitempty"`
	PayloadLowSpeed          string `json:"payload_low_speed,omitempty"`
	PayloadMediumSpeed       string `json:"payload_medium_speed,omitempty"`
	PayloadNotAvailable      string `json:"payload_not_available,omitempty"`
	PayloadOff               string `json:"payload_off,omitempty"`
	PayloadOn                string `json:"payload_on,omitempty"`
	PayloadOpen              string `json:"payload_open,omitempty"`
	PayloadOscillationOff    string `json:"payload_oscillation_off,omitempty"`
	PayloadOscillationOn     string `json:"payload_oscillation_on,omitempty"`
	PayloadStop              string `json:"payload_stop,omitempty"`
	PayloadUnlock            string `json:"payload_unlock,omitempty"`
	PowerCommandTopic        string `json:"power_command_topic,omitempty"`
	Retain                   string `json:"retain,omitempty"`
	RgbCommandTemplate       string `json:"rgb_command_template,omitempty"`
	RgbCommandTopic          string `json:"rgb_command_topic,omitempty"`
	RgbStateTopic            string `json:"rgb_state_topic,omitempty"`
	RgbValueTemplate         string `json:"rgb_value_template,omitempty"`
	SendIfOff                string `json:"send_if_off,omitempty"`
	SetPositionTemplate      string `json:"set_position_template,omitempty"`
	SetPositionTopic         string `json:"set_position_topic,omitempty"`
	SpeedCommandTopic        string `json:"speed_command_topic,omitempty"`
	SpeedStateTopic          string `json:"speed_state_topic,omitempty"`
	SpeedValueTemplate       string `json:"speed_value_template,omitempty"`
	Speeds                   string `json:"speeds,omitempty"`
	StateClosed              string `json:"state_closed,omitempty"`
	StateOff                 string `json:"state_off,omitempty"`
	StateOn                  string `json:"state_on,omitempty"`
	StateOpen                string `json:"state_open,omitempty"`
	StateTopic               string `json:"state_topic,omitempty"`
	StateValueTemplate       string `json:"state_value_template,omitempty"`
	SwingModeCommandTopic    string `json:"swing_mode_command_topic,omitempty"`
	SwingModeStateTemplate   string `json:"swing_mode_state_template,omitempty"`
	SwingModeStateTopic      string `json:"swing_mode_state_topic,omitempty"`
	TemperatureCommandTopic  string `json:"temperature_command_topic,omitempty"`
	TemperatureStateTemplate string `json:"temperature_state_template,omitempty"`
	TemperatureStateTopic    string `json:"temperature_state_topic,omitempty"`
	TiltClosedValue          string `json:"tilt_closed_value,omitempty"`
	TiltCommandTopic         string `json:"tilt_command_topic,omitempty"`
	TiltInvertState          string `json:"tilt_invert_state,omitempty"`
	TiltMax                  string `json:"tilt_max,omitempty"`
	TiltMin                  string `json:"tilt_min,omitempty"`
	TiltOpenedValue          string `json:"tilt_opened_value,omitempty"`
	TiltStatusOptimistic     string `json:"tilt_status_optimistic,omitempty"`
	TiltStatusTopic          string `json:"tilt_status_topic,omitempty"`
	Topic                    string `json:"topic,omitempty"`
	UniqueID                 string `json:"unique_id,omitempty"`
	UnitOfMeasurement        string `json:"unit_of_measurement,omitempty"`
	ValueTemplate            string `json:"value_template,omitempty"`
	WhiteValueCommandTopic   string `json:"white_value_command_topic,omitempty"`
	WhiteValueStateTopic     string `json:"white_value_state_topic,omitempty"`
	WhiteValueTemplate       string `json:"white_value_template,omitempty"`
	XyCommandTopic           string `json:"xy_command_topic,omitempty"`
	XyStateTopic             string `json:"xy_state_topic,omitempty"`
	XyValueTemplate          string `json:"xy_value_template,omitempty"`
}

func (mqd *MQTTDiscovery) topic() string {
	return fmt.Sprintf("%s/%s/%s/%s/config", mqd.DiscoveryPrefix, mqd.Component, mqd.NodeID, mqd.ObjectID)
}

// PublishDiscovery - Publish discovery on a given topic
func (mqd *MQTTDiscovery) PublishDiscovery(client mqtt.Client) {
	topic := mqd.topic()
	payloadBytes, _ := json.Marshal(mqd)
	payload := string(payloadBytes)
	retain := true
	if token := client.Publish(topic, 0, retain, payload); token.Wait() && token.Error() != nil {
		log.Printf("Publishing Error: %s", token.Error())
	}
	log.Print(fmt.Sprintf("Publishing - Topic:%s ; Payload: %s", topic, payload))
}

type newMQTTClientOptsFunc func() *mqtt.ClientOptions
type newMQTTClientFunc func(*mqtt.ClientOptions) mqtt.Client

// MQTTFuncWrapper - Wraps the functions needed to create a new MQTT client.
type MQTTFuncWrapper struct {
	clientOptsFunc newMQTTClientOptsFunc
	clientFunc     newMQTTClientFunc
}

// NewMQTTFuncWrapper - Returns a fancy new wrapper for the mqtt creation functions.
func NewMQTTFuncWrapper() *MQTTFuncWrapper {
	return &MQTTFuncWrapper{
		clientOptsFunc: mqtt.NewClientOptions,
		clientFunc:     mqtt.NewClient,
	}
}
