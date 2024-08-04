package diskutil

import (
	"encoding/json"
	"strings"
)

// PhysicalDriveStat is a struct to get the Physical Drive Stat of a RAID card.
type PhysicalDriveStat struct {
	EnclosureDeviceId      int    `json:"enclosure_device_id"`
	DeviceId               int    `json:"device_id"`
	SlotNumber             int    `json:"slot_number"`
	MediaErrorCount        int    `json:"media_error_count"`
	OtherErrorCount        int    `json:"other_error_count"`
	PredictiveFailureCount int    `json:"predictive_failure_count"`
	PdMediaType            string `json:"pd_media_type"`
	PdType                 string `json:"pd_type"`
	RawSize                string `json:"raw_size"`
	FirmwareState          string `json:"firmware_state"`
	Brand                  string `json:"brand"`
	Model                  string `json:"model"`
	SerialNumber           string `json:"serial_number"`
	DriveTemperature       string `json:"drive_emperature"`
}

// String() is used to get the print string.
func (p *PhysicalDriveStat) String() string {
	data, err := json.Marshal(p)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

// ToJson() is used to get the json encoded string.
func (p *PhysicalDriveStat) ToJson() (string, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (p *PhysicalDriveStat) parseLine(line string) error {
	if strings.HasPrefix(line, keyPdEnclosureDeviceId) {
		EnclosureDeviceId, err := parseFiled(line, keyPdEnclosureDeviceId, typeInt)
		if err != nil {
			return err
		}
		p.EnclosureDeviceId = EnclosureDeviceId.(int)
	} else if strings.HasPrefix(line, keyPdDeviceId) {
		deviceId, err := parseFiled(line, keyPdDeviceId, typeInt)
		if err != nil {
			return err
		}
		p.DeviceId = deviceId.(int)
	} else if strings.HasPrefix(line, keyPdSlotNumber) {
		slotNumber, err := parseFiled(line, keyPdSlotNumber, typeInt)
		if err != nil {
			return err
		}
		p.SlotNumber = slotNumber.(int)
	} else if strings.HasPrefix(line, keyPdMediaErrorCount) {
		mediaErrorCount, err := parseFiled(line, keyPdMediaErrorCount, typeInt)
		if err != nil {
			return err
		}
		p.MediaErrorCount = mediaErrorCount.(int)
	} else if strings.HasPrefix(line, keyPdOtherErrorCount) {
		otherErrorCount, err := parseFiled(line, keyPdOtherErrorCount, typeInt)
		if err != nil {
			return err
		}
		p.OtherErrorCount = otherErrorCount.(int)
	} else if strings.HasPrefix(line, keyPdPredictiveFailureCount) {
		predictiveFailureCount, err := parseFiled(line, keyPdPredictiveFailureCount, typeInt)
		if err != nil {
			return err
		}
		p.PredictiveFailureCount = predictiveFailureCount.(int)
	} else if strings.HasPrefix(line, keyPdMediaType) {
		pdMediaType, err := parseFiled(line, keyPdMediaType, typeString)
		if err != nil {
			return err
		}
		p.PdMediaType = pdMediaType.(string)
	} else if strings.HasPrefix(line, keyPdPdtype) {
		pdType, err := parseFiled(line, keyPdPdtype, typeString)
		if err != nil {
			return err
		}
		p.PdType = pdType.(string)
	} else if strings.HasPrefix(line, keyPdRawSize) {
		rawSize, err := parseFiled(line, keyPdRawSize, typeString)
		if err != nil {
			return err
		}
		p.RawSize = rawSize.(string)
	} else if strings.HasPrefix(line, keyPdFirmwareState) {
		firmwareState, err := parseFiled(line, keyPdFirmwareState, typeString)
		if err != nil {
			return err
		}
		p.FirmwareState = firmwareState.(string)
	} else if strings.HasPrefix(line, keyPdInquiryData) {
		inquiryData, err := parseFiled(line, keyPdInquiryData, typeString)
		if err != nil {
			return err
		}
		inquiryStr := inquiryData.(string)
		parts := strings.Fields(inquiryStr)
		switch len(parts) {
		case 4:
			p.SerialNumber = parts[3]
			p.Model = parts[1] + " " + parts[2]
			p.Brand = parts[0]
		case 3:
			p.SerialNumber = parts[2]
			p.Model = parts[1]
			p.Brand = parts[0]
		case 2:
			p.SerialNumber = parts[1]
			p.Model = parts[0]
		case 1:
			p.SerialNumber = parts[0]
		default:
			p.SerialNumber = "unknown"
		}
	} else if strings.HasPrefix(line, keyPdDriveTemperature) {
		driveTemperature, err := parseFiled(line, keyPdDriveTemperature, typeString)
		if err != nil {
			return err
		}
		p.DriveTemperature = driveTemperature.(string)
	}
	return nil
}
