# tunecon

No more ipod for me. I have a "AGPtEK A12 8GB Portable Clip Mp3 Player" which I got from Amazon because I'm tired of the terrible ipod interface. This is a program written in golang wich will synchronize files from your itunes folder to your generic mp3 player.

The main functionality is copying files that live in your itunes directory but not on your mp3 player to your mp3 player. This is much nicer than dragging and dropping the whole directory each time.

Example output:
```
17719555 bytes: #532_ The Wild West of the Internet.mp3
#############################################################################
16930072 bytes: #566_ The Zoo Economy.mp3
#############################################################################
19943202 bytes: #723_ The Risk Farmers.mp3
#############################################################################
```

etcetera.

To make this go, run it with `-src <directory>` and `-dst <directory>`.

I run mine like this:
```
./build/tuncon -src ~/Music/iTunes/iTunes\ Media/Podcasts/ -dst /Volumes/NO\ NAME/
```

In the future I hope to also sync out positional information so we can auto-remote already played files, etc.

