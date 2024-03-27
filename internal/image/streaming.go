package image

import (
	"io"
	"sync"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/frame"
)

//Based on https://github.com/suyashkumar/dicom/blob/main/cmd/dicomutil/main.go

// FrameBufferSize represents the size of the *Frame buffered channel for streaming calls
const FrameBufferSize = 100

func ParseWithStreaming(in io.Reader, size int64, filename string) (*dicom.Dataset, error) {
	fc := make(chan *frame.Frame, FrameBufferSize)

	// Go routine to process frames as they are sent to frameChannel
	var wg sync.WaitGroup
	wg.Add(1)
	go WriteStreamingFrames(fc, &wg, filename)

	ds, err := dicom.Parse(in, size, fc)
	if err != nil {
		return &ds, err
	}
	wg.Wait()

	return &ds, nil

}

func WriteStreamingFrames(frameChan chan *frame.Frame, doneWG *sync.WaitGroup, filename string) {
	count := 0 // may not correspond to frame number
	var wg sync.WaitGroup
	for fr := range frameChan {
		count++
		wg.Add(1)
		go generateImage(fr, &wg, filename)
	}
	wg.Wait()
	doneWG.Done()
}
