FROM alpine as base-rg

ARG VERSION=11.0.1
ARG RGNAME=ripgrep-${VERSION}-x86_64-unknown-linux-musl
ARG URL=https://github.com/BurntSushi/ripgrep/releases/download/${VERSION}/${RGNAME}.tar.gz
RUN wget ${URL} -O /tmp/rg.tar.gz && \
    tar -xzvf /tmp/rg.tar.gz -C /tmp/ && \
    cp /tmp/${RGNAME}/rg /rg && \
    rm -rf /tmp/

# ---
FROM alpine as base-git

RUN apk add git

# ---
FROM base-git

COPY --from=base-rg /rg /usr/local/bin/

ADD findIt.sh /usr/local/bin/findIt.sh
RUN chmod +x /usr/local/bin/findIt.sh

ADD toodler_linux /usr/local/bin/toodler

CMD ["toodler"]
