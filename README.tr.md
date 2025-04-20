# Go-Wordlistgen 🔑

[![en](https://img.shields.io/badge/lang-en-red.svg)](README.md)
[![tr](https://img.shields.io/badge/lang-tr-green.svg)](README.tr.md)

Go dilinde yazılmış, kişisel bilgilere dayalı özelleştirilmiş wordlist oluşturucu.

## Özellikler

- 🖥️ İki Arayüz: Terminal Kullanıcı Arayüzü (TUI) veya Komut Satırı Arayüzü (CLI) seçeneği
- 👤 Kişisel Bilgi Tabanlı: Şu bilgileri kullanarak wordlist oluşturma:
  - Ad ve soyad
  - Doğum tarihi
  - İlgili kelimeler
- 🔄 Gelişmiş Varyasyonlar:
  - Leet (1337) dönüşümleri
  - Büyük-küçük harf varyasyonları
  - Uzunluk sınırlamaları
- 💾 Özelleştirilebilir çıktı dosyası konumu

## Kurulum

```bash
go install github.com/efeaslansoyler/go-wordlistgen@latest
```

## Kullanım

### TUI Modu (Varsayılan)

Basitçe çalıştırın:
```bash
go-wordlistgen
```

### CLI Modu

```bash
go-wordlistgen --cli [seçenekler]

Seçenekler:
  -c, --cli                CLI modunda çalıştır
  -f, --firstname string   Ad (ve varsa ikinci ad)
  -l, --lastname string    Soyad
  -b, --birthday string    Doğum tarihi (GG/AA/YYYY formatında)
  -w, --words string       Virgülle ayrılmış ilgili kelimeler
      --min string        Minimum şifre uzunluğu (varsayılan "6")
      --max string        Maksimum şifre uzunluğu (varsayılan "12")
  -o, --output string     Çıktı dosyası yolu (varsayılan "wordlist.txt")
      --leet             Leet konuşma varyasyonlarını etkinleştir
      --caps             Büyük-küçük harf varyasyonlarını etkinleştir
```

Örnek:
```bash
go-wordlistgen --cli -f "Ahmet" -l "Yılmaz" -b "01/01/1990" -w "hobi,evcilhayvan,şehir" --leet --caps
```

## Lisans

Bu proje MIT Lisansı ile lisanslanmıştır - detaylar için [LICENSE](LICENSE) dosyasına bakınız.

## Yazar

Efe Aslan Söyler (efeaslan1703@gmail.com)
