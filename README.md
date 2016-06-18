Shade SDK
=========

A simple and easy to use 2.5D game SDK for the Go programming language.

NOTE: This SDK should be considered very experimental as it is still under development.  It is currently being modeled after some aspects of the PyGame SDK, but this will probably change some as it matures.  The project will not have its "experimental" status removed until it is easy to install, easy to use, and supports dynamic lighting.

While the above should work without needing to work with the OpenGL SDK, the packages of this SDK should be extendable such that more advanced uses are possible.

Installing
----------

Linux (Debian based) Specific Instructions
```
# Install Go 1.6
$ sudo apt-get install git-core libgl1-mesa-dev libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev
```

Windows Specific Instructions (work in progress)
```
https://git-scm.com/download/win
https://golang.org/doc/install
https://cygwin.com/install.html
http://www.glfw.org/
```

To install:

Dependencies

```
$ go get github.com/go-gl/gl/v{3.2,3.3,4.1,4.4,4.5}-{core,compatibility}/gl
$ go get github.com/go-gl/gl/v3.3-core/gl
```

NOTE: the first "go get" will produce an error because generated files are not generated yet.

```
$ go get github.com/aeonurutu/shade/...
package github.com/aeonurutu/shade/gen: cannot find package "github.com/aeonurutu/shade/gen" in any of:
	/usr/local/go/src/github.com/aeonurutu/shade/gen (from $GOROOT)
	/Users/hurricanerix/bin/usr/gocode/src/github.com/aeonurutu/shade/gen (from $GOPATH)
$ go generate github.com/aeonurutu/shade/...
$ go get github.com/aeonurutu/shade/...
```

```
cd $GOPATH/src/github.com/aeonurutu/examples/ex2-platform
go run main.go
```

Testing
-------

You can run tests using:

```
$ make test
ok  	github.com/aeonurutu/shade	0.006s	coverage: 0.0% of statements
?   	github.com/aeonurutu/shade/actor	[no test files]
?   	github.com/aeonurutu/shade/archive	[no test files]
?   	github.com/aeonurutu/shade/camera	[no test files]
?   	github.com/aeonurutu/shade/display	[no test files]
?   	github.com/aeonurutu/shade/engine	[no test files]
ok  	github.com/aeonurutu/shade/entity	0.006s	coverage: 10.3% of statements
?   	github.com/aeonurutu/shade/events	[no test files]
?   	github.com/aeonurutu/shade/fonts	[no test files]
?   	github.com/aeonurutu/shade/gen	[no test files]
?   	github.com/aeonurutu/shade/light	[no test files]
?   	github.com/aeonurutu/shade/scene	[no test files]
ok  	github.com/aeonurutu/shade/shapes	0.007s	coverage: 100.0% of statements
?   	github.com/aeonurutu/shade/splash	[no test files]
?   	github.com/aeonurutu/shade/splash/ghost	[no test files]
?   	github.com/aeonurutu/shade/sprite	[no test files]
?   	github.com/aeonurutu/shade/time	[no test files]
?   	github.com/aeonurutu/shade/time/clock	[no test files]
```

If you would like simplified statuses

```
$ make test | awk -f ci.awk
build passing
```

Attribution
-----------

This project was inspired by the article ["Normal Mapping with Javascript and Canvas"](https://29a.ch/2010/3/24/normal-mapping-with-javascript-and-canvas-tag).

Some aspects of the SDK are inspired by the [PyGame SDK](http://www.pygame.org/).

assets/gopher* adapted from [creations by Renee French under a Creative Commons Attributions 3.0 license.](https://golang.org/doc/gopher/)

Helpful Tools
-------------

[Pyxel Edit](http://pyxeledit.com/) - Very nice pixel art editor.

[Sprite DLight](https://www.kickstarter.com/projects/2dee/sprite-dlight-instant-normal-maps-for-2d-graphics) - Instant normal maps for 2D graphics

[Tiled](http://www.mapeditor.org/) - Your free, easy to use and flexible tile map editor.

[sfxr](http://www.drpetter.se/project_sfxr.html)/[cfxr](http://thirdcog.eu/apps/cfxr) - Simple means of getting basic sound effects.


Troubleshooing
--------------

#### cannot find package "github.com/aeonurutu/shade/gen"

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
