# use make and then a windows version of avrdude as windows COM works
make
./avrdude/avrdude.exe -Cavrdude/avrdude.conf  -v -V -patmega328p -carduino -PCOM3 -b115200 -D -Uflash:w:main.hex:i
