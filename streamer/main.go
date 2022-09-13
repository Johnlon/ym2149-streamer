package main

// based on https://github.com/FlorentFlament/ym2149-streamer.git
// and github.com/nguillaumin/ymtool

/*
# R00 = Channel A Pitch LO (8 bits)
# R01 = Channel A Pitch HI (4 bits)
# R02 = Channel B Pitch LO (8 bits)
# R03 = Channel B Pitch HI (4 bits)
# R04 = Channel C Pitch LO (8 bits)
# R05 = Channel C Pitch HI (4 bits)
# R06 = Noise Frequency    (5 bits)
# R07 = I/O & Mixer        (IOB|IOA|NoiseC|NoiseB|NoiseA|ToneC|ToneB|ToneA)
# R08 = Channel A Level    (M | 4 bits) (where M is mode)
# R09 = Channel B Level    (M | 4 bits)
# R10 = Channel C Level    (M | 4 bits)
# R11 = Envelope Freq LO   (8 bits)
# R12 = Envelope Freq HI   (8 bits)
# R13 = Envelope Shape     (CONT|ATT|ALT|HOLD)
*/

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nguillaumin/ymtool/ym"
	"github.com/tarm/serial"
	"go.bug.st/serial.v1/enumerator"
)

const exitCodeCmdParsing = 1
const exitCodeIOError = 2
const exitCodeUnsupportedVersion = 3

// Command Line Args
// EX: ./serial -com=COM3 -baud=57600
var (
	com  string
	baud int
)

// Set Defaults
func init() {
	ports := eports()
	var port = "COM1"
	if len(ports) == 1 {
		port = ports[0].Name
	}
	flag.StringVar(&com, "com", port, "The COM port to write data to")
	flag.IntVar(&baud, "baud", 57600, "The Baud Rate")
}

func main() {
	flag.Parse()

	if len(flag.Args()) < 1 {
		Usage()
		os.Exit(exitCodeCmdParsing)
	}

	command := flag.Arg(0)
	switch command {
	case "port":
		PortCmd()
	case "stream":
		StreamCmd()
	case "info":
		InfoCmd()
	default:
		Usage()
		os.Exit(exitCodeCmdParsing)
	}

}

// Usage prints the command line utility usage
func Usage() {
	fmt.Println("YM Tool")
	fmt.Println("")
	fmt.Println("Usage: ym <command> [arg...]")
	fmt.Println("")

	fmt.Println("Available commands:")
	fmt.Println("  port                       : Looks for a single port print to console")
	fmt.Println("  stream <file.ym>           : Stream YM file over the USB/Comm port")
	fmt.Println("  info <file.ym>             : Get metadata of a YM file")
	fmt.Println("")
	flag.PrintDefaults()

	fmt.Println("")
	fmt.Println("Run a command without argument to get more information.")
	fmt.Println("")

	fmt.Print("Supported YM versions: ")
	for version := range ym.SupportedYmVersions {
		fmt.Printf("%v ", version)
	}

	fmt.Println("")
	fmt.Println("")

	fmt.Println("Exit codes:")
	fmt.Printf("  - %v: Error parsing command line.\n", exitCodeCmdParsing)
	fmt.Printf("  - %v: I/O error reading or writing files.\n", exitCodeIOError)
	fmt.Printf("  - %v: Unsupported YM version.\n", exitCodeUnsupportedVersion)
	fmt.Println("")

	fmt.Println("WARNING: This tool doesn't unpack LHA/LZH compressed YM files.")
	fmt.Println("The YM file is expected to be already unpacked.")

}

func StreamCmd() {
	ports := eports()
	if ports == nil {
		fmt.Println("No serial ports found!")
		os.Exit(exitCodeCmdParsing)
	}

	if len(flag.Args()) < 2 {

		fmt.Println("YM Tool - Stream a song")
		fmt.Println("")
		fmt.Println("Usage: ym stream <file.ym>")

		os.Exit(exitCodeCmdParsing)
	}

	// open data file
	filePath := flag.Arg(1)
	ymFile, err := ym.NewFile(filePath, false)
	switch err := err.(type) {
	case nil:
	case ym.UnsupportedVersionError:
		fmt.Printf("Error opening file %v: %v\n", filePath, err)
		os.Exit(exitCodeUnsupportedVersion)
	default:
		fmt.Printf("Error opening file %v: %v\n", filePath, err)
		os.Exit(exitCodeIOError)
	}

	fmt.Printf("Data for %v = %d\n\n", filePath, len(ymFile.Frames))

	isInterleaved := ymFile.Header.Attributes&0x01 != 0
	if !isInterleaved {
		fmt.Printf("Only interleaved files permitted %v\n", filePath)
		os.Exit(exitCodeIOError)
	}

	// ******************
	// optionally, takes an "interleaved file" where the data is organised by register so all R0's first then all the R1's ... etc
	// and changes this so that the data is reorganised frame by frame so we have R0,R1.. of frame 0 then the same of frame 1 and so on.
	// ******************

	frameCount := int(ymFile.Header.FrameCount)
	fmt.Printf("FrameCount %d\n", frameCount)

	var dataByFrame [][16]byte
	if isInterleaved {

		// split to 16 equal chunks sequentially
		var regBlocks [16][]byte
		//for i := 0; i < len(ymFile.Frames); i += frameCount {
		for reg := 0; reg < 16; reg++ {
			start := reg*frameCount
			regData := ymFile.Frames[start:start+frameCount]
			//regBlocks = append(regBlocks, regData)
			regBlocks[reg] = regData
		}

		// zip together
		for frameNo := 0; frameNo < frameCount; frameNo++ {
			var frameData[16]byte

			for regNo := 0; regNo < 16; regNo++ {
				var regByte = regBlocks[regNo][frameNo]
				if regNo == 7 {
					// force IO ports as input - save power??
					regByte = regByte | 0b11000000
				}
				frameData[regNo] = regByte
			}
			dataByFrame = append(dataByFrame, frameData)
		}
	} else {
		os.Exit(1)
	}

	//***************************
	// stream data to the device
	//***************************
	// http://leonard.oxg.free.fr/ymformat.html
	// ftp://ftp.modland.com/pub/documents/format_documentation/Atari%20ST%20Sound%20Chip%20Emulator%20YM1-6%20(.ay,%20.ym).txt
	port := openCom()
	time.Sleep(2 * time.Second)

	log.Printf("Writing %s", com)

	// connecting USB I think does a reset - so give a little pause for that to occur

	// send data
	frameStart := time.Now()
	frameRate := ymFile.Header.PlayerFrame
	frameDurationMs := float64((1 / float64(frameRate)) * 1000)

	//timeSecs := ymFile.Header.FrameCount / uint32(ymFile.Header.PlayerFrame)
	//song_mins := timeSecs / 60
	//song_secs := timeSecs % 60

	//last_sec := 0


//		fmt.Printf("r%-2d=%08b ", regNo, b);
	for frameNo := 0; frameNo < frameCount; frameNo++ {
		frameStart = time.Now()

		frameData := dataByFrame[frameNo] 
		slice := frameData[:]
		port.Write(slice)
		port.Flush()

		for regNo:=0; regNo<16; regNo++ {
			fmt.Printf("%08b ", frameData[regNo])
		}
		fmt.Printf("\n")

		elapsedMs := time.Now().Sub(frameStart).Milliseconds()
		remainingMs := frameDurationMs - float64(elapsedMs)
		NanoSleep(time.Duration(remainingMs) * time.Millisecond)

		// print frame rate
		/*
		if 1==2 && last_sec != secs {
			e_timeSecs := (frameNo) / int(frameRate)
			mins := e_timeSecs / 60
			secs := e_timeSecs % 60
			log.Printf("Playing %02d:%02d of %02d:%02d\n", mins, secs, song_mins, song_secs)
			last_sec = secs
		}
		*/

	}
	port.Close()

}

func PortCmd() {
	ports := eports()
	if ports == nil {
		os.Exit(exitCodeCmdParsing)
	}
	fmt.Println(ports[0].Name)
}

// spin wait as it's more accurate than sleep
func NanoSleep(n time.Duration) {

	for target := time.Now().Add(n); time.Now().Before(target); {
		// pass
	}
}

// InfoCmd shows information about a YM file
func InfoCmd() {
	if len(flag.Args()) < 2 {

		fmt.Println("YM Tool - Show information about a song")
		fmt.Println("")
		fmt.Println("Usage: ym info <file.ym>")

		os.Exit(exitCodeCmdParsing)
	}

	filePath := flag.Arg(1)
	ymFile, err := ym.NewFile(filePath, false)
	switch err := err.(type) {
	case nil:
	case ym.UnsupportedVersionError:
		fmt.Printf("Error opening file %v: %v\n", filePath, err)
		os.Exit(exitCodeUnsupportedVersion)
	default:
		fmt.Printf("Error opening file %v: %v\n", filePath, err)
		os.Exit(exitCodeIOError)
	}

	fmt.Printf("Information for %v:\n\n", filePath)
	fmt.Println(ymFile.Header)
}

func openCom() *serial.Port {

	flag.Parse()

	// Program Info
	log.Printf("Starting with port %s at baud rate %d", com, baud)

	// Setup Serial Port
	c := &serial.Config{Name: com, Baud: baud}
	s, err := serial.OpenPort(c)
	if err != nil {
		fmt.Printf("Cannot open port: %s\n", com)
		listPorts()
		log.Fatal(err, " : ", com)
	}

	log.Printf("Opened %s", com)
	return s
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

/*
func ports() {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}
	for _, port := range ports {
		fmt.Printf("Found port: %v\n", port)
	}
}
*/

func eports() []*enumerator.PortDetails {

	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		return nil
	}
	return ports
}

func listPorts() {
	ports := eports()

	for _, port := range ports {
		fmt.Printf("Found port: %s\n", port.Name)
		if port.IsUSB {
			fmt.Printf("   USB ID     %s:%s\n", port.VID, port.PID)
			fmt.Printf("   USB serial %s\n", port.SerialNumber)
		}
	}
}
