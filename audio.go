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

// NewPlayer initializes the speaker with a sample rate of sr,
// and a buffer size of sr/buffRatio. Returns an error if there
// is an error initializing the speaker.
func NewPlayer(sr beep.SampleRate, buffRatio int) (player, error) {
	err := speaker.Init(sr, sr.N(time.Second)/buffRatio)
	if err != nil {
		return player{}, err
	}
	return player{}, nil
}

// playSound resamples the audio stream, stream, and plays it through the
// speaker. It takes an optional done channel on which a signal will be sent
// when the audio finishes playing.
func (p player) PlaySound(speakerSampleRate beep.SampleRate, stream audioStream, done chan<- bool) {
	speaker.Clear()
	stream.stream.Seek(stream.startPos)
	speaker.Play(beep.Seq(beep.Resample(3, stream.format.SampleRate, speakerSampleRate, stream.stream), beep.Callback(func() {
		if done != nil {
			done <- true
		}
	})))
}

// newAudioStream decodes an mp3 file at the given path and returns a new
// audioStream.
func (p player) NewAudioStream(filePath string) (audioStream, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return audioStream{}, err
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return audioStream{}, err
	}
	return audioStream{
		stream:   streamer,
		format:   format,
		startPos: streamer.Position(),
	}, nil
}

func (p player) NilAudioStream(sr beep.SampleRate) audioStream {
	return audioStream{
		stream: nilStream{},
		format: beep.Format{
			SampleRate:  sr,
			NumChannels: 2,
			Precision:   3,
		},
		startPos: 0,
	}
}

type nilStream struct{}

func (ns nilStream) Stream(samples [][2]float64) (n int, ok bool) {
	return 0, true
}

func (ns nilStream) Err() error {
	return nil
}

func (ns nilStream) Len() int {
	return 0
}

func (ns nilStream) Position() int {
	return 0
}

func (ns nilStream) Seek(p int) error {
	return nil
}

func (ns nilStream) Close() error {
	return nil
}
