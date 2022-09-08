#include "common.h"
#include <stdlib.h>


//FiSim bootloader simulator "magical" function; the name is important - please do not touch
//it always returns "!! Pwned boot !!" - this is the password that is beging guessed
//the name is the hardcoded in the FISim engine
#define readPassword flash_load_img
void __attribute__ ((noinline)) flash_load_img(void* base_addr, size_t* size_ptr) {
	*size_ptr = 16;
	
	// The engine will "magically" write some data to this location based on the hardware model
};


//FiSim function to indicate successful glitch
void __attribute__ ((noinline)) authentication_successful (void) {
	// Indicate we successfully bypassed the signature verification
	__SET_SIM_SUCCESS();
}

void entry() {
	 // So we fucking have a termination character
	unsigned char real_password[17]="Test Payload!!!!";
	volatile bool notCheck = false;
	
	//unsigned char provided_password[16]="Test Payload!!!!";
	unsigned char* provided_password = GUESSED_PASSWORD_MEMORY;
	size_t provided_size;

	if (notCheck) 
		goto authenticated;

	serial_puts("Start Password Checking...\n");
	
	readPassword(provided_password, &provided_size);
	
	int authenticate = 0;
	
	//please put your comparison code here:
	// My comparison code
	// I guess
	// Wow
	// It's so good
	// Right?
	
	int index = 0; // The index, big souprise
	while(real_password[index] != '\0' && provided_password[index] != '\0') { // We check if both the passwords are not yet the termination character
	
		// Just to be sure, this is not needed for normal C code, but it needs to be here.
		if (real_password[index] != provided_password[index]) { 
			break;
		}
	
		index += 1; // We increment
	}
	
	// Doing it twice prevents the fault injection
	if (real_password[index] != provided_password[index]) { // Comparison, wowzers
		authenticate = 1;
	}
	
	if (real_password[index] != provided_password[index]) { // Comparison, wowzers
		authenticate = 1;
	}
	
	//end

	if(authenticate) {
		serial_puts("Auth failed!\n");
		__SET_SIM_FAILED();
	}
	
	if(authenticate) {
		serial_puts("Auth failed!\n");
		__SET_SIM_FAILED();
	}


authenticated:
	serial_puts("Authenticated!\n");

	authentication_successful();
}

//if you find FI issues in the entry.S - please ignore it. It is an artifact of using FiSim for bootloaders. 