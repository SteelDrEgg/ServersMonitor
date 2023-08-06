package util

import (
	"Server/graph/model"
	"github.com/shirou/gopsutil/v3/mem"
	"reflect"
)

func RamStat() model.RAM {
	stat, _ := mem.VirtualMemory()
	ram := model.RAM{}
	ram.Total = int(stat.Total)
	ram.Available = int(stat.Available)
	ram.Used = int(stat.Used)
	ram.Free = int(stat.Free)
	ram.Active = int(stat.Active)
	ram.Inactive = int(stat.Inactive)
	ram.Wired = int(stat.Wired)
	return ram
}

func PrettyRam(decimal int) model.PrettyRAM {
	stat := RamStat()
	ram := model.PrettyRAM{}
	statVal := reflect.ValueOf(stat)
	statName := reflect.TypeOf(stat)
	ramVal := reflect.ValueOf(&ram).Elem()
	for i := 0; i < statName.NumField(); i++ {
		name := statName.Field(i).Name
		ramVal.FieldByName(name).SetString(ProperUnit(uint64(statVal.FieldByName(name).Int()), decimal))
	}
	return ram
}
