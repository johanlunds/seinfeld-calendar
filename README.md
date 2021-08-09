# Seinfeld Calendar

This is a [Seinfeld Calendar](https://lifehacker.com/jerry-seinfelds-productivity-secret-281626) I made.
I made my own version in Golang. It produces a PDF using the [draw2d](https://github.com/llgcode/draw2d) library.

It has 2 improvements I wanted that I couldn't find on the web:

- ability to set start date yourself (ie. day 1 = Aug 2 2021)
- include weekday, day of the month, and name of month

[Here's another good article explaining Seinfeld Calendars.](https://jamesclear.com/stop-procrastinating-seinfeld-strategy)

It uses the [Inter font family](https://rsms.me/inter/).

## How to run

```sh
make
open output/calendar.pdf
```

## Example output

[Here's a PDF example](reference/example-output.pdf).

## Generate fonts

To generate fonts to use with draw2d and gofpdf:

```sh
cd gofpdf/makefont/
go build
./makefont \
  --embed \
  --enc=../font/iso-8859-1.map \
  --dst=~/seinfeld-calendar/resource/font/ \
  "~/Inter-3.19/Inter Hinted for Windows/Desktop/Inter-Bold.ttf"
```

That will generate a .json and .z file. Take those files and the .ttf file, [add the appropriate suffix](https://github.com/llgcode/draw2d/blob/577c1ead272a7aad4e14b84427a948b2336bc088/font.go#L42-L63) to the filenames (for example "sr") and put them in the `resource/font/` folder.
