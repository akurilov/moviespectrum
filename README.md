# moviespectrum

A command line tool to generate the color spectrum for a given input video file.
| Source video | Result Spectrum |
|--------------|-----------------|
| <img src="examples/Screenshot_20201022_103659.png" width="360"/> | <img src="examples/KikoRiki%20Ep.%206%20-%20Season%203%20-%20The Border-3-HiGBJ7nJ0.mp4.png" width="360" />

## Usage

Build:
```
go build cmd/moviespectrum.go
```

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

## How it works

1. Convert the source video to 256x144 RGB frames
2. Convert each pixel to HSL
3. Spectrum X axis is pixel's Hue
4. Color weight = Saturation * Median lightness difference
5. Median lightness difference = chi square for the pixel's lightness and lightness range middle 

## Most popular movies

| Year | Source Video Link | Result Spectrum |
|------|-------------------|-----------------|
| 1975 | [Jaws - Theatrical Trailer (HD) (1975)](https://www.youtube.com/watch?v=4pxkU9GVAoA) | <img src="examples/Jaws%20-%20Theatrical%20Trailer%20(HD)%20(1975)-4pxkU9GVAoA.mp4.png"/> |
| 1977 | [Star Wars A New Hope 1977 Trailer](https://www.youtube.com/watch?v=1g3_CFmnU7k) | <img src="examples/Star%20Wars%20A%20New%20Hope%201977%20Trailer-1g3_CFmnU7k.mp4.png"/> |
|  | []() | <img src=""/> |