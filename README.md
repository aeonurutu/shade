# Shade Engine

A experimental game engine written in Go.

Primary goals:

* Prototype games quickly with Go
* Simple and easy to use
* Extendable allowing advanced customization
* Built-in support for advanced features, such as physics and dynamic lighting

NOTE: This SDK should be considered very experimental as it is still under development.  It is currently being modeled after some aspects of the PyGame SDK, but this will probably change some as it matures.  The project will not have its "experimental" status removed until it is easy to install and to use.

For usage/help using Shade, please see the [shade-engine](https://groups.google.com/forum/#!forum/shade-engine) google group.

## Installing

Install the [Go Programming Language](https://golang.org/doc/install) version 1.8.1 for your OS.

### Linux (Debian based):
```
$ sudo apt-get install git-core libgl1-mesa-dev libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev
```

### macOS

```
$ brew install pkgconfig sdl2
```

### macOS

```
$ brew install pkgconfig sdl2
```

Raspberry Pi users will also need to install
```
sudo apt-get install libgles2-mesa-dev
```

Windows:
TODO: document this.

### Dependencies

```
$ go get github.com/hurricanerix/shade/...
```

### To test your install:

```
go run -ldflags="-X github.com/hurricanerix/shade.ldDevBuild=true" $GOPATH/src/github.com/hurricanerix/shade/examples/demos/platformer/main.go --debug-server
```

## Contributing

Join the dev discussions via the [shade-engine-dev](https://groups.google.com/forum/#!forum/shade-engine-dev) google group.

Request access to shade-engine.slack.com

### Dev Build

To compile your app with the Shade's dev option available:

```
go build -ldflags="-X github.com/hurricanerix/shade.ldDevBuild=true" main.go
```

## Attribution

This project was inspired by the article ["Normal Mapping with Javascript and Canvas"](https://29a.ch/2010/3/24/normal-mapping-with-javascript-and-canvas-tag).

Some aspects of the SDK are inspired by the [PyGame SDK](http://www.pygame.org/).

assets/gopher* adapted from [creations by Renee French under a Creative Commons Attributions 3.0 license.](https://golang.org/doc/gopher/)

## Helpful Tools

[Pyxel Edit](http://pyxeledit.com/) - Very nice pixel art editor.

[Sprite DLight](https://www.kickstarter.com/projects/2dee/sprite-dlight-instant-normal-maps-for-2d-graphics) - Instant normal maps for 2D graphics

[Tiled](http://www.mapeditor.org/) - Your free, easy to use and flexible tile map editor.

[sfxr](http://www.drpetter.se/project_sfxr.html)/[cfxr](http://thirdcog.eu/apps/cfxr) - Simple means of getting basic sound effects.


## Troubleshooing

#### cannot find package "github.com/hurricanerix/shade/gen"

Some variables/assets are packaged into a generated code file, if you get this error run the bindata.sh script to generate that file.

#### Slow saves in VIM

Sometimes, go-imports inserts the the wrong things.  VIM was hanging for me for about 10~20 seconds after saving.  I'm not sure why, but switching from v2.1 to v4.1-core fixed the issue.

```
-       "github.com/go-gl/gl/v2.1/gl"
+       "github.com/go-gl/gl/v4.1-core/gl"
```

#### Error: "could not decode file player.png: image: unknown format"

The following line will not be auto-imported since our code does not call anything in the png, library directly.  In order to ensure that the library loads (so image.decode("somefile.png") works), we must have the following import to load the PNG package.

```
+	_ "image/png" // register PNG decode
```

#### ... has no field or method ...

And it just does not make sense why it is not found, check that go-imports imported the correct thing.
