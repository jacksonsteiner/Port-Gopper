# Port-Gopper

Proof-of-concept client to server communication over random emphemeral port sequence. Concept is demonstrated by breaking client message up into 10-byte udp datagrams and individually sent over a random ephemeral port by both the client and server.

Build with:
`
git clone git@github.com:jacksonsteiner/Port-Gopper.git
cd Port-Gopper
make
`
This builds 2 files in the new local `bin` directory. `Port-Gopper` is the server and `Port-Gopper-Client` is the client.