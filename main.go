package main

import (
	"bufio"
	"fmt"
	"github.com/go-audio/audio"
	"io"
	"os"
	"github.com/hraban/opus"
	"github.com/go-audio/wav"
	"github.com/viert/go-lame"
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

func wav2mp3(fileIn, fileOut string) {
	f_out, err := os.Create(fileOut)
	check(err)
	defer f_out.Close()
	enc := lame.NewEncoder(f_out)
	err = enc.SetNumChannels(1)
	err = enc.SetQuality(5)
	check(err)
	defer enc.Close()

	f_in, err := os.Open(fileIn)
	check(err)
	defer f_in.Close()

	r := bufio.NewReader(f_in)
	// simple way - just write whole file without modifications
	//_, err = r.WriteTo(enc)

	// but doing via Read because of a weird click at the beginning
	buf := make([]byte, 1024)
	bytesToSkip := 1024 * 5 // skip first glitch samples. 10 = 0.9sec?
	var bytesRead int
	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			check(err)
		}
		bytesRead += n
		if bytesRead < bytesToSkip {
			continue
		}
		pcm := buf[:n]
		_, err = enc.Write(pcm)
		check(err)
	}
}

func wawMix(fileIn1, fileIn2, fileOut string) error {

	return nil
}

func main() {
	// os.Args[0] - path to binary
	if len(os.Args) < 4 {
		// TODO print help
		fmt.Println("use: [ogg2wav,wav2mp3] file_in file_out")
		os.Exit(0)
	}
	argsCommand := os.Args[1]
	argsFileIn := os.Args[2]
	argsFileOut := os.Args[3]
	switch argsCommand {
	case "ogg2wav":
		ogg2wav(argsFileIn, argsFileOut)
	case "wav2mp3":
		wav2mp3(argsFileIn, argsFileOut)
	default:
		fmt.Printf("command %s not found", argsCommand)
	}
}