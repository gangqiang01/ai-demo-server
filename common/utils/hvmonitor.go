package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
)

type HwMonitor struct {
	VBAT       float64 `json:"vabt"`
	VSB        float64 `json:"vsb"`
	VIN        float64 `json:"vin"`
	VCORE      float64 `json:"vcore"`
	TempCPU    float64 `json:"tempcpu"`
	TempSystem float64 `json:"tempsystem"`
}

var hw HwMonitor

func GetHvMonitor() HwMonitor {

	filepathNames, err := filepath.Glob(filepath.Join("/sys/class/hwmon/*"))
	if err != nil {
		log.Fatal(err)
		return hw
	} else {
		re := regexp.MustCompile(`[0-9]+[0-9]`)
		length := len(filepathNames) - 1

	OuterLoop:
		for i := length; i >= 0; i-- {

			mt := fmt.Sprintf(filepathNames[i] + "/name")
			contents, _ := ioutil.ReadFile(mt)
			result := string(contents)

			mes := regexp.MustCompile("advhwmon")
			if mes.MatchString(result) {

				vabt, _ := ioutil.ReadFile(
					filepathNames[i] + "/in1_input")
				vabtResult := string(vabt)
				vbatValue := re.FindString(vabtResult)
				VABT, _ := strconv.ParseFloat(vbatValue, 64)
				hw.VBAT = VABT / 1000

				vsb, _ := ioutil.ReadFile(
					filepathNames[i] + "/in2_input")
				vsbResult := string(vsb)
				vsbValue := re.FindString(vsbResult)
				VSB, _ := strconv.ParseFloat(vsbValue, 64)
				hw.VSB = VSB / 1000

				vin, _ := ioutil.ReadFile(
					filepathNames[i] + "/in3_input")
				vinResult := string(vin)
				vinValue := re.FindString(vinResult)
				VIN, _ := strconv.ParseFloat(vinValue, 64)
				hw.VIN = VIN / 1000

				vcore, _ := ioutil.ReadFile(
					filepathNames[i] + "/in4_input")
				vcoreResult := string(vcore)
				vcoreValue := re.FindString(vcoreResult)
				VCORE, _ := strconv.ParseFloat(vcoreValue, 64)
				hw.VCORE = VCORE / 1000

				tempCpu, _ := ioutil.ReadFile(
					filepathNames[i] + "/temp2_input")
				tempCpuResult := string(tempCpu)
				tempCpuValue := re.FindString(tempCpuResult)
				TempCpu, _ := strconv.ParseFloat(tempCpuValue, 64)
				hw.TempCPU = TempCpu / 1000

				tempSystem, _ := ioutil.ReadFile(
					filepathNames[i] + "/temp3_input")
				tempSystemResult := string(tempSystem)
				tempSystemValue := re.FindString(tempSystemResult)
				TempSystem, _ := strconv.ParseFloat(tempSystemValue, 64)
				hw.TempSystem = TempSystem / 1000
				break OuterLoop
			}
		}
		return hw
	}

}
