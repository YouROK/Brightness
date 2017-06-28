package brightness

import (
	"fmt"
	"sync"
	"time"
)

type Options struct {
	CameraAlwaysOn   bool
	CameraDevicePath string
	CameraTimeout    int

	MaxBrightnessPath    string
	SetBrightnessPath    string
	ActualBrightnessPath string

	MinPercent          int
	MaxPercent          int
	AveragePercentCount int
}

func GetOptions() *Options {
	o := new(Options)

	o.CameraAlwaysOn = true
	o.CameraDevicePath = "/dev/video0"
	o.CameraTimeout = 1

	o.MaxBrightnessPath = "/sys/class/backlight/intel_backlight/max_brightness"
	o.SetBrightnessPath = "/sys/class/backlight/intel_backlight/brightness"
	o.ActualBrightnessPath = "/sys/class/backlight/intel_backlight/actual_brightness"

	o.MinPercent = 8
	o.MaxPercent = 60
	o.AveragePercentCount = 5

	return o
}

type Brightness struct {
	opt            *Options
	lastPercent    int
	lastCount      int
	averagePercent []int
	isStop         bool
}

func NewBrightness(opts *Options) *Brightness {
	b := new(Brightness)
	b.opt = opts
	if b.opt.MinPercent == b.opt.MaxPercent {
		if b.opt.MaxPercent == 0 {
			b.opt.MaxPercent = 100
		} else {
			b.opt.MinPercent = 0
		}
	}
	SetCameraDevice(opts.CameraDevicePath)
	return b
}

func (b *Brightness) RunAutoBrightness() error {
	var lock sync.Mutex

	CameraON()
	defer CameraOFF()
	if b.opt.CameraAlwaysOn {
		CameraStart()
	}
	stepCount := 0
	b.isStop = false
	for !b.isStop {
		frame, err := b.getFrame()
		if err != nil {
			return err
		}
		if len(frame) != 0 {
			prc := b.getBright(frame, true)
			aver := b.getAveragePercent(prc)
			//fmt.Println("aver:", aver, "curr:", prc)
			//SaveJPEG("./camera.jpeg", frame)
			if stepCount > 60 && b.opt.CameraAlwaysOn {
				CameraStop()
				time.Sleep(time.Millisecond * 100)
				CameraStart()
				stepCount = 0
			}
			if isSetBrightness(&b.lastPercent, &aver, &b.lastCount) {
				b.lastPercent = aver
				lock.Lock()
				go func() {
					b.setBrightness(aver)
					lock.Unlock()
				}()
			} else {
				if b.opt.CameraAlwaysOn {
					stepCount++
				}
				time.Sleep(time.Second)
			}
		}
	}
	return nil
}

func (b *Brightness) TestBrightness() error {
	var lock sync.Mutex

	CameraON()
	defer CameraOFF()
	if b.opt.CameraAlwaysOn {
		CameraStart()
	}

	b.isStop = false
	for !b.isStop {
		frame, err := b.getFrame()
		if err != nil {
			return err
		}
		if len(frame) != 0 {
			prc := b.getBright(frame, false)
			SaveJPEG("./camera.jpeg", frame)
			aver := b.getAveragePercent(prc)
			fmt.Println("Percent aver:", aver, "curr:", prc)
			lock.Lock()
			go func() {
				b.setBrightness(aver)
				lock.Unlock()
			}()
		}
	}
	return nil
}

func (b *Brightness) Stop() {
	b.isStop = true
}

func (b *Brightness) getFrame() ([]byte, error) {
	if !b.opt.CameraAlwaysOn {
		err := CameraStart()
		if err != nil {
			return nil, err
		}
	}

	frame, err := CameraGetFrame(b.opt.CameraTimeout)
	if err != nil {
		CameraOFF()
		return nil, err
	}

	if !b.opt.CameraAlwaysOn {
		err := CameraStop()
		if err != nil {
			return nil, err
		}
	}
	return frame, err
}

func (b *Brightness) getBright(frame []byte, setMinMax bool) int {
	yBytes := make([]byte, 0)

	for index, value := range frame {
		if index%4 == 2 {
			yBytes = append(yBytes, value)
		}
	}

	all := 0
	for _, y := range yBytes {
		all += int(y)
	}
	//среднее значение яркости по картинке
	clr := all / len(yBytes)
	//процент от яркости
	prc := (clr * 100 / 255)
	if setMinMax {
		prc -= b.opt.MinPercent
		prc = prc * 100 / (b.opt.MaxPercent - b.opt.MinPercent)
	}
	if prc > 100 {
		prc = 100
	}
	if prc < 1 { //0% is off display
		prc = 1
	}
	return prc
}

func (b *Brightness) getAveragePercent(percent int) int {
	b.averagePercent = append([]int{percent}, b.averagePercent...)

	if len(b.averagePercent) > b.opt.AveragePercentCount {
		b.averagePercent = b.averagePercent[:b.opt.AveragePercentCount]
	}
	all := 0
	for _, v := range b.averagePercent {
		all += v
	}
	return all / len(b.averagePercent)
}

func (b *Brightness) setBrightness(prc int) {
	maxBright := readFileInt(b.opt.MaxBrightnessPath)
	actBright := readFileInt(b.opt.ActualBrightnessPath)
	targBright := maxBright / 100 * prc
	var curr float64 = float64(actBright)
	step := float64(targBright-actBright) / 100
	for i := 0; i < 99; i++ {
		curr += step
		if actBright != int(curr) {
			writeFileInt(b.opt.SetBrightnessPath, int(curr))
			actBright = int(curr)
		}
		time.Sleep(time.Millisecond * 10)
	}
	writeFileInt(b.opt.SetBrightnessPath, int(targBright))
	time.Sleep(time.Millisecond * 10)
}
