


#include "uart.h"
#include "ym2149.h"


// Digital port 13 noted the board has an led on pin 13 already(should be yellow and be right below the pin)
#define LED PORTB5  
void set_led_out(void) {
  DDRB |= 1 << LED;
}

void clear_registers(void) {
  int i;
  for (i=0; i<14; i++) {
    send_data(i, 0);
  }
}

int main() {
  unsigned int i;
  unsigned char data[16];

  set_led_out();
  set_ym_clock();
  set_bus_ctl();
  initUART();
  clear_registers();

// JL
send_data(7, 0xf8); // bits 7:6=1 set IO ports as output, bits 5:0 turn off noise
send_data(8, 0x0f); // set all voices to constant tone - no envelope
send_data(9, 0x0f); // set all voices to constant tone - no envelope
send_data(10, 0x0f);// set all voices to constant tone - no envelope
// JL

  for/*ever*/(;;) {
    for (i=0; i<16; i++) {
      data[i] = getByte();
    }

    // Working around envelope issue (kind of). When writing on the
    // envelope shape register, it resets the envelope. This cannot be
    // properly expressed with the YM file format. So not using it.
    // Thanks sebdel: https://github.com/sebdel
    for (i=0; i<13; i++) {
      send_data(i, data[i]);
    }

    /*
    send_data(0, data[0]); // a fine
    send_data(1, data[1]); // a coarse
    send_data(2, data[2]); // b
    send_data(3, data[3]); // b
    send_data(4, data[4]); // c
    send_data(5, data[5]); // c
    send_data(6, data[6]); // noise freq
//    send_data(8, 0);
//    send_data(9, 0);
//    send_data(10, 0);

    send_data(11, data[11]); // freq of envelope fine
    send_data(12, data[12]); // freq of envelope rough
    //send_data(13, data[13]); // shape of env

    send_data(14, data[14]); // port a
    send_data(15, data[15]); // port b
*/

    // Have LED blink with noise (drums)
    //if (~data[7] & 0x38) {
    //  PORTB |=   1 << LED;
    //} else {
    //  PORTB &= ~(1 << LED);
   // }


    if (i%2==0) {
        PORTB |=   1 << 5;
      } else {
        PORTB &= ~(1 << 5);
      }
  }

  return 0;
}
