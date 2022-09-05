rem Found the right args by turning on the verbose logging in the Arduino IDE where I saw this..
rem C:\Program Files\WindowsApps\ArduinoLLC.ArduinoIDE_1.8.57.0_x86__mdqgnx93n4wtt\hardware\tools\avr/bin/avrdude -CC:\Program Files\WindowsApps\ArduinoLLC.ArduinoIDE_1.8.57.0_x86__mdqgnx93n4wtt\hardware\tools\avr/etc/avrdude.conf -v -patmega328p -carduino -PCOM3 -b115200 -D -Uflash:w:C:\Users\johnl\AppData\Local\Temp\arduino_build_90162/integrated-circuit-tester.ino.hex:i 


.\avrdude\avrdude.exe -Cavrdude\avrdude.conf  -v -V -patmega328p -carduino -PCOM3 -b115200 -D -Uflash:w:main.hex:i
