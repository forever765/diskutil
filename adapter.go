package diskutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// AdapterStat is a struct to get the Adapter Stat of a RAID card.
// AdapterStat has VirtualDriveStats and PhysicalDriveStats in itself.
type AdapterStat struct {
	AdapterId          int                 `json:"adapter_id"`
	VirtualDriveStats  []VirtualDriveStat  `json:"virtual_drive_stats"`
	PhysicalDriveStats []PhysicalDriveStat `json:"physical_drive_stats"`
}

// String() is used to get the print string.
func (a *AdapterStat) String() string {
	data, err := json.Marshal(a)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

// ToJson() is used to get the json encoded string.
func (a *AdapterStat) ToJson() (string, error) {
	data, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (a *AdapterStat) parseMegaRaidVdInfo(info string, common string, adapterId string) error {
	if info == "" {
		return errors.New("mageRaid vd info nil")
	}

	vds := make([]VirtualDriveStat, 0)

	parts := strings.Split(info, keyVdVirtualDrive)
	// vdinfo为megacli响应的完整vd信息，遍历所有VD
	for _, vdinfo := range parts {
		if strings.Contains(vdinfo, keyVdTargetId) {
			vdinfo = keyVdVirtualDrive + vdinfo
			vd := VirtualDriveStat{}
			lines := strings.Split(vdinfo, "\n")
			for _, line := range lines {
				err := vd.parseLine(line)
				if err != nil {
					return err
				}
			}
			osPath := getMegaRaidAdapterInfo(common, vd.VirtualDrive, adapterId)
			vd.OsPath = osPath
			vds = append(vds, vd)
		}
	}

	a.VirtualDriveStats = vds
	return nil
}

func (a *AdapterStat) getMegaRaidVdInfo(command string) error {
	adapterId := strconv.Itoa(a.AdapterId)
	args := "-ldinfo -lall -a" + adapterId + " -NoLog"

	output, err := execCmd(command, args)
	if err != nil {
		return err
	}
	parts := strings.SplitN(output, keyExitResult, 2)
	if len(parts) != 2 {
		return errors.New("megaCli output illegal")
	}
	result := strings.TrimSpace(parts[1])
	if result != "0x00" {
		return errors.New("megaCli return error: " + result)
	}

	err = a.parseMegaRaidVdInfo(output, command, adapterId)
	if err != nil {
		return err
	}
	return nil
}

func (a *AdapterStat) parseMegaRaidPdInfo(info string) error {
	if info == "" {
		return errors.New("mageRaid pd info nil")
	}

	pds := make([]PhysicalDriveStat, 0)

	parts := strings.Split(info, keyPdEnclosureDeviceId)
	for _, pdinfo := range parts {
		if strings.Contains(pdinfo, keyPdSlotNumber) {
			pdinfo = keyPdEnclosureDeviceId + pdinfo
			pd := PhysicalDriveStat{}
			lines := strings.Split(pdinfo, "\n")
			for _, line := range lines {
				err := pd.parseLine(line)
				if err != nil {
					return err
				}
			}
			pds = append(pds, pd)
		}
	}

	a.PhysicalDriveStats = pds
	return nil
}

func (a *AdapterStat) getMegaRaidPdInfo(command string) error {
	args := "-pdlist -a" + strconv.Itoa(a.AdapterId) + " -NoLog"

	output, err := execCmd(command, args)
	if err != nil {
		return err
	}
	parts := strings.SplitN(output, keyExitResult, 2)
	if len(parts) != 2 {
		return errors.New("megaCli output illegal")
	}
	result := strings.TrimSpace(parts[1])
	if result != "0x00" {
		return errors.New("megaCli return error: " + result)
	}

	err = a.parseMegaRaidPdInfo(output)
	if err != nil {
		return err
	}
	return nil
}

// 提取HBAPCIInfo
func parseHBAPCIInfo(output string) string {
	busprefix := "0000"
	var busid, devid, functionid, pcipath string

	busRegex := regexp.MustCompile(`^Bus Number.*:.*$`)
	deviceRegex := regexp.MustCompile(`^Device Number.*:.*$`)
	functionRegex := regexp.MustCompile(`^Function Number.*:.*$`)

	lines := strings.Split(output, "\n")

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if busRegex.MatchString(trimmedLine) {
			parts := strings.Split(trimmedLine, ":")
			busid = fmt.Sprintf("%02s", strings.TrimSpace(parts[1]))
		}
		if deviceRegex.MatchString(trimmedLine) {
			parts := strings.Split(trimmedLine, ":")
			devid = fmt.Sprintf("%02s", strings.TrimSpace(parts[1]))
		}
		if functionRegex.MatchString(trimmedLine) {
			parts := strings.Split(trimmedLine, ":")
			functionid = fmt.Sprintf("%01s", strings.TrimSpace(parts[1]))
		}
	}

	if busid != "" {
		pcipath = fmt.Sprintf("%s:%s:%s.%s", busprefix, busid, devid, functionid)
		//fmt.Println("Array PCI path :", pcipath)
		return pcipath
	}

	return ""
}

// 解析raid卡pcie信息
func (a *AdapterStat) parseMegaRaidAdapterInfo(info string) error {
	if info == "" {
		return errors.New("mageRaid pcie info nil")
	}

	pds := make([]PhysicalDriveStat, 0)

	parts := strings.Split(info, keyPdEnclosureDeviceId)
	for _, pdinfo := range parts {
		if strings.Contains(pdinfo, keyPdSlotNumber) {
			pdinfo = keyPdEnclosureDeviceId + pdinfo
			pd := PhysicalDriveStat{}
			lines := strings.Split(pdinfo, "\n")
			for _, line := range lines {
				err := pd.parseLine(line)
				if err != nil {
					return err
				}
			}
			pds = append(pds, pd)
		}
	}

	a.PhysicalDriveStats = pds
	return nil
}

// GetMegaRaidAdapterInfo 获取vd对应的系统盘符
func getMegaRaidAdapterInfo(command string, virtualDriveId int, adapterId string) string {
	var (
		args   = "-AdpGetPciInfo -a" + adapterId + " -NoLog"
		osPath = "N/A"
	)

	output, err := execCmd(command, args)
	if err != nil {
		return osPath
	}
	parts := strings.SplitN(output, keyExitResult, 2)
	if len(parts) != 2 {
		fmt.Printf("megaCli output illegal")
		return osPath
	}
	result := strings.TrimSpace(parts[1])
	if result != "0x00" {
		fmt.Printf("megaCli return error: " + result)
		return osPath
	}

	pciPath := parseHBAPCIInfo(output)

	//vDId := 0
	if pciPath != "" {
		diskPrefix := "/dev/disk/by-path/pci-" + pciPath + "-scsi-0:"

		// RAID disks are usually with a channel of '2', JBOD disks with a channel of '0'
		for j := 1; j < 8; j++ {
			diskPath := diskPrefix + fmt.Sprintf("%d:%d:0", j, virtualDriveId)
			//fmt.Println("Looking for DISKpath : " + diskpath)
			if _, err := os.Stat(diskPath); err == nil {
				if realpath, err := filepath.EvalSymlinks(diskPath); err == nil {
					osPath = realpath
					//fmt.Println("Found DISK match: " + diskpath + " -> " + tempResult)
					break
				}
			}
		}
	}
	//fmt.Println("got real os path: ", osPath)
	return osPath
}
