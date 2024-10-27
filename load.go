package name_gen

import (
    "fmt"
    "bytes"
    "os"
    "encoding/binary"
    "encoding/base64"
)

func load() error {
    if _fioData == nil || NameBytesMax <= 0 {
        return fmt.Errorf("bad input")
    }

    for country, itemData := range *_fioData {
        for _, item := range []*FioDataCountryItem{
            itemData.Names,
            itemData.Surnames,
            itemData.Midnames,
        } {
            if err := loadItemFromImportEnc(item, NameBytesMax); err != nil {
                return fmt.Errorf("failed load [country:%v]: %w", country, err)
            }
            // If need: load from binary encoded source
            //if err := loadItemFromPathEnc(item, NameBytesMax); err != nil {
            //    return fmt.Errorf("failed load [country:%v]: %v", country, err)
            //}
        }
    }
    return nil
}

func loadItemFromImportEnc(item *FioDataCountryItem, elementMaxBytes int) error {
    if item == nil {
        return fmt.Errorf("bad input")
    }

    listEncB64 := item.GetEncB64F()
    if len(listEncB64) <= 0 {
        return fmt.Errorf("empty encoded data")
    }

    listEncByt, err := base64.StdEncoding.DecodeString(listEncB64)
    if err != nil {
        return fmt.Errorf("failed base64 decode: %w", err)
    }

    listEncBytBuf := bytes.NewBuffer(listEncByt)
    listEnc := make(ListEncoded, len(listEncByt) / elementMaxBytes)
    if err := binary.Read(listEncBytBuf, binary.LittleEndian, &listEnc); err != nil {
        return fmt.Errorf("failed binary read: %w", err)
    }

    item.List = make(FioDataCountryItemList, len(listEnc))
    for i, _ := range listEnc {
        decodeElement(&item.List[i], listEnc[i])
    }

    return nil
}

func loadItemFromPathEnc(item *FioDataCountryItem, elementMaxBytes int) error {
    if item == nil {
        return fmt.Errorf("bad input")
    }

    fileEnc, err := os.OpenFile(item.PathEnc, os.O_RDONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed open file: %w", err)
    }
    defer fileEnc.Close()

    if _, err = fileEnc.Seek(0, 0); err != nil {
        return fmt.Errorf("failed file seek: %w", err)
    }

    fInfo, err := fileEnc.Stat()
    if err != nil {
        return fmt.Errorf("failed file stat: %w", err)
    }

    fsize := int(fInfo.Size())
    if fsize <= 0 {
        return fmt.Errorf("empty encoded data")
    }

    listEnc := make(ListEncoded, fsize / elementMaxBytes)
    if err := binary.Read(fileEnc, binary.LittleEndian, &listEnc); err != nil {
        return fmt.Errorf("failed binary read: %w", err)
    }

    item.List = make(FioDataCountryItemList, len(listEnc))
    for i, _ := range listEnc {
        decodeElement(&item.List[i], listEnc[i])
    }

    return nil
}

func decodeElement(dest *string, src ListEncodedElement) {
    n := bytes.IndexByte(src[:], 0) // Search for zero byte (i.e. end of string)
    if n == -1 {
        *dest = string(src[:])
    } else {
        *dest = string(src[:n])
    }
}
