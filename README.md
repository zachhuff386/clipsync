# clipsync: Clipboard sync

[Clipsync](https://github.com/zachhuff386/clipsync)
is a high performance clipboard sharing application for linux. The clipboard is
shared over the network with support for string clipboard content up to several
megabytes in size. Works with both Wayland and X servers.

## Building

Requires Golang, X11 and XFixes libraries.

```bash
sudo dnf install libX11-devel libXfixes-devel xorg-x11-proto-devel
go get github.com/zachhuff386/clipsync
cd ~/go/src/github.com/zachhuff386/clipsync
CGO_ENABLED=1 go build -v
```

## Computer 111.111.111.111

```bash
./clipsync generate-key
public_key=oDKlgt9NsfJObrQu+Xp2GTLY80EpGkW0Hr09bBwsUTI

nano clipsync.conf
{
  "bind": "0.0.0.0:9774",
  "private_key": "3zyjOTKaXO0zFOIx2cOaeYBmQ8bSsQjTr9dLGBHTNto",
  "public_key": "cjZ5vR4R2t3QI8xzMz0Jw2lGvnim3nBlsmiViyM0iWo",
  "clients": [
    {
      "address": "222.222.222.222:9774",
      "public_key": "C4HZ1DkOIbG3u2zqC4mL8JPhliOfjex0h3E3XoKfJhw"
    }
  ]
}

./clipsync start
```

## Computer 222.222.222.222

```bash
./clipsync generate-key
public_key=C4HZ1DkOIbG3u2zqC4mL8JPhliOfjex0h3E3XoKfJhw

nano clipsync.conf
{
  "bind": "0.0.0.0:9774",
  "private_key": "Vc7BVAyVFdtmvtOv5uhm/2/EoAZlOXvsL/QgCUVlVAg",
  "public_key": "C4HZ1DkOIbG3u2zqC4mL8JPhliOfjex0h3E3XoKfJhw",
  "clients": [
    {
      "address": "111.111.111.111:9774",
      "public_key": "cjZ5vR4R2t3QI8xzMz0Jw2lGvnim3nBlsmiViyM0iWo"
    }
  ]
}

./clipsync start
```
