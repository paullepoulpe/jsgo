// Oracle tool: vectors for internal/jutil/sax_parser.go. Parses every *.xml file of a directory
// (sorted by name) with the JDK SAX parser configured exactly as XtbMessageBundle configures it,
// and prints the ContentHandler events, one per line, strings escaped as in ParsingDump.java.
// Documents that are not well-formed print the exception class instead of the remaining events.
//
//   javac -d build JutilSaxVectors.java && java -cp build JutilSaxVectors <dir> > sax.txt
import java.io.ByteArrayInputStream;
import java.io.StringReader;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;
import java.util.stream.Stream;
import javax.xml.XMLConstants;
import javax.xml.parsers.SAXParser;
import javax.xml.parsers.SAXParserFactory;
import org.xml.sax.Attributes;
import org.xml.sax.ContentHandler;
import org.xml.sax.InputSource;
import org.xml.sax.Locator;
import org.xml.sax.XMLReader;

public final class JutilSaxVectors {
  static String esc(String s) {
    if (s == null) {
      return "<null>";
    }
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < s.length(); i++) {
      char c = s.charAt(i);
      if (c == '\\') {
        sb.append("\\\\");
      } else if (c == '\t') {
        sb.append("\\t");
      } else if (c == '\n') {
        sb.append("\\n");
      } else if (c == '\r') {
        sb.append("\\r");
      } else if (c < 0x20 || c > 0x7e) {
        sb.append(String.format("\\u%04x", (int) c));
      } else {
        sb.append(c);
      }
    }
    return sb.toString();
  }

  public static void main(String[] args) throws Exception {
    List<Path> files = new ArrayList<>();
    try (Stream<Path> s = Files.list(Paths.get(args[0]))) {
      s.filter(p -> p.toString().endsWith(".xml")).sorted().forEach(files::add);
    }
    for (Path f : files) {
      StringBuilder out = new StringBuilder();
      out.append("== ").append(f.getFileName()).append('\n');
      try {
        SAXParserFactory factory = SAXParserFactory.newInstance();
        factory.setValidating(false);
        factory.setXIncludeAware(false);
        factory.setFeature("http://xml.org/sax/features/external-general-entities", false);
        factory.setFeature("http://xml.org/sax/features/external-parameter-entities", false);
        factory.setFeature("http://apache.org/xml/features/nonvalidating/load-external-dtd", false);
        factory.setFeature(XMLConstants.FEATURE_SECURE_PROCESSING, true);
        SAXParser parser = factory.newSAXParser();
        XMLReader reader = parser.getXMLReader();
        reader.setEntityResolver((publicId, systemId) -> new InputSource(new StringReader("")));
        reader.setContentHandler(
            new ContentHandler() {
              @Override
              public void setDocumentLocator(Locator locator) {}

              @Override
              public void startDocument() {
                out.append("startDocument\n");
              }

              @Override
              public void endDocument() {
                out.append("endDocument\n");
              }

              @Override
              public void startPrefixMapping(String prefix, String uri) {
                out.append("startPrefixMapping ").append(esc(prefix)).append('\n');
              }

              @Override
              public void endPrefixMapping(String prefix) {
                out.append("endPrefixMapping ").append(esc(prefix)).append('\n');
              }

              @Override
              public void startElement(String uri, String localName, String qName, Attributes atts) {
                out.append("start ").append(esc(uri)).append('|').append(esc(localName)).append('|')
                    .append(esc(qName));
                for (int i = 0; i < atts.getLength(); i++) {
                  out.append(' ').append(esc(atts.getQName(i))).append('=').append(esc(atts.getValue(i)));
                }
                out.append('\n');
              }

              @Override
              public void endElement(String uri, String localName, String qName) {
                out.append("end ").append(esc(uri)).append('|').append(esc(localName)).append('|')
                    .append(esc(qName)).append('\n');
              }

              @Override
              public void characters(char[] ch, int start, int length) {
                out.append("chars ").append(esc(new String(ch, start, length))).append('\n');
              }

              @Override
              public void ignorableWhitespace(char[] ch, int start, int length) {
                out.append("ws ").append(esc(new String(ch, start, length))).append('\n');
              }

              @Override
              public void processingInstruction(String target, String data) {
                out.append("pi ").append(esc(target)).append('|').append(esc(data)).append('\n');
              }

              @Override
              public void skippedEntity(String name) {
                out.append("skipped ").append(esc(name)).append('\n');
              }
            });
        reader.parse(new InputSource(new ByteArrayInputStream(Files.readAllBytes(f))));
      } catch (Exception e) {
        out.append("error ").append(e.getClass().getName()).append('\n');
      }
      System.out.print(out);
    }
  }
}
