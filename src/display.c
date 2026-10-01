#include "main.h"

void show_menu(void) {
    printf(BOLD CYN
        "\n  +=============================================================+\n"
        "  |              BGP SIMULATION                |\n"
        "  +=============================================================+\n" RST);
    printf(CYN "  |" RST BOLD "  VIEW ──────────────────────────────────────────────────-   " CYN "|\n" RST);
    printf(CYN "  |" RST "   1)  Network Topology                                           " CYN "|\n" RST);
    printf(CYN "  |" RST "   2)  BGP Session Status                                         " CYN "|\n" RST);
    printf(CYN "  |" RST "   3)  BGP Routing Tables                                         " CYN "|\n" RST);
    printf(CYN "  |" RST "   4)  Event Log                                                  " CYN "|\n" RST);
    printf(CYN "  |" RST BOLD "  OPERATIONS ────────────────────────────────────────-----   " CYN "|\n" RST);
    printf(CYN "  |" RST "   5)  Establish BGP Sessions                                     " CYN "|\n" RST);
    printf(CYN "  |" RST "   6)  Run Route Advertisement                                    " CYN "|\n" RST);
    printf(CYN "  |" RST BOLD "  PACKET SIMULATION ──────────────────────────────────────   " CYN "|\n" RST);
    printf(CYN "  |" RST "   7)  Send a Packet                                              " CYN "|\n" RST);
    printf(CYN "  |" RST BOLD "  ATTACKS -----------─────────────────────────────────────   " CYN "|\n" RST);
    printf(CYN "  |" RST RED  "   8)  Simulate BGP Route Hijack                             " CYN "|\n" RST);
    printf(CYN "  |" RST YLW  "   9)  Simulate BGP Blackhole                                " CYN "|\n" RST);
    printf(CYN "  |" RST BOLD "  SYS ----------──────────────────────────────────────────   " CYN "|\n" RST);
    printf(CYN "  |" RST "  10)  Restore Network to Normal State                            " CYN "|\n" RST);
    printf(CYN "  |" RST "  11)  Exit                                                       " CYN "|\n" RST);
    printf(CYN "  +=============================================================+\n" RST);
    printf(BOLD "  Enter choice: " RST);
}
