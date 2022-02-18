package metadata

import "time"

type FileMetadata struct {
	Name           string
	RepositoryPath string
	FileUrl        string
	Created        time.Time
}

/*
	Repository: docker-remote
	Blob 		/repo/docker-remote/v2/postgres/blobs/sha256:da2cb49d7a8d1416cfc2ec6fb47b60112b3a2f276bcf7439ef18e7c505b83fc6
	Blob 		/repo/docker-remote/v2/postgres/blobs/sha256:daa0467a6c4883c02b241fe5f4f1703245f43ccbe5bcd56a3dceddef285bf31e
	Blob 		/repo/docker-remote/v2/postgres/blobs/sha256:bb8afcc973b2a3de0135018fd8bd12fc2f56ef30a64e6e503d1a20cf35e340f3
	Blob 		/repo/docker-remote/v2/postgres/blobs/sha256:c74bf40d29ee27431deec24f6d21d1a09f178335b7c25aa6fd26850bec90252a
	Manifest	/repo/docker-remote/v2/postgres/manifests/latest
*/

/*
	v docker-remote
		v hello-world
			> _uploads
			> v1.0

	v docker-remote
		v hello-world
			> _uploads
			v v1.0
				da2cb49d7a8d1416cfc2ec6fb47b60112b3a2f276bcf7439ef18e7c505b83fc6
				daa0467a6c4883c02b241fe5f4f1703245f43ccbe5bcd56a3dceddef285bf31e
				manifest.json
		v postgres
			> latest
				0da2cb49d7a8d1416cfc2ec6fb47b60112b3a2f276bcf7439ef18e7c505b83fc6
				daa0467a6c4883c02b241fe5f4f1703245f43ccbe5bcd56a3dceddef285bf31e
				manifest.json
*/
