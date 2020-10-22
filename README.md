# moviespectrum

A command line tool to generate the color spectrum for a given input video file.
![video](Screenshot_20201022_103659.png)
![spectrum](KikoRiki Ep. 6 - Season 3 - The Border-3-HiGBJ7nJ0.mp4.png)

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
