/*
 * Copyright 2007 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
// Port of com/google/javascript/jscomp/LoggerErrorManager.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.LoggerErrorManager
type LoggerErrorManager struct {
	BasicErrorManagerBase
	formatter MessageFormatter `gson:"formatter"`
	logger *j2gort.Ext_java_util_logging_Logger `gson:"logger"`
}

func (this_ *LoggerErrorManager) BindSelf_(s any) {
	_ = s
	this_.BasicErrorManagerBase.BindSelf_(s)
}
func (this_ *LoggerErrorManager) IsJavaNil_() bool { return this_ == nil }
func (*LoggerErrorManager) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.LoggerErrorManager", "LoggerErrorManager" }
func (this_ *LoggerErrorManager) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.BasicErrorManagerBase.IdentityHash_()
}
func (*LoggerErrorManager) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.LoggerErrorManager", "com.google.javascript.jscomp.LoggerErrorManager", func(v_ any) bool {
		return j2gort.InstanceOf[*LoggerErrorManager](v_)
	})
}
func (*LoggerErrorManager) Implements_jscomp_BasicErrorManager() {}
func (*LoggerErrorManager) Implements_jscomp_SortingErrorManager() {}
func (*LoggerErrorManager) Implements_jscomp_ErrorManager() {}
func (*LoggerErrorManager) Implements_jscomp_ErrorHandler() {}

func (this_ *LoggerErrorManager) initFields_LoggerErrorManager() {
	_ = this_
}

// Java constructor: LoggerErrorManager(com.google.javascript.jscomp.MessageFormatter formatter, java.util.logging.Logger logger)
func NewLoggerErrorManager_2(formatter MessageFormatter, logger *j2gort.Ext_java_util_logging_Logger) *LoggerErrorManager {
	this_ := &LoggerErrorManager{}
	this_.BindSelf_(this_)
	InitLoggerErrorManager_2(this_, formatter, logger)
	return this_
}

func InitLoggerErrorManager_2(this_ *LoggerErrorManager, formatter MessageFormatter, logger *j2gort.Ext_java_util_logging_Logger) {
	_ = formatter
	_ = logger
	InitBasicErrorManager(&this_.BasicErrorManagerBase)
	this_.initFields_LoggerErrorManager()
	this_.formatter = formatter
	this_.logger = logger
}

// Java constructor: LoggerErrorManager(java.util.logging.Logger logger)
func NewLoggerErrorManager_1(logger *j2gort.Ext_java_util_logging_Logger) *LoggerErrorManager {
	this_ := &LoggerErrorManager{}
	this_.BindSelf_(this_)
	InitLoggerErrorManager_1(this_, logger)
	return this_
}

func InitLoggerErrorManager_1(this_ *LoggerErrorManager, logger *j2gort.Ext_java_util_logging_Logger) {
	_ = logger
	InitLoggerErrorManager_2(this_, ToIface_MessageFormatter(LightweightMessageFormatter_WithoutSource()), logger)
}

// Java method: LoggerErrorManager.println(com.google.javascript.jscomp.CheckLevel level, com.google.javascript.jscomp.JSError error)
func (this_ *LoggerErrorManager) Println(level *CheckLevel, error_ *JSError) {
	_ = level
	_ = error_
	switch j2gort.CheckNotNullPtr(level).Ordinal() {
	case int32(0) /* CheckLevel.ERROR */:
		this_.logger.Severe_1_string(error_.Format(level, this_.formatter))
	case int32(1) /* CheckLevel.WARNING */:
		this_.logger.Warning_1_string(error_.Format(level, this_.formatter))
	case int32(2) /* CheckLevel.OFF */:
	}
}

// Java method: LoggerErrorManager.printSummary()
func (this_ *LoggerErrorManager) PrintSummary() {
	var level *j2gort.Ext_java_util_logging_Level
	_ = level
	if (((this_.GetErrorCount() + this_.GetWarningCount()) == int32(0))) {
		level = j2gort.Ext_java_util_logging_Level_F_INFO
	} else {
		level = j2gort.Ext_java_util_logging_Level_F_WARNING
	}
	if ((this_.GetTypedPercent() > float64(0.0))) {
		this_.logger.Log_2_Ext_java_util_logging_Level_string(level, j2gort.Ext_java_lang_String_Format_2("%d error(s), %d warning(s), %.1f%% typed", []any{j2gort.ValueOfInt(this_.GetErrorCount()), j2gort.ValueOfInt(this_.GetWarningCount()), j2gort.ValueOfDouble(this_.GetTypedPercent())}))
	} else if (((this_.GetErrorCount() + this_.GetWarningCount()) > int32(0))) {
		this_.logger.Log_2_Ext_java_util_logging_Level_string(level, j2gort.Ext_java_lang_String_Format_2("%d error(s), %d warning(s)", []any{j2gort.ValueOfInt(this_.GetErrorCount()), j2gort.ValueOfInt(this_.GetWarningCount())}))
	}
}

