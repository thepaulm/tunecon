package sync

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const copyCount = 1024 * 10
const lineHashes = 78

func HashCopyFile(shortName string, size int64, src string, dst string) error {
	fmt.Printf("%d bytes: %s\n", size, shortName)
	from, err := os.Open(src)
	if err != nil {
		return err
	}
	defer from.Close()
	to, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer to.Sync()
	defer to.Close()

	var sofar int64
	var hashes = 0
	for sofar < size {
		written, err := io.CopyN(to, from, copyCount)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		sofar += written
		needHashes := int(float32(sofar) / float32(size) * lineHashes)
		for ; hashes < needHashes; hashes++ {
			fmt.Printf("#")
		}
	}
	fmt.Printf("\n")
	return nil
}

func Files(src string, dst string) error {
	srcFiles, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	for _, srcF := range srcFiles {
		dstF := fmt.Sprintf("%s/%s", dst, srcF.Name())
		_, err := os.Stat(dstF)
		if err != nil {
			err = HashCopyFile(srcF.Name(), srcF.Size(),
				fmt.Sprintf("%s/%s", src, srcF.Name()), dstF)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Dirs(src string, dst string) error {
	srcFiles, err := ioutil.ReadDir(src)
	if err != nil {
		return fmt.Errorf("src dir: %v", err)
	}
	dstFiles, err := ioutil.ReadDir(dst)
	if err != nil {
		return fmt.Errorf("dst dir: %v", err)
	}

	for _, srcD := range srcFiles {
		if srcD.IsDir() {
			fmt.Printf("D: %s\n", srcD.Name())
			dstD := fmt.Sprintf("%s/%s", dst, srcD.Name())
			dStat, err := os.Stat(dstD)
			if err != nil {
				fmt.Printf("Skipping missing dest: %s\n", dstD)
				continue
			}
			if !dStat.IsDir() {
				fmt.Printf("Skipping non directory dest: %s\n", dstD)
				continue
			}
			err = Files(fmt.Sprintf("%s/%s", src, srcD.Name()), dstD)
			if err != nil {
				fmt.Printf("dstD: %v", err)
				return err
			}
		}
	}
	_ = dstFiles
	return nil
}
