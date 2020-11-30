package main

// typedef unsigned char Uint8;
// void GetStdin(void *userdata, Uint8 *stream, int len);
import "C"
import (
	"bufio"
	"flag"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"unsafe"

	log "github.com/sirupsen/logrus"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	stdin = make(chan C.Uint8)
	quit  = make(chan os.Signal, 1)

	DefaultFrequency *int64
	Default16Bit     *bool
	DefaultChannels  *uint64
	DefaultSamples   *uint64
	DebugMode        *bool
)

func init() {
	DefaultFrequency = flag.Int64("f", 44000, "Sampling frequency")
	DefaultChannels = flag.Uint64("c", 1, "Number of channels")
	DefaultSamples = flag.Uint64("s", 512, "Buffer size in bytes")
	Default16Bit = flag.Bool("16", false, "Use 16-bit samples")
	DebugMode = flag.Bool("d", false, "Debug mode")

	flag.Parse()
}

//export GetStdin
func GetStdin(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	n := int(length)
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(stream)), Len: n, Cap: n}
	buf := *(*[]C.Uint8)(unsafe.Pointer(&hdr))

	for i := 0; i < n; i++ {
		select {
		case buf[i] = <-stdin:
		case <-quit:
			quit <- syscall.SIGUSR1
			log.Debugf("Signalled. Quitting.")
			return
		}
	}
}

func main() {
	var format sdl.AudioFormat
	if *Default16Bit {
		format = sdl.AUDIO_S16
	} else {
		format = sdl.AUDIO_S8
	}

	if *DebugMode {
		log.SetLevel(log.DebugLevel)

		log.Infof("Frequency: %d", *DefaultFrequency)
		log.Infof("Channels:  %d", *DefaultChannels)
		if *Default16Bit {
			log.Infof("Bits:      %d", 16)
		} else {
			log.Infof("Bits:      %d", 8)
		}

		log.Infof("Buffer:    %d", *DefaultSamples)
	} else {
		log.SetLevel(log.ErrorLevel)
	}

	var dev sdl.AudioDeviceID
	var err error

	// Initialize SDL2
	if err := sdl.Init(sdl.INIT_AUDIO); err != nil {
		log.Println(err)
		return
	}
	defer sdl.Quit()

	go func() {
		reader := bufio.NewReader(os.Stdin)
		defer log.Debugf("Quitting reader goroutine")
		for {
			bb, err := reader.ReadByte()
			if err != nil {
				quit <- syscall.SIGUSR1
				log.Debug("Done reading.")
				return
			}
			stdin <- C.Uint8(bb)
		}
	}()

	spec := sdl.AudioSpec{
		Freq:     int32(*DefaultFrequency),
		Format:   format,
		Channels: uint8(*DefaultChannels),
		Samples:  uint16(*DefaultSamples),
		Callback: sdl.AudioCallback(C.GetStdin),
	}

	// Open default playback device
	if dev, err = sdl.OpenAudioDevice("", false, &spec, nil, 0); err != nil {
		log.Println(err)
		return
	}
	defer sdl.CloseAudioDevice(dev)

	// Start playback audio
	sdl.PauseAudioDevice(dev, false)

	// Listen to OS signals

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Run infinite loop until we receive SIGINT or SIGTERM or we are done.
	running := true
	for running {
		select {
		case <-quit:
			log.Debug("Exiting loop.")
			quit <- syscall.SIGUSR1
			running = false
		}
	}
}
