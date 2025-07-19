# 📖 ANITR-CLI Kullanım Kılavuzu

Bu kılavuz, anitr-cli'nin yeni özelliklerini nasıl kullanacağınızı detaylı olarak açıklar.

## 🚀 Hızlı Başlangıç

1. **Uygulamayı başlatın**:
   ```bash
   anitr-cli
   ```

2. **Ana menüden seçim yapın**:
   - Anime Ara
   - Favoriler
   - İzleme Geçmişi
   - Gelişmiş Arama
   - Çıkış

## ⭐ Favori Sistemi

### Favori Ekleme
1. Anime ara menüsünden bir anime seçin
2. İzleme menüsünde "Favorilere Ekle" seçeneğini seçin
3. Anime favorilerinize eklenir

### Favori Çıkarma
1. Favori bir animenin izleme menüsünde "Favorilerden Çıkar" seçeneğini seçin
2. Anime favorilerden kaldırılır

### Favorileri Görüntüleme
1. Ana menüden "Favoriler" seçeneğini seçin
2. Favori animelerinizin listesini görün
3. İzlemek istediğiniz animeyi seçin

### Kaldığı Yerden Devam Etme
Favori bir anime seçtiğinizde:
- Eğer daha önce izlemişseniz, "Kaldığı yerden devam et" seçeneği görünür
- "Baştan başla" seçeneği ile sıfırdan başlayabilirsiniz

## 📚 İzleme Geçmişi

### Otomatik Kayıt
- Her izlediğiniz bölüm otomatik olarak geçmişe kaydedilir
- Anime adı, bölüm adı, izleme tarihi kaydedilir

### Geçmişi Görüntüleme
1. Ana menüden "İzleme Geçmişi" seçeneğini seçin
2. Son 20 izlediğiniz bölümü görün
3. Tekrar izlemek istediğiniz bölümü seçin

### Geçmişi Temizleme
1. İzleme geçmişi menüsünde "--- Geçmişi Temizle ---" seçeneğini seçin
2. Onay verin
3. Tüm geçmiş silinir

## 🔍 Anime Arama

### Basit Arama
1. Ana menüden "Anime Ara" seçeneğini seçin
2. Anime adını yazın
3. Sonuçlardan seçim yapın

### Arama İpuçları
- Türkçe veya İngilizce anime adları kullanabilirsiniz
- Kısmi isimler de çalışır (örn: "naruto" → "Naruto Shippuden")
- Büyük/küçük harf duyarlı değildir

## 🎬 İzleme Özellikleri

### Çözünürlük Seçimi
1. İzleme menüsünde "Çözünürlük seç" seçeneğini seçin
2. Mevcut kalite seçeneklerini görün (720p, 1080p, vb.)
3. İstediğiniz kaliteyi seçin

### Bölüm Navigasyonu
- **Sonraki bölüm**: Bir sonraki bölüme geç
- **Önceki bölüm**: Bir önceki bölüme geç
- **Bölüm seç**: Belirli bir bölümü seç

### Discord Rich Presence
- Otomatik olarak Discord'da ne izlediğinizi gösterir
- `--disable-rpc` parametresi ile kapatabilirsiniz

## 🖥️ Arayüz Seçenekleri

### Terminal (TUI) Modu
```bash
anitr-cli
```

### Rofi Modu
```bash
anitr-cli --rofi
```

### Rofi ile Özel Ayarlar
```bash
anitr-cli --rofi --rofi-flags "-theme ~/.config/rofi/anime.rasi"
```

## 📁 Veri Yönetimi

### Veri Konumları
- **Favoriler**: `~/.config/anitr-cli/favorites.json`
- **Geçmiş**: `~/.config/anitr-cli/history.json`
- **Loglar**: Uygulama çalışma dizininde

### Yedekleme
```bash
# Favori ve geçmiş verilerinizi yedekleyin
cp -r ~/.config/anitr-cli ~/anitr-cli-backup
```

### Geri Yükleme
```bash
# Yedekten geri yükleyin
cp -r ~/anitr-cli-backup ~/.config/anitr-cli
```

## 🔧 Sorun Giderme

### Yaygın Sorunlar

#### Favoriler Görünmüyor
- Config dizininin var olduğunu kontrol edin: `ls ~/.config/anitr-cli/`
- Dosya izinlerini kontrol edin: `ls -la ~/.config/anitr-cli/`

#### Geçmiş Kaydedilmiyor
- Disk alanınızı kontrol edin
- Yazma izinlerinizi kontrol edin

#### Anime Bulunamıyor
- İnternet bağlantınızı kontrol edin
- Farklı arama terimleri deneyin
- VPN kullanıyorsanız kapatmayı deneyin

### Log Dosyaları
Hata durumunda log dosyalarını kontrol edin:
```bash
ls -la *.log
```

### Temiz Kurulum
Tüm verileri sıfırlamak için:
```bash
rm -rf ~/.config/anitr-cli/
```

## 🎯 İpuçları ve Püf Noktaları

1. **Hızlı Erişim**: Sık izlediğiniz animeleri favorilere ekleyin
2. **Kalite Ayarı**: İnternet hızınıza göre çözünürlük seçin
3. **Geçmiş Takibi**: Uzun seriler için geçmiş özelliğini kullanın
4. **Rofi Kullanımı**: Daha hızlı navigasyon için Rofi modunu deneyin
5. **Yedekleme**: Önemli favori listelerinizi düzenli yedekleyin

## 🔮 Gelecek Özellikler

- **Gelişmiş Filtreler**: Tür, yıl, puan bazlı arama
- **Entegrasyon**: MyAnimeList, AniList bağlantısı
- **Bildirimler**: Yeni bölüm bildirimleri
- **Temalar**: Farklı renk temaları

## 📞 Destek

Sorun yaşıyorsanız:
1. Bu kılavuzu kontrol edin
2. [GitHub Issues](https://github.com/xeyossr/anitr-cli/issues) sayfasını ziyaret edin
3. Yeni bir issue açın (detaylı açıklama ile)

---

**İyi Seyirler! 🍿**
