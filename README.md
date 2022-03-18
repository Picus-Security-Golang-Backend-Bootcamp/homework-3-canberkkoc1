# Go KİTAP ARAMA , LİSTELEME , SİLME VE GÜNCELLEME UYGULAMASI
Bu uygulama Patika ve Picus Security iş birliğinde gerçekleşen Golang Backend Web Development Bootcamp kapsamındaki ilk ödevi içermektedir.  
> Projenin çalışırlığına [Uygulama Çıktısı](#uygulama-%C3%A7%C4%B1kt%C4%B1s%C4%B1)'ndan ulaşabilirsiniz.  


<br>  
<br>  

# İçindekiler
- [Uygulama İsterleri](#uygulama-i%CC%87sterleri)  
- [Uygulama Öncesi Hazırlık](#uygulama-%C3%B6ncesi-haz%C4%B1rl%C4%B1k)  
- [Uygulama Aşamaları](#uygulama-a%C5%9Famalar%C4%B1)  
- [Uygulama Çıktısı](#uygulama-%C3%A7%C4%B1kt%C4%B1s%C4%B1) 


<br>  
<br>  

# Uygulama İsterleri
- Kullanıcı yapmak istediği işlemleri `list`, `search`, `delete` ve `buy` komutlarını kullanarak belirtecek.    
- Uygulama içinde oluşturulmuş csv dosyası ile veritabanına kayıt işlemleri yapılmaktadır. Kullanıcının istekler veritabanında gerçekleştirilecektir.
- *Komutlar:*  
	+ **search:** Arama yapmak için kullanılacak bu komuttan sonraki gelen argümanlar birleştirilerek arama yaparken kullanılacak değeri içermelidir.  
	> örn: search star wars  
	+ **list:** Ekli olan kitapları listelemek için kullanılmalıdır.  
	+ **delete** Id belirterek veritabanında o id'ye ait olan kayıt silinecektir.
	> örn: delete 1
	+ **buy** Id ve alınacak kitap sayısı girerek kullanıcı satın alma işlemi yapabilir. Ilk değer id ikinci değer ise kitap sayısıdır.
	> buy 1 5


<br>  
<br>  


# Uygulama Öncesi Hazırlık

- **Proje içerisinde kullanılan yapılar/fonkiyonlar:** if/else, for-range, slice,struct,switch
- **Proje içeriğinde kullanılan paketler:** `strigs`, `os`, `ftm`,`gorm`,`encoding/csv`,`log`
	`os` | `helper`,`models`(kendi oluşturduğum paket) 


<br>  
<br>  


# Uygulama Aşamaları

Bu projede komut satırına girilen verileri okuyup girdiye göre kitapların listelenmesi (list) ,aranması (search), silinmesi (delete) ve güncellenmesi (update) işlemlerini gerçekleştirdim.

```
func init() {

	var err error

	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=* password=*")

	if err != nil {
		panic(err)
	}

	booksRepo = models.NewBookRepo(db)

	booksRepo.Setup()
	 booksRepo.InsertData()

}

```

- Öncelikle init fonskiyonu içinde veritabanı ile ilgili bağlantı işlemlerini yaptım.

```
func (g *BookDB) Setup() {
	g.db.AutoMigrate(&Book{})

}

```

- Models içinde bulunan Setup fonksiyonu ile migration işlemini gerçekleştirdim.



<br>


```
func (g *BookDB) InsertData() {

	var books []Book

	file, err := os.Open("/media/canberk/hdd1/HW/HW-3/booklist.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	text, err := csv.NewReader(file).ReadAll()

	if err != nil {

		log.Fatal(err)

	}

	for _, record := range text {

		books = append(books, Book{

			StockNumber: helper.RandomNumber(1, 100),
			PageNumber:  helper.RandomNumber(1, 100),
			Price:       helper.RandomFloat(1, 100),
			Name:        record[3],
			StockCode:   helper.RandomString(5),
			Isbn:        helper.RandomString(5),
			Author:      record[6],
		})

	}

	for _, v := range books[1:] {
		g.db.Create(&v)
	}

}

```

- Yine models içinde bulunan InsertData fonksiyonunda csv dosyasında bilgileri okuyarak struct üzerinden bu verileri veritabanına kaydettim.


```
input := os.Args

	firstArg := strings.ToLower(input[1])

	switch firstArg {
	case "search":
		secondArg := strings.ToLower(input[2])
		bookResult := booksRepo.GetBookByName(secondArg)

		for _, v := range bookResult {
			fmt.Println(v.Name)
		}

	case "list":
		bookList := booksRepo.GetAllBook()
		for _, v := range bookList {
			fmt.Println(v.Name)
		}
	case "delete":
		secondArg, _ := strconv.Atoi(input[2])
		booksRepo.DeleteByID(secondArg)

	case "buy":
		secondArg, _ := strconv.Atoi(input[2])
		thirdArg, _ := strconv.Atoi(input[3])
		booksRepo.UpdateStock(secondArg, thirdArg, books)

	}


```

- Bu işlemlerden sonra main içinde kullanıcı tarafından girilen değerleri okuyarak bu değerlere göre işlemler yaptım. 

<br>

```
func (g *BookDB) GetAllBook() []Book {

	var bookList []Book
	g.db.Find(&bookList)

	return bookList

}
func (g *BookDB) DeleteByID(id int) {
	var books []Book

	var n []int

	g.db.Model(&books).Pluck("id", &n)

	g.db.Unscoped().Delete(&books, id)

	isDeleted := helper.CheckSlice(n, id)

	if !isDeleted {
		panic("id not found")
	}

}

func (g *BookDB) GetBookByName(name string) []Book {

	var book []Book
	g.db.Where(" Name LIKE ?", "%"+name+"%").Find(&book)

	return book

}

func (g *BookDB) UpdateStock(id, stock int, book Book) {

	var n []int

	g.db.Model(&book).Pluck("stock_number", &n)

	stoc_num := n[id-1]

	if stoc_num <= 0 || stoc_num < stock {
		panic("not in stock")
	}

	newStock := stoc_num - stock

	g.db.Model(&book).Where("id = ?", id).Update("stock_number", newStock)

}


```

- Main içinde kullanılan fonksiyonları models içinden aldım. Burada öncellikle **listeme** işlemi için GetAllBook fonskiyonunu kullandım `Find` ile db ye kayıtlı bütün verileri getirdim.
- **Delete** işlemi için DeleteByID fonksiyonunu kullandım burada `Unscoped` ile girilen id'ye ait kaydın tamamen veritabanından sildim. Girilen id'ye ait değerin daha önceden silinip silinmediğini veya veritabanında bulunup bulunmadığını kontrol edebilmek için `g.db.Model(&books).Pluck("id", &n)` sorgusunu kullanarak `n` isimli slice'a veritabanında bulunan bütün idleri aktardım. Ardından helper içinde tanımladığım `CheckSlice` fonksiyonunu kullanarak dbden gelen idler ile girilen idyi karşılaştırdım eğer bu id mevcut değilse panic() ile bir hata oluşturdum.
- **Search** işlemi için ise GetBookByName fonksiyonunu kullandım. Bu fonksiyonda `Where(" Name LIKE ?", "%"+name+"%").Find(&book)` kullanarak aranan girdi ile dbde sorgulama yaptım ve `Book` olarak return ettim.
- **Buy** ile update işlemlerini ise UpdateStock fonksiyonu içinde gerçekleştirdim. Öncelikle `n` isimli bir slice tanımladım ve `g.db.Model(&book).Pluck("stock_number", &n)` kullanarak dbde bulunan bütün `stock_num` değerlerini getirdim. Ardından bu slice içinde bulunan değerleri kullanarak ürünün yeterli sayıda olup olmadığını kontrol ettim ardından eğer yeterli sayıda ürün var ise update işlemini gerçekleştirdim.

# Uygulama Çıktısı
Aşağıda uygulamanın çalışırlığını gözlemleyebilirsiniz.

<br>

![CRUD]()  







