// package store_test memastika MemoryStore mengimplementasikan ItemStore dengan benar.
// Test disini tidak butuh databse, jadi kita bisa pakai MemoryStore.
package store_test

// func TestNewMemoryStore memastikan MemoryStore dibuat dengan benar.
// magsud dari t *testing.T adalah untuk menjalankan test ini dengan parameter t yang akan digunakan untuk melakukan assertion dan reporting error.
// Tujuannya adalah untuk memastikan bahwa MemoryStore dibuat dengan benar dan dapat digunakan untuk menyimpan dan mengambil data.
func TestNewMemoryStore(t *testing.T) {
	S := store.NewMemoryStore()
    // Jika Variable S benilai nil atau tidak ada
	if S == nil {
		// t.Fatal mermakna test gagal dan program berhenti.
		// ("NewMemoryStore() tidak boleh nil atau tidak ada") adalah pesan error yang akan ditampilkan jika test gagal.
		t.Fatal("NewMemoryStore() tidak boleh nil atau tidak ada")
	}
	items, err := s.List()
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("List() awal harus kosong, got %d items", len(items))
	}

}