box2perm - inav 2.5.0 AUX update
================================

# Introduction

ianv 2.5 changes the `aux` mode ID from using a `boxId` to a `permanentId`.

* `boxid` are not guaranteed to be constant between releases
* `permanentId` are guaranteed to be constant between releases (the clue is in the name).

**This is a one-off change that is intrusive and disruptive for inav 2.5 but will result in a more consistent user experience in the future.**

# `diff` / `dump` updates required

If inav CLI `diff` or `dump` files are used to migrate setting from an earlier firmware revision to inav 2.5, then the user will have to either:

* Remove all `aux` settings and manually recreate them; or
* Update the source `diff` or `dump`, either:
  * using `box2perm` a cross-platform, command line tool to automate the process; or
  * manually edit the file, using the information in the [table below](#additional-information).

# `box2perm`

`box2perm` is a simple command line application that migrates pre-2.5 `diff` or `dump` files to 2.5 format.

Binaries are provides for Linux (ia32, x86-64, arm7), FreeBSD (x86-64), Windows (win32) and MacOS (darwin) in the [release area](https://github.com/stronnag/box2perm/releases). `box2perm` has no external dependencies and should build /run on any platform where Go is available.

`box2perm` takes two parameters, the input file and the output file. If either of these is `-` or missing, then `stdin` / `stdout` are used.

## example

* Unzip the release archive for your platform
* The command line application (`box2perm`, `box2perm.exe` etc.) is extracted to a `box2perm` directory.
* Run the `box2perm` application specifying with the old dump/diff file name and the new output file name:

```
## migrate an inav 2.4 diff (or dump) to 2.5
$ cd box2perm
$ ./box2perm mydiff-2.4.txt mydiff-2.5.txt
```

## sample output

Given the 2.4 aux statements:

```
# aux
aux 0 0 5 1500 2100
aux 1 3 0 1300 1700
aux 2 3 2 1300 1700
aux 3 9 0 1300 1700
aux 4 8 0 1700 2100
aux 5 19 1 1700 2100
aux 6 35 2 1300 1700
aux 7 10 6 1450 2100
aux 8 14 3 1600 2100
```

The 2.5 `box2perm` output is:

```
# aux
aux 0 0 5 1500 2100         #  0 =>  0 (ARM)
aux 1 3 0 1300 1700         #  3 =>  3 (NAV ALTHOLD)
aux 2 3 2 1300 1700         #  3 =>  3 (NAV ALTHOLD)
aux 3 11 0 1300 1700        #  9 => 11 (NAV POSHOLD)
aux 4 10 0 1700 2100        #  8 => 10 (NAV RTH)
aux 5 28 1 1700 2100        # 19 => 28 (NAV WP)
aux 6 45 2 1300 1700        # 35 => 45 (NAV CRUISE)
aux 7 12 6 1450 2100        # 10 => 12 (MANUAL)
aux 8 36 3 1600 2100        # 14 => 36 (NAV LAUNCH)
```
The addition text after the upper range will be ignored by the inav CLI and provides verification to the user that the mode lines have been updated. The final line of `box2perm` output will be as:

```
### inav 2.5 aux conversion by box2perm 2020-05-01T14:08:03+0100 ###
```

# Author  and licence

`box2perm` is (c) Jonathan Hudson 2020

License : MIT / Public domain / BSD / WTF i.e. what ever is the most permissive in your locale.

E&OE. No warranty

# Additional information

The table illustrates the conversion between `boxId` and `permanentId`. This may be used to manually update a `diff` or `dump`.

| Box Id | Name            | Permanent Id |
| ------ | --------------- | ------------ |
| 0 | ARM | 0 |
| 1 | ANGLE | 1 |
| 2 | HORIZON | 2 |
| 3 | NAV ALTHOLD | 3 |
| 4 | HEADING HOLD | 5 |
| 5 | HEADFREE | 6 |
| 6 | HEADADJ | 7 |
| 7 | CAMSTAB | 8 |
| 8 | NAV RTH | 10 |
| 9 | NAV POSHOLD | 11 |
| 10 | MANUAL | 12 |
| 11 | BEEPER | 13 |
| 12 | LEDLOW | 15 |
| 13 | LIGHTS | 16 |
| 15 | OSD SW | 19 |
| 16 | TELEMETRY | 20 |
| 28 | AUTO TUNE | 21 |
| 17 | BLACKBOX | 26 |
| 18 | FAILSAFE | 27 |
| 19 | NAV WP | 28 |
| 20 | AIR MODE | 29 |
| 21 | HOME RESET | 30 |
| 22 | GCS NAV | 31 |
| 39 | FPV ANGLE MIX | 32 |
| 24 | SURFACE | 33 |
| 25 | FLAPERON | 34 |
| 26 | TURN ASSIST | 35 |
| 14 | NAV LAUNCH | 36 |
| 27 | SERVO AUTOTRIM | 37 |
| 23 | KILLSWITCH | 38 |
| 29 | CAMERA CONTROL 1 | 39 |
| 30 | CAMERA CONTROL 2 | 40 |
| 31 | CAMERA CONTROL 3 | 41 |
| 32 | OSD ALT 1 | 42 |
| 33 | OSD ALT 2 | 43 |
| 34 | OSD ALT 3 | 44 |
| 35 | NAV CRUISE | 45 |
| 36 | MC BRAKING | 46 |
| 37 | USER1 | 47 |
| 38 | USER2 | 48 |
| 40 | LOITER CHANGE | 49 |
| 41 | MSP RC OVERRIDE | 50 |
