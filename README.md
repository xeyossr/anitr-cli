<div>
 
 # ANITR-CLI
  **anitr-cli:** Hızlı bir şekilde anime araması yapabileceğiniz ve istediğiniz animeyi Türkçe altyazılı izleyebileceğiniz terminal aracıdır 💫 Anime severler için hafif, pratik ve kullanışlı bir çözüm sunar 🚀

  <p>
    <a href="https://github.com/xeyossr/anitr-cli/releases">
      <img src="https://img.shields.io/github/v/release/xeyossr/anitr-cli?style=for-the-badge&include_prereleases&label=GitHub%20Release" alt="GitHub Release">
    </a>
    <a href="https://github.com/mstsecurity/anitr-cli-windows">
      <img src="https://img.shields.io/github/v/release/mstsecurity/anitr-cli-windows?include_prereleases&display_name=release&label=Windows%20Fork&style=for-the-badge" alt="Windows Fork">
    </a>
    <a href="https://aur.archlinux.org/packages/anitr-cli">
      <img src="https://img.shields.io/aur/version/anitr-cli?style=for-the-badge" alt="AUR">
    </a>
  </p>
</div>

## 💻 Kurulum

### 🐧 Linux Kullanıcıları

Eğer Arch tabanlı bir dağıtım kullanıyorsanız, [AUR](https://aur.archlinux.org/packages/anitr-cli) üzerinden tek bir komut ile indirebilirsiniz:

```bash
yay -S anitr-cli
```

Eğer Arch tabanlı olmayan bir dağıtım kullanıyorsanız, **en son sürümü** indirmek için aşağıdaki komutları kullanabilirsiniz:
```bash
curl -L -o /tmp/anitr-cli https://github.com/xeyossr/anitr-cli/releases/latest/download/anitr-cli

sudo mv /tmp/anitr-cli /usr/bin/anitr-cli
sudo chmod +x /usr/bin/anitr-cli
```

#### Güncelleme

Her çalıştırdığınızda yeni bir güncelleme olup olmadığı denetlenecektir. Eğer güncelleme mevcutsa, şu komutla güncelleyebilirsiniz:

- **AUR** üzerinden kurulum yaptıysanız:
```bash
yay -Sy anitr-cli
```

- **Manuel** kurulum yaptıysanız:
> Eğer manuel kurulum yaptıysanız, güncellemeleri manuel olarak yapmanız gerekmektedir. Alternatif olarak, [anitrcli-update](anitrcli-update) betiğini indirip `/usr/bin` dizinine yerleştirerek otomatik güncelleme yapabilirsiniz.


### 🪟 Windows Kullanıcıları

Bu proje Linux için geliştirilmiştir. **Windows kullanıcıları**, [anitr-cli-windows](https://github.com/mstsecurity/anitr-cli-windows) forkunu kullanabilirler:

> 🔗 [https://github.com/mstsecurity/anitr-cli-windows](https://github.com/mstsecurity/anitr-cli-windows)

## ✨ Özellikler

- **🔍 Hızlı Anime Arama**: Anime adı ile hızlı arama
- **🇹🇷 Türkçe Altyazı Desteği**: Otomatik Türkçe altyazı bulma
- **⭐ Favori Sistemi**: Sevdiğiniz animeleri favorilere ekleyin
- **📚 İzleme Geçmişi**: İzlediğiniz bölümleri takip edin
- **🎮 Discord Rich Presence**: Discord'da ne izlediğinizi gösterin
- **🖥️ Çoklu UI Desteği**: Modern Bubble Tea TUI ve Rofi arayüzleri
- **📺 Çoklu Çözünürlük**: 720p, 1080p seçenekleri
- **🔄 Otomatik Güncelleme**: En son sürüme otomatik güncelleme
- **🎨 Modern UI**: Bubble Tea framework ile geliştirilmiş interaktif arayüz
- **⌨️ Klavye Navigasyonu**: Vim benzeri klavye kısayolları (hjkl)

## 👾 Kullanım

```bash
Usage of ./anitr-cli:
  -disable-rpc
    	Discord Rich Presence özelliğini devre dışı bırakır.
  -rofi
    	Rofi arayüzü ile başlatır.
  -rofi-flags string
    	Rofi için flag'ler
  -update
    	anitr-cli aracını en son sürüme günceller.
  -version
    	versiyon
```

### 🎯 Ana Menü Seçenekleri

1. **Anime Ara**: Yeni anime aramak için
2. **Favoriler**: Favori animelerinizi görüntüleyin ve izleyin
3. **İzleme Geçmişi**: Son izlediğiniz animeleri görün
4. **Gelişmiş Arama**: Filtreleme seçenekleri (yakında)
5. **Çıkış**: Uygulamadan çık

### 📁 Veri Depolama

Uygulama verileriniz şu konumda saklanır:
- **Linux/macOS**: `~/.config/anitr-cli/`
  - `favorites.json`: Favori animeleriniz
  - `history.json`: İzleme geçmişiniz

### 🎮 Klavye Kısayolları

#### Ana Menü
- `↑/↓` veya `j/k`: Yukarı/aşağı hareket
- `Enter` veya `Space`: Seçim yapma
- `q` veya `Ctrl+C`: Çıkış

#### Arama Ekranı
- `Yazma`: Arama terimi girme
- `Enter`: Arama yapma
- `Backspace`: Son karakteri silme
- `Esc`: Ana menüye dönme

#### Favoriler/Geçmiş
- `↑/↓` veya `j/k`: Liste navigasyonu
- `Enter` veya `Space`: Anime seçme
- `Esc`: Ana menüye dönme

#### Genel
- `Ctrl+C`: Uygulamadan çıkış
- `q`: Hızlı çıkış (ana menüde)

## Sorunlar

Eğer bir sorunla karşılaştıysanız lütfen bir [**issue**](https://github.com/xeyossr/anitr-cli/issue) açarak karşılaştığınız problemi detaylı bir şekilde açıklayın.

## Lisans

Bu proje GNU General Public License v3.0 (GPL-3) altında lisanslanmıştır. Yazılımı bu lisansın koşulları altında kullanmakta, değiştirmekte ve dağıtmakta özgürsünüz. Daha fazla ayrıntı için lütfen [LICENSE](LICENSE) dosyasına bakın.
