Simple Movie Application To learn with basic golang, grpc, crdb.


Some points to remember :
1. crdb issue : 
    ````
    var db *gorm.DB
    var err errOR
   db, err = gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")+"&application_name=$ docs_simplecrud_gorm"), &gorm.Config{})
   if err != nil { 
        log.Fatal(err) 
   }
    Here you have declared db *gorm.DB as global so when initialing you have to not do like this 
    db, err := gorm.....
   
    It will throw this error
    panic: runtime error: 
   invalid memory address or nil pointer dereference 
   [signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x14adba2] 
   goroutine 1 [running]: gorm.io/gorm.(*DB).Create(0x1853423?, {0x170b640?, 0xc0000e22d0?}) 
   /Users/avinashreddyk/go/pkg/mod/gorm.io/gorm@v1.24.3/finisher_api.go:18 +0x22
   ````
   
2. You can use Postman to check responses check path carefully.
3. Run coackroachDb by 
    ````  
    cockroach start-single-node --advertise-addr 'localhost' --insecure
    ````
    to run locally.
4. Remember to know what is this error 
    ````
   undefined grpc.supportpackageisversion7
   it is still coming even if all packages are in latest check this, but still it is working...
   ````
   

