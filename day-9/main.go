package main

import (
	"fmt"
	"strings"
)

const (
	BREAK   = true
	NOBREAK = false
)

type File struct {
	ID   int
	Size int
	Next *File
	Prev *File
}

func (f *File) Copy() *File {
	ptr := f
	var head, prev, cur *File
	for ptr != nil {
		cur = &File{ID: ptr.ID, Size: ptr.Size}
		if prev != nil {
			cur.Prev = prev
			prev.Next = cur
		}
		prev = cur
		ptr = ptr.Next
	}
	return head
}

func expand2(s string) *File {
	if len(s) == 0 {
		return nil
	}
	var head, prev, cur *File
	fid := 0
	empty := false
	ix := 0
	for ix < len(s) {
		cur = &File{ID: -1, Size: int(s[ix] - byte('0'))}
		if !empty {
			cur.ID = fid
			fid++
		}
		cur.Prev = prev
		if prev != nil {
			prev.Next = cur
		}
		if head == nil {
			head = cur
		}
		prev = cur
		empty = !empty
		ix++
	}

	return head
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

func defrag(head *File, brk bool) *File {
	ptr := head
	var tail *File
	for ptr != nil {
		tail = ptr
		ptr = ptr.Next
	}

	for tail != head {
		if tail.ID < 0 {
			tail = tail.Prev
			continue
		}
		chead := head
		for chead != tail {
			if chead.ID >= 0 {
				chead = chead.Next
				continue
			}
			if chead.Size < tail.Size {
				chead = chead.Next
				continue
			}
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
			// Compact adjacent empty files
			if tail.Next != nil && tail.Next.ID < 0 {
				tail.Size += tail.Next.Size
				tail.Next = tail.Next.Next
				if tail.Next != nil {
					tail.Next.Prev = tail
				}
			}
			if tail.Prev != nil && tail.Prev.ID < 0 {
				tail.Prev.Size += tail.Size
				tail.Prev.Next = tail.Next
				if tail.Next != nil {
					tail.Next.Prev = tail.Prev
				}
			}
			break
		}
		tail = tail.Prev
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

	files1 := expand(lines[0])
	compacted1 := compact(files1)
	printf("checksum1: %d", checksum(compacted1))

	files2 := expand2(lines[0])
	compacted2 := defrag(files2, NOBREAK)
	printf("checksum2: %d", checksum2(compacted2))
}
