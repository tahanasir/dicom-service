package image

import (
	"fmt"
	"image/png"
	"os"
	"sync"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/frame"
)

//Based on https://github.com/suyashkumar/dicom/blob/main/cmd/dicomutil/main.go

func WritePixelDataElement(info dicom.PixelDataInfo, filename string) {
	for _, f := range info.Frames {
		generateImage(f, nil, filename)
	}
}

func generateImage(fr *frame.Frame, wg *sync.WaitGroup, filename string) error {
	image, err := fr.GetImage()
	if err != nil {
		return err
	}

	name := fmt.Sprintf("%s.png", filename)
	f, err := os.Create(name)
	if err != nil {
		return err
	}

	err = png.Encode(f, image)
	if err != nil {
		return err
	}

	if err = f.Close(); err != nil {
		return err
	}

	if wg != nil {
		wg.Done()
	}

	return nil
}
