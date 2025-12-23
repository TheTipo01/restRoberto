FROM golang:alpine AS build

COPY . /restRoberto

WORKDIR /restRoberto
RUN GOOS=windows GOARCH=386 CGO_ENABLED=0 go mod download
RUN GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o "restRoberto.exe"

FROM debian:trixie-slim

RUN mkdir -pm755 /etc/apt/keyrings

RUN dpkg --add-architecture i386 && apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates wget curl gpg && \
    curl -fsSL https://dl.winehq.org/wine-builds/winehq.key | gpg --dearmor -o /etc/apt/keyrings/winehq-archive.key && \
    wget -NP /etc/apt/sources.list.d/ https://dl.winehq.org/wine-builds/debian/dists/trixie/winehq-trixie.sources && \
    apt-get update && \
    apt-get install -y --no-install-recommends winehq-stable xvfb xauth unzip && \
    rm -rf /var/lib/apt/lists/* && \
    apt-get clean

RUN mkdir /roberto

COPY ./roberto_setup.exe /roberto

# Supresses wine gui popups
ENV WINEARCH=win32
ENV WINEDLLOVERRIDES="mscoree,mshtml="

# Install Speech SDK dependency (for SAPI support)
RUN curl -o /usr/bin/winetricks https://raw.githubusercontent.com/Winetricks/winetricks/master/src/winetricks && \
    chmod +x /usr/bin/winetricks
RUN xvfb-run -a winetricks -q speechsdk

# Install Loquendo Roberto and patch the DLL
RUN xvfb-run -a wine /roberto/roberto_setup.exe /SILENT && wineserver -w && rm /roberto/roberto_setup.exe && rm "/root/.wine/drive_c/Program Files/Loquendo/LTTS/LoqTTS6.dll"
ARG target="/root/.wine/drive_c/Program Files/Loquendo/LTTS/LoqTTS6.dll"
COPY ./LoqTTS6.dll ${target}

RUN mkdir /root/.wine/drive_c/roberto
WORKDIR /root/.wine/drive_c/roberto

# Dependencies: Balcon
RUN wget https://www.cross-plus-a.com/balcon.zip && unzip balcon.zip && rm balcon.zip && rm *.txt

COPY --from=build /restRoberto/restRoberto.exe /root/.wine/drive_c/roberto
RUN chmod +x /root/.wine/drive_c/roberto/restRoberto.exe

CMD ["sh", "-c", "cd /root/.wine/drive_c/roberto && xvfb-run -a wine cmd /c C:\\\\roberto\\\\restRoberto.exe"]