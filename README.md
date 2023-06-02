# Port-Gopper

A computer network implementation of [Frequency Hopping](https://en.wikipedia.org/wiki/Frequency-hopping_spread_spectrum).

Proof-of-concept client to server communication over random emphemeral port sequence. Concept is demonstrated by breaking client message up into 10-byte udp datagrams and individually sent over a random ephemeral port by both the client and server.

Build with:

    git clone https://github.com/jacksonsteiner/Port-Gopper.git
    cd Port-Gopper
    make

This builds 2 files in the new local `bin` directory. `Port-Gopper` is the server and `Port-Gopper-Client` is the client.