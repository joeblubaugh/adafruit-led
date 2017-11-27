A Go Library supporting Adafruit Backpacks on systems with a devfs / sysfs I2C bus.

Includes support for:
- 8x8 LED Grid packs
- Generic HT16K33 LED driver packs

Planned New Features:
- 2x12 (aka 24 :-P) LED Bar Graph Packs.
- 4-character 7 segment display packs.

Planned Enhancements:
- Remove application-specific code ("scroll chars") from library code.
- Add support for 8x8 Grid Rotation
- Consider removing local i2c package in favor of golang.org/x/exp/io/i2c
- Remove support for MMA7455 accelerometer / move into a separate package.

Go code written by George McBay (george.mcbay@gmail.com) 
and Joe Blubaugh (joe.blubaugh@gmail.com) under BSD License.  

Portions adapted from Adafruit's Raspberry-Pi Python Code Library.

This library is forked from the one hosted at: https://bitbucket.org/gmcbay/i2c

---------------------------------------------------------------------------------------

Adafruit's Raspberry-Pi Python Code Library
============
  Here is a growing collection of libraries and example python scripts
  for controlling a variety of Adafruit electronics with a Raspberry Pi
  
  In progress!

  Adafruit invests time and resources providing this open source code,
  please support Adafruit and open-source hardware by purchasing
  products from Adafruit!

  Written by Limor Fried, Kevin Townsend and Mikey Sklar for Adafruit Industries.
  BSD license, all text above must be included in any redistribution
  
  To download, we suggest logging into your Pi with Internet accessibility and typing:
  git clone https://github.com/adafruit/Adafruit-Raspberry-Pi-Python-Code.git
