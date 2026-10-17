// Helper of ClosureModulesPassDump (same oracle): ClosureRequireProcessor is package-private in
// com.google.javascript.jscomp.modules, so its results are rendered here.
package com.google.javascript.jscomp.modules;

import com.google.javascript.rhino.Node;
import java.util.ArrayList;
import java.util.List;

public final class ClosureModulesPassDumpRequires {
  private ClosureModulesPassDumpRequires() {}

  private static String nodeRef(Node n) {
    if (n == null) {
      return "null";
    }
    return n.getToken() + "@" + n.getLineno() + ":" + n.getCharno();
  }

  /** One line per ClosureRequireProcessor.getAllRequires(nameDeclaration) entry. */
  public static List<String> requires(Node nameDeclaration) {
    List<String> out = new ArrayList<>();
    for (ClosureRequireProcessor.Require r :
        ClosureRequireProcessor.getAllRequires(nameDeclaration)) {
      out.add(
          "require "
              + r.localName()
              + " "
              + r.importRecord().moduleRequest()
              + " "
              + r.importRecord().importName()
              + " "
              + r.importRecord().localName()
              + " "
              + nodeRef(r.importRecord().importNode())
              + " "
              + nodeRef(r.importRecord().nameNode())
              + " "
              + r.createdBy());
    }
    return out;
  }
}
