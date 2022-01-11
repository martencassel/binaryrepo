package memory

import "testing"

func TestMemory(t *testing.T) {
	t.Run("WriteFile", func(t *testing.T) {
		fs := NewFileStore("/tmp/memory")
		b := []byte("hello world")
		digest, err := fs.WriteFile(b)
		if err != nil {
			t.Fatal(err)
		}
		if !fs.Exists(digest) {
			t.Fatal("file not found")
		}
	})
}
