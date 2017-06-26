package main

import (
	"brightness"
	"flag"
	"fmt"
)

var (
	opt    = brightness.GetOptions()
	isTest bool
)

func init() {
	flag.BoolVar(&isTest, "test", false, "Testing program to obtain percent")

	flag.BoolVar(&opt.CameraAlwaysOn, "cao", true, "Camera is always on")
	flag.StringVar(&opt.CameraDevicePath, "cdp", "/dev/video0", "Camera device path")
	flag.IntVar(&opt.CameraTimeout, "ct", 1, "Camera timeout")

	flag.StringVar(&opt.MaxBrightnessPath, "mbp", "/sys/class/backlight/intel_backlight/max_brightness", "Path to max brightness")
	flag.StringVar(&opt.SetBrightnessPath, "sbp", "/sys/class/backlight/intel_backlight/brightness", "Path to set brightness")
	flag.StringVar(&opt.ActualBrightnessPath, "abp", "/sys/class/backlight/intel_backlight/actual_brightness", "Path to actual brightness")

	flag.IntVar(&opt.MinPercent, "min", 8, "Percentage with minimum illumination")
	flag.IntVar(&opt.MaxPercent, "max", 60, "Percentage with maximum illumination")
	flag.IntVar(&opt.AveragePercentCount, "apc", 5, "Average percent count")

	flag.Parse()
}

func main() {
	bright := brightness.NewBrightness(opt)
	var err error
	if !isTest {
		err = bright.RunAutoBrightness()
	} else {
		err = bright.TestBrightness()
	}
	if err != nil {
		fmt.Println("Error:", err)
	}
}
