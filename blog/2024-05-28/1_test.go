package roku

import (
   "fmt"
   "testing"
)

func TestOne(t *testing.T) {
   var one one_response
   err := one.New()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", one)
}
