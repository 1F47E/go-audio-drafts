package main

import (
	"fmt"
	"github.com/go-audio/audio"
	"io"
	"os"
	"github.com/hraban/opus"
	"github.com/go-audio/wav"
)

var (
	sampleRate  = 48000
	numChannels = 1
	bitDepth    = 16
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ogg2wav(fileIn, fileOut string) {
	file_in, err := os.Open(fileIn)
	check(err)
	stream, err := opus.NewStream(file_in)
	check(err)
	defer stream.Close()

	// create wav buffer
	file_out, err := os.Create(fileOut)
	check(err)

	// wav encoder
	encoder_out := wav.NewEncoder(file_out, sampleRate, bitDepth, numChannels, 1)

	frameSize := numChannels * 60 * sampleRate / 1000
	oggBuffer := make([]int16, frameSize)

	// create audio.IntBuffer
	audioBuffer := audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: numChannels,
			SampleRate:  sampleRate,
		},
		SourceBitDepth: bitDepth,
	}

	for {
		cnt, err := stream.Read(oggBuffer)
		if err == io.EOF {
			break
		} else if err != nil {
			check(err)
		}
		pcm := oggBuffer[:cnt*numChannels]


		// TODO write actual decoded data
		for _, b := range pcm {
			audioBuffer.Data = append(audioBuffer.Data, int(b))
		}


	}
	fmt.Printf("audioBuffer len: %d\n", len(audioBuffer.Data))
	// Write buffer to output file. This writes a RIFF header and the PCM chunks from the audio.IntBuffer.

	if err := encoder_out.Write(&audioBuffer); err != nil {
		check(err)
	}
	// close file + write header
	if err := encoder_out.Close(); err != nil {
		check(err)
	}
	fmt.Printf("encoder_out WrittenBytes: %d\n", encoder_out.WrittenBytes)
}

func wav2mp3(fileIn, fileOut string) error {

	return nil
}

func wawMix(fileIn1, fileIn2, fileOut string) error {

	return nil
}

func main() {
	ogg2wav("in.ogg", "out.wav")
}