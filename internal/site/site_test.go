package site

import (
 "os"
 "testing"
)

func TestSlugs(t *testing.T) {
 for _,s:=range []string{"../private","Upper","a/b",""} {if ValidSlug(s){t.Fatalf("accepted unsafe slug %q",s)}}
 if !ValidSlug("find-my-dev"){t.Fatal("valid slug rejected")}
}
func TestLocales(t *testing.T) {
 if err:=os.Chdir("../..");err!=nil{t.Fatal(err)}
 for _,locale:=range []string{"pt","en"}{c,err:=Load(locale);if err!=nil{t.Fatal(err)};if len(c.Projects)!=4{t.Fatalf("%s: expected four projects",locale)}}
 if _,err:=Load("../private");err==nil{t.Fatal("unsafe locale accepted")}
}
