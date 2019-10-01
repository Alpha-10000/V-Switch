// gcc send.c -o send

#include <arpa/inet.h>
#include <linux/if_packet.h>
#include <net/ethernet.h>
#include <net/if.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>

#define MAC_A 0x9C
#define MAC_B 0xB6
#define MAC_C 0xD0
#define MAC_D 0x9A
#define MAC_E 0xD7
#define MAC_F 0x45

#define MAX_SZ 1024

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
    ssize_t read = 0;
    while ( (read = getline(&line, &len, fp)) != -1 ) {
    }
    fclose(fp);

    unsigned char mac[6];
    sscanf(line, "%hhx:%hhx:%hhx:%hhx:%hhx:%hhx",
	   &mac[0], &mac[1], &mac[2], &mac[3], &mac[4], &mac[5]);

    char eth_buf[MAX_SZ];
    int eth_len = 0;
    memset(eth_buf, 0, MAX_SZ);
    struct ethhdr *eh = (struct ethhdr *)eth_buf;

    eh->h_dest[0] = mac[0];
    eh->h_dest[1] = mac[1];
    eh->h_dest[2] = mac[2];
    eh->h_dest[3] = mac[3];
    eh->h_dest[4] = mac[4];
    eh->h_dest[5] = mac[5];

    eh->h_source[0] = MAC_A;
    eh->h_source[1] = MAC_B;
    eh->h_source[2] = MAC_C;
    eh->h_source[3] = MAC_D;
    eh->h_source[4] = MAC_E;
    eh->h_source[5] = MAC_F;

    eh->h_proto = htons(ETH_P_ALL);

    eth_len += sizeof (struct ethhdr);
    eth_buf[eth_len++] = 0xAA;
    eth_buf[eth_len++] = 0xBB;
    eth_buf[eth_len++] = 0xCC;
    eth_buf[eth_len++] = 0xDD;
    eth_buf[eth_len++] = 0xEE;
    eth_buf[eth_len++] = 0xFF;

    struct sockaddr_ll sockaddr;
    memset(&sockaddr, 0, sizeof (struct sockaddr_ll));

    sockaddr.sll_family = AF_PACKET;
    sockaddr.sll_protocol = htons(ETH_P_ALL);
    sockaddr.sll_ifindex = if_nametoindex(if_name);
    sockaddr.sll_halen = ETH_ALEN;

    if ( sendto(sockfd, eth_buf, eth_len, 0, (struct sockaddr*)&sockaddr, sizeof(struct sockaddr_ll)) < 0 ) {
	printf("Impossible to send ethernet message\n");
	exit(1);
    }
    printf("[SENT] Message of size (%d)\n", eth_len);

    return 0;
}
