Title: ym2149-streamer

The ym2149-streamer allows sending Atari ST YM files to a YM2149 chip
driven by an Arduino.


I'm using a Windows machine so things are a little complicated.
USB support isn't native in WSL2 so there's messing around I can't be bothered with.
So I'll use WSL2 for the bits that are best done on Linux and use Windows for all the USB/Comm port stuff.

How to build the Arduino Hex file
=======

Built on Linux - I use WSL2.

Requirements
------------

The following libraries are required:

* gcc-avr
* avr-libc
* avrdude

Build Hex file
------------

    $ make

Upload to Arduino Nano
================

A Windows copy of 'avrdude' is sitting in this dir for convenience and the script "_program.sh_" is configured to upload the hex file to an Arduino nano.

Find out what Comm port your Arduino is on and edit program.sh accordingly.

    $ ./program.sh



More information
===================

More information can be found at blog:

* [Streaming music to YM2149F][1]
* [Driving YM2149F sound chip with an Arduino][2]
* [Arduino Hello World without IDE][3]

Besides, a video showing the [YM2149 & Arduino circuit playing a tune][4] is
available.


[1]: http://www.florentflament.com/blog/streaming-music-to-ym2149f.html
[2]: http://www.florentflament.com/blog/driving-ym2149f-sound-chip-with-an-arduino.html
[3]: http://www.florentflament.com/blog/arduino-hello-world-without-ide.html
[4]: https://www.youtube.com/watch?v=MTRJdDbY048
