package diskutil

import (
	"encoding/json"
	"errors"
	"strings"
)

// VirtualDriveStat is a struct to get the Virtual Drive Stat of a RAID card.
type VirtualDriveStat struct {
	VirtualDrive   int    `json:"virtual_drive"`
	Name           string `json:"name"`
	Size           string `json:"size"`
	State          string `json:"state"`
	NumberOfDrives int    `json:"number_of_drives"`
	Encryptiontype string `json:"encryption_type"`
	OsPath         string `json:"os_path"`
}

// String() is used to get the print string.
func (v *VirtualDriveStat) String() string {
	data, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

// ToJson() is used to get the json encoded string.
func (v *VirtualDriveStat) ToJson() (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (v *VirtualDriveStat) parseLine(line string) error {
	if strings.HasPrefix(line, keyVdVirtualDrive) {
		parts := strings.SplitN(line, "(", 2)
		if len(parts) != 2 {
			return errors.New("format illegal: " + line)
		}
		virtualDrive, err := parseFiled(parts[0], keyVdVirtualDrive, typeInt)
		if err != nil {
			return err
		}
		v.VirtualDrive = virtualDrive.(int)
	} else if strings.HasPrefix(line, keyVdName) {
		name, err := parseFiled(line, keyVdName, typeString)
		if err != nil {
			return err
		}
		v.Name = name.(string)
	} else if strings.HasPrefix(line, keyVdSize) {
		size, err := parseFiled(line, keyVdSize, typeString)
		if err != nil {
			return err
		}
		v.Size = size.(string)
	} else if strings.HasPrefix(line, keyVdState) {
		state, err := parseFiled(line, keyVdState, typeString)
		if err != nil {
			return err
		}
		v.State = state.(string)
	} else if strings.HasPrefix(line, keyVdNumberOfDrives) {
		numberOfDrives, err := parseFiled(line, keyVdNumberOfDrives, typeInt)
		if err != nil {
			return err
		}
		v.NumberOfDrives = numberOfDrives.(int)
	} else if strings.HasPrefix(line, keyVdEncryptiontype) {
		encryptiontype, err := parseFiled(line, keyVdEncryptiontype, typeString)
		if err != nil {
			return err
		}
		v.Encryptiontype = encryptiontype.(string)
	}
	return nil
}
