ym2149-streamer
===========

Should also work for ay3-8913 / ay3-8910 chips.

Adapted to work on Windows and WSL2.

This project provides a means to stream "YM" ChipTune files at a YM2149/AY3-8910 chip via an Arduino adapter that provides USB connectivity.

There are two parts to the project 

- an Arduino adaptor component that listens for messages on the USB Comm port and then writes these to a YM2149 chip

- a Windows (golang) Streamer component that reads "YM" formatted files and streams them to the Arduino via a USB Comms port


Arduino
-------

The Arduino component is based on https://github.com/FlorentFlament/ym2149-streamer/ which includes an Arduino component and also a streaming component written in Python..
I've included only the Arduino component here and left out the Python. The reason for this is that I couldn't get it to work in my environment and I dont' enjoy tht aspect of Python.
I found that the program seemed to work in Py2 on Linux but not at all on Windows. However I can't using Linux in my dev env as it's a WSL2 Linux and the USB comm ports don't work natively and trivially yet in WSL2. So everything to do with USB/Comm needs to happen in Windows.

So the instructions in the Arduino submodule explain building the sketch using a commment line compiler (no need for Arduino IDE) in Linux and then flashing the Ardino using Windows.


Streamer
------

The streamer component present in https://github.com/FlorentFlament/ym2149-streamer/ didn't work at all for me and so I wrote an alternative using Go.
This has been tested on Windows.



