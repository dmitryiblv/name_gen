package name_gen

import (
    "fmt"
    "strings"
    "os"
    "io"
    "io/ioutil"
    "encoding/binary"
    "encoding/base64"
)

func encode() error {
    _fioData = configure()
    if _fioData == nil || NameBytesMax <= 0 {
        return fmt.Errorf("bad input")
    }

    for country, itemData := range *_fioData {
        for _, item := range []*FioDataCountryItem{
            itemData.Names,
            itemData.Surnames,
            itemData.Midnames,
        } {
            if err := encodeItem(item, NameBytesMax, country); err != nil {
                return fmt.Errorf("failed encode item [country:%v]: %w", country, err)
            }
        }
    }
    return nil
}

func encodeItem(item *FioDataCountryItem, elementMaxBytes int, country Country) error {
    if item == nil {
        return fmt.Errorf("bad input")
    }

    contentTxt, err := ioutil.ReadFile(item.PathTxt)
    if err != nil {
        return fmt.Errorf("failed read file: %w", err)
    }

    listEnc := make(ListEncoded, 0, 32)
    for _, row := range strings.Split(string(contentTxt), "\n") {
        if row != "" {
            row = strings.TrimSpace(row)
            if l := len(row); l > elementMaxBytes {
                return fmt.Errorf("row '%v' is too long: %v; path: %v", row, l, item.PathTxt)
            }
            el := ListEncodedElement{}
            copy(el[:], row)
            listEnc = append(listEnc, el)
        }
    }

    // PathEnc
    fileEnc, err := os.OpenFile(item.PathEnc, os.O_RDWR | os.O_CREATE, 0644)
    if err != nil {
        return fmt.Errorf("failed open file: %w", err)
    }
    defer fileEnc.Close()

    err = binary.Write(fileEnc, binary.LittleEndian, listEnc)
    if err != nil {
        return fmt.Errorf("failed binary write: %w", err)
    }

    // PathGo
    fileGo, err := os.OpenFile(item.PathGo, os.O_RDWR | os.O_CREATE, 0644)
    if err != nil {
        return fmt.Errorf("failed open file: %w", err)
    }
    defer fileGo.Close()

    if _, err := fileEnc.Seek(0, 0); err != nil {
        return fmt.Errorf("failed file seek: %w", err)
    }

    listEncBytes, err := io.ReadAll(fileEnc)
    if err != nil {
        return fmt.Errorf("failed file read all: %w", err)
    }

    listEncB64 := make([]byte, base64.StdEncoding.EncodedLen(len(listEncBytes)))
    base64.StdEncoding.Encode(listEncB64, listEncBytes)

    contentEnc := "package " + string(country) + "\n" +
        "func " + item.GetEncB64N + "() string { return `" + string(listEncB64) + "` }"

    if n, err := fmt.Fprintf(fileGo, contentEnc); err != nil {
        return fmt.Errorf("failed write file: %w", err)
    } else if n <= 0 {
        return fmt.Errorf("failed write file '%v': zero bytes written", item.PathGo)
    }

    return nil
}
