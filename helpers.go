package main
import (
  "os"
  "html/template"
  "strings"
  "github.com/jszroberto/kindle-words/kindledb"
)


func exportToFolders(path string,extension string,tmpl string,words []kindledb.Word,joined bool) error {

  outputs := map[string][]Word{}

  for _,word := range words{
    outputs[word.Book]=append(outputs[word.Book],Word{strings.ToLower(word.Value), template.HTML(word.Usage),word.IsEnglish(),word.GetLanguage()})
  }

  if err := os.MkdirAll(path, 0755); err != nil {
    panic("Unable to create directory for tagfile! - " + err.Error())
  }

  for key,value := range outputs {

    if joined {
      // t := template.New("Book template")
      file, err := os.Create(path+"/"+key+"." + extension)
      if err != nil {
        return err
      }
      defer file.Close()

      temp, err := template.ParseFiles(tmpl)
      if err != nil {
        return err
      }

      book:= Book{key,value}
      err = temp.Execute(file, book)
      if err != nil {
        return err
      }

      }else{
        if err := os.MkdirAll(path+"/"+value[0].Language+"/"+key, 0755); err != nil {
          panic("Unable to create directory for tagfile! - " + err.Error())
        }

        for _,word := range value {
          file, err := os.Create(path+"/"+word.Language+"/"+key+"/"+word.Value+ "." + extension)
          if err != nil {
            return err
          }

          temp, err := template.ParseFiles(tmpl)
          if err != nil {
            return err
          }
          err = temp.Execute(file, word)
          if err != nil {
            return err
          }
          file.Close()
        }
      }


    }
    return nil
  }
