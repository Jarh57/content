package content

// import (
//   "database/sql/driver"
//   "encoding/json"
//   "fmt"
// )

// type Provider map[string]string

// // LangChain runs the lang chain over the contents of the container.
// func (content Provider) Chain(lang string) string {
//   return LangChain(content, lang)
// }

// func (content Provider) Value() (driver.Value, error) {
//   serialized, err := json.Marshal(content)
//   if err != nil {
//     return nil, fmt.Errorf("content/translated: cannot serialize value: %s", err)
//   }

//   return serialized, nil
// }

// func (content *Provider) Scan(value interface{}) error {
//   b, ok := value.([]byte)
//   if !ok {
//     return fmt.Errorf("content/translated: cannot scan type into bytes: %T", value)
//   }

//   if err := json.Unmarshal(b, content); err != nil {
//     return fmt.Errorf("content/translated: cannot scan value: %s", err)
//   }

//   return nil
// }

// // LangChain returns a value following a prearranged chain of preference.
// //
// // The chain gives preference to the requested lang, then english and finally spanish
// // if no other translation is available.
// func LangChain(translations map[string]string, lang string) string {
//   if translations[lang] != "" {
//     return translations[lang]
//   }
//   if translations["en"] != "" {
//     return translations["en"]
//   }
//   return translations["es"]
// }
