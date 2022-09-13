
make
port=$(/mnt/c/Users/johnl/work/github/johnlon/ym2149-streamer/streamer/ym.exe port)

if [ $? -ne 0 ]; then
	  echo cant find port
	    exit 1
fi
./avrdude/avrdude.exe -Cavrdude/avrdude.conf  -v -V -patmega328p -carduino -P$port -b115200 -D -Uflash:w:main.hex:i
