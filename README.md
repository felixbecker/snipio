# SNIPO

SNIPIO is a cross platform single binary tool to work with [draw io files](https://app.diagrams.net/)

With `snipio` you can do the following operations:

- list all layers in a drawing
- delete a layer from a drawing
- export a layer from a drawing
- unpack a mxfile
- classify a drawing as draft

## How to use it

`snipio show layers` will show all available layers in a drawio document. It requires you to specify the document file with the parameter `-f` or `--file`.

Example: `snipio show layers -f mydrawiofile.xml`

`snipio delete layer` will delete a layer from the draw.io document. it requires the input file with the parameter `-f` or `--file`. Also the name of the layer has to be provided via the flag `-n` or `--name`. Optionally you can specify the output file with the flag `-o` or `--output`. If no output file is provided the results will be written to a file named `export.xml`.

Example: `snipio delete layer -f mydrawiofile.xml -n layer1 -o precessedfile.xml`

`snipio export layer` will export the selected layer from the draw.io document to a new file. It requires the input file with the parameter `-f` or `--file`. Also the name of the layer has to be provided via the flag `-n` or `--name`. Optionally you can specify the output file with the flag `-o` or `--output`. If no output file is provided the results will be written to a file named `export.xml`.

Example: `snipio export layer -n mydrawiofile.xml -n layer1 -o precessedfile.xml`

`snipio unpack` will convert a mxfile to a xml file. It requires the input file with the parameter `-f` or `--file`. Optionally you can specify the output file with the flag `-o` or `--output`. If no output file is provided the results will be written to the console.

Example: `snipio unpack -f myMxFile.xml`

`snipio classify draft` will add a layer on top of the document with a draft watermark. It requires the input file with the parameter `-f` or `--file`. Optionally you can specify the output file with the flag `-o` or `--output`. If no output file is provided the results will be written into a file named `export.xml`.

Example: `snipio classify draft -f mydrawiofile.xm`

## Installation

Right now there are no pre-builds to download. Might come in the future.

- Please install the go programming language. Follow the instructions at [https://go.dev/](https://go.dev/).
- Clone the repo
- run `go install` in the directory of the repo. That's it.
- you could also use `go run main.go` to run it. Or `go build main.go` to build a binary for yourself.

That's it! Have fun.
