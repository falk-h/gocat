package main

import (
        goopt "github.com/droundy/goopt"
        "fmt"
        "bufio"
        "time"
        "os"
        "math/rand"
)

var output *bufio.Writer
var startColor uint8
var freq = uint8(2)
var offset = uint8(2)
var printPlain = false
var printNumbers = false
var printNonblankNumbers = false
var lineNumber = 1
const resetSequence = "\u001B[0m" // Escape sequence to reset terminal colors.
const invertSequence = "\u001B[7m" // Escape sequence to invert terminal colors.
const backSequence = "\u001B[F" // Escape sequence to go back to the first collumn.
const version = "0.1"

func main() {
        info, _ := os.Stdout.Stat()
        printPlain = (info.Mode() & os.ModeCharDevice) != os.ModeCharDevice // True if output is redirected.

        output = bufio.NewWriter(os.Stdout) // Use buffered output.
        defer output.Flush() // Flush the output later so no characters get lost.

        goopt.RequireOrder = true // Require flags to come before filenames.
        goopt.Usage = func() string {return rainbowStrings(helpText)}
        goopt.Version = rainbowStrings([]string{version})

        // Flags
        freqFlag := goopt.Int([]string{"-F", "--freq"}, 2, "color frequency")
        offsetFlag := goopt.Int([]string{"-O", "--offset"}, 2, "vertical offset")
        invertFlag := goopt.Flag([]string{"-i", "--invert"}, []string{}, "invert output", "")
        animateFlag := goopt.Flag([]string{"-a", "--animate"}, []string{}, "animate output", "")
        seedFlag := goopt.Int([]string{"-S", "--seed"}, 0, "RNG seed")
        durationFlag := goopt.Int([]string{"-d", "--duration"}, 12, "animation duration")
        speedFlag := goopt.Int([]string{"-s", "--speed"}, 20, "animation speed")
        numberFlag := goopt.Flag([]string{"-n", "--number"}, []string{}, "number output lines", "")
        numberNonblankFlag := goopt.Flag([]string{"-b", "--number-nonblank"}, []string{}, "number nonempty output lines", "")
        forceFlag := goopt.Flag([]string{"-f", "--force"}, []string{}, "force color output", "")
        //forceFlag := goopt.Flag([]string{"-f", "--force"}, []string{}, "force color output", "")

        goopt.Parse(func() []string { return []string{} }) // Parse with no extra arguments.

        printPlain = !*forceFlag && printPlain // Add in the -f flag.
        printNumbers = *numberFlag
        printNonblankNumbers = *numberNonblankFlag

        if *seedFlag == 0 {
                rand.Seed(time.Now().UnixNano()) // Seed RNG with current time.
        } else {
                rand.Seed(int64(*seedFlag)) // Seed with specified seed.
        }

        startColor = uint8(rand.Intn(256)) // Start at random color.

        if *invertFlag {
                output.WriteString(invertSequence) // Invert output.
        }

        freq = uint8(*freqFlag % 256)
        offset = uint8(*offsetFlag % 256)
        var filenames []string

        if  len(goopt.Args) == 0 {
                filenames = []string{"-"} // Read from stdin if no files are specified.
        } else {
                filenames = goopt.Args
        }

        files := make([]*bufio.Reader, len(filenames))

        for i,filename := range filenames {
                if filename == "-" {
                        files[i] = bufio.NewReader(os.Stdin) // Read from stdin on -
                } else {
                        file, err := os.Open(filename)
                        defer file.Close()
                        // Print err to sdterr if the file can't be opened.
                        if err != nil {
                                fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err.Error())
                        }
                        files[i] = bufio.NewReader(file)
                }
        }

        // Print plainly if stdout is redirected
        if printPlain {
                for _,file := range files {
                        output.ReadFrom(file)
                }
                return
        }

        // Loop over files.
        for _,file := range files {
                // Loop over lines in the file.
                for {
                        line,err := file.ReadString('\n') // Read the whole line.
                        if err != nil {
                                break // Break on EOF.
                        }

                        lineStart := startColor
                        printLine(line, lineStart)

                        // Animate if -a was passed.
                        if *animateFlag && line != "\n" {
                                for i := 0; i < *durationFlag; i++ {
                                        output.Flush() // Flush the output to avoid printing things in the wrong order.
                                        time.Sleep(time.Second / time.Duration(*speedFlag))
                                        output.WriteString(backSequence) // Go back to the first collumn.
                                        lineStart = startColor + uint8(i) * freq * offset * 2
                                        printLine(line, lineStart)
                                }
                        }

                        startColor += offset
                }
        }

        output.WriteString(resetSequence) // Reset output color.
}

var colors [256]string

// Converts a uint8 from HSV to RGB and generates the escape sequence.
// n corresponds to the HVS hue.
func color(n uint8) string {
        if colors[n] != "" {
                return colors[n] // Return cached result, if it exists.
        } else {
                h := n / 43
                f := n - 43 * h
                t := f * 6
                q := 255 - t

                switch h {
                case 0:
                        colors[n] = fmt.Sprintf("\u001B[38;2;255;%d;0m", t)
                case 1:
                        colors[n] = fmt.Sprintf("\u001B[38;2;%d;255;0m", q)
                case 2:
                        colors[n] = fmt.Sprintf("\u001B[38;2;0;255;%dm", t)
                case 3:
                        colors[n] = fmt.Sprintf("\u001B[38;2;0;%d;255m", q)
                case 4:
                        colors[n] = fmt.Sprintf("\u001B[38;2;%d;0;255m", t)
                default:
                        colors[n] = fmt.Sprintf("\u001B[38;2;255;0;%dm", q)
                }
        }
        return colors[n]
}

// Prints a string to stdout, starting at the specified color.
func printLine(line string, startColor uint8) {
        if printNonblankNumbers && line != "\n" {
                line = numberString(lineNumber) + line
                lineNumber++
        } else if printNumbers && !printNonblankNumbers {
                line = numberString(lineNumber) + line
                lineNumber++
        }
        for _,r := range line {
                output.WriteString(color(startColor)) // Print escape sequence
                output.WriteRune(r) // Print rune
                startColor += freq; // Increment color
        }
}

func numberString(n int) string {
       number := fmt.Sprintf("%d", lineNumber)
       for i := 8; i > len(number); i-- {
               number = " " + number
       }
       return number + "  "
}

// Rainbowifies a slice of strings.
// Used for help and version text.
func rainbowStrings(s []string) string {
        text := ""
        if printPlain {
                for _,line := range s {
                        text = text + line
                }
        } else {
                rand.Seed(time.Now().UnixNano())
                startColor = uint8(rand.Intn(256))
                for _,line := range s {
                        lineStart := startColor
                        for _,r := range line {
                                text = text + color(lineStart) + fmt.Sprintf("%c", r)
                                lineStart++
                        }
                        startColor++
                }
        }
        text = text + resetSequence
        return text
}

// Help text.
var helpText = []string{
"Usage: gocat [OPTION]... [FILE]...\n",
"\n",
"With no FILE, or when FILE is -, read standard input.\n",
"\n",
"  -a, --animate            animate the output\n",
"  -d, --duration=<d>       animation duration (default: 12)\n",
"  -f, --force              force color output\n",
"  -F, --freq=<f>           ranbow frequency (default: 2)\n",
"  -i, --invert             invert the output\n",
"  -n, --number             number all output lines\n",
"  -O, --offset=<o>         vertical offset (default: 2)\n",
"  -S, --seed=<s>           RNG seed, 0 means random (default: 2)\n",
"  -s, --speed=<s>          animation speed (default: 20)\n",
"  -h, --help               display this help text\n",
"      --version            display version information\n"}
