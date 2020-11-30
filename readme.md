# DevSnd

Golang SDL emulation of the `> /dev/snd` pipe in Linux (for MacOs users).

## Installation

```sh
$ go get -u github.com/ganglio/devsnd
```

## Usage

`devsnd` accepts a stream from `stdin` and plays it.

Example:

```sh
$ ps ax | devsnd
```

My usage is to follow slow logs with lots of identical entries and rare occurrences of different lines.

The idea comes from the soviet nuclear submarines. In the control rooms for the nuclear reactors there were lots of dials showing different measures.
Because of the very large numbers of them it was impossible for the engineers to follow all of them at the same time.
So each dial used to click at a different rate/pitch depending on its value.
This was the control room was filled with a regular cacophony of beeps and clicks.
As it was regular, the human brain filters it out after a while (like the noise on a plane) but, if one of the dials changed it's value, the global symphony changed and the engineers spotted it.

Same for the logs. You listen to it and it always sounds the same, but, if all of a sudden it changes, you notice it straight away.

## Flags

`devsnd` accepts the following flags:


```
  -16
        Use 16-bit samples
  -c uint
        Number of channels (default 1)
  -d    Debug mode
  -f int
        Sampling frequency (default 44000)
  -s uint
        Buffer size in bytes (default 512)
```
