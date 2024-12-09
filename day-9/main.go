package main

import (
	"fmt"
	"strings"
)

type File struct {
	ID   int
	Size int
	Next *File
	Prev *File
}

func expand(s string) []File {
	files := make([]File, 0)
	empty := false
	ix := 0
	fid := 0
	for ix < len(s) {
		size := int(s[ix] - byte('0'))
		file := File{ID: -1, Size: size}
		if !empty {
			file.ID = fid
			fid++
		}
		ix++
		empty = !empty
		files = append(files, file)
	}

	return files
}

func compact(files []File) []File {
	compacted := make([]File, 0, len(files))
	head, tail := 0, len(files)-1
	reclaimed := 0
	for head <= tail {
		if files[head].ID >= 0 {
			compacted = append(compacted, files[head])
			head++
			continue
		}
		if files[tail].ID < 0 {
			tail--
			continue
		}
		// head points at a space, tail points at a file
		if files[head].Size >= files[tail].Size {
			// trailing file fits into the space
			compacted = append(compacted, files[tail])
			reclaimed += files[tail].Size
			files[head].Size -= files[tail].Size
			tail--
		} else {
			if files[head].Size > 0 {
				compacted = append(compacted, File{ID: files[tail].ID, Size: files[head].Size})
				files[tail].Size -= files[head].Size
			}
			head++
		}
	}

	res := make([]File, 0, len(compacted))
	res = append(res, compacted[0])
	ix := 1
	for ix < len(compacted) {
		if compacted[ix].ID == res[len(res)-1].ID {
			res[len(res)-1].Size += compacted[ix].Size
			ix++
			continue
		}
		res = append(res, compacted[ix])
		ix++
	}

	return res
}

func printFiles(head *File) string {
	ptr := head
	var buf strings.Builder
	for ptr != nil {
		if buf.Len() > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprintf("File{ID:%d, Size:%d}", ptr.ID, ptr.Size))
		ptr = ptr.Next
	}

	return buf.String()
}

func compact2(files []File) *File {
	var head, tail, prev *File
	for _, file := range files {
		node := &File{ID: file.ID, Size: file.Size}
		if prev != nil {
			node.Prev = prev
			prev.Next = node
		}
		if head == nil {
			head = node
		}
		tail = node
		prev = node
	}

	for tail != head {
		if tail.ID < 0 {
			tail = tail.Prev
			continue
		}
		chead := head
		for chead != tail {
			if chead.ID >= 0 || chead.Size < tail.Size {
				chead = chead.Next
				continue
			}
			//debugf("Found free block of size %d for node %d/%d between %d and %d", chead.Size, tail.ID, tail.Size, chead.Prev.ID, chead.Next.ID)
			delta := chead.Size - tail.Size
			chead.ID = tail.ID
			chead.Size = tail.Size
			if delta > 0 {
				node := &File{ID: -1, Size: delta}
				oldnext := chead.Next
				chead.Next = node
				node.Next = oldnext
				node.Prev = chead
				if oldnext != nil {
					oldnext.Prev = node
				}
			}
			tail.ID = -1
			break
		}
		tail = tail.Prev

		prev = head
		cur := head.Next
		for cur != nil {
			if cur.ID < 0 && prev.ID < 0 {
				cur.Prev = prev.Prev
				cur.Size += prev.Size
				if cur.Prev != nil {
					cur.Prev.Next = cur
				}
			}
			prev = cur
			cur = cur.Next
		}
	}

	return head
}

func checksum(files []File) uint64 {
	sum := uint64(0)
	ix := 0
	for _, file := range files {
		if file.ID > 0 {
			for i := 0; i < file.Size; i++ {
				sum += uint64(file.ID) * uint64(ix)
				ix++
			}
		} else {
			ix += file.Size
		}
	}
	return sum
}

func checksum2(head *File) uint64 {
	sum := uint64(0)
	ix := 0
	for ptr := head; ptr != nil; ptr = ptr.Next {
		if ptr.ID > 0 {
			for i := 0; i < ptr.Size; i++ {
				sum += uint64(ptr.ID) * uint64(ix)
				ix++
			}
		} else {
			ix += ptr.Size
		}
	}
	return sum
}

func main() {
	lines := input()
	files := expand(lines[0])
	filescp := make([]File, len(files))
	copy(filescp, files)
	compacted := compact(files)
	checksum := checksum(compacted)
	printf("checksum: %d", checksum)

	compacted2 := compact2(filescp)
	printf("checksum2: %d", checksum2(compacted2))
}
