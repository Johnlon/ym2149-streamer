
ym streamer 
===========

This component reads "YM" formatted files and sends them over USB/Comm port to the Arduino.

The script build.bat will build the program "ym.exe".

Unpack the compressed "YM" file using 7Zip (for example) before using this streamer and inside it you'll find a file called something with a "YM" suffix too.
Use this file as the data for the streamer.


Example usage
---------

   $ ym  --com COM5 stream ChipTuneFile.YM

