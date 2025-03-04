package consts

const (
	PCISysPath   = "/sys/bus/pci/devices"
	PCIInfoRegex = `^([^:]+): \{(.*)\}$`
)

var (
	PCIDevicesInfoPath = []string{
		"/usr/share/misc/pci.ids",
		"/usr/share/hwdata/pci.ids",
	}
)

//  "00:00.0": {
//     "device_id": "8086:6f00",
//     "driver": "intel_uncore",
//     "module_aliases": "pci:v00008086d00006F00sv*sd*bc*sc*i*",
//     "udev_aliases": "primary-bridge",
//     "sys_path": "/sys/bus/pci/devices/0000:00:00.0"
//   },

type PCIInfo struct {
	PCIClass    string `json:"pci_class,omitempty"`     // PCI类
	PCIID       string `json:"pci_id,omitempty"`        // PCIID
	PCISubsysID string `json:"pci_subsys_id,omitempty"` // PCI子系统ID
	PCISlotName string `json:"pci_slot_name,omitempty"` // PCI插槽名称
	MODALIAS    string `json:"modalias,omitempty"`      // 模块别名
	Driver      string `json:"driver"`                  // 驱动
	SysPath     string `json:"sys_path,omitempty"`      // 系统路径
	DeviceName  string `json:"device_name,omitempty"`   // 设备名称
	VendorName  string `json:"vendor_name,omitempty"`   // 厂商名称
	DeviceID    string `json:"device_id,omitempty"`     // 设备ID
	VendorID    string `json:"vendor_id,omitempty"`     // 厂商ID
}

// GetPCIInfoListParam 获取PCI信息列表参数
type GetPCIInfoListParam struct {
	Async       bool `json:"async,omitempty"`       // 是否异步
	Concurrency int  `json:"concurrency,omitempty"` // 并发数
}

// GetPCIInfoListResult 获取PCI信息列表结果
type GetPCIInfoListResult struct {
	PCIInfoList []PCIInfo `json:"pci_info_list,omitempty"` // PCI信息列表
}

// GetPCIInfoByPCISlotNameParam 获取PCI信息参数
type GetPCIInfoByPCISlotNameParam struct {
	PCISlotName string `json:"pci_slot_name,omitempty"` // PCI插槽名称
	// PCIDevicesInfoPath string `json:"pci_devices_info_path,omitempty"` // PCI设备信息路径
}

// GetPCIInfoByPCISlotNameResult 获取PCI信息结果
type GetPCIInfoByPCISlotNameResult struct {
	PCIInfo `json:"pci_info,omitempty"` // PCI信息
}

// GeneratePCIDeviceMapResult 生成PCI设备Map结果
// type GeneratePCIDeviceMapResult struct {
// 	PCISubsysMap map[string]string `json:"pci_subsys_map",omitempty` // PCI子系统Map
// 	DeviceMap    map[string]string `json:"device_map",omitempty`     // 设备Map
// 	VendorMap    map[string]string `json:"vendor_map",omitempty`     // 厂商Map
// }
