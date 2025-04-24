package interp

import (
	"fmt"
	"go/token"
)

type InterpreterArchive struct {
	nindex int64

	fsetFilesLen int

	frameDataLen int

	universeChildLen int
	universeTypesLen int
	universeSymLen   int

	scopesLen int

	srcPkgLen int

	pkgNamesLen int

	rootsLen int
}

func (interp *Interpreter) SaveArchive() *InterpreterArchive {
	archive := &InterpreterArchive{}

	archive.nindex = interp.nindex

	fileSet := interp.FileSet()
	archive.fsetFilesLen = 0
	fileSet.Iterate(func(f *token.File) bool {
		fmt.Print(f.Name())
		archive.fsetFilesLen++
		return true
	})

	archive.frameDataLen = len(interp.frame.data)

	archive.universeChildLen = len(interp.universe.child)
	archive.universeTypesLen = len(interp.universe.types)
	archive.universeSymLen = len(interp.universe.sym)

	archive.scopesLen = len(interp.scopes)

	archive.srcPkgLen = len(interp.srcPkg)

	archive.pkgNamesLen = len(interp.pkgNames)

	archive.rootsLen = len(interp.roots)
	return archive
}

func (interp *Interpreter) RestoreArchive(archive *InterpreterArchive, packNameToRmv string) {
	interp.nindex = archive.nindex

	counter := 0
	fileSet := interp.FileSet()
	fileNeedToRmv := make([]*token.File, 0)
	fileSet.Iterate(func(f *token.File) bool {
		counter++
		if counter > archive.fsetFilesLen {
			fileNeedToRmv = append(fileNeedToRmv, f)
		}
		return true
	})
	for _, file := range fileNeedToRmv {
		fileSet.RemoveFile(file)
	}

	interp.frame.data = interp.frame.data[:archive.frameDataLen]

	interp.universe.child = interp.universe.child[:archive.universeChildLen]
	interp.universe.types = interp.universe.types[:archive.universeTypesLen]
	delete(interp.universe.sym, packNameToRmv)

	delete(interp.scopes, packNameToRmv)

	delete(interp.srcPkg, packNameToRmv)

	delete(interp.pkgNames, packNameToRmv)

	interp.roots = interp.roots[:archive.rootsLen]
}
