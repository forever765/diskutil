## 2024-08-03 forked from [buaazp/diskutil](https://github.com/buaazp/diskutil)

1. feat: 获取VirtualDisk对应的盘符OsPath
2. 新增字段  
   1. pds.PdMediaType (标识 HDD or SSD)
   2. pds.PdDiskGroup（物理磁盘所在的VD组）
   3. pds.PdArm（物理磁盘所在的VD组中的序号）
2. improve: 适配Enclosure Device ID=N/A的情况
3. improve: 禁用MegaCli64的日志输出文件
4. improve: complete example code
5. fixed: RawSize获取异常、物理硬盘型号获取失败的bug

---

## Golang MegaCli Utils

This package is used for go codes to get MegaRaid stat.

### Usage

*At first, you need install MegaRAID in your servers.*

Create a DiskStatus struct by calling `diskutil.NewDiskStatus()`. You need provide the MegaCli binary path and the count of RAID card in your server.

```
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
	fmt.Println(ds)
```

After calling `Get()`, you can visit any stat in the DiskStatus like this:

```
	for i, ads := range ds.AdapterStats {
		fmt.Printf("adapter #%d \n", i)
		for j, pds := range ads.PhysicalDriveStats {
			pdStatus := pds.FirmwareState
			pdName := []string{pds.Brand, pds.Model, pds.SerialNumber}
			pdSN := strings.Join(pdName, " ")
			fmt.Printf("PD%d: %s status: %s\n", j, pdSN, pdStatus)
		}
		fmt.Printf("\n")
	}
```

If you focus on the disk which is broken, you can use `ListBrokenDrive()` to get them:

```
	brokenVds, brokenPds, err := ds.ListBrokenDrive()
	if err != nil {
		fmt.Fprintf(os.Stderr, "DiskStatus ListBrokenDrive error: %v\n", err)
		return
	}
	for _, bvd := range brokenVds {
		fmt.Println(bvd)
	}
	for _, bpd := range brokenPds {
		fmt.Println(bpd)
	}
```

Or you can print the DiskStatus in json format by calling `ToJson()`:

```
	jsonStatus, err := ds.ToJson()
	if err != nil {
		fmt.Fprintf(os.Stderr, "DiskStatus ToJson error: %v\n", err)
		return
	}
	fmt.Println(jsonStatus)

	{
		"adapter_stats": [
			{
				"id": 0, 
				"virtual_drive_stats": [
					{
						"virtual_drive": 0, 
						"name": "", 
						"size": "278.875 GB", 
						"state": "Optimal", 
						"number_of_drives": 1, 
						"encryption_type": "None"
					}
				], 
				"physical_drive_stats": [
					{
						"enclosure_device_id": 64, 
						"device_id": 8, 
						"slot_number": 0, 
						"media_error_count": 0, 
						"other_error_count": 0, 
						"predictive_failure_count": 0, 
						"pd_type": "SAS", 
						"raw_size": "279.396 GB [0x22ecb25c Sectors]", 
						"firmware_state": "Online, Spun Up", 
						"brand": "SEAGATE", 
						"model": "ST9300605SS", 
						"serial_number": "00046XP4MQNJ", 
						"drive_emperature": "65C (149.00 F)"
					}
				]
			}
		]
	}
```

Full sample code is in /examples. Try it to test this package:

```
go build -v examples/printDiskStat.go
sudo ./printDiskStat
```
![example-image-1](https://github.com/forever765/diskutil/blob/master/images/example-1.png)  
![example-image-2](https://github.com/forever765/diskutil/blob/master/images/example-2.png)

### GoDoc

Visit Godoc to get full api documents:

[https://godoc.org/github.com/buaazp/diskutil](https://godoc.org/github.com/buaazp/diskutil) 

### Issue

If you meet some problems in your servers, please create a github [issue](https://github.com/buaazp/diskutil/issues) or contact me:

### Original Author
weibo: [@招牌疯子](http://weibo.com/buaazp)  
mail: zp@buaa.us


