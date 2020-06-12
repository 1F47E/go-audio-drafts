package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/go-audio-drafts/audiotools"
	"io/ioutil"
	"os"
)

func main() {
	// read OGG file
	f, err := os.Open("data/response.ogg")
	if err != nil {
		panic(err)
	}
	oggReader := bufio.NewReader(f)

	// convert to WAV bytes
	var wavData []byte
	wavData, err = audiotools.OggToWav(oggReader)
	if err != nil {
		panic(err)
	}
	fmt.Printf("got %d bytes of wav audio\n", len(wavData))

	// SAVE BYTES TO WAV FILE
	fmt.Println("writing wav file...")
	err = ioutil.WriteFile("data/response.wav", wavData, 0755)
	if err != nil {
		panic(err)
	}
	fmt.Printf("done\n")

	wavReader := bytes.NewReader(wavData)
	dataMp3, err := audiotools.WavToMp3(wavReader)
	if err != nil {
		panic(err)
	}

	// SAVE BYTES TO MP3 FILE
	fmt.Println("writing mp3 file...")
	err = ioutil.WriteFile("data/response.mp3", dataMp3, 0755)
	if err != nil {
		panic(err)
	}
	fmt.Printf("done\n")
}