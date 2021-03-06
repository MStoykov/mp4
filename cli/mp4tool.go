package main

import (
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/MStoykov/mp4"
	"github.com/MStoykov/mp4/clip"
	cli "github.com/jawher/mow.cli"
)

func main() {
	go func() {
		fmt.Println(http.ListenAndServe(":6060", nil))
	}()

	cmd := cli.App("mp4tool", "MP4 command line tool")

	cmd.Command("info", "Displays information about a media", func(cmd *cli.Cmd) {
		file := cmd.StringArg("FILE", "", "the file to display")
		cmd.Action = func() {
			rr := &fileRangeReader{fileName: *file}
			v, err := mp4.Decode(rr)
			if err != nil {
				fmt.Println(err)
			}
			v.Dump()
		}
	})

	cmd.Command("clip", "Generates a clip", func(cmd *cli.Cmd) {
		start := cmd.IntOpt("s start", 0, "start time (sec)")
		src := cmd.StringArg("SRC", "", "the source file name")
		dst := cmd.StringArg("DST", "", "the destination file name")
		cmd.Action = func() {
			rr := &fileRangeReader{fileName: *src}
			v, err := mp4.Decode(rr)
			if err != nil {
				fmt.Println(err)
				return
			}
			out, err := os.Create(*dst)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer closeWithPrintingError(out, "error on closing output file")
			clip, err := clip.New(v, time.Duration(*start)*time.Second, rr)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Scanln()
			size, err := clip.WriteTo(out)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("wrote", size)
		}
	})

	cmd.Command("copy", "Decodes a media and reencodes it to another file", func(cmd *cli.Cmd) {
		src := cmd.StringArg("SRC", "", "the source file name")
		dst := cmd.StringArg("DST", "", "the destination file name")
		cmd.Action = func() {
			rr := &fileRangeReader{fileName: *src}
			v, err := mp4.Decode(rr)
			if err != nil {
				fmt.Println(err)
				return
			}
			out, err := os.Create(*dst)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer closeWithPrintingError(out, "error on closing output file")
			if err := v.Encode(out); err != nil {
				fmt.Println("error on encoding", err)
			}
		}
	})
	if err := cmd.Run(os.Args); err != nil {
		fmt.Println("error form cmd.Run", err)
	}
	fmt.Println("press return to exit the program")
	fmt.Scanln()
}

type fileRangeReader struct {
	fileName string
}

func (frr *fileRangeReader) RangeRead(start, length uint64) (io.ReadCloser, error) {
	file, err := os.Open(frr.fileName)
	if err != nil {
		return nil, err
	}
	if _, err := file.Seek(int64(start), os.SEEK_SET); err != nil {
		return nil, err
	}
	r := &readCloser{Reader: io.LimitReader(file, int64(length)), closeFunc: file.Close}
	return r, nil
}

type readCloser struct {
	io.Reader
	closeFunc func() error
}

func (r *readCloser) Close() error {
	return r.closeFunc()
}

// close the closer and if it returns error print it preceded by the message
func closeWithPrintingError(closer io.Closer, message string) {
	if err := closer.Close(); err != nil {
		fmt.Println(message, err)
	}
}
