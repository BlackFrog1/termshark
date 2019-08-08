# Install Packages

Here's how to install termshark on various OSes and with various package managers.

## Arch Linux

- [termshark-bin](https://aur.archlinux.org/packages/termshark-bin): binary
  package which simply copies the released binary to install directory. Made by
  [jerry73204](https://github.com/jerry73204)
- [termshark-git](https://aur.archlinux.org/packages/termshark-git): Compiles
  from source, made by [Thann](https://github.com/Thann) 

## Debian

Termshark is only available in unstable/sid at the moment.

```bash
apt update
apt install termshark
```

## FreeBSD

Thanks to [Ryan Steinmetz](https://github.com/zi0r)

Termshark is in the FreeBSD ports tree!  To install the package, run:

```pkg install termshark```

To build/install the port, run:

```cd /usr/ports/net/termshark/ && make install clean```

## Homebrew

```bash
brew update
brew install termshark
```
## Kali Linux

```bash
apt update
apt install termshark
```

## SnapCraft

Thanks to [mharjac](https://github.com/mharjac)

Termshark can be easily installed on almost all major distros just by issuing: 

```bash
snap install termshark
```

After installation, it requires some additional permissions:

```bash
snap connect termshark:network-control
snap connect termshark:bluetooth-control
snap connect termshark:firewall-control
snap connect termshark:ppp
snap connect termshark:raw-usb
snap connect termshark:removable-media
```

## Termux (Android)

```bash
pkg install root-repo
pkg install termshark
```
Note that termshark does not require a rooted phone to inspect a pcap, but it does depend on tshark which is itself in Termux's root-repo for programs that do work best on a rooted phone.

If you would like to use termshark's copy-mode to copy sections of packets to your Android clipboard, you will also need [Termux:API](https://play.google.com/store/apps/details?id=com.termux.api&hl=en_US). Install from the Play Store, then from termux, type:

```bash
pkg install termux-api
```

![device art](https://drive.google.com/uc?export=view&id=1RzilBvj5YFsSqv72kO6yOD0Oil88mwp3)

## Ubuntu

Thanks to [Nicolai Søberg](https://github.com/NicolaiSoeborg)

You can use the PPA *nicolais/termshark* to install termshark:

```bash
sudo add-apt-repository --update ppa:nicolais/termshark
sudo apt install termshark
```


