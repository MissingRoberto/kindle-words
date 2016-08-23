package kindledb


type Word struct{
    Value string  `sql:"word"`
    Usage string  `sql:"usage"`
    Book string  `sql:"book"`
    Language string  `sql:"language"`
    Count int  `sql:"count_usage"`
    Mastered int `sql:"mastered"`
}

type BookInfo struct{
    Guid string  `sql:"word"`
    Language string  `sql:"lang"`
    Title string  `sql:"title"`
    Authors string  `sql:"authors"`
}


func (w *Word) IsEnglish() bool {
    return w.Language == "en" || w.Language == "en-GB"
}

func (w *Word) IsMastered() bool {
    return w.Mastered == 100
}

func (w *Word) GetLanguage() string {
    if w.Language == "en" || w.Language == "en-GB"{
      return "english"
    }else if w.Language == "de"{
      return "german"
    }else if w.Language == "es"{
      return "spanish"
    }
    return ""
}
