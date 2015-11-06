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

    var nearest float64 = float64(phi_input)

    for _, kc := range K {
        kc.delta = math.Abs(kc.phi - float64(phi_input))
        if nearest > kc.delta {
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
    fmt.Println("Text language: ", DoPhiTest(buffer))
}

