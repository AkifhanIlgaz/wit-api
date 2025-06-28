
# Wit API

Wit API, kullanıcıların kıyafet kombinlerini paylaşıp, diğer kullanıcılardan ilham alabileceği bir sosyal medya platformunun backend servisidir. Firebase ile entegre çalışarak kimlik doğrulama, veri depolama ve dosya yönetimi gibi işlemleri gerçekleştirir.

## Özellikler

* **Kullanıcı Yönetimi:**
    * Firebase Authentication ile güvenli kullanıcı kaydı ve girişi.
    * Kullanıcı profili oluşturma ve güncelleme (profil fotoğrafı, kullanıcı adı vb.).
    * Diğer kullanıcıları takip etme ve takibi bırakma.
* **Kıyafet (Outfit) Paylaşımı:**
    * Kullanıcıların kendi kıyafet kombinlerini (outfit) fotoğraflarıyla birlikte paylaşabilmesi.
    * Paylaşılan kıyafetlere ürün linkleri ekleyebilme.
    * Anasayfa akışında takip edilen kullanıcıların paylaşımlarını görme.
* **Sosyal Etkileşim:**
    * Paylaşılan kıyafetleri beğenme ve beğeniyi geri alma.
    * Beğenilen veya ilham alınan kıyafetleri kaydetme.
* **Veri Yönetimi:**
    * Firestore ile kullanıcı, kıyafet ve sosyal etkileşim verilerinin yönetimi.
    * Cloud Storage for Firebase ile kullanıcıların yüklediği fotoğrafların saklanması ve güvenli bir şekilde erişilmesi.

## Teknolojiler

* **Go:** Backend servisinin geliştirildiği ana dil.
* **Firebase:**
    * **Authentication:** Kullanıcı kimlik doğrulama işlemleri.
    * **Firestore:** NoSQL veritabanı.
    * **Cloud Storage:** Dosya depolama.
* **Chi:** Hafif ve hızlı bir Go web framework'ü.

## Kurulum

1. **Projeyi klonlayın:**
   ```bash
   git clone https://github.com/AkifhanIlgaz/wit-api.git
   ```
2. **Gerekli bağımlılıkları yükleyin:**
   ```bash
   go mod tidy
   ```
3. **Firebase projenizi oluşturun ve konfigürasyon dosyalarınızı ayarlayın:**
    * Firebase konsolundan yeni bir proje oluşturun.
    * Authentication, Firestore ve Storage servislerini aktif edin.
    * Proje ayarlarından servis hesabı (service account) için özel bir anahtar (private key) oluşturun ve `GOOGLE_APPLICATION_CREDENTIALS` ortam değişkenine bu dosyanın yolunu atayın.
4. **Uygulamayı çalıştırın:**
   ```bash
   go run main.go
   ```

## API Endpoints

Tüm endpointler `setup/routes.go` dosyasında tanımlanmıştır. Başlıca endpointler şunlardır:

* `POST /user/new`: Yeni kullanıcı oluşturur.
* `PUT /user/update`: Kullanıcı bilgilerini günceller.
* `POST /outfit/new`: Yeni kıyafet paylaşımı yapar.
* `GET /outfit/home`: Takip edilen kullanıcıların kıyafetlerini listeler.
* `PUT /outfit/like`: Bir kıyafeti beğenir.
* `PUT /user/follow`: Bir kullanıcıyı takip eder.

Detaylı bilgi için `setup/routes.go` dosyasını inceleyebilirsiniz.
