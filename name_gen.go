package name_gen

import (
    "sync"
)

type Country string

const (
    CountryRu Country = "ru"
)

const (
    NameBytesMax = 32 // Max len of name or surname or middle name
                      // NOTE: In bytes, not in symbols
                      //       Onchange: re-encode *.enc, *.go
)

type FioData map[Country]*FioDataCountry

type FioDataCountry struct {
    Names       *FioDataCountryItem
    Surnames    *FioDataCountryItem
    Midnames    *FioDataCountryItem
}

type FioDataCountryItem struct {
    PathTxt     string
    PathEnc     string
    PathGo      string
    GetEncB64N  string
    GetEncB64F  func() string
    List        FioDataCountryItemList
}

type FioDataCountryItemList []string

type ListEncoded []ListEncodedElement

type ListEncodedElement [NameBytesMax]byte

var _fioData *FioData

var _once sync.Once
