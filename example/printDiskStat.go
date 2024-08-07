package main

import (
	"flag"
	"fmt"
	"github.com/buaazp/diskutil"
	"os"
	"strings"
	"unicode"
)

var (
	megaPath     string
	adapterCount int
)

func init() {
	flag.StringVar(&megaPath, "mega-path", "/opt/MegaRAID/MegaCli/MegaCli64", "megaCli binary path")
	flag.IntVar(&adapterCount, "adapter-count", 1, "adapter count in your server")
}

func keepUppercaseLetters(input string) string {
	var result []rune
	for _, char := range input {
		if unicode.IsUpper(char) {
			result = append(result, char)
		}
	}
	return string(result)
}

func main() {
	flag.Parse()
	ds, err := diskutil.NewDiskStatus(megaPath, adapterCount)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DiskStatus New error: %v\n", err)
		return
	}

	err = ds.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "DiskStatus Get error: %v\n", err)
		return
	}

	for _, ads := range ds.AdapterStats {
		for _, vds := range ads.VirtualDriveStats {
			vdStatus := vds.State
			fmt.Printf("VD-%d: status: %s, size: %s, NumberOfDrives:%v, OsPath: %s\n", vds.VirtualDrive, vdStatus, vds.Size, vds.NumberOfDrives, vds.OsPath)
		}
		fmt.Printf("\n")

		for num, pds := range ads.PhysicalDriveStats {

			pdStatus := pds.FirmwareState
			pdName := []string{pds.Brand, pds.Model, pds.SerialNumber}
			pdSN := strings.Join(pdName, " ")
			diskGroup := "null"
			if pds.PdDiskGroup != "" || pds.PdArm != "" {
				diskGroup = pds.PdDiskGroup + "-" + pds.PdArm
			}
			fmt.Printf("PD-%d: %s, Size: %s, status: %s, PdType: %s %s, DiskGroup: %s\n", num, pdSN, pds.RawSize, pdStatus, pds.PdType, keepUppercaseLetters(pds.PdMediaType), diskGroup)
		}
		fmt.Printf("\n")
	}

}
