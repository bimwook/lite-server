package woo;
import (
  "bytes"
);

//AboutMe 关于
func AboutMe() string{
  var buf bytes.Buffer;
  buf.WriteString(`  Lite-Server/1.1` + "\r\n");
  buf.WriteString(`-----------------------------------------` + "\r\n");
  buf.WriteString(`  - Name: 杨波 (Bamboo Young)` + "\r\n");
  buf.WriteString(`  - QQ: 3262706` + "\r\n");
  buf.WriteString(`  - E-Mail: woo@omuen.com` + "\r\n");
  buf.WriteString(`-----------------------------------------` + "\r\n");
  return buf.String();  
}