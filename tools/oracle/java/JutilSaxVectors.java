// Oracle tool: vectors for internal/jutil/sax_parser.go, from the JDK's SAXParser configured
// exactly as XtbMessageBundle configures it (non-validating, namespace-unaware, no external
// entities or DTD, secure processing, an EntityResolver returning an empty source).
//
//   javac -d build JutilSaxVectors.java
//   java -cp build JutilSaxVectors <testdata dir> [seed [count]]
//
// Writes two files:
//   - <dir>/sax/expected.txt: the events of every <dir>/sax/*.xml (sorted by name), one per
//     line, without positions; a document that is not well-formed prints the exception class
//     instead of the remaining events. Replayed by TestSAXParserMatchesJava.
//   - <dir>/sax_vectors.txt: hand-written documents plus <count> (default 400) documents generated
//     from <seed> (default 20260415), well-formed and malformed, in UTF-8, UTF-16 and
//     ISO-8859-1. Each block is "== <name>", "doc <bytes>", then one line per event with the
//     Locator position at the time of the call ("@line:column"), the exception (class, line,
//     column, message) and whatever the parser printed on System.err. Gzip it to
//     sax_vectors.txt.gz; replayed by TestSAXParserVectors (SAX_VECTORS=<file> replays another).
//
// Bytes are printed as-is for printable ASCII, "\\" for a backslash and "\xHH" otherwise;
// strings in events are escaped as in ParsingDump.java.
import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.PrintStream;
import java.io.StringReader;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;
import java.util.Random;
import java.util.stream.Stream;
import javax.xml.XMLConstants;
import javax.xml.parsers.SAXParser;
import javax.xml.parsers.SAXParserFactory;
import org.xml.sax.Attributes;
import org.xml.sax.ContentHandler;
import org.xml.sax.InputSource;
import org.xml.sax.Locator;
import org.xml.sax.SAXParseException;
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

  static String escBytes(byte[] b) {
    StringBuilder sb = new StringBuilder();
    for (byte x : b) {
      int c = x & 0xff;
      if (c == '\\') {
        sb.append("\\\\");
      } else if (c >= 0x20 && c <= 0x7e) {
        sb.append((char) c);
      } else {
        sb.append(String.format("\\x%02x", c));
      }
    }
    return sb.toString();
  }

  static SAXParserFactory newFactory() throws Exception {
    SAXParserFactory factory = SAXParserFactory.newInstance();
    factory.setValidating(false);
    factory.setXIncludeAware(false);
    factory.setFeature("http://xml.org/sax/features/external-general-entities", false);
    factory.setFeature("http://xml.org/sax/features/external-parameter-entities", false);
    factory.setFeature("http://apache.org/xml/features/nonvalidating/load-external-dtd", false);
    factory.setFeature(XMLConstants.FEATURE_SECURE_PROCESSING, true);
    return factory;
  }

  // ---------------------------------------------------------------------------------------------
  // Old format: testdata/sax/*.xml -> expected.txt (events without positions).

  static String runPlain(byte[] data) {
    StringBuilder out = new StringBuilder();
    try {
      SAXParser parser = newFactory().newSAXParser();
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
      reader.parse(new InputSource(new ByteArrayInputStream(data)));
    } catch (Exception e) {
      out.append("error ").append(e.getClass().getName()).append('\n');
    }
    return out.toString();
  }

  // ---------------------------------------------------------------------------------------------
  // New format: events with Locator positions, exceptions with messages, System.err.

  static String runFull(byte[] data, ByteArrayOutputStream errBuf) {
    StringBuilder out = new StringBuilder();
    Locator[] loc = new Locator[1];
    errBuf.reset();
    try {
      SAXParser parser = newFactory().newSAXParser();
      XMLReader reader = parser.getXMLReader();
      reader.setEntityResolver(
          (publicId, systemId) -> {
            out.append("resolveEntity ").append(esc(publicId)).append('|').append(esc(systemId))
                .append(pos(loc[0])).append('\n');
            return new InputSource(new StringReader(""));
          });
      reader.setContentHandler(
          new ContentHandler() {
            @Override
            public void setDocumentLocator(Locator locator) {
              loc[0] = locator;
              out.append("locator").append(pos(locator)).append('\n');
            }

            @Override
            public void startDocument() {
              out.append("startDocument").append(pos(loc[0])).append('\n');
            }

            @Override
            public void endDocument() {
              out.append("endDocument").append(pos(loc[0])).append('\n');
            }

            @Override
            public void startPrefixMapping(String prefix, String uri) {
              out.append("startPrefixMapping ").append(esc(prefix)).append('|').append(esc(uri))
                  .append(pos(loc[0])).append('\n');
            }

            @Override
            public void endPrefixMapping(String prefix) {
              out.append("endPrefixMapping ").append(esc(prefix)).append(pos(loc[0])).append('\n');
            }

            @Override
            public void startElement(String uri, String localName, String qName, Attributes atts) {
              out.append("start ").append(esc(uri)).append('|').append(esc(localName)).append('|')
                  .append(esc(qName));
              for (int i = 0; i < atts.getLength(); i++) {
                out.append(' ').append(esc(atts.getQName(i))).append('=').append(esc(atts.getValue(i)));
              }
              out.append(pos(loc[0])).append('\n');
            }

            @Override
            public void endElement(String uri, String localName, String qName) {
              out.append("end ").append(esc(uri)).append('|').append(esc(localName)).append('|')
                  .append(esc(qName)).append(pos(loc[0])).append('\n');
            }

            @Override
            public void characters(char[] ch, int start, int length) {
              out.append("chars ").append(esc(new String(ch, start, length))).append(pos(loc[0]))
                  .append('\n');
            }

            @Override
            public void ignorableWhitespace(char[] ch, int start, int length) {
              out.append("ws ").append(esc(new String(ch, start, length))).append(pos(loc[0]))
                  .append('\n');
            }

            @Override
            public void processingInstruction(String target, String data) {
              out.append("pi ").append(esc(target)).append('|').append(esc(data)).append(pos(loc[0]))
                  .append('\n');
            }

            @Override
            public void skippedEntity(String name) {
              out.append("skipped ").append(esc(name)).append(pos(loc[0])).append('\n');
            }
          });
      reader.parse(new InputSource(new ByteArrayInputStream(data)));
    } catch (SAXParseException e) {
      out.append("error ").append(e.getClass().getName()).append(' ').append(e.getLineNumber())
          .append(':').append(e.getColumnNumber()).append(' ').append(esc(e.getPublicId()))
          .append('|').append(esc(e.getSystemId())).append(' ').append(esc(e.getMessage()))
          .append('\n');
    } catch (Exception e) {
      out.append("error ").append(e.getClass().getName()).append(' ').append(esc(e.getMessage()))
          .append('\n');
    }
    System.err.flush();
    String err = errBuf.toString(StandardCharsets.UTF_8);
    if (!err.isEmpty()) {
      out.append("stderr ").append(esc(err)).append('\n');
    }
    return out.toString();
  }

  static String pos(Locator l) {
    if (l == null) {
      return " @null";
    }
    String ids = "";
    if (l.getPublicId() != null || l.getSystemId() != null) {
      ids = " [" + esc(l.getPublicId()) + "|" + esc(l.getSystemId()) + "]";
    }
    return " @" + l.getLineNumber() + ":" + l.getColumnNumber() + ids;
  }

  // ---------------------------------------------------------------------------------------------
  // Hand-written documents.

  static final List<String[]> HAND = new ArrayList<>();

  static void hand(String name, String doc) {
    HAND.add(new String[] {name, doc});
  }

  static {
    // XTB-shaped documents.
    String head =
        "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"
            + "<!DOCTYPE translationbundle SYSTEM \"translationbundle.dtd\">\n";
    hand("xtb_basic", head + "<translationbundle lang=\"zh-HK\">\n"
        + "<translation id=\"7639678437384034548\">descargar</translation>\n"
        + "<translation id=\"2398375912250604550\">Se han\nignorado <ph name=\"NUM\"/> conversaciones.</translation>\n"
        + "<translation id=\"3945720239421293834\"></translation>\n"
        + "</translationbundle>");
    hand("xtb_entities", head + "<translationbundle lang=\"fr\">\n"
        + "<translation id=\"1\">a &amp; b &lt;c&gt; d &quot;e&apos; &#233;&#x1F600; &unknown; x</translation>\n"
        + "</translationbundle>\n");
    hand("xtb_crlf", head.replace("\n", "\r\n") + "<translationbundle lang=\"de\">\r\n"
        + "<translation id=\"10\">eins\r\nzwei\r\ndrei</translation>\r\n</translationbundle>\r\n");
    hand("xtb_gender", head + "<translationbundle lang=\"es_ES\">\n<translation id=\"7\">\n"
        + "<branch variants=\"variants { grammatical_gender_variant { grammatical_gender_case:\nMASCULINE } }\">"
        + "Bienvenido, <ph name=\"NAME\"/></branch></translation>\n</translationbundle>\n");
    hand("xtb_malformed", head + "<translationbundle lang=\"fr\">\n<translation id=\"1\">a</translatio>\n");
    hand("xtb_unclosed", head + "<translationbundle lang=\"fr\">\n<translation id=\"1\">a");
    // Character data chunking around newlines, brackets and supplementary characters.
    hand("chunk_newlines", "<r>a\nb\nc\nd\ne\nf\n\n\ng\r\nh\ri\r\r\nj</r>");
    hand("chunk_brackets", "<r>x]y]z]w ]a]]b]]]c\n]d ]]]x]] ]]</r>");
    hand("chunk_supplementary", "<r>a😀b😀😀c d\n😀\n]😀]</r>");
    hand("chunk_long", "<r>" + "word ".repeat(7000) + "</r>");
    hand("chunk_long_lines", "<r>" + "line of text\n".repeat(2000) + "</r>");
    hand("chunk_long_cdata", "<r><![CDATA[" + "cdata ]] text\n".repeat(1500) + "]]></r>");
    hand("chunk_long_comment", "<r>a<!--" + "comment text ".repeat(1500) + "-->b</r>");
    hand("chunk_long_attr", "<r a=\"" + "attr value ".repeat(1500) + "\">t</r>");
    // Prolog.
    hand("decl_standalone", "<?xml version=\"1.0\" standalone=\"yes\"?><a/>");
    hand("decl_spaces", "<?xml   version = '1.0'   encoding = 'utf-8'   standalone = 'no'   ?>\n<a/>");
    hand("decl_long", "<?xml version='1.0'" + " ".repeat(100) + "?><a>" + "x".repeat(20000) + "</a>");
    hand("decl_latin1", "<?xml version=\"1.0\" encoding=\"ISO-8859-1\"?><a>café</a>");
    hand("doctype_public", "<!DOCTYPE a PUBLIC \"-//X//DTD Y//EN\" \"y.dtd\"><a>&e;</a>");
    hand("doctype_subset", "<!DOCTYPE a SYSTEM \"a.dtd\" [\n<!-- c --><?pi x?>\n<!ELEMENT a ANY>\n"
        + "<!ENTITY ext SYSTEM \"ext.xml\">\n<!ENTITY % pe \"x\">\n<!NOTATION n PUBLIC \"n\">\n]>\n<a>&ext;&undeclared;</a>");
    hand("doctype_noid", "<!DOCTYPE a><a>&amp;</a>");
    hand("misc", "<!-- c1 -->\n<?pi data ?>\n<a/>\n<!-- c2 --><?q?>\n");
    // Errors.
    hand("err_undeclared", "<a>&e;</a>");
    hand("err_undeclared_standalone", "<?xml version=\"1.0\" standalone=\"yes\"?><!DOCTYPE a SYSTEM \"x\"><a>&e;</a>");
    hand("err_mismatch", "<a><b></c></a>");
    hand("err_prolog", "x<a/>");
    hand("err_trailing", "<a/>x");
    hand("err_eof", "<a>");
    hand("err_empty", "");
    hand("err_cdend", "<a>]]></a>");
    hand("err_lt_attr", "<a b=\"<\"/>");
    hand("err_dup_attr", "<a b=\"1\" b=\"2\"/>");
    hand("err_charref", "<a>&#0;</a>");
    hand("err_control", "<a>\u0001</a>");
    hand("err_fffe", "<a>￾</a>");
    hand("err_comment", "<a><!-- a -- b --></a>");
    hand("err_pi_xml", "<a><?xml x?></a>");
    hand("err_version", "<?xml version=\"2.0\"?><a/>");

    // Where characters() calls end around the end of the first buffer (32 bytes).
    String[] boundaryTails = {
      "\n\nyy", "\r\nyy", "\ryy", "]yy", "]]yyyy", "😀yy", "\n<e/>yy", "\n&amp;yy", "\n]yy",
      "\r\n<e/>", "\r\nyyy", "\r\r\nyyyy<e/>", "\r\ryyyy<e/>", "\n\ryyy", "\n\r\u0085yy", "<e/>yy",
      "</r>", "&amp;yy"
    };
    for (int t = 0; t < boundaryTails.length; t++) {
      for (int k = 20; k <= 31; k++) {
        String tail = boundaryTails[t];
        String doc = "<r>" + "x".repeat(k) + tail + (tail.endsWith("</r>") ? "" : "</r>");
        hand("boundary" + t + "_" + k, doc);
      }
    }
    // A ']' in an attribute value ends a step (seen with a decoding error
    // right after the first read).
    for (String t : new String[] {"]]>", "]>", "]]]", "]x", "x]", "]\t", "]\n"}) {
      String doc = "<a b='" + "x".repeat(26 - t.length()) + t;
      hand("attr_bracket_" + Integer.toHexString(t.hashCode()), doc + "\u0001'/>");
    }
    for (int k = 18; k <= 23; k++) {
      hand("boundary_markup_" + k, "<r><e/>" + "z".repeat(k) + "<f>\n\r\nyy<e/>&amp;\r\r</f></r>");
    }
    for (int k = 8180; k <= 8232; k += 4) {
      hand("boundary_second_" + k, "<r>" + "y".repeat(k) + "\n\nab]]c\r\n<e/>" + "x".repeat(9000) + "</r>");
    }
    // The first read splits after the four bytes used for encoding detection.
    hand("split_first_read", "<a>é" + "x".repeat(40) + "</a>");
    hand("split_first_read2", "<?pé\nÿ!>abc?><f.g>&#233; \n#\n\r\u0085&#x3c; </f.g>");
    // Columns after lone CRs.
    for (String t : new String[] {"x\r\ny", "x\ry", "x\r\ry", "x\r\r\ny", "\r\ny", "\ry", "\r\r\ny", "x\n\ry",
        "x\r\r\r\r\ny", "\r\ry]z", "\r\ry&amp;z", "x\r yz"}) {
      hand("cr_" + Integer.toHexString(t.hashCode()), "<r>" + t + "<e/></r>");
    }
    hand("cr_attr", "<r a=\"x\r\ry\" b=\"x\r\r\ny\" c=\"\r\r&amp;y\"/>");
    hand("cr_comment", "<r>x<!--\r\r-->y<![CDATA[q\r\rz]]>y<?p \r\r?>y</r>");
    hand("cr_eof", "<r>a\r\r");
    // The XML declaration and its pre-scan.
    for (String d : new String[] {
      "<?xml\nversion=\"1.0\"?><?p?>\n<a/>", "<?xml \n version=\"1.0\"?><?p ?>\n<a/>",
      "<?xml version\n=\n\"1.0\"?><a/>", "<?xml\r\n version\r\n=\r\n\"1.0\"?><a/>",
      "<?xml   version='1.0'?><?p?><a/>", "<?xml version  ='1.0'?><?p?><a/>",
      "<?xml-stylesheet href=\"a\"?><a/>", "<?xmlfoo?>\n<a/>", "<?xml?><a/>",
      "<?xml version = x?><a/>", "<?xml  version\t=.0'?><a/>", "<?xml\nversion='abcdef'?><a/>",
      "<?xml\n\nversion='abc&def'?><a/>", "<?xml\r\nversion\r\n=\r\n'1&#0;'?><a/>",
      "<?xml version='1.0\u000b'?><a/>", "<?xml version=\"1.0\" encoding=\"UT\u0001F\"?><a/>",
      "<?xml version=\"1.0\" foo=\"x\"?><a/>", "<?xml version=\"1.0\"x?><a/>",
      "<?xml version=\"1.0\"?x><a/>", "<?xml version=\"1.0\" version=\"1.0\"?><a/>",
      "<?xml version=\"1.0\" encoding=\"UTF-8\" encoding=\"UTF-8\"?><a/>",
      "<?xml version=\"1.0\" standalone=\"no\" standalone=\"no\"?><a/>",
      "<?xml version=\"1.0\"standalone=\"no\"?><a/>", "<?xml encoding=\"UTF-8\"?><a/>",
      "<?xml version=\"1.0\" standalone=\"maybe\"?><a/>", "<?xml version=\"1.0\" encoding=\"1FOO\"?><a/>",
      "<?xml version=\"1.0\" encoding=\"\"?><a/>", "<?xml version=\"1.0\" encoding=\"a\nb\"?><a/>",
      "<?xml version=\"1", "<?xml version=\"1.0", "<?xml version", "<?xml ", "<?xm", "<?xml version=\"1.0\" ?",
      "<?xml version=\"1.0\"?>", "\n\n", " <?xml version=\"1.0\"?><a/>", "<?xml ?><a/>", "<?xml \n?>",
      "<?xml   ?x", "<?xml version =\"1.0\"standalone='yes'?><a/>", "<?xml\r\nversion='1.0'encoding=\"UTF-8\"?><a/>",
      "<?xml version ='1.0'standalone='no'encoding=\"UTF-8\"?><a/>", "<?xml version=\"1.0\"encoding=\"UTF-8\"?><a/>"
    }) {
      hand("decl_" + Integer.toHexString(d.hashCode()), d);
    }
    // Error messages and positions.
    for (String d : new String[] {
      "<a b></a>", "<a b=></a>", "<a b=x></a>", "<a b=\"x\"c=\"y\"></a>", "<a b=\"x\" c></a>", "<a b=\"x></a>",
      "<a/ >", "<a></ab>", "<ab></a>", "<a></a x>", "<a></>", "<a><b></a></b>", "<a>< b/></a>", "<a><!X></a>",
      "<a><!-x></a>", "<a><!--x---></a>", "<a><![CDAT[x]]></a>", "<a><!DOCTYPE a></a>", "<a><?></a>",
      "<a><?p</a>", "<a b=\"&#0;\"/>", "<a b=\"&amp\"/>", "<a b=\"& x\"/>", "<a b=\"&undef;\"/>", "<a>x]", "<a>x]]",
      "<a><!-", "<a></", "<a>&am", "<a>&#1", "<a>]]]></a>", "<a>&#x;</a>", "<a>&#X41;</a>", "<a>&#12a;</a>",
      "<a>&#xD800;</a>", "<a>&#x110000;</a>", "<a>&#99999999999;</a>", "<a/><b/>", "<a/>&amp;", "&amp;<a/>",
      "<a/><![CDATA[x]]>", "<a/><!-", "<a>x\ny", "<a>\n"
    }) {
      hand("err_" + Integer.toHexString(d.hashCode()), d);
    }
    // The document type declaration.
    for (String d : new String[] {
      "<!DOCTYPE a [<?pi in subset?><!-- c -->]><a/>", "<!DOCTYPE a [<!ELEMENT a ANY>]><a> <b> </b> </a>",
      "<!DOCTYPE a [<!ELEMENT a (#PCDATA|b)*><!ELEMENT b (#PCDATA)>]><a> <b> </b> </a>",
      "<!DOCTYPE a [<!ENTITY e SYSTEM \"e.xml\" NDATA n>]><a>&e;</a>",
      "<!DOCTYPE a [<!ENTITY e PUBLIC \"p\" \"e.xml\">]><a>&e;</a>",
      "<!DOCTYPE a [<!ENTITY e SYSTEM \"e.xml\">]><a b=\"&e;\"/>", "<!DOCTYPE a SYSTEM \"x.dtd\"><a b=\"&e;\"/>",
      "<!DOCTYPE a [<!ENTITY amp \"&#38;#38;\">]><a>&amp;</a>", "<!DOCTYPE a [%p;]><a/>",
      "<!DOCTYPE><a/>", "<!DOCTYPE 1a><a/>", "<!DOCTYPE a SYSTEM><a/>", "<!DOCTYPE a SYSTEM x><a/>",
      "<!DOCTYPE a PUBLIC><a/>", "<!DOCTYPE a PUBLIC \"p\"><a/>", "<!DOCTYPE a PUBLIC \"p{\" \"s\"><a/>",
      "<!DOCTYPE a PUBLIC x><a/>", "<!DOCTYPE a FOO \"x\"><a/>", "<!DOCTYPE a [<!FOO>]><a/>",
      "<!DOCTYPE a [<!ELEMENT>]><a/>", "<!DOCTYPE a [<!ELEMENT a>]><a/>", "<!DOCTYPE a [<!ELEMENT a FOO>]><a/>",
      "<!DOCTYPE a [<!ELEMENT a (#PCDATA|b)>]><a/>", "<!DOCTYPE a [<!ENTITY e>]><a/>",
      "<!DOCTYPE a [<!ENTITY e x>]><a/>", "<!DOCTYPE a [<!ENTITY e SXSTEM \"x\">]><a/>",
      "<!DOCTYPE a [<!ENTITY e \"x>]><a/>", "<!DOCTYPE a [<!ENTITY % >]><a/>", "<!DOCTYPE a [<!ENTITY e \"&#0;\">]><a/>",
      "<!DOCTYPE a [<!ENTITY e \"%p;\">]><a/>", "<!DOCTYPE a [<!ENTITY e \"\u0001\">]><a/>",
      "<!DOCTYPE a [<!NOTATION n>]><a/>", "<!DOCTYPE a [<!NOTATION n FOO>]><a/>", "<!DOCTYPE a [<!ATTLIST a>]><a/>",
      "<!DOCTYPE a [ x ]><a/>", "<!DOCTYPE a [ ] x><a/>", "<!DOCTYPE a [", "<!DOCTYPE a [<!NOTATION n SYSTEM 'n'>",
      "<!DOCTYPE a [<!NOTATION n PUBLIC '-/", "<!DOCTYPE a SYSTEM \"x\"", "<!DOCTYPE a><!DOCTYPE a><a/>",
      "<!DOCTYPE a [<![INCLUDE[<!ELEMENT a ANY>]]>]><a/>", "<!DOCTYPE a [<![",
      "<!DOCTYPE a [<!-- x -->]", "<!DOCTYPE a [] \n", "<!DOCTYPE a [<!ELEMENT a (#PCDATA|x>]><a/>",
      "<!DOCTYPE a [<!ELEMENT a ( #PCDATA | x y)*>]><a/>", "<!DOCTYPE a [<!ELEMENT a (#PCDATA|x )x>]><a/>"
    }) {
      hand("dtd_" + Integer.toHexString(d.hashCode()), d);
    }
    // Public identifiers: the char after a line end counts two columns.
    for (String d : new String[] {
      "<!DOCTYPE a PUBLIC \"-//E \n 'e.xml'> ]><a/>", "<!DOCTYPE a PUBLIC \"a\r\nb\" \"x>\"><a/>",
      "<!DOCTYPE a PUBLIC \"a\rb\n\n>\" \"x\"><a/>", "<!DOCTYPE a PUBLIC \"x\n", "<!DOCTYPE a PUBLIC \"x\n\" \"y>\"><a/>",
      "<!DOCTYPE a [<!NOTATION n PUBLIC \"x\n\"\u0001>]><a/>"
    }) {
      hand("pubid_" + Integer.toHexString(d.hashCode()), d);
    }
    // A CDATA section reports what it has when a run of line ends stops
    // one char before the end of the buffer.
    for (String run : new String[] {"\n\n\n\n", "\r\n\r\n\r\n", "\r\r\r\r", "\n\n\nx"}) {
      for (int k = 8205; k <= 8210; k++) {
        hand("cdata_run_" + Integer.toHexString(run.hashCode()) + "_" + k,
            "<a><![CDATA[" + "x".repeat(k) + run + "y]]></a>");
      }
    }
    // Entity values: a line end inside a step makes the next char count two
    // columns; references to parameter entities are errors.
    for (String v : new String[] {"a\n%", "a\n\n%", "\n%", "a\r\r%", "a\r\r\n%", "a\r%", "&#65;\n%", "a&#10;\n%",
        "a\nb\u0001", "%x", "%", "% x", "a\n%x;", "&\nx;"}) {
      hand("entity_value_" + Integer.toHexString(v.hashCode()), "<!DOCTYPE a [<!ENTITY e \"" + v + "\">]><a/>");
    }
    hand("pe_ref_no_semicolon", "<!DOCTYPE a [%x ]><a/>");
    hand("pe_ref_no_name", "<!DOCTYPE a [\n%\n]><a/>");
    // Secure-processing limits.
    hand("limit_name_1000", "<" + "a".repeat(1000) + "/>");
    hand("limit_name_1001", "<" + "a".repeat(1001) + "/>");
    hand("limit_name_20000", "<" + "a".repeat(20000) + "/>");
    hand("limit_entity_name", "<a>&" + "b".repeat(1001) + ";</a>");
    // The name limit is checked when a name reaches the end of the buffer.
    for (int before : new int[] {1000, 1001, 2000}) {
      hand("limit_name_boundary_" + before, "<r>" + "y".repeat(8220 - before) + "<" + "n".repeat(before + 5) + "/></r>");
    }
    StringBuilder manyAttrs = new StringBuilder("<a");
    for (int i = 0; i <= 10000; i++) {
      manyAttrs.append(" x").append(i).append("='1'");
    }
    hand("limit_attributes", manyAttrs.append("/>").toString());
  }

  // Hand-written documents given as bytes (hex).

  static String hx(String s, String charset) {
    StringBuilder b = new StringBuilder();
    for (byte x : s.getBytes(java.nio.charset.Charset.forName(charset))) {
      b.append(String.format("%02x", x & 0xff));
    }
    return b.toString();
  }
  static final String[][] HAND_BYTES = {
    {"utf8_invalid_lead", "3c613e" + "78".repeat(40) + "ff3c2f613e"},
    {"utf8_invalid_2", "3c613e" + "78".repeat(40) + "c3283c2f613e"},
    {"utf8_invalid_3", "3c613e" + "78".repeat(40) + "e282283c2f613e"},
    {"utf8_overlong", "3c613e78e080803c2f613e"},
    {"utf8_surrogate", "3c613e78eda0803c2f613e"},
    {"utf8_planes", "3c613e" + "78".repeat(40) + "f49080803c2f613e"},
    {"utf8_truncated", "3c613e78f09f98"},
    {"utf8_fffe", "3c613e78efbfbe3c2f613e"},
    {"utf8_first_read", "3c613eff3c2f613e"},
    {"utf8_in_decl", "3c3f786d6c200a2076657273696f6e3de28222312e30223f3e3c612f3e"},
    // encoding='UTF8' takes the lenient reader.
    {"utf8_alias_lenient", "3c3f786d6c2076657273696f6e3d27312e302720656e636f64696e673d2755544638273f3e3c613e78"
        + "ff78c32878e2827878f490808078eda08078c0af78f09f983c2f613e"},
    {"utf8_alias_eof", hx("<?xml version='1.0' encoding='UTF8'?><a>xx", "UTF-8") + "e697"},
    {"ascii_strict", "3c3f786d6c2076657273696f6e3d27312e302720656e636f64696e673d2755532d4153434949273f3e3c613e78e93c2f613e"},
    {"latin1", "3c3f786d6c2076657273696f6e3d27312e302720656e636f64696e673d274c4154494e31273f3e3c613e78e9ff3c2f613e"},
    {"utf16le_bom", "fffe3c0061003e007800e90041d800de3c002f0061003e00"},
    {"utf16be_nobom", "003c003f0078006d006c002000760065007200730069006f006e003d00270031002e003000270020"
        + "0065006e0063006f00640069006e0067003d0027005500540046002d003100360027003f003e003c0061002f003e"},
    {"utf16_lone_surrogate", "feff003c0061003e0078d800003c002f0061003e"},
    // Java-only aliases of UTF-16 and US-ASCII take lenient readers.
    {"utf16_alias_be", "feff" + hx("<?xml version=\"1.0\" encoding=\"UTF16\"?><a>x", "UTF-16BE")
        + "d8000079dc00d800d800dc00" + hx("y</a>", "UTF-16BE") + "00"},
    {"utf16_alias_swapped_bom", "fffe" + hx("<?xml version=\"1.0\" encoding=\"UnicodeLittle\"?>", "UTF-16LE")
        + "feff" + hx("<a>x</a>", "UTF-16BE")},
    {"utf16_alias_long", "fffe" + hx("<?xml version=\"1.0\" encoding=\"utf-16le\"?><a>" + "x".repeat(9000)
        + "\r\n" + "y".repeat(5000) + "\ud83d\ude00</a>", "UTF-16LE")},
    {"utf16_alias_eof", "feff" + hx("<?xml version=\"1.0\" encoding=\"UTF_16\"?><a>xy", "UTF-16BE") + "00"},
    {"utf16_entity_value_surrogate", "feff" + hx("<!DOCTYPE a [<!ENTITY e \"a", "UTF-16BE") + "d800"
        + hx("b\"><!ENTITY f 'x'>]><a/>", "UTF-16BE")},
    {"utf16_entity_value_low_surrogate", "feff" + hx("<!DOCTYPE a [<!ENTITY e \"a", "UTF-16BE") + "dc00"
        + hx("b\">]><a/>", "UTF-16BE")},
    {"utf16_high_surrogate_eof", "feff" + hx("<a><?pi x", "UTF-16BE") + "dbff"},
    {"utf16_high_surrogate_eof_text", "feff" + hx("<a>x", "UTF-16BE") + "dbff"},
    {"ascii7_lenient", hx("<?xml version='1.0' encoding='ascii7'?><a>x", "UTF-8") + "e9ff" + hx("</a>", "UTF-8")},
    {"attr_bracket_invalid_byte", hx("<a b='xxxxxxxxxxxxxxxxxxxxxxx]]>", "UTF-8") + "ff" + hx("'/>", "UTF-8")},
    {"mixed_pcdata_invalid_byte", hx("<!DOCTYPE a [<!ELEMENT a (#PCDATA)", "UTF-8") + "ff" + hx(">]><a/>", "UTF-8")},
    {"pubid_invalid_byte", hx("<!DOCTYPE a PUBLIC \"a-", "UTF-8") + "ff" + hx("\" \"x\"><a/>", "UTF-8")},
    {"utf16_declares_utf8", "feff003c003f0078006d006c002000760065007200730069006f006e003d00220031002e00300022"
        + "00200065006e0063006f00640069006e0067003d0022005500540046002d00380022003f003e003c0061002f003e"},
  };

  // ---------------------------------------------------------------------------------------------
  // Random documents.

  static final class Gen {
    final Random r;
    final StringBuilder sb = new StringBuilder();
    boolean externalSubset; // undeclared entities are skipped instead of fatal
    boolean latin1;
    String encName;
    int budget;
    List<String> extEntities = new ArrayList<>();

    Gen(Random r) {
      this.r = r;
    }

    int n(int bound) {
      return r.nextInt(bound);
    }

    boolean p(int percent) {
      return r.nextInt(100) < percent;
    }

    String pick(String... options) {
      return options[n(options.length)];
    }

    String space() {
      switch (n(10)) {
        case 0:
          return "\n";
        case 1:
          return "\t";
        case 2:
          return "\r\n";
        case 3:
          return "  ";
        case 4:
          return " \n ";
        default:
          return " ";
      }
    }

    String optSpace() {
      return p(70) ? "" : space();
    }

    static final String[] NAMES = {
      "a", "b", "c", "translation", "ph", "x1", "d-e", "f.g", "h_i", "_u", "A", "B9", "ns:p", "q:",
      "café", "日本", "Àb", "long_element_name_with_many_characters_in_it"
    };

    String name() {
      return NAMES[n(NAMES.length)];
    }

    String nonAscii() {
      return pick("é", "日本", " ", "\u0085", "�", " ", "ÿ",
          "😀", "𐀀", "􏿿");
    }

    String charRef() {
      return pick("&#65;", "&#x42;", "&#10;", "&#13;", "&#9;", "&#x1F600;", "&#xD7FF;", "&#65533;",
          "&#x10FFFF;", "&#32;", "&#x3c;", "&#38;", "&#0000065;", "&#x00000041;", "&#233;");
    }

    String entRef() {
      return pick("&amp;", "&lt;", "&gt;", "&quot;", "&apos;");
    }

    // Text with the characters that matter to the scanner: newlines, brackets, '>', non-ASCII.
    void text(StringBuilder out, int len) {
      int target = out.length() + len;
      while (out.length() < target) {
        int k = n(40);
        String frag;
        if (k < 12) {
          frag = pick("word ", "x", "abc", "Hello, world", "a b", "z");
        } else if (k < 18) {
          frag = pick("\n", "\r\n", "\r", "\n\n", "\r\r\n", "\t", "\n\r");
        } else if (k < 22) {
          frag = pick("]", "]]", "]]]", "] ]", "]x]");
        } else if (k < 24) {
          frag = ">";
        } else if (k < 28) {
          frag = nonAscii();
        } else if (k < 30) {
          frag = "x".repeat(1 + n(60));
        } else {
          frag = pick(" ", "  ", "-", "?", "'", "\"", "=", "/", "!", "#", ";");
        }
        if (frag.startsWith(">") && endsWithBrackets(out)) {
          out.append(' ');
        }
        if (latin1) {
          frag = toLatin1(frag);
        }
        out.append(frag);
      }
    }

    static boolean endsWithBrackets(StringBuilder s) {
      int n = s.length();
      return n >= 2 && s.charAt(n - 1) == ']' && s.charAt(n - 2) == ']';
    }

    static String toLatin1(String s) {
      StringBuilder b = new StringBuilder();
      for (int i = 0; i < s.length(); i++) {
        char c = s.charAt(i);
        b.append(c <= 0xff ? c : 'Y');
      }
      return b.toString();
    }

    int textLen() {
      int k = n(100);
      if (k < 50) {
        return 1 + n(12);
      } else if (k < 80) {
        return 1 + n(80);
      } else if (k < 95) {
        return 1 + n(600);
      }
      return 1 + n(12000);
    }

    void comment(StringBuilder out) {
      out.append("<!--");
      StringBuilder t = new StringBuilder();
      text(t, n(4) == 0 ? textLen() : n(20));
      String s = t.toString().replace("--", "- -");
      if (s.endsWith("-")) {
        s += " ";
      }
      out.append(s).append("-->");
    }

    void pi(StringBuilder out) {
      out.append("<?").append(pick("pi", "app", "x-y", "target", "xml-stylesheet", "pé"));
      if (p(70)) {
        out.append(space());
        StringBuilder t = new StringBuilder();
        text(t, n(4) == 0 ? textLen() : n(20));
        out.append(t.toString().replace("?>", "? >"));
        if (p(30)) {
          out.append(optSpace());
        }
      }
      out.append("?>");
    }

    void cdata(StringBuilder out) {
      out.append("<![CDATA[");
      StringBuilder t = new StringBuilder();
      text(t, n(3) == 0 ? textLen() : n(20));
      if (p(30)) {
        t.append(pick("<b>", "&amp;", "<", "&", "]]"));
      }
      out.append(t.toString().replace("]]>", "]] >")).append("]]>");
    }

    void misc(StringBuilder out) {
      int count = n(4);
      for (int i = 0; i < count; i++) {
        switch (n(3)) {
          case 0:
            comment(out);
            break;
          case 1:
            pi(out);
            break;
          default:
            out.append(space());
        }
      }
      if (p(50)) {
        out.append(p(50) ? "\n" : space());
      }
    }

    String quoted(String s) {
      return p(50) ? "\"" + s + "\"" : "'" + s + "'";
    }

    void doctype(StringBuilder out) {
      out.append("<!DOCTYPE").append(space()).append(name());
      int id = n(3);
      if (id == 1) {
        out.append(space()).append("SYSTEM").append(space()).append(quoted(pick("x.dtd", "translationbundle.dtd", "", "a b")));
        externalSubset = true;
      } else if (id == 2) {
        out.append(space()).append("PUBLIC").append(space()).append(quoted(pick("-//A//B", "p", "")))
            .append(space()).append(quoted("y.dtd"));
        externalSubset = true;
      }
      if (p(40)) {
        out.append(optSpace()).append('[');
        int count = n(5);
        for (int i = 0; i < count; i++) {
          out.append(optSpace());
          switch (n(6)) {
            case 0:
              comment(out);
              break;
            case 1:
              pi(out);
              break;
            case 2:
              out.append("<!ELEMENT").append(space()).append(name()).append(space())
                  .append(pick("ANY", "EMPTY", "(#PCDATA)", "(#PCDATA|a|b)*", "( #PCDATA | x )*"))
                  .append(optSpace()).append('>');
              break;
            case 3:
              {
                String e = pick("ext", "ext2", "e");
                extEntities.add(e);
                out.append("<!ENTITY").append(space()).append(e).append(space())
                    .append(p(50) ? "SYSTEM" + space() + quoted("e.xml")
                        : "PUBLIC" + space() + quoted("-//E") + space() + quoted("e.xml"))
                    .append(optSpace()).append('>');
                break;
              }
            case 4:
              out.append("<!ENTITY").append(space()).append('%').append(space()).append(pick("pe", "p2"))
                  .append(space())
                  .append(quoted(pick("text", "", "<!ELEMENT a ANY>", "a\nb c", "a\r\nb\r\rc", "x&#65;\ny", "\n\n",
                      "a&#10;\n")))
                  .append(optSpace()).append('>');
              break;
            default:
              out.append("<!NOTATION").append(space()).append(pick("n", "gif")).append(space())
                  .append(p(50) ? "SYSTEM" + space() + quoted("n") : "PUBLIC" + space() + quoted(pick("-//N", "-//N\n x", "-//N\r\n\r\n y")))
                  .append(optSpace()).append('>');
          }
        }
        out.append(optSpace()).append(']');
      }
      out.append(optSpace()).append('>');
    }

    void xmlDecl(StringBuilder out, String encoding) {
      out.append("<?xml").append(space()).append("version").append(optSpace()).append('=')
          .append(optSpace()).append(quoted("1.0"));
      if (encoding != null) {
        out.append(space()).append("encoding").append(optSpace()).append('=').append(optSpace())
            .append(quoted(encoding));
      }
      if (p(25)) {
        out.append(space()).append("standalone").append(optSpace()).append('=').append(optSpace())
            .append(quoted(p(30) ? "yes" : "no"));
      }
      if (p(30)) {
        out.append(optSpace());
      }
      if (p(5)) {
        out.append(" ".repeat(n(60)));
      }
      out.append("?>");
    }

    void attrValue(StringBuilder out, char quote) {
      int count = n(6);
      for (int i = 0; i < count; i++) {
        int k = n(20);
        String frag;
        if (k < 6) {
          frag = pick("v", "value", "1", "a b", "x y z");
        } else if (k < 9) {
          frag = pick("\t", "\n", "\r\n", "\r", "  ");
        } else if (k < 11) {
          frag = entRef();
        } else if (k < 13) {
          frag = charRef();
        } else if (k < 15) {
          frag = nonAscii();
        } else if (k < 16) {
          frag = pick(">", "]]>", quote == '"' ? "'" : "\"");
        } else if (k < 17 && externalSubset) {
          frag = "&undeclared;";
        } else if (k < 18) {
          frag = "x".repeat(n(3) == 0 ? textLen() : 1 + n(40));
        } else {
          frag = pick("=", "/", "?", "!", "#", "-");
        }
        if (latin1) {
          frag = toLatin1(frag);
        }
        out.append(frag);
      }
    }

    void element(StringBuilder out, int depth) {
      String nm = name();
      out.append('<').append(nm);
      List<String> used = new ArrayList<>();
      int attrs = n(10) < 6 ? 0 : n(4) + 1;
      for (int i = 0; i < attrs; i++) {
        String an = pick("id", "lang", "name", "x", "a-b", "c.d", "e_f", "xmlns", "xmlns:p", "p:q", "été");
        if (used.contains(an)) {
          continue;
        }
        used.add(an);
        char q = p(50) ? '"' : '\'';
        out.append(space()).append(an).append(optSpace()).append('=').append(optSpace()).append(q);
        attrValue(out, q);
        out.append(q);
      }
      if (p(20)) {
        out.append(optSpace());
      }
      if (depth > 6 || p(15)) {
        out.append("/>");
        return;
      }
      out.append('>');
      int items = n(8);
      for (int i = 0; i < items && out.length() < budget; i++) {
        int k = n(30);
        if (k < 12) {
          text(out, textLen());
        } else if (k < 16) {
          element(out, depth + 1);
        } else if (k < 18) {
          out.append(entRef());
        } else if (k < 20) {
          out.append(charRef());
        } else if (k < 22) {
          cdata(out);
        } else if (k < 24) {
          comment(out);
        } else if (k < 26) {
          pi(out);
        } else if (k < 27 && externalSubset) {
          out.append("&undeclared;");
        } else if (k < 28 && !extEntities.isEmpty()) {
          out.append('&').append(extEntities.get(n(extEntities.size()))).append(';');
        } else {
          out.append(space());
        }
      }
      out.append("</").append(nm).append(p(15) ? optSpace() : "").append('>');
    }

    byte[] document() {
      int enc = n(100);
      String declared = null;
      boolean bom = false;
      boolean forceDecl = false;
      if (enc < 75) {
        encName = "UTF-8";
        if (p(30)) {
          declared = pick("UTF-8", "utf-8", "UTF8");
        }
        bom = p(8);
      } else if (enc < 85) {
        encName = p(50) ? "UTF-16LE" : "UTF-16BE";
        bom = p(70);
        if (!bom) {
          forceDecl = true;
        }
        if (forceDecl || p(50)) {
          declared = p(85) ? "UTF-16"
              : pick("utf-16", "UTF16", encName, encName.toLowerCase(),
                  encName.equals("UTF-16LE") ? "UnicodeLittle" : "UnicodeBig");
        }
      } else {
        encName = "ISO-8859-1";
        latin1 = true;
        declared = pick("ISO-8859-1", "iso-8859-1", "latin1");
      }
      budget = p(70) ? 400 : 30000;
      StringBuilder out = new StringBuilder();
      if (declared != null || forceDecl || p(40)) {
        xmlDecl(out, declared);
      }
      misc(out);
      if (p(45)) {
        doctype(out);
        misc(out);
      }
      element(out, 0);
      misc(out);
      String s = out.toString();
      if (latin1) {
        s = toLatin1(s);
      }
      byte[] body = s.getBytes(encName.equals("ISO-8859-1") ? StandardCharsets.ISO_8859_1
          : encName.equals("UTF-16LE") ? StandardCharsets.UTF_16LE
          : encName.equals("UTF-16BE") ? StandardCharsets.UTF_16BE : StandardCharsets.UTF_8);
      if (bom) {
        byte[] b = encName.equals("UTF-8") ? new byte[] {(byte) 0xef, (byte) 0xbb, (byte) 0xbf}
            : encName.equals("UTF-16LE") ? new byte[] {(byte) 0xff, (byte) 0xfe}
            : new byte[] {(byte) 0xfe, (byte) 0xff};
        byte[] all = new byte[b.length + body.length];
        System.arraycopy(b, 0, all, 0, b.length);
        System.arraycopy(body, 0, all, b.length, body.length);
        body = all;
      }
      return body;
    }
  }

  // Malformations, applied to the bytes of a UTF-8 document.
  static byte[] mutate(Random r, byte[] doc) {
    String[] inserts = {
      "<", "&", "]]>", "\u0001", "￾", "</x>", "<x", "\"", "'", "=", "&#0;", "&#xD800;", "&bogus",
      "<!--", "--", "<?", "?>", "<![CDATA[", "<!DOCTYPE a>", "<?xml version=\"1.0\"?>", "&#x;", "&#12",
      "&amp", "& ", "<a b c>", "</", "<a/ >", "<1>", "\u000b", "\u001f", ">", "&#65;", " ", "\n", "/>"
    };
    byte[][] rawInserts = {
      {(byte) 0xff}, {(byte) 0x80}, {(byte) 0xc3}, {(byte) 0xc3, 0x28}, {(byte) 0xe2, (byte) 0x82},
      {(byte) 0xf0, (byte) 0x9f, (byte) 0x98}, {(byte) 0xed, (byte) 0xa0, (byte) 0x80},
      {(byte) 0xf4, (byte) 0x90, (byte) 0x80, (byte) 0x80}, {(byte) 0xc0, (byte) 0xaf}, {0}
    };
    int at = doc.length == 0 ? 0 : r.nextInt(doc.length + 1);
    ByteArrayOutputStream b = new ByteArrayOutputStream();
    switch (r.nextInt(5)) {
      case 0: // delete one to three bytes
        {
          int len = 1 + r.nextInt(3);
          b.write(doc, 0, at);
          if (at + len < doc.length) {
            b.write(doc, at + len, doc.length - at - len);
          }
          break;
        }
      case 1: // truncate
        b.write(doc, 0, at);
        break;
      case 2: // raw bytes
        {
          byte[] ins = rawInserts[r.nextInt(rawInserts.length)];
          b.write(doc, 0, at);
          b.write(ins, 0, ins.length);
          b.write(doc, at, doc.length - at);
          break;
        }
      default: // markup
        {
          byte[] ins = inserts[r.nextInt(inserts.length)].getBytes(StandardCharsets.UTF_8);
          b.write(doc, 0, at);
          b.write(ins, 0, ins.length);
          b.write(doc, at, doc.length - at);
        }
    }
    return b.toByteArray();
  }

  // Malformations of a UTF-16 document (whole code units; the byte order mark is kept).
  static byte[] mutate16(Random r, byte[] doc, boolean littleEndian) {
    int units = (doc.length - 2) / 2;
    int at = 2 + 2 * (units == 0 ? 0 : r.nextInt(units + 1));
    String[] inserts = {"<", "&", "]]>", "\u0001", "\ufffe", "\ud800", "\udc00", "\ud800x", "</x>", "\"", "&#0;", "<!--"};
    ByteArrayOutputStream b = new ByteArrayOutputStream();
    switch (r.nextInt(3)) {
      case 0:
        b.write(doc, 0, at);
        if (at + 2 < doc.length) {
          b.write(doc, at + 2, doc.length - at - 2);
        }
        break;
      case 1:
        b.write(doc, 0, at);
        if (r.nextBoolean() && at < doc.length) {
          b.write(doc[at]); // an odd number of bytes
        }
        break;
      default:
        String ins = inserts[r.nextInt(inserts.length)];
        b.write(doc, 0, at);
        for (int i = 0; i < ins.length(); i++) {
          char c = ins.charAt(i);
          if (littleEndian) {
            b.write(c & 0xff);
            b.write(c >> 8);
          } else {
            b.write(c >> 8);
            b.write(c & 0xff);
          }
        }
        b.write(doc, at, doc.length - at);
    }
    return b.toByteArray();
  }

  public static void main(String[] args) throws Exception {
    Path dir = Paths.get(args[0]);
    long seed = args.length > 1 ? Long.parseLong(args[1]) : 20260415L;
    int count = args.length > 2 ? Integer.parseInt(args[2]) : 400;
    PrintStream realErr = System.err;
    ByteArrayOutputStream errBuf = new ByteArrayOutputStream();
    System.setErr(new PrintStream(errBuf, true, StandardCharsets.UTF_8));

    // expected.txt for testdata/sax/*.xml.
    Path saxDir = dir.resolve("sax");
    if (Files.isDirectory(saxDir)) {
      List<Path> files = new ArrayList<>();
      try (Stream<Path> s = Files.list(saxDir)) {
        s.filter(p -> p.toString().endsWith(".xml")).sorted().forEach(files::add);
      }
      StringBuilder out = new StringBuilder();
      for (Path f : files) {
        out.append("== ").append(f.getFileName()).append('\n');
        out.append(runPlain(Files.readAllBytes(f)));
      }
      Files.writeString(saxDir.resolve("expected.txt"), out.toString(), StandardCharsets.UTF_8);
    }

    // sax_vectors.txt.
    StringBuilder out = new StringBuilder();
    out.append("# JutilSaxVectors seed=").append(seed).append(" count=").append(count).append('\n');
    List<Object[]> cases = new ArrayList<>();
    if (Files.isDirectory(saxDir)) {
      try (Stream<Path> s = Files.list(saxDir)) {
        for (Path f : (Iterable<Path>) s.filter(p -> p.toString().endsWith(".xml")).sorted()::iterator) {
          cases.add(new Object[] {"file_" + f.getFileName(), Files.readAllBytes(f)});
        }
      }
    }
    for (String[] h : HAND) {
      cases.add(new Object[] {"hand_" + h[0], h[1].getBytes(StandardCharsets.UTF_8)});
    }
    for (String[] h : HAND_BYTES) {
      byte[] b = new byte[h[1].length() / 2];
      for (int i = 0; i < b.length; i++) {
        b[i] = (byte) Integer.parseInt(h[1].substring(2 * i, 2 * i + 2), 16);
      }
      cases.add(new Object[] {"bytes_" + h[0], b});
    }
    Random r = new Random(seed);
    for (int i = 0; i < count; i++) {
      Gen g = new Gen(r);
      byte[] doc = g.document();
      String kind = "ok";
      if (r.nextInt(100) < 35 && g.encName.equals("UTF-8")) {
        int times = 1 + (r.nextInt(4) == 0 ? 1 : 0);
        for (int t = 0; t < times; t++) {
          doc = mutate(r, doc);
        }
        kind = "mut";
      } else if (r.nextInt(100) < 35 && g.encName.startsWith("UTF-16") && doc.length > 2
          && (doc[0] == (byte) 0xff || doc[0] == (byte) 0xfe)) {
        doc = mutate16(r, doc, g.encName.equals("UTF-16LE"));
        kind = "mut16";
      }
      cases.add(new Object[] {"rand" + i + "_" + kind, doc});
    }
    for (Object[] c : cases) {
      byte[] doc = (byte[]) c[1];
      out.append("== ").append(c[0]).append('\n');
      out.append("doc ").append(escBytes(doc)).append('\n');
      out.append(runFull(doc, errBuf));
    }
    Files.writeString(dir.resolve("sax_vectors.txt"), out.toString(), StandardCharsets.UTF_8);
    System.setErr(realErr);
  }
}
