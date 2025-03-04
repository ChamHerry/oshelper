package controller

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/ChamHerry/oshelper/consts"
	"github.com/ChamHerry/oshelper/utils"
)

// GetPCIInfoList 获取PCI信息列表
func (s *Controller) GetPCIInfoList(ctx context.Context, in consts.GetPCIInfoListParam) (out consts.GetPCIInfoListResult, err error) {
	command := "ls " + consts.PCISysPath
	pciInfo, err := s.RunCommand(consts.RunCommandConfig{
		Command:                command,
		RunCommandFailedCounts: consts.DefaultRunCommandFailedCounts,
	})
	if err != nil {
		return out, err
	}
	pciInfoList := strings.Split(pciInfo, "\n")
	out.PCIInfoList = make([]consts.PCIInfo, 0)
	if in.Async {
		// g.Log().Debug(ctx, "异步处理")
		if in.Concurrency <= 0 {
			in.Concurrency = 10
		}
		// 异步处理
		interfaceSlice := utils.ConvertSliceToInterfaceSlice(pciInfoList)
		asyncCallResultList, err := utils.AsyncCall(ctx, consts.AsyncCallParam{
			Items: interfaceSlice,
			Operation: func(ctx context.Context, item interface{}) (interface{}, error) {
				return s.GetPCIInfoByPCISlotName(ctx, consts.GetPCIInfoByPCISlotNameParam{
					PCISlotName: item.(string),
				})
			},
			Concurrency: in.Concurrency,
		})
		if err != nil {
			return out, err
		}
		for _, asyncCallResult := range asyncCallResultList.RetList {
			out.PCIInfoList = append(out.PCIInfoList, asyncCallResult.Ret.(consts.GetPCIInfoByPCISlotNameResult).PCIInfo)
		}
	} else {
		// g.Log().Debug(ctx, "同步处理")
		// 同步处理
		for _, PCIID := range pciInfoList {
			// g.Log().Debugf(ctx, "PCIID: %s", PCIID)
			pciInfo := consts.PCIInfo{}
			pciInfo.PCIID = PCIID
			getPCIInfoByPCISlotNameResult, err := s.GetPCIInfoByPCISlotName(ctx, consts.GetPCIInfoByPCISlotNameParam{
				PCISlotName: PCIID,
			})
			if err != nil {
				return out, err
			}
			pciInfo = getPCIInfoByPCISlotNameResult.PCIInfo
			// g.Log().Debug(ctx, "pciInfo:", pciInfo)
			out.PCIInfoList = append(out.PCIInfoList, pciInfo)
		}
	}
	return out, nil
}

// GetPCIInfoByPCISlotName 获取指定PCI设备信息
func (s *Controller) GetPCIInfoByPCISlotName(ctx context.Context, in consts.GetPCIInfoByPCISlotNameParam) (out consts.GetPCIInfoByPCISlotNameResult, err error) {
	if strings.HasPrefix(in.PCISlotName, "0000") {
		deviceIDList := strings.Split(in.PCISlotName, ":")
		in.PCISlotName = strings.Join(deviceIDList[1:], ":")
	}
	SysPath := consts.PCISysPath + "/0000:" + in.PCISlotName
	if !s.IsFileOrDirExist(ctx, SysPath) {
		SysPath = consts.PCISysPath + "/" + in.PCISlotName
		if !s.IsFileOrDirExist(ctx, SysPath) {
			return out, errors.New("deviceID not found")
		}
	}
	out.SysPath = SysPath
	out.PCISlotName = in.PCISlotName
	ueventPath := SysPath + "/uevent"
	// g.Log().Debug(ctx, "ueventPath:", ueventPath)
	uevent, err := s.GetFileContent(ctx, ueventPath)
	if err != nil {
		// g.Log().Error(ctx, "GetFileContent error:", err)
		return out, err
	}
	// g.Log().Debug(ctx, "uevent:\n", uevent)
	ueventMap, err := utils.ParseToJSON(uevent, "=")
	if err != nil {
		return out, err
	}
	// g.Log().Debug(ctx, "ueventMap:", ueventMap)
	json.Unmarshal([]byte(ueventMap), &out.PCIInfo)
	// if s.IsFileExist(ctx, ueventPath) {

	// }
	// g.Log().Debug(ctx, "-------------------------------------------")
	deviceIdPath := SysPath + "/device"
	vendorIdPath := SysPath + "/vendor"
	if s.IsFileExist(ctx, deviceIdPath) {
		deviceId, err := s.GetFileContent(ctx, deviceIdPath)
		if err != nil {
			return out, err
		}
		out.DeviceID = deviceId
	}
	if s.IsFileExist(ctx, vendorIdPath) {
		vendorId, err := s.GetFileContent(ctx, vendorIdPath)
		if err != nil {
			return out, err
		}
		out.VendorID = vendorId
	}
	return out, nil
}
