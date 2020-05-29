package main

import "fmt"
import "github.com/go-audio/audio"

var (
	sampleRate  = 48000
	numChannels = 1
	bitDepth    = 16
)


func ogg2wav(fileIn, fileOut string) error {

	audioBuffer := audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: numChannels,
			SampleRate:  sampleRate,
		},
		SourceBitDepth: bitDepth,
	}
	fmt.Printf("audioBuffer len: %d\n", len(audioBuffer.Data))

	return nil
}

func wav2mp3(fileIn, fileOut string) error {

	return nil
}

func wawMix(fileIn1, fileIn2, fileOut string) error {

	return nil
}

func main() {

}