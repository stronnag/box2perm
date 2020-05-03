/* bpx2perm
 * Convert pre-inav 2.5 diff or dump files to inav 2.5
 * updates the aux lines fron boxid to permid
 * this only needs to be done once
 *
 * https://github.com/iNavFlight/inav/pull/5654
 *
 * (c) Jonathan Hudson 2020
 * License : Public domain / BSD / MIT / WTF
 *           i.e. what ever is the most permissive in your locale
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"io"
	"log"
	"bufio"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const FLAGDAY = 4
const FLAGMON = 5

var box2perm = []struct {
        permid byte
        name string
}{
  { permid: 0, name: "ARM"},		// 0
  { permid: 1, name: "ANGLE"},		// 1
  { permid: 2, name: "HORIZON"},		// 2
  { permid: 3, name: "NAV ALTHOLD"},		// 3
  { permid: 5, name: "HEADING HOLD"},		// 4
  { permid: 6, name: "HEADFREE"},		// 5
  { permid: 7, name: "HEADADJ"},		// 6
  { permid: 8, name: "CAMSTAB"},		// 7
  { permid: 10, name: "NAV RTH"},		// 8
  { permid: 11, name: "NAV POSHOLD"},		// 9
  { permid: 12, name: "MANUAL"},		// 10
  { permid: 13, name: "BEEPER"},		// 11
  { permid: 15, name: "LEDLOW"},		// 12
  { permid: 16, name: "LIGHTS"},		// 13
  { permid: 36, name: "NAV LAUNCH"},		// 14
  { permid: 19, name: "OSD SW"},		// 15
  { permid: 20, name: "TELEMETRY"},		// 16
  { permid: 26, name: "BLACKBOX"},		// 17
  { permid: 27, name: "FAILSAFE"},		// 18
  { permid: 28, name: "NAV WP"},		// 19
  { permid: 29, name: "AIR MODE"},		// 20
  { permid: 30, name: "HOME RESET"},		// 21
  { permid: 31, name: "GCS NAV"},		// 22
  { permid: 38, name: "KILLSWITCH"},		// 23
  { permid: 33, name: "SURFACE"},		// 24
  { permid: 34, name: "FLAPERON"},		// 25
  { permid: 35, name: "TURN ASSIST"},		// 26
  { permid: 37, name: "SERVO AUTOTRIM"},		// 27
  { permid: 21, name: "AUTO TUNE"},		// 28
  { permid: 39, name: "CAMERA CONTROL 1"},		// 29
  { permid: 40, name: "CAMERA CONTROL 2"},		// 30
  { permid: 41, name: "CAMERA CONTROL 3"},		// 31
  { permid: 42, name: "OSD ALT 1"},		// 32
  { permid: 43, name: "OSD ALT 2"},		// 33
  { permid: 44, name: "OSD ALT 3"},		// 34
  { permid: 45, name: "NAV CRUISE"},		// 35
  { permid: 46, name: "MC BRAKING"},		// 36
  { permid: 47, name: "USER1"},		// 37
  { permid: 48, name: "USER2"},		// 38
  { permid: 32, name: "FPV ANGLE MIX"},		// 39
  { permid: 49, name: "LOITER CHANGE"},		// 40
  { permid: 50, name: "MSP RC OVERRIDE"},		// 41
}

var monstrings = map[string]int{
	"Jan": 1, "Feb": 2, "Mar": 3, "Apr": 4, "May": 5, "Jun": 6,
	"Jul": 7, "Aug": 8, "Sep": 9, "Oct": 10, "Nov": 11, "Dec": 12}


func openStdinOrFile(path string) (io.ReadCloser, error) {
        var err error
        var r io.ReadCloser

        if len(path) == 0 || path == "-" {
                r = os.Stdin
        } else {
                r, err = os.Open(path)
        }
        return r, err
}

func openStdoutOrFile(path string) (io.WriteCloser, error) {
        var err error
        var w io.WriteCloser

        if len(path) == 0 || path == "-" {
                w = os.Stdout
        } else {
                w, err = os.Create(path)
        }
        return w, err
}

var (
	force = flag.Bool("force", false, "Forces update regardless of version")
)

func main() {
	var inpfn=""
	var outfn=""

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: box2perm {infile|-} [outfile|-]\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	nargs := len(flag.Args())
	switch nargs {
	case 2:
		outfn = flag.Args()[1]
		fallthrough
	case 1:
		inpfn = flag.Args()[0]
	}

	input,err := openStdinOrFile(inpfn)
	if err == nil {
		defer input.Close()
	} else {
		log.Fatal("Can't open input file\n")
	}

	sb := strings.Builder{}

	r := regexp.MustCompile(`^# INAV\/\S+\s+(\d+)\.(\d+)\.\d+\s+(\S+)\s+(\d+)\s+\d+ `)

	init := false
	doconv := *force
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		if !doconv && !init {
			m := r.FindAllStringSubmatch(line,-1)
			if len(m) > 0 {
				major,_ := strconv.Atoi(m[0][1])
				minor,_ := strconv.Atoi(m[0][2])
				mon := monstrings[m[0][3]]
				day,_ := strconv.Atoi(m[0][4])
				doconv = (major == 2 && (minor < 5 || (minor == 5 && (mon < FLAGMON || (mon == FLAGMON && day < FLAGDAY)))))
				init = true
				if !doconv {
					fmt.Fprintln(os.Stderr,"No conversion required")
					return
				}
			}
		}
		if init && doconv && strings.HasPrefix(line, "aux") {
			a := strings.Split(line, " ")
			boxid,_ :=  strconv.Atoi(a[2])
			if boxid >= len(box2perm) {
				fmt.Fprintf(os.Stderr, "Invalid ID %d at \"%s\", no conversion performed\n",
					boxid, line)
				return
			}
			// only updates lines with 'set' values
			if !(a[4] == "900" && a[5] == "900") {
				p := box2perm[boxid]
				line = fmt.Sprintf("aux %s %d %s %s %s\t\t# %2d => %2d (%s)",
					a[1], p.permid, a[3], a[4], a[5], boxid, p.permid, p.name)
			}
		}
		fmt.Fprintln(&sb, line)
	}
	if doconv {
		fmt.Fprintf(&sb, "### inav 2.5 aux conversion by box2perm %s ###\n",
			time.Now().Format("2006-01-02T15:04:05-0700"))
	}

	output,err := openStdoutOrFile(outfn)
	if err == nil {
		output.Write([]byte(sb.String()))
		defer output.Close()
	}  else {
		log.Fatal("Can't open output file\n")
	}


}
