# moviespectrum

A command line tool to generate the color spectrum for a given input video file.
| Source video | Result |
|--------------|--------|
| <img src="Screenshot_20201022_103659.png" width="360"/> | <img src="KikoRiki%20Ep.%206%20-%20Season%203%20-%20The Border-3-HiGBJ7nJ0.mp4.png" width="360" />

## Usage

Build:
```go build cmd/moviespectrum.go```

Prepare a video file, e.g. download it from youtube:
```
youtube-dl --format 160 https://www.youtube.com/watch?v=3-HiGBJ7nJ0
```

Run the tool:
```
./moviespectrum 'KikoRiki Ep. 6 - Season 3 - The Border-3-HiGBJ7nJ0.mp4'
```

The resulting spectrum image file is saved to the file with the same file name as the source video with additional 
".png" extension.
