package main

import (
    "fmt"
    "log"
    "net/http"
    // "io/os.
    "strings"
    "strconv"
    "math/rand"
    "time"
    "os"
)

func randfilename()  string {
    rand.Seed(time.Now().UnixNano())
    
    for i:=0; i<1000000;i++ {
        f := strconv.FormatInt(int64(rand.Intn(1000000)), 10)
        if _, err := os.Stat("pastes/" + f); os.IsNotExist(err) {
            // file doesn exist, which is good
            return f
        }
    }

    // if we exaulst those options, start using very very long numbers
    f := strconv.FormatInt(int64(rand.Intn(1000000000) + 1000000), 10)
    return f

}

func newpaste(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        fmt.Println("error:", err)
    }
    fmt.Println("form:",r.Form)
    // fmt.Fprintf(w, "data: %s", r.Form["d"])
    filename := ""
    if len(r.Form["name"]) > 0 {
        filename = strings.Replace(r.Form["name"][0],"/","",-1)
    } else {
        filename = randfilename()
        fmt.Println("random filename:", filename)
    }
    if _, err := os.Stat("pastes/" + filename); os.IsNotExist(err) {
        // file doesn currently exist
        err := os.WriteFile("pastes/"+filename, []byte(r.Form["d"][0]), 0644)
        if err != nil {
            fmt.Fprintf(w, "Error: %s", err)
            return
        }
        fmt.Fprintf(w, "OK /%s",filename)
    }
    // r.Form["d"]
}
 

func homePage(w http.ResponseWriter, r *http.Request){
    // fmt.Fprintf(w, "Nothing to see here")

    if r.Method == "GET" {
        if r.URL.Path == "/" || r.URL.Path == "/index.html" || r.URL.Path == "" {
            content, err := os.ReadFile("static/index.html") 
            if err != nil {
                fmt.Println(err)
            }
            fmt.Fprintf(w, string(content))   
        } else {
            pastepath := strings.Replace(r.URL.Path, "/", "", -1)
            fmt.Println("loading paste:", pastepath)
            content, err := os.ReadFile("pastes/"+pastepath) 
            if err != nil {
                fmt.Println(err)
                http.NotFound(w, r)
                return
            }
            fmt.Fprintf(w, string(content))
        }
        
    } else if (r.Method == "POST" || r.Method == "PUT") {
        //TODO new paste
        newpaste(w, r)
        return
    }
    
    // fmt.Println("Endpoint Hit: homePage, path:", r.URL)


}


func handleRequests() {
    http.HandleFunc("/", homePage)
    // http.HandleFunc("/new", newpaste)
    log.Fatal(http.ListenAndServe(":9000", nil))
}
 
func main() {
    handleRequests()
}