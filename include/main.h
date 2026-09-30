#pragma once
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// Platform clear-screen command
#ifdef _WIN32
#  define CLRCMD "cls"
#else
#  define CLRCMD "clear"
#endif

// Colors
#define RST  "\033[0m"
#define BOLD "\033[1m"
#define DIM  "\033[2m"
#define RED  "\033[91m"
#define GRN  "\033[92m"
#define YLW  "\033[93m"
#define BLU  "\033[94m"
#define MAG  "\033[95m"
#define CYN  "\033[96m"
#define WHT  "\033[97m"

// Compile-time limits
#define MAX_AS     10   // Total autonomus systems allows
#define MAX_PEERS   8   // Maximum BGP peers per AS
#define MAX_ROUTES 60   // Maximum routing table entries per AS
#define MAX_PATH   10   // Maximum AS_PATH length
#define PFXSZ      26   // Prefix string size e.g. "203.0.113.0/24"
#define NMSZ       40   // Name string size
#define PFXMAX      5   // Owned prefixes per AS
#define LOGSZ     200   // Event log capacity

void show_menu(void);

// Packet being traced through the network
typedef struct {
    char src[PFXSZ];     // Source IP
    char dst[PFXSZ];     // Destination prefix                 
    int  src_asn;        // Source AS number                   
    int  hops[MAX_PATH]; // AS numbers visited along the path  
    int  hopcnt;         // Total hops taken                   
    int  dropped;        // 1 = dropped, 0 = delivered         
    char reason[100];    // Reason for drop 
} Packet;
