// Java reference for jsgo's hidden --max_optimization_loop_iterations=N flag: CommandLineRunner
// with CompilerOptions.setMaxOptimizationLoopIterations(N), an API option that Java's CLI does
// not expose. It takes the same arguments as `java -jar closure-compiler.jar`, plus that flag, so
// one argv drives both compilers and their outputs can be compared byte for byte
// (tools/bench/maxloop_oracle.py).
//
//   javac -cp closure-compiler.jar -d build MaxLoopRunner.java
//   java -cp closure-compiler.jar:build MaxLoopRunner --max_optimization_loop_iterations=2 <args>
import com.google.javascript.jscomp.CommandLineRunner;
import com.google.javascript.jscomp.CompilerOptions;
import java.util.ArrayList;
import java.util.List;

public final class MaxLoopRunner extends CommandLineRunner {
  private static final String FLAG = "--max_optimization_loop_iterations=";

  private final int maxIterations;

  private MaxLoopRunner(String[] args, int maxIterations) {
    super(args);
    this.maxIterations = maxIterations;
  }

  @Override
  protected CompilerOptions createOptions() {
    CompilerOptions options = super.createOptions();
    if (maxIterations > 0) {
      options.setMaxOptimizationLoopIterations(maxIterations);
    }
    return options;
  }

  public static void main(String[] args) {
    int maxIterations = 0;
    List<String> rest = new ArrayList<>();
    for (String arg : args) {
      if (arg.startsWith(FLAG)) {
        maxIterations = Integer.parseInt(arg.substring(FLAG.length()));
      } else {
        rest.add(arg);
      }
    }
    MaxLoopRunner runner = new MaxLoopRunner(rest.toArray(new String[0]), maxIterations);
    if (runner.shouldRunCompiler()) {
      runner.run(); // exits the JVM with CommandLineRunner's status
    }
    System.exit(runner.hasErrors() ? -1 : 0);
  }
}
