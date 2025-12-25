FROM golang:alpine AS build

COPY . /restRoberto

WORKDIR /restRoberto
RUN GOOS=windows GOARCH=386 CGO_ENABLED=0 go mod download
RUN GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o "restRoberto.exe"

FROM debian:trixie-slim

RUN mkdir -pm755 /etc/apt/keyrings && \
    dpkg --add-architecture i386 && apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates wget gpg && \
    wget -qO - https://dl.winehq.org/wine-builds/winehq.key | gpg --dearmor -o /etc/apt/keyrings/winehq-archive.key && \
    wget -NP /etc/apt/sources.list.d/ https://dl.winehq.org/wine-builds/debian/dists/trixie/winehq-trixie.sources && \
    apt-get update && \
    apt-get install -y --no-install-recommends winehq-stable xvfb xauth unzip && \
    apt-get purge -y gpg && \
    rm -rf /var/lib/apt/lists/*

RUN mkdir /roberto

# Supresses wine gui popups and sets architecture to win32
ENV WINEARCH=win32
ENV WINEDLLOVERRIDES="mscoree,mshtml="

# Install Speech SDK dependency (for SAPI support)
RUN wget -O /usr/bin/winetricks https://raw.githubusercontent.com/Winetricks/winetricks/master/src/winetricks && \
    chmod +x /usr/bin/winetricks && \
    xvfb-run -a winetricks -q speechsdk && \
    rm /usr/bin/winetricks

# Install Loquendo Roberto, Loquendo Paola voices and apply the DLL patch
RUN wget -O /roberto/roberto_setup.exe https://archive.org/download/loquendo-6-voices-pack-multilanguage-with-crack-dll/Languages%20and%20voices/Raffaele%20-%20Italian.exe &&  \
    xvfb-run -a wine /roberto/roberto_setup.exe /SILENT &&  \
    wineserver -w &&  \
    wget -O /roberto/paola_setup.exe https://archive.org/download/loquendo-6-voices-pack-multilanguage-with-crack-dll/Languages%20and%20voices/Paola%20-%20Italian.exe &&  \
    xvfb-run -a wine /roberto/paola_setup.exe /SILENT &&  \
    wineserver -w &&  \
    rm /roberto/paola_setup.exe /roberto/roberto_setup.exe &&  \
    rm "/root/.wine/drive_c/Program Files/Loquendo/LTTS/LoqTTS6.dll" && \
    wget -O "/root/.wine/drive_c/Program Files/Loquendo/LTTS/LoqTTS6.dll" https://archive.org/download/loquendo-6-voices-pack-multilanguage-with-crack-dll/Crack%20dll/LoqTTS6.dll

RUN mkdir /root/.wine/drive_c/roberto
WORKDIR /root/.wine/drive_c/roberto

# Download and install Balcon
RUN wget https://www.cross-plus-a.com/balcon.zip && unzip balcon.zip && rm balcon.zip && rm *.txt

COPY --from=build /restRoberto/restRoberto.exe /root/.wine/drive_c/roberto
RUN chmod +x /root/.wine/drive_c/roberto/restRoberto.exe

CMD ["sh", "-c", "cd /root/.wine/drive_c/roberto && xvfb-run -a wine cmd /c C:\\\\roberto\\\\restRoberto.exe"]