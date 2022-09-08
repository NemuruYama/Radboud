#include "common.h"

void __attribute__ ((noinline)) serial_putc(char c) {
	*(char*)UART_OUT_BUF_ADDR = c;
}

void serial_puts (char* s) {
	while (*s) {
		serial_putc(*s);
		
		s++;
	}
}



