# Port-Gopper

A computer network implementation of [Frequency Hopping](https://en.wikipedia.org/wiki/Frequency-hopping_spread_spectrum).

More specifically, a proof-of-concept, client-server communication model over a random emphemeral port sequence. This method of communication is demonstrated by this application by breaking the client message up into 10-byte UDP datagrams and sending them individually over a random ephemeral port known by the server at connection time.

Build with:

    git clone https://github.com/jacksonsteiner/Port-Gopper.git
    cd Port-Gopper
    make

This builds both the client and server into a new `bin` directory within the Port-Gopper folder. Run `Port-Gopper` to first start the server, then `Port-Gopper-Client` to specify the server IP, your desired port range to hop over, and a message to send. Both the client and server will output how many bytes they sent and over which port before the message is displayed in its entirety by the server.