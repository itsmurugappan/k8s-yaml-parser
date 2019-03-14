package main

import (
  "fmt"
  "os"
  "strings"
  "io/ioutil"
  "log"
  "bufio"
  "io"

  "github.com/ghodss/yaml"
)


func main() {

  filePath, _ := os.LookupEnv("file_path")

  file, _ := os.Open(filePath)

  reader := bufio.NewReader(file)

  var str strings.Builder

  destDir, _ := os.LookupEnv("dir")

  for {
    bytesRead, errs := reader.ReadString('\n')
    if strings.TrimSpace(bytesRead) != "---" {
      str.WriteString(bytesRead)
    } else {
        if strings.TrimSpace(str.String()) == "" {
          str.Reset()
        } else {
          var result map[string]interface{}
          err := yaml.Unmarshal([]byte(str.String()), &result)
          if err != nil {
            log.Fatalf("error: %v", err)
          }

          kind := fmt.Sprintf("%v",result["kind"])

          mdi := result["metadata"]

          md, _ := mdi.(map[string]interface{})

          newName := fmt.Sprintf("%v",md["name"])

          if _, err := os.Stat(destDir + "/" + kind); os.IsNotExist(err) {
            os.MkdirAll(destDir + "/" + kind, 0777)
          }

          err2 := ioutil.WriteFile(destDir + "/" + kind + "/" + newName + ".yaml", []byte(str.String()), 0777)

          if err2 != nil {
                  log.Fatalf("error: %v", err2)
          }
          str.Reset()
        }
    }
    if errs != nil {
      if errs == io.EOF {
        fmt.Printf("success")
      } else
      {
        log.Fatalf("error: %v", errs)
      }
      break
    }
  }
}
