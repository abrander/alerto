BINARY=alerto
SOURCES := $(shell find . -name '*.go')

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	go build -o ${BINARY} .

install:
	install --preserve-context alerto /usr/sbin/
	setcap cap_net_raw=ep /usr/sbin/alerto # allow icmp from unprivileged users

clean:
	rm -f ${BINARY}

uninstall:
	rm /usr/sbin/alerto

run: $(BINARY)
	sudo setcap cap_net_raw=ep alerto
	DEBUG=* ./alerto
