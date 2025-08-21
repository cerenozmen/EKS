# ETKİNLİK KAYIT SİSTEMİ

#### Bu proje, mikroservis yapısını kullanarak hazırlanmış basit bir etkinlik yönetim ve kayıt sistemidir. Kullanıcılar sisteme üye olabilir, giriş yapabilir, yeni etkinlikler oluşturabilir ve mevcut etkinlikleri listeleyebilir. Ayrıca istedikleri etkinliklere kayıt olabilir ya da daha önce yaptıkları kayıtları iptal edebilirler. Projenin amacı, küçük parçalar halinde çalışan servislerle daha düzenli, kolay geliştirilebilir ve esnek bir uygulama ortaya koymaktır.

## PROJE AÇIKLAMASI

#### Bu sistem, üç ana mikroservisten oluşmaktadır:
- #### User Service: Kullanıcı kayıt, giriş ve kimlik doğrulama işlemlerinden sorumludur. JWT (JSON Web Token) kullanarak güvenli bir kimlik doğrulama mekanizması sağlar.
- #### Event Service: Etkinliklerin oluşturulması, listelenmesi ve detaylarının görüntülenmesi gibi işlevleri yönetir. Veritabanı ile etkileşim kurar ve performans artışı için Redis'i bir önbellek katmanı olarak kullanır.
- #### Booking Service: Kullanıcıların etkinliklere kayıt yapmasını ve iptal etmesini sağlar. Rezervasyon öncesi etkinlik kontenjanını ve aktiflik durumunu kontrol etmek için Event Service ile, kullanıcı kimliğini doğrulamak için ise User Service ile iletişim kurar.

#### Bu servisler, HTTP üzerinden birbirleriyle haberleşir ve her biri kendi veritabanı veya önbellek kaynağını kullanabilir.

## GEREKSİNİMLER
#### Projeyi yerel makinenizde çalıştırmak için aşağıdaki yazılımların kurulu olması gerekir:

- #### Go (versiyon 1.18 veya üstü)
- #### Docker ve Docker Compose

## KURULUM

#### 1. Projeyi klonlayın:
 > git clone "proje-repo-adresi"

 > cd "proje-klasoru"

#### 2. Veritabanı ve Redis konteynerlerini başlatın:

 > docker-compose up -d

#### 3.Her bir mikroservisin bağımlılıklarını indirin:

> cd user-service && go mod tidy && cd ..

> cd event-service && go mod tidy && cd ..

> cd booking-service && go mod tidy && cd ..


#### 4.Her servisi ayrı ayrı çalıştırın:

> cd user-service

> go run main.go

> cd event-service

> go run main.go

> cd booking-service

> go run main.go

## API DÖKÜMANTASYONU

#### User Service (:3031)

| Endpoint | Metot | Açıklama |
| :--- | :--: | --- |
| `/register` | `POST` | Yeni bir kullanıcı kaydı oluşturur. |
| `/login` | `POST` | Kullanıcı girişi yapar ve JWT token döndürür. |
| `/me` | `GET` | Gönderilen token'ı doğrulayarak kullanıcı bilgilerini getirir. |

#### Event Service (:8081)

| Endpoint | Metot | Açıklama |
| :--- | :--- | :--- |
| `/events` | `POST` | Yeni bir etkinlik oluşturur. |
| `/events` | `GET` | Tüm etkinlikleri listeler. İsteğe bağlı olarak `isActive` sorgu parametresiyle filtreleme yapılabilir. |
| `/events/:id` | `GET` | Belirtilen ID'ye sahip etkinliğin detaylarını getirir. |

#### Booking Service (:4041)

| Endpoint | Metot | Açıklama |
| :--- | :--- | :--- |
| `/bookings` | `POST` | Yetkili bir kullanıcının bir etkinliğe kayıt olmasını sağlar. Bu işlem öncesinde kullanıcı ve etkinlik bilgileri doğrulanır, ayrıca etkinlik kontenjanı kontrol edilir. |
| `/bookings` | `DELETE` | Yetkili bir kullanıcının daha önce yaptığı bir rezervasyonu iptal etmesini sağlar. |