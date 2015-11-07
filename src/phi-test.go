/*
 *                               Copyright (C) 2015 by Rafael Santiago
 *
 * This is a free software. You can redistribute it and/or modify under
 * the terms of the GNU General Public License version 2.
 *
 */
package main

import (
    "strings"
    "io/ioutil"
    "fmt"
    "os"
    "math"
)

type KappaCtx struct {
    k float64
    phi float64
    delta float64
}

func GetBufferFromFile(filepath string) (string, error) {
    var retval []byte
    var err error
    retval, err  = ioutil.ReadFile(filepath);
    return string(retval), err
}

func NormalizeBuffer(buffer string) string {
    //  WARN(Santiago): Unicode sucks as hell... I am sleepy and by now just jumping out this corral.
    //                  "maybe tomorrow... it is such a shame to waste your time away like this." :)
    var ascii []byte
    if len(buffer) < 2048 {
        ascii = []byte(strings.ToUpper(buffer))
    } else {
        ascii = []byte(strings.ToUpper(buffer[:2048]))
    }
    retval := ""
    for _, a := range ascii {
        switch (a) {
            case 0xc0, 0xc1, 0xc2, 0xc3, 0xc4, 0xc5,
                 0xe0, 0xe1, 0xe2, 0xe3, 0xe4, 0xe5:
                retval += "A"
                break

            case 0xc8, 0xc9, 0xca, 0xcb,
                 0xe8, 0xe9, 0xea, 0xeb:
                retval += "E"
                break

            case 0xcc, 0xcd, 0xce, 0xcf,
                 0xec, 0xed, 0xee, 0xef:
                retval += "I"
                break

            case 0xd2, 0xd3, 0xd4, 0xd5, 0xd6,
                 0xf2, 0xf3, 0xf4, 0xf5, 0xf6:
                retval += "O"
                break

            case 0xd9, 0xda, 0xdb, 0xdc,
                 0xf9, 0xfa, 0xfb, 0xfc:
                retval += "U"
                break

            case 0xdd,
                 0xfd,
                 0x9f:
                retval += "Y"
                break

            case 0xc7,
                 0xe7:
                retval += "C"
                break

            case 0xd1,
                 0xf1:
                retval += "N"
                break

            default:
                if a >= 'A' && a <= 'Z' {
                    retval += string(a)
                }
                break
        }
    }
    return retval
}

func GetOption(option, stdvalue string) string {
    for _, arg := range os.Args {
        if strings.HasPrefix(arg, "--" + option + "=") {
            return arg[len("--" + option + "="):]
        }
        if strings.HasPrefix(arg, "--" + option) {
            return "1"
        }
    }
    return stdvalue
}

func DoPhiTest(buffer string) []string {

    if len(buffer) == 0 {
        return make([]string,0)
    }

    var K map[string]*KappaCtx
    K = make(map[string]*KappaCtx)
    K["Random"] = &KappaCtx{0.0385, 0, 0}
    K["Portuguese"] = &KappaCtx{0.0781, 0, 0}
    K["French"] = &KappaCtx{0.0778, 0, 0}
    K["Spanish"] = &KappaCtx{0.0775, 0, 0}
    K["English"] = &KappaCtx{0.0667, 0, 0}

    l := len(buffer)

    for _, kc := range K {
        kc.phi = kc.k * float64(l) * (float64(l-1))
    }

    var phi_alpha map[string]int
    phi_alpha = make(map[string]int)
    for _, b := range buffer {
        phi_alpha[string(b)]++
    }

    phi_input := 0
    for _, b := range buffer {
        f := phi_alpha[string(b)]
        phi_input += f
    }

    var nearest float64 = -1

    for _, kc := range K {
        kc.delta = math.Abs(kc.phi - float64(phi_input))
        if nearest == -1 || nearest > kc.delta {
            nearest = kc.delta
        }
    }

    var retval []string

    retval = make([]string, 0)
    for k, kc := range K {
        if kc.delta == nearest {
            retval = append(retval, k)
        }
    }

    return retval
}

func main() {
    if GetOption("help", "0") == "1" {
        fmt.Println("usage: phi-test --from-file=<filepath>|--buffer=<buffer>")
        os.Exit(1)
    }
    var buffer string
    var err error
    if option := GetOption("buffer", ""); option != "" {
        buffer = option
    } else if option := GetOption("from-file", ""); option != "" {
        buffer, err = GetBufferFromFile(option)
        if err != nil {
            fmt.Println(err.Error())
            os.Exit(1)
        }
    } else {
        fmt.Println("duh: No buffer to guess about.")
        os.Exit(1)
    }
    fmt.Println("Text language: ", DoPhiTest(NormalizeBuffer(buffer)))
}
