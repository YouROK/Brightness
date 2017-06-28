package brightness

import (
	"errors"
	"github.com/blackjack/webcam"
	"image"
	"image/jpeg"
	"os"
	"sort"
	"strings"
)

var (
	cam     *webcam.Webcam
	camPath string
	width   int
	height  int
)

func SetCameraDevice(path string) {
	camPath = path
}

func CameraON() error {
	var err error
	cam, err = webcam.Open(camPath)
	if err != nil {
		return err
	}

	err = setFormat()
	if err != nil {
		cam.Close()
		return err
	}

	return nil
}

func CameraOFF() error {
	if cam != nil {
		return cam.Close()
	}
	return nil
}

func CameraStart() error {
	return cam.StartStreaming()
}

func CameraStop() error {
	return cam.StopStreaming()
}

func CameraGetFrame(timeout int) ([]byte, error) {
	err := cam.WaitForFrame(uint32(timeout))
	if err != nil {
		return nil, err
	}

	frame, err := cam.ReadFrame()
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, len(frame))
	copy(buffer, frame)
	for {
		frame, err := cam.ReadFrame()
		if err != nil {
			break
		}
		copy(buffer, frame)
	}
	return buffer, nil
}

func setFormat() error {
	format_desc := cam.GetSupportedFormats()
	var format webcam.PixelFormat

	for f := range format_desc {
		if strings.Contains(format_desc[f], "YUYV") {
			format = f
			break
		}
	}

	if format == 0 {
		return errors.New("Support format not found")
	}

	frames := cam.GetSupportedFrameSizes(format)
	if len(frames) == 0 {
		return errors.New("Supported frame sizes not found")
	}

	sort.Slice(frames, func(i, j int) bool {
		ls := frames[i].MaxWidth * frames[i].MaxHeight
		rs := frames[j].MaxWidth * frames[j].MaxHeight
		return ls < rs
	})

	size := frames[0]

	_, w, h, err := cam.SetImageFormat(format, uint32(size.MaxWidth), uint32(size.MaxHeight))
	width = int(w)
	height = int(h)
	return err
}

func SaveJPEG(name string, frame []byte) error {
	yBytes := make([]byte, 0)
	cbBytes := make([]byte, 0)
	crBytes := make([]byte, 0)

	for index, value := range frame {
		switch index % 4 {
		case 0:
			fallthrough
		case 2:
			// Y values
			yBytes = append(yBytes, value)
		case 1:
			// U values
			cbBytes = append(cbBytes, value) // TODO: U == cb? think so...
		case 3:
			// V values
			crBytes = append(crBytes, value) // TODO: V == cr? think so...
		}
	}

	img := &image.YCbCr{
		Y:              yBytes,
		Cb:             cbBytes,
		Cr:             crBytes,
		SubsampleRatio: image.YCbCrSubsampleRatio422,
		YStride:        width,
		CStride:        height / 2,
		Rect:           image.Rect(0, 0, width, height),
	}

	out, err := os.Create(name)
	if err != nil {
		return err
	}
	defer out.Close()

	options := &jpeg.Options{Quality: 90}
	return jpeg.Encode(out, img, options)
}
