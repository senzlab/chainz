package main

import (
    "os"
    "bytes"
    "text/template"
    "path/filepath"
)

func read(la LienAdd) {
    // format template path 
    cwd, _ := os.Getwd()
    tp := filepath.Join(cwd, "./template/lienadd.xml")
    println(tp)

    // template from file
    t, err := template.ParseFiles(tp)
    if err != nil {
        println(err.Error())
    }

    // parse template
    var buf bytes.Buffer
    err = t.Execute(&buf, la)
    if err != nil {
        println(err.Error())
    }

    println(buf.String())
}
