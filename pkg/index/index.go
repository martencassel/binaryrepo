package index

type FileIndex struct{}

type FileInfo struct {
	ID       int
	RepoKey  string
	Type     string
	Path     string
	Name     string
	Size     int64
	Checksum string
}

func (*FileIndex) Get(reponame string, path string) FileInfo {
	return FileInfo{}
}

/*var files []FileInfo = []FileInfo{
	{
		ID:       1,
		RepoKey:  "golang-remote",
		Path:     "github.com/gorilla/mux/@v/v1.8.0.zip",
		Type:     "zip",
		Name:     "file1.txt",
		Size:     250,
		Checksum: "1590ef6d972aea060ed538e5055d1c9e713eb818f1fac6332bfd4942b64eb825",
	},
	{
		ID:       2,
		RepoKey:  "golang-remote",
		Path:     "github.com/gorilla/mux/@v/v1.8.0.mod",
		Name:     "file2.txt",
		Type:     "mod",
		Size:     20,
		Checksum: "1590ef6d972aea060ed538e5055d1c9e713eb818f1fac6332bfd4942b64eb825",
	},
}
*/
