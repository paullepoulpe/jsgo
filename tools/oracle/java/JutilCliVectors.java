// Vector generator for jutil (black-box JDK output): OutputStreamWriter
// encoding per standard charset, Charset.isSupported and Level.parse.
// Output: internal/jutil/testdata/cli_charset_vectors.tsv
import java.io.*;
import java.nio.charset.*;
import java.util.logging.Level;
public final class JutilCliVectors {
  static String hex(byte[] b){StringBuilder sb=new StringBuilder();for(byte x:b)sb.append(String.format("%02x",x));return sb.toString();}
  public static void main(String[] a) throws Exception {
    String[] cs={"US-ASCII","ISO-8859-1","UTF-8","UTF-16","UTF-16BE","UTF-16LE"};
    String[] in={"aéb","x𝒳y","p\ud835q","\udcb3z","€",""};
    for(String c:cs) for(String s:in){
      ByteArrayOutputStream bo=new ByteArrayOutputStream();
      Writer w=new BufferedWriter(new OutputStreamWriter(bo,Charset.forName(c)));
      w.write(s); w.write(s); w.close();
      System.out.println(c+"\t"+hex(s.getBytes(StandardCharsets.UTF_8))+"\t"+hex(bo.toByteArray()));
    }
    for(String n:new String[]{"","a b","-x","x-","utf8","UTF8","latin1","windows-1252","foo","1a","a:b","a_b","é"}){
      String r; try{ r=String.valueOf(Charset.isSupported(n)); }catch(Exception e){ r=e.toString(); }
      System.out.println("isSupported\t"+n+"\t"+r);
    }
    for(String n:new String[]{"WARNING","warning","OFF","ALL","FINEST","1000","900","123","-5","x","", " INFO"}){
      String r; try{ Level l=Level.parse(n); r=l.getName()+"/"+l.intValue(); }catch(Exception e){ r=e.toString(); }
      System.out.println("level\t"+n+"\t"+r);
    }
  }
}
