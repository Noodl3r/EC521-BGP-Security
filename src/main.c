#include "main.h"

/* Discard remaining characters up to and including the newline */
static void flush_buf(void) {
    int c;
    while ((c = getchar()) != '\n' && c != EOF);
}

static void pause_key(void) {
    printf(DIM "\n  [Press ENTER to continue...]" RST);
    flush_buf();
}

int main(void) {
    // init_sim();
    printf(GRN "Simulation initialised: 5 Autonomous Systems ready.\n" RST);
    printf("BGP sessions not yet established.  No routes advertised yet.\n\n");
    pause_key();

    int choice, running = 1;
    while (running) {
        show_menu();

        char input[100];

        if (fgets(input, sizeof(input), stdin) == NULL) {
            printf(RED "  Invalid input.\n" RST);
            pause_key();
            continue;
        }
        if (sscanf(input, "%d", &choice) != 1) {
            printf(RED "  Invalid input.\n" RST);
            pause_key();
            continue;
        }
        flush_buf();

        switch (choice) {
        case 1:
            pause_key();
            break;
        case 2:
            pause_key();
            break;
        case 3:
            pause_key();
            break;
        case 4:
            pause_key();
            break;

        case 5:
            pause_key();
            break;

        case 6:
            pause_key();
            break;

        case 7:
            pause_key();
            break;
        case 8:
            pause_key();
            break;
        case 9:
            pause_key();
            break;

        case 10:
            break;

        case 11:
            printf(CYN "\n  Goodbye!  BGP Simulation terminated.\n\n" RST);
            running = 0;
            break;

        default:
            printf(RED "  Unknown option.\n" RST);
            pause_key();
        }
    }
    return 0;
}