package name_gen

import (
    "fmt"
    "time"
    "math/rand"
)

func Generate(country Country) (string, error) {
    _once.Do(func() { // For lazy load on call Generate()
        _fioData = configure()
        if err := load(); err != nil { }
    })
    if _fioData == nil {
        return "", fmt.Errorf("bad input")
    }

    fio, ok := (*_fioData)[country]
    if !ok {
        return "", fmt.Errorf("bad input: country: %v", country)
    }

    rand.Seed(time.Now().UTC().UnixNano())
    name    := fio.Names.List[ rand.Intn(len(fio.Names.List)) ]
    surname := fio.Surnames.List[ rand.Intn(len(fio.Surnames.List)) ]
    midname := fio.Midnames.List[ rand.Intn(len(fio.Midnames.List)) ]

    return surname + " " + name + " " + midname, nil
}
