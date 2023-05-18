package paramount

import (
   "fmt"
   "testing"
   "time"
)

func Test_Item(t *testing.T) {
   token, err := New_App_Token()
   if err != nil {
      t.Fatal(err)
   }
   for _, test := range tests {
      item, err := token.Item(test.content_ID)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(item)
      name, err := item.Name()
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(name)
      time.Sleep(time.Second)
   }
}
