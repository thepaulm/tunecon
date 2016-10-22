package sync

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const copyCount = 1024 * 10
const lineHashes = 78

// HashCopyFile copies src to dst and prints out # status
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

// PushFiles copies all files from src dir to dst dir
func PushFiles(src string, dst string) error {
	srcFiles, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	for _, srcF := range srcFiles {
		dstF := fmt.Sprintf("%s/%s", dst, srcF.Name())
		_, err := os.Stat(dstF)
		if err != nil {
			// src not in dst, copy to dst
			err = HashCopyFile(srcF.Name(), srcF.Size(),
				fmt.Sprintf("%s/%s", src, srcF.Name()), dstF)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CleanFiles removes files from dst not in src
func CleanFiles(src string, dst string) error {
	dstFiles, err := ioutil.ReadDir(dst)
	if err != nil {
		return err
	}
	for _, dstF := range dstFiles {
		srcF := fmt.Sprintf("%s/%s", src, dstF.Name())
		_, err := os.Stat(srcF)
		if err != nil {
			// dst not in src, remove from src
			err = os.Remove(dstF.Name())
			if err != nil {
				fmt.Printf("Error removing file %s: %v\n", dstF.Name(), err)
			}
		}
	}
	return nil
}

func srcDstFiles(src string, dst string) ([]os.FileInfo, []os.FileInfo, error) {
	srcFiles, err := ioutil.ReadDir(src)
	if err != nil {
		return nil, nil, fmt.Errorf("src dir: %v", err)
	}
	dstFiles, err := ioutil.ReadDir(dst)
	if err != nil {
		return nil, nil, fmt.Errorf("dst dir: %v", err)
	}
	return srcFiles, dstFiles, nil
}

func cleanDst(src string, dst string) error {
	fmt.Printf("cleaning: %s <- %s", src, dst)
	srcFiles, dstFiles, err := srcDstFiles(src, dst)
	if err != nil {
		return err
	}

	for _, dstD := range dstFiles {
		if dstD.IsDir() {
			fmt.Printf("D: %s\n", dstD.Name())
			srcD := fmt.Sprintf("%s/%s", src, dstD.Name())
			sStat, err := os.Stat(srcD)
			if err != nil {
				fmt.Printf("Skipping missing src: %s\n", srcD)
				fmt.Printf("(you may want to check this yourself)\n")
				continue
			}
			if !sStat.IsDir() {
				fmt.Printf("Skipping non directory src: %s\n", srcD)
				fmt.Printf("(you may want to check this yourself)\n")
				continue
			}
			err = CleanFiles(fmt.Sprintf("%s/%s", dst, dstD.Name()), srcD)
			if err != nil {
				fmt.Printf("srcD: %v", err)
				return err
			}
		}
	}
	_ = srcFiles
	return nil
}

func pushSrc(src string, dst string) error {
	fmt.Printf("Pusing: %s -> %s\n", src, dst)
	srcFiles, dstFiles, err := srcDstFiles(src, dst)
	if err != nil {
		return err
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
			err = PushFiles(fmt.Sprintf("%s/%s", src, srcD.Name()), dstD)
			if err != nil {
				fmt.Printf("dstD: %v", err)
				return err
			}
		}
	}
	_ = dstFiles
	return nil
}

// Dirs syncs src and dst subdirectories
func Dirs(src string, dst string) error {
	e1 := cleanDst(src, dst)
	if e1 != nil {
		fmt.Printf("clean error: %v\n", e1)
	}
	e2 := pushSrc(src, dst)
	if e2 != nil {
		fmt.Printf("sync error: %v\n", e2)
	}
	if e1 != nil || e2 != nil {
		return fmt.Errorf("clean: %v, sync: %v", e1, e2)
	}
	return nil
}
