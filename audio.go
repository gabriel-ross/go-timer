package timer

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type audioStream struct {
	stream   beep.StreamSeekCloser
	format   beep.Format
	startPos int
}

type player struct{}

// initAndPlay is a shortcut for initializing the speaker and playing an
// an audio stream. Takes in an optional done channel on which a signal
// will be sent when the audio finishes playing.
func (p player) initAndPlay(stream audioStream, done chan<- bool) {
	speaker.Clear()
	p.initSpeaker(stream.format)
	stream.stream.Seek(stream.startPos)
	speaker.Play(beep.Seq(stream.stream, beep.Callback(func() {
		if done != nil {
			done <- true
		}
	})))
}

func (p player) initSpeaker(format beep.Format) error {
	return speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
}

func (p player) playSound(stream audioStream, done chan<- bool) {
	speaker.Clear()
	stream.stream.Seek(stream.startPos)
	speaker.Play(beep.Seq(stream.stream, beep.Callback(func() {
		if done != nil {
			done <- true
		}
	})))
}

func (p player) newAudioStream(filePath string) (audioStream, error) {
	streamer, format, err := p.decodeSampleFromFile(filePath)
	if err != nil {
		return audioStream{}, err
	}
	return audioStream{
		stream:   streamer,
		format:   format,
		startPos: streamer.Position(),
	}, nil
}

func (p player) decodeSampleFromFile(path string) (beep.StreamSeekCloser, beep.Format, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, beep.Format{}, err
	}
	defer f.Close()
	return p.decodeSample(f)
}

func (p player) decodeSample(f *os.File) (beep.StreamSeekCloser, beep.Format, error) {
	return mp3.Decode(f)
}
