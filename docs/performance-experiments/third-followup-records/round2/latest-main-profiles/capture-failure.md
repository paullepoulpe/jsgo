The first capture attempt exited 0 and matched every Java digest on typed-app, but produced an empty CPU file and no allocs profile. No valid performance profile was captured, and the remaining four tasks were not run.

Cause: AbstractCommandLineRunnerBase.Run invokes exitCodeReceiver for status 0 as well as failures. Its default receiver calls System_Exit, so deferred main cleanup cannot run. This is independent of compiler performance.

The replacement capture in ../latest-main-profiles-v2 uses only a cmd/jsgo/main.go overlay. It reconstructs the original CommandLineRunner_Main startup and wraps the public exit-code receiver: successful profile finalization runs immediately before the original exported normalized exit routine. The original Run exception handling, stock compiler construction, logger level, and main throwable handling remain intact. No library source changes were made. Failed first-capture artifacts remain untouched.
