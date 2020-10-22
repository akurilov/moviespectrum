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
| 1985 | [Back To The Future (1985) Theatrical Trailer - Michael J. Fox Movie HD](https://www.youtube.com/watch?v=qvsgGtivCgs) | <img src="examples/Back%20To%20The%20Future%20(1985)%20Theatrical%20Trailer%20-%20Michael%20J.%20Fox%20Movie%20HD-qvsgGtivCgs.mp4.png"/> |
| 1991 | [Terminator 2: Judgment Day (1991) Trailer #1 Movieclips Classic Trailers](https://www.youtube.com/watch?v=CRRlbK5w8AE) | <img src="examples/Terminator%202%20-%20Judgment%20Day%20%281991%29%20Trailer%20%231%20_%20Movieclips%20Classic%20Trailers-CRRlbK5w8AE.mp4.png"/> |
| 1993 | [Jurassic Park Trailer](https://www.youtube.com/watch?v=lc0UehYemQA) | <img src="examples/Jurassic%20Park%20Trailer-lc0UehYemQA.mp4.png"/> |
| 1994 | [Pulp Fiction (1994) Official Trailer - Samuel L. Jackson, John Travolta Movie HD](https://www.youtube.com/watch?v=5ZAhzsi1ybM) | <img src="examples/Pulp%20Fiction%20%281994%29%20Official%20Trailer%20-%20Samuel%20L.%20Jackson%2C%20John%20Travolta%20Movie%20HD-5ZAhzsi1ybM.mp4.png"/> |
| 1994 | [The Lion King - Original Release Trailer (1994)](https://www.youtube.com/watch?v=hY7xBISLBIA) | <img src="examples/The%20Lion%20King%20-%20Original%20Release%20Trailer%20%281994%29-hY7xBISLBIA.mp4.png"/> |
| | []() | <img src="examples/"/> |
| | []() | <img src="examples/"/> |
| | []() | <img src="examples/"/> |
| | []() | <img src="examples/"/> |
| | []() | <img src="examples/"/> |

