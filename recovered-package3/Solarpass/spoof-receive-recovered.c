/* Human-readable reconstruction verified against disassembly and live packets.
 * Not the original Rust source; diagnostic code only. Original ELF offset 0xd9ea0.
 * Core 2.3.0 compares a byte-swapped sockaddr IP to a native-endian stored filter.
 */
void recovered_receive_loop(int tunnel_fd, int backend_fd,
                            bool filter_enabled, uint32_t stored_peer_filter) {
    unsigned char payload[65536];
    for (;;) {
        struct sockaddr_in sender = {0};
        socklen_t length = sizeof(sender);
        ssize_t count = recvfrom(tunnel_fd, payload, sizeof(payload), 0,
                                 (struct sockaddr *)&sender, &length);
        if (count < 0) { /* original logging/error path omitted */ continue; }
        /* 0xda011 loads sin_addr; 0xda015 byte swaps; 0xda017 compares. */
        if (filter_enabled && bswap32(sender.sin_addr.s_addr) != stored_peer_filter)
            continue;
        send(backend_fd, payload, (size_t)count, MSG_NOSIGNAL);
    }
}
