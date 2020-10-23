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

## Some popular movies

| Year | Source Video Link | Result Spectrum |
|------|-------------------|-----------------|
| 1975 | [Jaws - Theatrical Trailer (HD) (1975)](https://www.youtube.com/watch?v=4pxkU9GVAoA) | <img src="examples/Jaws%20-%20Theatrical%20Trailer%20(HD)%20(1975)-4pxkU9GVAoA.mp4.png"/> |
| 1977 | [Star Wars A New Hope 1977 Trailer](https://www.youtube.com/watch?v=1g3_CFmnU7k) | <img src="examples/Star%20Wars%20A%20New%20Hope%201977%20Trailer-1g3_CFmnU7k.mp4.png"/> |
| 1985 | [Back To The Future (1985) Theatrical Trailer - Michael J. Fox Movie HD](https://www.youtube.com/watch?v=qvsgGtivCgs) | <img src="examples/Back%20To%20The%20Future%20(1985)%20Theatrical%20Trailer%20-%20Michael%20J.%20Fox%20Movie%20HD-qvsgGtivCgs.mp4.png"/> |
| 1991 | [Terminator 2: Judgment Day (1991) Trailer #1 Movieclips Classic Trailers](https://www.youtube.com/watch?v=CRRlbK5w8AE) | <img src="examples/Terminator%202%20-%20Judgment%20Day%20%281991%29%20Trailer%20%231%20_%20Movieclips%20Classic%20Trailers-CRRlbK5w8AE.mp4.png"/> |
| 1993 | [Jurassic Park Trailer](https://www.youtube.com/watch?v=lc0UehYemQA) | <img src="examples/Jurassic%20Park%20Trailer-lc0UehYemQA.mp4.png"/> |
| 1994 | [Pulp Fiction (1994) Official Trailer - Samuel L. Jackson, John Travolta Movie HD](https://www.youtube.com/watch?v=5ZAhzsi1ybM) | <img src="examples/Pulp%20Fiction%20%281994%29%20Official%20Trailer%20-%20Samuel%20L.%20Jackson%2C%20John%20Travolta%20Movie%20HD-5ZAhzsi1ybM.mp4.png"/> |
| 1994 | [The Lion King - Original Release Trailer (1994)](https://www.youtube.com/watch?v=hY7xBISLBIA) | <img src="examples/The%20Lion%20King%20-%20Original%20Release%20Trailer%20%281994%29-hY7xBISLBIA.mp4.png"/> |
| 1994 | [Forrest Gump - Trailer](https://www.youtube.com/watch?v=bLvqoHBptjg) | <img src="examples/Forrest%20Gump%20-%20Trailer-bLvqoHBptjg.mp4.png"/> |
| 1995 | [Braveheart Trailer - 1995 HQ](https://www.youtube.com/watch?v=1cnoM8EiGGU) | <img src="examples/Braveheart%20Trailer%20-%201995%20HQ-1cnoM8EiGGU.mp4.png"/> |
| 1997 | [Titanic (1997) - Original Trailer](https://www.youtube.com/watch?v=jUm88F3MEbQ) | <img src="examples/Titanic%20%281997%29%20-%20Original%20Trailer-jUm88F3MEbQ.mp4.png"/> |
| 1999 | [Matrix Trailer HD (1999)](https://www.youtube.com/watch?v=m8e-FF8MsqU) | <img src="examples/Matrix%20Trailer%20HD%20%281999%29-m8e-FF8MsqU.mp4.png"/> |
| 2001 | [Shrek 2001 Official Trailer](https://www.youtube.com/watch?v=ooJJX3R42WM) | <img src="examples/Shrek%202001%20Official%20Trailer-ooJJX3R42WM.mp4.png"/> |
| 2003 | [The Lord of the Rings: The Two Towers (2002) Official Trailer #2 - Orlando Bloom Movie HD](https://www.youtube.com/watch?v=LbfMDwc4azU) | <img src="examples/The%20Lord%20of%20the%20Rings%20-%20The%20Two%20Towers%20%282002%29%20Official%20Trailer%20%232%20-%20Orlando%20Bloom%20Movie%20HD-LbfMDwc4azU.mp4.png"/> |
| 2003 | [Pirates of the Caribbean: The Curse of the Black Pearl Official Trailer 1 (2003) HD](https://www.youtube.com/watch?v=naQr0uTrH_s) | <img src="examples/Pirates%20of%20the%20Caribbean%20-%20The%20Curse%20of%20the%20Black%20Pearl%20Official%20Trailer%201%20%282003%29%20HD-naQr0uTrH_s.mp4.png"/>
| 2007 | [Transformers (2007) - Full Trailer HD](https://www.youtube.com/watch?v=dxQxgAfNzyE) | <img src="examples/Transformers%20%282007%29%20-%20Full%20Trailer%20%5BHD%5D-dxQxgAfNzyE.mp4.png"/> |
| 2008 | [The Dark Knight - Official Trailer 2 HD](https://www.youtube.com/watch?v=TQfATDZY5Y4) | <img src="examples/The%20Dark%20Knight%20-%20Official%20Trailer%202%20%5BHD%5D-TQfATDZY5Y4.mp4.png"/> |
| 2012 | [Marvel's The Avengers Trailer 2 (OFFICIAL)](https://www.youtube.com/watch?v=hIR8Ar-Z4hw) | <img src="examples/Marvel%27s%20The%20Avengers%20Trailer%202%20%28OFFICIAL%29-hIR8Ar-Z4hw.mp4.png"/> |
| | []() | <img src="examples/"/> |
| | []() | <img src="examples/"/> |
| | []() | <img src="examples/"/> |