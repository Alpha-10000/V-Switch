// gcc send.c -o send

#include <arpa/inet.h>
#include <linux/if_packet.h>
#include <net/ethernet.h>
#include <net/if.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>

unsigned char MAC[6] = { 0xEA, 0xC0, 0x52, 0x37, 0xB5, 0x40 };
#define MAX_SZ 1024

union ethframe
{
  struct
  {
    struct ethhdr    header;
    unsigned char    data[ETH_DATA_LEN];
  } field;
  unsigned char    buffer[ETH_FRAME_LEN];
};

int main (int argc, char *argv[])
{
    if ( argc < 2 ) {
	printf("Usage: ./send INTF_NAME\n");
	exit(1);
    }
    char *if_name = argv[1];
    int sockfd = socket(AF_PACKET, SOCK_RAW, htons(ETH_P_ALL));
    if ( sockfd == -1 ) {
	printf("Can't create socket for interface %s\n", if_name);
	exit(1);
    }

    char addr_path[MAX_SZ];
    memset(&addr_path, 0, MAX_SZ);
    snprintf(addr_path, MAX_SZ, "/sys/class/net/%s/address", if_name);
    FILE *fp = fopen(addr_path, "r");
    if ( fp == NULL ) {
	printf("Couldn't open device address file\n");
        exit(1);
    }
    char *line = NULL;
    size_t len = 0;
    if ( getline(&line, &len, fp) == -1 ) {
	printf("Couldn't read mac address\n");
	exit(1);
    }
    fclose(fp);
    unsigned char mac[6];
    sscanf(line, "%hhx:%hhx:%hhx:%hhx:%hhx:%hhx",
	   &mac[0], &mac[1], &mac[2], &mac[3], &mac[4], &mac[5]);


    union ethframe frame;
    memcpy(frame.field.header.h_dest, MAC, sizeof (MAC));
    memcpy(frame.field.header.h_source, mac, sizeof (mac));
    frame.field.header.h_proto = htons(ETH_P_802_2);
    unsigned char payload[] = { 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF };
    memcpy(frame.field.data, payload, sizeof (payload));


    struct sockaddr_ll sockaddr;
    memset(&sockaddr, 0, sizeof (struct sockaddr_ll));

    sockaddr.sll_family = AF_PACKET;
    //sockaddr.sll_protocol = htons(ETH_P_ALL);
    sockaddr.sll_ifindex = if_nametoindex(if_name);
    sockaddr.sll_halen = ETH_ALEN;

    if ( sendto(sockfd, frame.buffer, ETH_FRAME_LEN, 0, (struct sockaddr*)&sockaddr, sizeof(struct sockaddr_ll)) < 0 ) {
	printf("Impossible to send ethernet message\n");
	exit(1);
    }

    int eth_len = sizeof (payload)+ ETH_ALEN;
    printf("[SENT] Message of size (%d)\n", eth_len);

    return 0;
}
