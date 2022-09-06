
ym streamer 
===========

This component reads "YM" formatted files and sends them over USB/Comm port to the Arduino.

The script build.bat will build the program "ym.exe".

Unpack the compressed "YM" file using 7Zip (for example) before using this streamer and inside it you'll find a file called something with a "YM" suffix too.
Use this file as the data for the streamer.

Build it
--------

I built it using Go on Windows.

     build.bat


Example usage
---------

    ym.exe --com com5 stream TheJetsons.YM

    Data for TheJetsons.YM = 61792

    FrameCount 3862
    2022/09/06 01:22:59 Starting with port com5 at baud rate 115200
    2022/09/06 01:22:59 Opened com5
    2022/09/06 01:22:59 Writing com5
    2022/09/06 01:23:02 Playing 00:01 of 01:17
    2022/09/06 01:23:04 Playing 00:02 of 01:17
    2022/09/06 01:23:05 Playing 00:03 of 01:17
