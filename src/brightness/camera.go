package brightness

import (
	"github.com/blackjack/webcam"
	"image"
	"image/jpeg"
	"os"
)

var (
	cam     *webcam.Webcam
	camPath string
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
	return cam.StartStreaming()
}

func CameraOFF() error {
	if cam != nil {
		return cam.Close()
	}
	return nil
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
		YStride:        320,
		CStride:        160,
		Rect:           image.Rect(0, 0, 320, 240),
	}

	out, err := os.Create(name)
	if err != nil {
		return err
	}
	defer out.Close()

	options := &jpeg.Options{Quality: 90}
	return jpeg.Encode(out, img, options)
}
