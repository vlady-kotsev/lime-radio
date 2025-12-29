package radio

import (
	"encoding/binary"
	"fmt"
	"time"
)

func CreateWAVHeader(sampleRate int) ([]byte, error) {
	header := make([]byte, 44)
	copy(header[0:4], "RIFF")
	binary.LittleEndian.PutUint32(header[4:8], 0x7FFFFFFF)
	copy(header[8:12], "WAVE")
	copy(header[12:16], "fmt ")
	// PCM format size
	binary.LittleEndian.PutUint32(header[16:20], 16)
	// PCM format
	binary.LittleEndian.PutUint16(header[20:22], 1)
	// Stereo
	binary.LittleEndian.PutUint16(header[22:24], uint16(2))
	if sampleRate < 0 || sampleRate > 0x3FFFFFFF {
		return []byte{}, fmt.Errorf("invalid sample rate")
	}
	// Sample rate
	binary.LittleEndian.PutUint32(header[24:28], uint32(sampleRate))
	// Byte rate
	binary.LittleEndian.PutUint32(header[28:32], uint32(sampleRate)*4)
	// Block align
	binary.LittleEndian.PutUint16(header[32:34], uint16(4))
	// Bits per sample
	binary.LittleEndian.PutUint16(header[34:36], uint16(16))
	copy(header[36:40], "data")
	binary.LittleEndian.PutUint32(header[40:44], 0x7FFFFFFF)
	return header, nil
}

// SpeedUp is used because browser falls behind
const SpeedUp float64 = 0.75

func calculateStreamIntervalForBytes(sampleRate int, numBytes int) time.Duration {
	// 16-bit stereo
	bytesPerSecond := float64(sampleRate * 2 * 2)
	intervalMs := (float64(numBytes) * 1000.0) / bytesPerSecond
	return time.Duration((intervalMs - SpeedUp) * float64(time.Millisecond))
}
