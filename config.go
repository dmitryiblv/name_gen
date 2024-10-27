package name_gen

import (
    "name_gen/ru"
)

func configure() *FioData {
    return &FioData{
        CountryRu: &FioDataCountry{
            Names: &FioDataCountryItem{
                PathTxt:    "ru/names.txt",
                PathEnc:    "ru/names.enc",
                PathGo:     "ru/names.go",
                GetEncB64N: "GetNamesEncB64",
                GetEncB64F: ru.GetNamesEncB64,
                List:       FioDataCountryItemList{},
            },
            Surnames: &FioDataCountryItem{
                PathTxt:    "ru/surnames.txt",
                PathEnc:    "ru/surnames.enc",
                PathGo:     "ru/surnames.go",
                GetEncB64N: "GetSurnamesEncB64",
                GetEncB64F: ru.GetSurnamesEncB64,
                List:       FioDataCountryItemList{},
            },
            Midnames: &FioDataCountryItem{
                PathTxt:    "ru/midnames.txt",
                PathEnc:    "ru/midnames.enc",
                PathGo:     "ru/midnames.go",
                GetEncB64N: "GetMidnamesEncB64",
                GetEncB64F: ru.GetMidnamesEncB64,
                List:       FioDataCountryItemList{},
            },
        },
    }
}
