package model

type Book struct {
    Judul  string `bson:"judul"`
    Harga string    `bson:"harga"`
    Deskripsi string    `bson:"deskripsi"`
}

type BookUpdate struct {
    ID string `bson:"id"`
    Judul  string `bson:"judul"`
    Harga string    `bson:"harga"`
    Deskripsi string    `bson:"deskripsi"`
}