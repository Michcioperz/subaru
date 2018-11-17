Subaru
======

**Subaru** is a software solution to a graphics driver problem.

Backstory
---------

I recently (as of time of writing this document) acquired a Pinebook laptop.
Those shiny MacBook-styled Chinese laptops sport an Allwinner aarch64 CPU and
Mali graphics. My laptop, specifically, is among the first series of the model
to have a 1080p display instead of 1366x768.

I don't know whether or not the new display is relevant to the case, but
software support for hardware accelerated H264 decoding slightly sucks. In
particular, using `mpv --vo=vdpau $file`, which yields pretty nice framerate at
low CPU usage, disables rendering of subtitles and any graphical user interface
mpv usually has.

I could dig into why things don't work as expected (and quite possibly run into
legal barriers), or I could bodge a workaround for this particular thing.
So I decided to outsource subtitle displaying away from the player.

The first step was to obtain synched subtitles _anyhow_. As it quickly turned
out, mpv has a Lua API. But lazy as I am, I wasn't interested in learning
the API. there had to some code I could scavenge.

Through this reasoning, I quickly found [User Scripts page on mpv Github
Wiki](https://github.com/mpv-player/mpv/wiki/User-Scripts). The next logical
step was to see if anyone's already solved my problem, but this wasn't
successful.

However, a quick search for `subt` on the linked page brought me to
[a repo with speed-transition script](https://github.com/zenyd/mpv-scripts),
which speeds the video up and down depending on whether or not there are
subtitles currently on the screen. That answered my question of how I could
capture "displayed" subtitles in real time.

Afterwards, I had to display the subtitles somewhere somehow. The first attempt
was just to `print(sub)` and that yielded a very barebones solution to the whole
problem. I could just read the subtitles off a terminal next to the video. It
wasn't what I wanted though, because my default terminal font is pretty small.

I figured out I had to build a separate GUI app to display the subtitles. But
first, I needed a mechanism to pipe the subtitles to it. That was easily
solvable with `lua-socket` and a bit of TCP. It's probably not the simplest
solution – it is not a Unix socket, for example – but it is still a lot simpler
protocol than HTTP. However, unlike in IRC, I couldn't separate subtitles with
newlines, because the subtitles could very well contain newlines themselves.

I quickly verified with my roommate that it was a safe thing to do and used null
bytes as separators.

Now I had to decide what I expected the GUI app to do. I wanted it to display
text, and scale it up and down to fit the window. I knew that I could get that
done using the excellent [Love2D](http://love2d.org) game engine, but a quick
run made it clear that the same issues that had me running mpv with VDPAU made
Love2D consume 350% CPU on Pinebook. Surely someone else could do better.

I went with Go language. I knew from my earlier encounters with its API for
sockets that to split socket input by null bytes I just had to use
a `bufio.Scanner` giving it a slightly changed copy-paste of `bufio.ScanLines`
function as a splitter.

There's this GUI toolkit that everyone always forgets about, save for 5 people:
[libui](https://github.com/andlabs/libui). It advertises itself as
cross-platform, because deep inside it's pretty much a wrapper for whatever the
dominant GUI framework on the given platform is. On Linux, for better or worse,
that means GTK+3.

So I quickly put together a TCP listener in Go, and all that network stuff, and
then it turned out something about libui didn't quite compile. So, since I was
expecting to use GTK+3 indirectly anyway, I replaced libui with GTK bindings
themselves.

And it worked pretty fast, even if during the initial build I found myself with
enough time to go out and do some grocery shopping. One bug slowed me down for
a good while, because I was checking some boolean the wrong way, but once
I finally noticed it, it was quick to fix.

I had a working solution.

But I only did because I threw away the idea of scaling the text. As it turned
out, `gtk.Label` (at least in the Go bindings) doesn't provide a way to measure
its width or height. I needed that to implement the binary search for _That
Destined Font Size_.

I guess that means there's still room for improvement, but I decided against
pushing forward for now, and instead went to enjoy some anime using my solution.

[Here's a fedi message where I posted a screenshot](https://m.atm.pl/notice/479900)

Installation
------------

1. Make sure your mpv's Lua has access to `lua-socket` library.
   On my Pinebook's KDE Neon, it meant running `apt install lua-socket`.

2. Put `subaru.lua` in your mpv scripts folder.
   Usually, that is `~/.config/mpv/scripts`

3. Make sure you have Go toolchain and GTK+3 development header available.
   On my Pinebook's KDE Neon, it meant running
   `apt install golang-go libgtk-3-dev`.

4. From the project directory, run `go get -d .`. This will download GTK+
   bindings for Go to your `$GOPATH` – if not changed, this is `$HOME/go`.

5. Build the GUI app by running `go build` in the project directory.
   You can then run it: `./subaru`.
