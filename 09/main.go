package main

import (
	"aoc-2024/internal"
	"fmt"
	"github.com/samber/lo"
	"io"
	"slices"
	"strings"

	_ "aoc-2024/internal"
)

const DemoInput = `
2333133121414131402
`

var Input string

func init() {
	//Input = DemoInput
	Input = internal.Download(2024, 9)
}

type sized struct {
	size int
}

func (s *sized) Size() int {
	return s.size
}

type File struct {
	sized
	id *int
}
type Space struct {
	sized
}

func (f *File) Id() *int {
	return f.id
}

func (s *Space) Id() *int {
	return nil
}

type Object interface {
	Size() int
	Id() *int
}

type Objects []Object

func (o Objects) debug(writer io.Writer) {
	lo.Must(fmt.Fprintln(writer, strings.Join(lo.Map(o.serialize(), func(o int, _ int) string {
		if o != -1 {
			return string(rune(o + '0'))
		} else {
			return "."
		}
	}), "")))
}
func (o Objects) serialize() []int {
	return lo.FlatMap(o, func(o Object, _ int) []int {
		if id := o.Id(); id != nil {
			return slices.Repeat([]int{*id}, o.Size())
		} else {
			return slices.Repeat([]int{-1}, o.Size())
		}
	})
}

func (o Objects) checksum() int {
	return lo.Sum(lo.Map(o.serialize(), func(obj int, index int) int {
		if obj == -1 {
			return 0
		}
		return obj * index
	}))
}

func (o Objects) Copy() Objects {
	// copy each object
	return lo.Map(o, func(o Object, _ int) Object {
		if o.Id() == nil {
			return NewSpace(o.Size())
		}
		id := *o.Id()
		return NewFile(o.Size(), &id)
	})
}

func NewFile(size int, id *int) Object {
	return &File{id: id, sized: sized{size: size}}
}

func NewSpace(size int) Object {
	return &Space{sized: sized{size: size}}
}

var _ Object = (*File)(nil)
var _ Object = (*Space)(nil)

func main() {
	objects := Objects(lo.Map([]rune(strings.TrimSpace(Input)), func(n rune, i int) Object {
		size := int(n - '0')
		if i%2 == 0 {
			return NewFile(size, lo.ToPtr(i/2))
		} else {
			return NewSpace(size)
		}
	}))

	endSpace := NewSpace(0).(*Space)
	objects = append(objects, endSpace)

	//objects.debug(os.Stderr)
	byBlocks := compactByBlocks(objects.Copy(), endSpace)
	//objects.debug(os.Stderr)
	fmt.Println(byBlocks.checksum())

	//objects.debug(os.Stderr)
	byFiles := compactByFiles(objects.Copy())
	//objects.debug(os.Stderr)
	fmt.Println(byFiles.checksum())
}

func compactByFiles(objs Objects) Objects {
	rf, _ := lo.Must2(lo.FindLastIndexOf(objs, func(o Object) bool {
		return o.Id() != nil
	}))
	cfid := *rf.Id()

	for cfid > 0 {
		//slog.Info("looking for space", "cfid", cfid)
		//objs.debug(os.Stderr)

		// take file with cfid
		fileToMove, fileIdx := lo.Must2(lo.FindIndexOf(objs, func(item Object) bool {
			id := item.Id()
			return id != nil && *id == cfid
		}))

		fileSize := fileToMove.Size()

		space, spaceIdx, found := lo.FindIndexOf(objs, func(o Object) bool {
			return o.Id() == nil && o.Size() >= fileSize
		})

		if !found || spaceIdx > fileIdx {
			// no suitable space found, try next
			// don't want to move file to the right, try next
			cfid--
			continue
		}

		s := space.(*Space)
		// file moved to space
		objs = lo.Splice(objs, spaceIdx, fileToMove)
		// decreased space
		s.size -= fileSize
		fileIdx++ // something added before the file

		objBeforeFile := objs[fileIdx-1]

		if objBeforeFile.Id() == nil {
			// is space
			s := objBeforeFile.(*Space)
			s.size += fileSize

			cfid--
			objs = lo.DropByIndex(objs, fileIdx)
			continue
		}

		objAfterFile, _ := lo.Nth(objs, fileIdx+1)
		if objAfterFile != nil && objAfterFile.Id() == nil {
			//	has something and it's space
			s := objBeforeFile.(*Space)
			s.size += fileSize

			cfid--
			objs = lo.DropByIndex(objs, fileIdx)
			continue
		}

		// files from both sides, insert space
		objs = lo.DropByIndex(objs, fileIdx)
		objs = lo.Splice(objs, fileIdx, NewSpace(fileSize))
	}
	return objs
}

func compactByBlocks(objs Objects, endSpace *Space) Objects {

	for {
		ls, lsidx := lo.Must2(lo.FindIndexOf(objs, func(o Object) bool {
			return o.Id() == nil
		}))

		rf, rfidx := lo.Must2(lo.FindLastIndexOf(objs, func(o Object) bool {
			return o.Id() != nil
		}))

		if rfidx < lsidx {
			// defragmented
			break
		}
		s := ls.(*Space)
		f := rf.(*File)

		ss := s.size
		fs := f.size

		switch {
		case ss == fs:
			// complete match
			objs = lo.DropByIndex(objs, lsidx, rfidx)
			endSpace.size += f.size

			objs = lo.Splice(objs, lsidx, rf)
		case ss > fs:
			//	space is bigger the file
			objs = lo.DropByIndex(objs, rfidx)
			endSpace.size += f.size
			objs = lo.Splice(objs, lsidx, rf)
			s.size -= f.size
		case fs > ss:
			// file bigger then space
			objs = lo.DropByIndex(objs, lsidx)
			endSpace.size += s.size
			f.size -= s.size
			objs = lo.Splice(objs, lsidx, NewFile(s.size, f.id))
		}

		//objs.debug(os.Stderr)
		//slog.Info("compacted", "diff", rfidx-lsidx)
	}
	return objs
}
