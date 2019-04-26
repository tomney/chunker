# chunker
This application downloads a file in 4 chunks and writes the chunks downloaded to disk. It is a simplified version of a “download booster”, which speeds up downloads by requesting
files in multiple pieces simultaneously (saturating the network), then reassembling the pieces.

## Running the program
After downloading the `chunker` package you should be able to execute it by running 
```
$ ./chunker
```
You can then enter in the URL address you intend to download a file from. After doing this you will be prompted to enter the filename (though you can leave this blank and recieve a randomly generated filename).

## Running tests
To run the tests, navigate to the `chunker` directory and run 
```
$ go test
```

## Installing go dependencies
If there are third party go dependecies which you need to run locally use
```
$ go get <third party library>
```