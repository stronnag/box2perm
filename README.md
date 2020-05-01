box2perm - inav 2.5.0 AUX update
================================

# Introduction

ianv 2.5 changes the mode ID from using a `boxId` to a `permanentId`.

* `boxid`s are not guaranteed to be constant between releases
* `permid`s are guaranteed to be constant between releases (clue's in the name).

Thus, this is a one-off change that is intrusive and disruptive for inav 2.5 but will result in a better user experience in the future.

# Update required

If inav CLI `diff` or `dump` files are used to migrate setting from an earlier firmware revision to inav 2.5, then the user will have to either:

* Remove all `aux` settings and manually recreate them; or
* Update the source `diff` or `dump`

# using `box2perm`

`box2perm` is a simple command line application that migrates pre-2.5 `diff` or `dump` files to 2.5 format.

Binaries are provides for Linux (ia32, x86-64, arm7), FreeBSD (x86-64), Windows (win32) and MacOS (darwin) in the release area. `box2perm` has no external dependencies and should build /run on any platform where Go is available.

`box2perm` takes two parameters, the input file and the output file. If either of these is `-` or missing, then `stdin` / `stdout` are used.

## example

```
## migrate an inav 2.4 diff (or dump) to 2.5
$ box2perm mydiff-2.4.txt mydiff-2.5.txt
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
The addition text will be ignored by the inav CLI and provides verification to the user that the mode lines have been updated. The final line of output will be as:

```
### inav 2.5 aux conversion by box2perm 2020-05-01T14:08:03+0100 ###
```

# Author  and licence

`box2perm` is (c) Jonathan Hudson 2020

License : MIT / Public domain / BSD / WTF i.e. what ever is the most permissive in your locale.

E&OE. No warranty
