# mp4

A encoder/decoder class, io.Reader and io.Writer compatible, usable for HTTP pseudo streaming

For the complete MP4 specifications, see http://standards.iso.org/ittf/PubliclyAvailableStandards/c061988_ISO_IEC_14496-12_2012.zip and http://standards.iso.org/ittf/PubliclyAvailableStandards/c061989_ISO_IEC_15444-12_2012.zip

## Doc

See http://godoc.org/github.com/MStoykov/mp4 and http://godoc.org/github.com/MStoykov/mp4/clip

## CLI

A CLI can be found in cli/mp4tool.go

It can :

* Display info about a media
```
mp4tool info file.mp4
```
* Copy a video (decode it and reencode it to another file, useful for debugging)
```
mp4tool copy in.mp4 out.mp4
```
* Generate a clip
```
mp4tool clip --start 10 in.mp4 out.mp4
```

(if you really want to generate a clip, you should use ffmpeg, you will ge better results)

## LICENSE

See LICENSE
