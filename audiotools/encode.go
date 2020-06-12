package audiotools

import (
	"encoding/binary"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/hraban/opus"
	"github.com/viert/go-lame"
	"io"
)

var (
	sampleRate  = 48000
	numChannels = 1
	bitDepth    = 16
)

func OggToWav(r io.Reader) ([]byte, error) {
	stream, err := opus.NewStream(r)
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	// create empty bytes buffer
	fb := NewFileBuffer()

	// create wav encoder with our file buffer that supports WriteSeeker interface
	encoderOut := wav.NewEncoder(fb, sampleRate, bitDepth, numChannels, 1)

	frameSize := numChannels * 60 * sampleRate / 1000
	oggBuffer := make([]int16, frameSize)

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
			return nil, err
		}
		pcm := oggBuffer[:cnt*numChannels]

		for _, b := range pcm {
			audioBuffer.Data = append(audioBuffer.Data, int(b))
		}
	}
	// Write buffer to output file. This writes a RIFF header and the PCM chunks from the audio.IntBuffer.
	if err := encoderOut.Write(&audioBuffer); err != nil {
		return nil, err
	}

	// MANUALLY FIXING HEADER because of a broken seek
	// position 4 (8 bytes) and pos 40 (8 bytes) is file length by wav file specs
	headerPositions := []int{4, 40}
	audioBytes := fb.Bytes()
	for _, pos := range headerPositions {
		binary.LittleEndian.PutUint32(audioBytes[pos:], uint32(encoderOut.WrittenBytes))
	}

	return audioBytes, nil
}

func WavToMp3(r io.Reader) ([]byte, error) {
	// create empty bytes buffer and write to it instead of file
	fb := NewFileBuffer()

	enc := lame.NewEncoder(fb)
	err := enc.SetNumChannels(1)
	if err != nil {
		return nil, err
	}
	err = enc.SetQuality(5)
	if err != nil {
		return nil, err
	}
	defer enc.Close()

	// but doing via Read because of a weird click at the beginning
	buf := make([]byte, 1024)
	bytesToSkip := 1024 * 5 // skip first glitch samples.
	var bytesRead int
	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		bytesRead += n
		if bytesRead < bytesToSkip {
			continue
		}
		pcm := buf[:n]
		_, err = enc.Write(pcm)
		if err != nil {
			return nil, err
		}
	}
	return fb.Bytes(), nil
}

