/*
 * Copyright 2011 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/regex/CaseCanonicalize.java (generated mechanically by j2go).

package regex

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.regex.CaseCanonicalize
type CaseCanonicalize struct {
	idHash_ atomic.Int32
}

func (this_ *CaseCanonicalize) BindSelf_(s any) {
	_ = s
}
func (this_ *CaseCanonicalize) IsJavaNil_() bool { return this_ == nil }
func (*CaseCanonicalize) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.regex.CaseCanonicalize", "CaseCanonicalize"
}
func (this_ *CaseCanonicalize) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CaseCanonicalize) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.regex.CaseCanonicalize", "com.google.javascript.jscomp.regex.CaseCanonicalize", func(v_ any) bool {
		return j2gort.InstanceOf[*CaseCanonicalize](v_)
	})
}

var static_CaseCanonicalize_CASE_SENSITIVE *CharRanges

func StaticGet_CaseCanonicalize_CASE_SENSITIVE() *CharRanges {
	clinit_CaseCanonicalize()
	return static_CaseCanonicalize_CASE_SENSITIVE
}
func StaticSet_CaseCanonicalize_CASE_SENSITIVE(v *CharRanges) {
	clinit_CaseCanonicalize()
	static_CaseCanonicalize_CASE_SENSITIVE = v
}

var static_CaseCanonicalize_UCASE_ASCII_LETTERS *CharRanges

func StaticGet_CaseCanonicalize_UCASE_ASCII_LETTERS() *CharRanges {
	clinit_CaseCanonicalize()
	return static_CaseCanonicalize_UCASE_ASCII_LETTERS
}
func StaticSet_CaseCanonicalize_UCASE_ASCII_LETTERS(v *CharRanges) {
	clinit_CaseCanonicalize()
	static_CaseCanonicalize_UCASE_ASCII_LETTERS = v
}

var static_CaseCanonicalize_DELTA_SETS j2gort.Ext_com_google_common_collect_ImmutableList

func StaticGet_CaseCanonicalize_DELTA_SETS() j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_CaseCanonicalize()
	return static_CaseCanonicalize_DELTA_SETS
}
func StaticSet_CaseCanonicalize_DELTA_SETS(v j2gort.Ext_com_google_common_collect_ImmutableList) {
	clinit_CaseCanonicalize()
	static_CaseCanonicalize_DELTA_SETS = v
}

var static_CaseCanonicalize_CANON_DELTA_SETS j2gort.Ext_com_google_common_collect_ImmutableList

func StaticGet_CaseCanonicalize_CANON_DELTA_SETS() j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_CaseCanonicalize()
	return static_CaseCanonicalize_CANON_DELTA_SETS
}
func StaticSet_CaseCanonicalize_CANON_DELTA_SETS(v j2gort.Ext_com_google_common_collect_ImmutableList) {
	clinit_CaseCanonicalize()
	static_CaseCanonicalize_CANON_DELTA_SETS = v
}

var clinitDone_CaseCanonicalize uint32
var clinitState_CaseCanonicalize j2gort.ClinitState

func clinit_CaseCanonicalize() {
	if atomic.LoadUint32(&clinitDone_CaseCanonicalize) != 0 {
		return
	}
	clinitSlow_CaseCanonicalize()
}

func clinitSlow_CaseCanonicalize() {
	if !j2gort.EnterClinit(&clinitState_CaseCanonicalize, &clinitDone_CaseCanonicalize) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CaseCanonicalize, &clinitDone_CaseCanonicalize)
	static_CaseCanonicalize_CASE_SENSITIVE = CharRanges_WithRanges([]int32{int32(65), int32(91), int32(97), int32(123), int32(181), int32(182), int32(192), int32(215), int32(216), int32(223), int32(224), int32(247), int32(248), int32(304), int32(306), int32(312), int32(313), int32(329), int32(330), int32(383), int32(384), int32(397), int32(398), int32(411), int32(412), int32(426), int32(428), int32(442), int32(444), int32(446), int32(447), int32(448), int32(452), int32(496), int32(497), int32(545), int32(546), int32(564), int32(570), int32(575), int32(577), int32(592), int32(595), int32(597), int32(598), int32(600), int32(601), int32(602), int32(603), int32(604), int32(608), int32(609), int32(611), int32(612), int32(616), int32(618), int32(619), int32(620), int32(623), int32(624), int32(626), int32(627), int32(629), int32(630), int32(637), int32(638), int32(640), int32(641), int32(643), int32(644), int32(648), int32(653), int32(658), int32(659), int32(837), int32(838), int32(891), int32(894), int32(902), int32(903), int32(904), int32(907), int32(908), int32(909), int32(910), int32(912), int32(913), int32(930), int32(931), int32(944), int32(945), int32(975), int32(976), int32(978), int32(981), int32(983), int32(984), int32(1011), int32(1013), int32(1014), int32(1015), int32(1020), int32(1021), int32(1154), int32(1162), int32(1300), int32(1329), int32(1367), int32(1377), int32(1415), int32(4256), int32(4294), int32(7549), int32(7550), int32(7680), int32(7830), int32(7835), int32(7836), int32(7840), int32(7930), int32(7936), int32(7958), int32(7960), int32(7966), int32(7968), int32(8006), int32(8008), int32(8014), int32(8017), int32(8018), int32(8019), int32(8020), int32(8021), int32(8022), int32(8023), int32(8024), int32(8025), int32(8026), int32(8027), int32(8028), int32(8029), int32(8030), int32(8031), int32(8062), int32(8112), int32(8114), int32(8120), int32(8124), int32(8126), int32(8127), int32(8136), int32(8140), int32(8144), int32(8146), int32(8152), int32(8156), int32(8160), int32(8162), int32(8165), int32(8166), int32(8168), int32(8173), int32(8184), int32(8188), int32(8498), int32(8499), int32(8526), int32(8527), int32(8544), int32(8576), int32(8579), int32(8581), int32(9398), int32(9450), int32(11264), int32(11311), int32(11312), int32(11359), int32(11360), int32(11373), int32(11381), int32(11383), int32(11392), int32(11492), int32(11520), int32(11558), int32(65313), int32(65339), int32(65345), int32(65371)})
	static_CaseCanonicalize_UCASE_ASCII_LETTERS = CharRanges_Inclusive(int32(uint16(65)), int32(uint16(90)))
	static_CaseCanonicalize_DELTA_SETS = j2gort.Ext_com_google_common_collect_ImmutableList_Of_13(NewCaseCanonicalize_DeltaSet(int32(-10795), CharRanges_WithMembers([]int32{int32(570)})), NewCaseCanonicalize_DeltaSet(int32(-10792), CharRanges_WithMembers([]int32{int32(574)})), NewCaseCanonicalize_DeltaSet(int32(-10743), CharRanges_WithMembers([]int32{int32(619)})), NewCaseCanonicalize_DeltaSet(int32(-10727), CharRanges_WithMembers([]int32{int32(637)})), NewCaseCanonicalize_DeltaSet(int32(-7264), CharRanges_WithRanges([]int32{int32(4256), int32(4294)})), NewCaseCanonicalize_DeltaSet(int32(-7205), CharRanges_WithMembers([]int32{int32(921)})), NewCaseCanonicalize_DeltaSet(int32(-3814), CharRanges_WithMembers([]int32{int32(7549)})), NewCaseCanonicalize_DeltaSet(int32(-743), CharRanges_WithMembers([]int32{int32(181)})), NewCaseCanonicalize_DeltaSet(int32(-219), CharRanges_WithMembers([]int32{int32(439)})), NewCaseCanonicalize_DeltaSet(int32(-218), CharRanges_WithMembers([]int32{int32(422), int32(425), int32(430)})), NewCaseCanonicalize_DeltaSet(int32(-217), CharRanges_WithRanges([]int32{int32(433), int32(435)})), NewCaseCanonicalize_DeltaSet(int32(-214), CharRanges_WithMembers([]int32{int32(415)})), []any{NewCaseCanonicalize_DeltaSet(int32(-213), CharRanges_WithMembers([]int32{int32(413)})), NewCaseCanonicalize_DeltaSet(int32(-211), CharRanges_WithMembers([]int32{int32(406), int32(412)})), NewCaseCanonicalize_DeltaSet(int32(-210), CharRanges_WithMembers([]int32{int32(385)})), NewCaseCanonicalize_DeltaSet(int32(-209), CharRanges_WithMembers([]int32{int32(407)})), NewCaseCanonicalize_DeltaSet(int32(-207), CharRanges_WithMembers([]int32{int32(404)})), NewCaseCanonicalize_DeltaSet(int32(-206), CharRanges_WithMembers([]int32{int32(390)})), NewCaseCanonicalize_DeltaSet(int32(-205), CharRanges_WithRanges([]int32{int32(393), int32(395), int32(403), int32(404)})), NewCaseCanonicalize_DeltaSet(int32(-203), CharRanges_WithMembers([]int32{int32(400)})), NewCaseCanonicalize_DeltaSet(int32(-202), CharRanges_WithMembers([]int32{int32(399)})), NewCaseCanonicalize_DeltaSet(int32(-195), CharRanges_WithMembers([]int32{int32(384)})), NewCaseCanonicalize_DeltaSet(int32(-163), CharRanges_WithMembers([]int32{int32(410)})), NewCaseCanonicalize_DeltaSet(int32(-130), CharRanges_WithRanges([]int32{int32(414), int32(415), int32(891), int32(894)})), NewCaseCanonicalize_DeltaSet(int32(-128), CharRanges_WithRanges([]int32{int32(8056), int32(8058)})), NewCaseCanonicalize_DeltaSet(int32(-126), CharRanges_WithRanges([]int32{int32(8060), int32(8062)})), NewCaseCanonicalize_DeltaSet(int32(-121), CharRanges_WithMembers([]int32{int32(255)})), NewCaseCanonicalize_DeltaSet(int32(-112), CharRanges_WithRanges([]int32{int32(8058), int32(8060)})), NewCaseCanonicalize_DeltaSet(int32(-100), CharRanges_WithRanges([]int32{int32(8054), int32(8056)})), NewCaseCanonicalize_DeltaSet(int32(-97), CharRanges_WithMembers([]int32{int32(405)})), NewCaseCanonicalize_DeltaSet(int32(-96), CharRanges_WithMembers([]int32{int32(917)})), NewCaseCanonicalize_DeltaSet(int32(-86), CharRanges_WithRanges([]int32{int32(922), int32(923), int32(8050), int32(8054)})), NewCaseCanonicalize_DeltaSet(int32(-84), CharRanges_WithMembers([]int32{int32(837)})), NewCaseCanonicalize_DeltaSet(int32(-80), CharRanges_WithRanges([]int32{int32(929), int32(930), int32(1024), int32(1040)})), NewCaseCanonicalize_DeltaSet(int32(-79), CharRanges_WithMembers([]int32{int32(398)})), NewCaseCanonicalize_DeltaSet(int32(-74), CharRanges_WithRanges([]int32{int32(8048), int32(8050)})), NewCaseCanonicalize_DeltaSet(int32(-71), CharRanges_WithMembers([]int32{int32(581)})), NewCaseCanonicalize_DeltaSet(int32(-69), CharRanges_WithMembers([]int32{int32(580)})), NewCaseCanonicalize_DeltaSet(int32(-64), CharRanges_WithMembers([]int32{int32(908)})), NewCaseCanonicalize_DeltaSet(int32(-63), CharRanges_WithRanges([]int32{int32(910), int32(912)})), NewCaseCanonicalize_DeltaSet(int32(-62), CharRanges_WithMembers([]int32{int32(914)})), NewCaseCanonicalize_DeltaSet(int32(-59), CharRanges_WithMembers([]int32{int32(7776)})), NewCaseCanonicalize_DeltaSet(int32(-57), CharRanges_WithMembers([]int32{int32(920)})), NewCaseCanonicalize_DeltaSet(int32(-56), CharRanges_WithMembers([]int32{int32(447)})), NewCaseCanonicalize_DeltaSet(int32(-54), CharRanges_WithMembers([]int32{int32(928)})), NewCaseCanonicalize_DeltaSet(int32(-48), CharRanges_WithRanges([]int32{int32(1329), int32(1367), int32(11264), int32(11311)})), NewCaseCanonicalize_DeltaSet(int32(-47), CharRanges_WithMembers([]int32{int32(934)})), NewCaseCanonicalize_DeltaSet(int32(-38), CharRanges_WithMembers([]int32{int32(902)})), NewCaseCanonicalize_DeltaSet(int32(-37), CharRanges_WithRanges([]int32{int32(904), int32(907)})), NewCaseCanonicalize_DeltaSet(int32(-32), CharRanges_WithRanges([]int32{int32(65), int32(91), int32(192), int32(215), int32(216), int32(223), int32(913), int32(930), int32(931), int32(940), int32(1040), int32(1072), int32(65313), int32(65339)})), NewCaseCanonicalize_DeltaSet(int32(-31), CharRanges_WithMembers([]int32{int32(931)})), NewCaseCanonicalize_DeltaSet(int32(-28), CharRanges_WithMembers([]int32{int32(8498)})), NewCaseCanonicalize_DeltaSet(int32(-26), CharRanges_WithRanges([]int32{int32(9398), int32(9424)})), NewCaseCanonicalize_DeltaSet(int32(-16), CharRanges_WithRanges([]int32{int32(8544), int32(8560)})), NewCaseCanonicalize_DeltaSet(int32(-15), CharRanges_WithMembers([]int32{int32(1216)})), NewCaseCanonicalize_DeltaSet(int32(-8), CharRanges_WithRanges([]int32{int32(7936), int32(7944), int32(7952), int32(7958), int32(7968), int32(7976), int32(7984), int32(7992), int32(8000), int32(8006), int32(8017), int32(8018), int32(8019), int32(8020), int32(8021), int32(8022), int32(8023), int32(8024), int32(8032), int32(8040), int32(8112), int32(8114), int32(8144), int32(8146), int32(8160), int32(8162)})), NewCaseCanonicalize_DeltaSet(int32(-7), CharRanges_WithMembers([]int32{int32(1010), int32(8165)})), NewCaseCanonicalize_DeltaSet(int32(-2), CharRanges_WithMembers([]int32{int32(452), int32(455), int32(458), int32(497)})), NewCaseCanonicalize_DeltaSet(int32(-1), CharRanges_WithMembers([]int32{int32(256), int32(258), int32(260), int32(262), int32(264), int32(266), int32(268), int32(270), int32(272), int32(274), int32(276), int32(278), int32(280), int32(282), int32(284), int32(286), int32(288), int32(290), int32(292), int32(294), int32(296), int32(298), int32(300), int32(302), int32(306), int32(308), int32(310), int32(313), int32(315), int32(317), int32(319), int32(321), int32(323), int32(325), int32(327), int32(330), int32(332), int32(334), int32(336), int32(338), int32(340), int32(342), int32(344), int32(346), int32(348), int32(350), int32(352), int32(354), int32(356), int32(358), int32(360), int32(362), int32(364), int32(366), int32(368), int32(370), int32(372), int32(374), int32(377), int32(379), int32(381), int32(386), int32(388), int32(391), int32(395), int32(401), int32(408), int32(416), int32(418), int32(420), int32(423), int32(428), int32(431), int32(435), int32(437), int32(440), int32(444), int32(452), int32(455), int32(458), int32(461), int32(463), int32(465), int32(467), int32(469), int32(471), int32(473), int32(475), int32(478), int32(480), int32(482), int32(484), int32(486), int32(488), int32(490), int32(492), int32(494), int32(497), int32(500), int32(504), int32(506), int32(508), int32(510), int32(512), int32(514), int32(516), int32(518), int32(520), int32(522), int32(524), int32(526), int32(528), int32(530), int32(532), int32(534), int32(536), int32(538), int32(540), int32(542), int32(546), int32(548), int32(550), int32(552), int32(554), int32(556), int32(558), int32(560), int32(562), int32(571), int32(577), int32(582), int32(584), int32(586), int32(588), int32(590), int32(984), int32(986), int32(988), int32(990), int32(992), int32(994), int32(996), int32(998), int32(1000), int32(1002), int32(1004), int32(1006), int32(1015), int32(1018), int32(1120), int32(1122), int32(1124), int32(1126), int32(1128), int32(1130), int32(1132), int32(1134), int32(1136), int32(1138), int32(1140), int32(1142), int32(1144), int32(1146), int32(1148), int32(1150), int32(1152), int32(1162), int32(1164), int32(1166), int32(1168), int32(1170), int32(1172), int32(1174), int32(1176), int32(1178), int32(1180), int32(1182), int32(1184), int32(1186), int32(1188), int32(1190), int32(1192), int32(1194), int32(1196), int32(1198), int32(1200), int32(1202), int32(1204), int32(1206), int32(1208), int32(1210), int32(1212), int32(1214), int32(1217), int32(1219), int32(1221), int32(1223), int32(1225), int32(1227), int32(1229), int32(1232), int32(1234), int32(1236), int32(1238), int32(1240), int32(1242), int32(1244), int32(1246), int32(1248), int32(1250), int32(1252), int32(1254), int32(1256), int32(1258), int32(1260), int32(1262), int32(1264), int32(1266), int32(1268), int32(1270), int32(1272), int32(1274), int32(1276), int32(1278), int32(1280), int32(1282), int32(1284), int32(1286), int32(1288), int32(1290), int32(1292), int32(1294), int32(1296), int32(1298), int32(7680), int32(7682), int32(7684), int32(7686), int32(7688), int32(7690), int32(7692), int32(7694), int32(7696), int32(7698), int32(7700), int32(7702), int32(7704), int32(7706), int32(7708), int32(7710), int32(7712), int32(7714), int32(7716), int32(7718), int32(7720), int32(7722), int32(7724), int32(7726), int32(7728), int32(7730), int32(7732), int32(7734), int32(7736), int32(7738), int32(7740), int32(7742), int32(7744), int32(7746), int32(7748), int32(7750), int32(7752), int32(7754), int32(7756), int32(7758), int32(7760), int32(7762), int32(7764), int32(7766), int32(7768), int32(7770), int32(7772), int32(7774), int32(7776), int32(7778), int32(7780), int32(7782), int32(7784), int32(7786), int32(7788), int32(7790), int32(7792), int32(7794), int32(7796), int32(7798), int32(7800), int32(7802), int32(7804), int32(7806), int32(7808), int32(7810), int32(7812), int32(7814), int32(7816), int32(7818), int32(7820), int32(7822), int32(7824), int32(7826), int32(7828), int32(7840), int32(7842), int32(7844), int32(7846), int32(7848), int32(7850), int32(7852), int32(7854), int32(7856), int32(7858), int32(7860), int32(7862), int32(7864), int32(7866), int32(7868), int32(7870), int32(7872), int32(7874), int32(7876), int32(7878), int32(7880), int32(7882), int32(7884), int32(7886), int32(7888), int32(7890), int32(7892), int32(7894), int32(7896), int32(7898), int32(7900), int32(7902), int32(7904), int32(7906), int32(7908), int32(7910), int32(7912), int32(7914), int32(7916), int32(7918), int32(7920), int32(7922), int32(7924), int32(7926), int32(7928), int32(8579), int32(11360), int32(11367), int32(11369), int32(11371), int32(11381), int32(11392), int32(11394), int32(11396), int32(11398), int32(11400), int32(11402), int32(11404), int32(11406), int32(11408), int32(11410), int32(11412), int32(11414), int32(11416), int32(11418), int32(11420), int32(11422), int32(11424), int32(11426), int32(11428), int32(11430), int32(11432), int32(11434), int32(11436), int32(11438), int32(11440), int32(11442), int32(11444), int32(11446), int32(11448), int32(11450), int32(11452), int32(11454), int32(11456), int32(11458), int32(11460), int32(11462), int32(11464), int32(11466), int32(11468), int32(11470), int32(11472), int32(11474), int32(11476), int32(11478), int32(11480), int32(11482), int32(11484), int32(11486), int32(11488), int32(11490)})), NewCaseCanonicalize_DeltaSet(int32(1), CharRanges_WithMembers([]int32{int32(257), int32(259), int32(261), int32(263), int32(265), int32(267), int32(269), int32(271), int32(273), int32(275), int32(277), int32(279), int32(281), int32(283), int32(285), int32(287), int32(289), int32(291), int32(293), int32(295), int32(297), int32(299), int32(301), int32(303), int32(307), int32(309), int32(311), int32(314), int32(316), int32(318), int32(320), int32(322), int32(324), int32(326), int32(328), int32(331), int32(333), int32(335), int32(337), int32(339), int32(341), int32(343), int32(345), int32(347), int32(349), int32(351), int32(353), int32(355), int32(357), int32(359), int32(361), int32(363), int32(365), int32(367), int32(369), int32(371), int32(373), int32(375), int32(378), int32(380), int32(382), int32(387), int32(389), int32(392), int32(396), int32(402), int32(409), int32(417), int32(419), int32(421), int32(424), int32(429), int32(432), int32(436), int32(438), int32(441), int32(445), int32(453), int32(456), int32(459), int32(462), int32(464), int32(466), int32(468), int32(470), int32(472), int32(474), int32(476), int32(479), int32(481), int32(483), int32(485), int32(487), int32(489), int32(491), int32(493), int32(495), int32(498), int32(501), int32(505), int32(507), int32(509), int32(511), int32(513), int32(515), int32(517), int32(519), int32(521), int32(523), int32(525), int32(527), int32(529), int32(531), int32(533), int32(535), int32(537), int32(539), int32(541), int32(543), int32(547), int32(549), int32(551), int32(553), int32(555), int32(557), int32(559), int32(561), int32(563), int32(572), int32(578), int32(583), int32(585), int32(587), int32(589), int32(591), int32(985), int32(987), int32(989), int32(991), int32(993), int32(995), int32(997), int32(999), int32(1001), int32(1003), int32(1005), int32(1007), int32(1016), int32(1019), int32(1121), int32(1123), int32(1125), int32(1127), int32(1129), int32(1131), int32(1133), int32(1135), int32(1137), int32(1139), int32(1141), int32(1143), int32(1145), int32(1147), int32(1149), int32(1151), int32(1153), int32(1163), int32(1165), int32(1167), int32(1169), int32(1171), int32(1173), int32(1175), int32(1177), int32(1179), int32(1181), int32(1183), int32(1185), int32(1187), int32(1189), int32(1191), int32(1193), int32(1195), int32(1197), int32(1199), int32(1201), int32(1203), int32(1205), int32(1207), int32(1209), int32(1211), int32(1213), int32(1215), int32(1218), int32(1220), int32(1222), int32(1224), int32(1226), int32(1228), int32(1230), int32(1233), int32(1235), int32(1237), int32(1239), int32(1241), int32(1243), int32(1245), int32(1247), int32(1249), int32(1251), int32(1253), int32(1255), int32(1257), int32(1259), int32(1261), int32(1263), int32(1265), int32(1267), int32(1269), int32(1271), int32(1273), int32(1275), int32(1277), int32(1279), int32(1281), int32(1283), int32(1285), int32(1287), int32(1289), int32(1291), int32(1293), int32(1295), int32(1297), int32(1299), int32(7681), int32(7683), int32(7685), int32(7687), int32(7689), int32(7691), int32(7693), int32(7695), int32(7697), int32(7699), int32(7701), int32(7703), int32(7705), int32(7707), int32(7709), int32(7711), int32(7713), int32(7715), int32(7717), int32(7719), int32(7721), int32(7723), int32(7725), int32(7727), int32(7729), int32(7731), int32(7733), int32(7735), int32(7737), int32(7739), int32(7741), int32(7743), int32(7745), int32(7747), int32(7749), int32(7751), int32(7753), int32(7755), int32(7757), int32(7759), int32(7761), int32(7763), int32(7765), int32(7767), int32(7769), int32(7771), int32(7773), int32(7775), int32(7777), int32(7779), int32(7781), int32(7783), int32(7785), int32(7787), int32(7789), int32(7791), int32(7793), int32(7795), int32(7797), int32(7799), int32(7801), int32(7803), int32(7805), int32(7807), int32(7809), int32(7811), int32(7813), int32(7815), int32(7817), int32(7819), int32(7821), int32(7823), int32(7825), int32(7827), int32(7829), int32(7841), int32(7843), int32(7845), int32(7847), int32(7849), int32(7851), int32(7853), int32(7855), int32(7857), int32(7859), int32(7861), int32(7863), int32(7865), int32(7867), int32(7869), int32(7871), int32(7873), int32(7875), int32(7877), int32(7879), int32(7881), int32(7883), int32(7885), int32(7887), int32(7889), int32(7891), int32(7893), int32(7895), int32(7897), int32(7899), int32(7901), int32(7903), int32(7905), int32(7907), int32(7909), int32(7911), int32(7913), int32(7915), int32(7917), int32(7919), int32(7921), int32(7923), int32(7925), int32(7927), int32(7929), int32(8580), int32(11361), int32(11368), int32(11370), int32(11372), int32(11382), int32(11393), int32(11395), int32(11397), int32(11399), int32(11401), int32(11403), int32(11405), int32(11407), int32(11409), int32(11411), int32(11413), int32(11415), int32(11417), int32(11419), int32(11421), int32(11423), int32(11425), int32(11427), int32(11429), int32(11431), int32(11433), int32(11435), int32(11437), int32(11439), int32(11441), int32(11443), int32(11445), int32(11447), int32(11449), int32(11451), int32(11453), int32(11455), int32(11457), int32(11459), int32(11461), int32(11463), int32(11465), int32(11467), int32(11469), int32(11471), int32(11473), int32(11475), int32(11477), int32(11479), int32(11481), int32(11483), int32(11485), int32(11487), int32(11489), int32(11491)})), NewCaseCanonicalize_DeltaSet(int32(2), CharRanges_WithMembers([]int32{int32(454), int32(457), int32(460), int32(499)})), NewCaseCanonicalize_DeltaSet(int32(7), CharRanges_WithMembers([]int32{int32(1017), int32(8172)})), NewCaseCanonicalize_DeltaSet(int32(8), CharRanges_WithRanges([]int32{int32(7944), int32(7952), int32(7960), int32(7966), int32(7976), int32(7984), int32(7992), int32(8000), int32(8008), int32(8014), int32(8025), int32(8026), int32(8027), int32(8028), int32(8029), int32(8030), int32(8031), int32(8032), int32(8040), int32(8048), int32(8120), int32(8122), int32(8152), int32(8154), int32(8168), int32(8170)})), NewCaseCanonicalize_DeltaSet(int32(15), CharRanges_WithMembers([]int32{int32(1231)})), NewCaseCanonicalize_DeltaSet(int32(16), CharRanges_WithRanges([]int32{int32(8560), int32(8576)})), NewCaseCanonicalize_DeltaSet(int32(26), CharRanges_WithRanges([]int32{int32(9424), int32(9450)})), NewCaseCanonicalize_DeltaSet(int32(28), CharRanges_WithMembers([]int32{int32(8526)})), NewCaseCanonicalize_DeltaSet(int32(31), CharRanges_WithMembers([]int32{int32(962)})), NewCaseCanonicalize_DeltaSet(int32(32), CharRanges_WithRanges([]int32{int32(97), int32(123), int32(224), int32(247), int32(248), int32(255), int32(945), int32(962), int32(963), int32(972), int32(1072), int32(1104), int32(65345), int32(65371)})), NewCaseCanonicalize_DeltaSet(int32(37), CharRanges_WithRanges([]int32{int32(941), int32(944)})), NewCaseCanonicalize_DeltaSet(int32(38), CharRanges_WithMembers([]int32{int32(940)})), NewCaseCanonicalize_DeltaSet(int32(47), CharRanges_WithMembers([]int32{int32(981)})), NewCaseCanonicalize_DeltaSet(int32(48), CharRanges_WithRanges([]int32{int32(1377), int32(1415), int32(11312), int32(11359)})), NewCaseCanonicalize_DeltaSet(int32(54), CharRanges_WithMembers([]int32{int32(982)})), NewCaseCanonicalize_DeltaSet(int32(56), CharRanges_WithMembers([]int32{int32(503)})), NewCaseCanonicalize_DeltaSet(int32(57), CharRanges_WithMembers([]int32{int32(977)})), NewCaseCanonicalize_DeltaSet(int32(59), CharRanges_WithMembers([]int32{int32(7835)})), NewCaseCanonicalize_DeltaSet(int32(62), CharRanges_WithMembers([]int32{int32(976)})), NewCaseCanonicalize_DeltaSet(int32(63), CharRanges_WithRanges([]int32{int32(973), int32(975)})), NewCaseCanonicalize_DeltaSet(int32(64), CharRanges_WithMembers([]int32{int32(972)})), NewCaseCanonicalize_DeltaSet(int32(69), CharRanges_WithMembers([]int32{int32(649)})), NewCaseCanonicalize_DeltaSet(int32(71), CharRanges_WithMembers([]int32{int32(652)})), NewCaseCanonicalize_DeltaSet(int32(74), CharRanges_WithRanges([]int32{int32(8122), int32(8124)})), NewCaseCanonicalize_DeltaSet(int32(79), CharRanges_WithMembers([]int32{int32(477)})), NewCaseCanonicalize_DeltaSet(int32(80), CharRanges_WithRanges([]int32{int32(1009), int32(1010), int32(1104), int32(1120)})), NewCaseCanonicalize_DeltaSet(int32(84), CharRanges_WithMembers([]int32{int32(921)})), NewCaseCanonicalize_DeltaSet(int32(86), CharRanges_WithRanges([]int32{int32(1008), int32(1009), int32(8136), int32(8140)})), NewCaseCanonicalize_DeltaSet(int32(96), CharRanges_WithMembers([]int32{int32(1013)})), NewCaseCanonicalize_DeltaSet(int32(97), CharRanges_WithMembers([]int32{int32(502)})), NewCaseCanonicalize_DeltaSet(int32(100), CharRanges_WithRanges([]int32{int32(8154), int32(8156)})), NewCaseCanonicalize_DeltaSet(int32(112), CharRanges_WithRanges([]int32{int32(8170), int32(8172)})), NewCaseCanonicalize_DeltaSet(int32(121), CharRanges_WithMembers([]int32{int32(376)})), NewCaseCanonicalize_DeltaSet(int32(126), CharRanges_WithRanges([]int32{int32(8186), int32(8188)})), NewCaseCanonicalize_DeltaSet(int32(128), CharRanges_WithRanges([]int32{int32(8184), int32(8186)})), NewCaseCanonicalize_DeltaSet(int32(130), CharRanges_WithRanges([]int32{int32(544), int32(545), int32(1021), int32(1024)})), NewCaseCanonicalize_DeltaSet(int32(163), CharRanges_WithMembers([]int32{int32(573)})), NewCaseCanonicalize_DeltaSet(int32(195), CharRanges_WithMembers([]int32{int32(579)})), NewCaseCanonicalize_DeltaSet(int32(202), CharRanges_WithMembers([]int32{int32(601)})), NewCaseCanonicalize_DeltaSet(int32(203), CharRanges_WithMembers([]int32{int32(603)})), NewCaseCanonicalize_DeltaSet(int32(205), CharRanges_WithRanges([]int32{int32(598), int32(600), int32(608), int32(609)})), NewCaseCanonicalize_DeltaSet(int32(206), CharRanges_WithMembers([]int32{int32(596)})), NewCaseCanonicalize_DeltaSet(int32(207), CharRanges_WithMembers([]int32{int32(611)})), NewCaseCanonicalize_DeltaSet(int32(209), CharRanges_WithMembers([]int32{int32(616)})), NewCaseCanonicalize_DeltaSet(int32(210), CharRanges_WithMembers([]int32{int32(595)})), NewCaseCanonicalize_DeltaSet(int32(211), CharRanges_WithMembers([]int32{int32(617), int32(623)})), NewCaseCanonicalize_DeltaSet(int32(213), CharRanges_WithMembers([]int32{int32(626)})), NewCaseCanonicalize_DeltaSet(int32(214), CharRanges_WithMembers([]int32{int32(629)})), NewCaseCanonicalize_DeltaSet(int32(217), CharRanges_WithRanges([]int32{int32(650), int32(652)})), NewCaseCanonicalize_DeltaSet(int32(218), CharRanges_WithMembers([]int32{int32(640), int32(643), int32(648)})), NewCaseCanonicalize_DeltaSet(int32(219), CharRanges_WithMembers([]int32{int32(658)})), NewCaseCanonicalize_DeltaSet(int32(743), CharRanges_WithMembers([]int32{int32(924)})), NewCaseCanonicalize_DeltaSet(int32(3814), CharRanges_WithMembers([]int32{int32(11363)})), NewCaseCanonicalize_DeltaSet(int32(7205), CharRanges_WithMembers([]int32{int32(8126)})), NewCaseCanonicalize_DeltaSet(int32(7264), CharRanges_WithRanges([]int32{int32(11520), int32(11558)})), NewCaseCanonicalize_DeltaSet(int32(10727), CharRanges_WithMembers([]int32{int32(11364)})), NewCaseCanonicalize_DeltaSet(int32(10743), CharRanges_WithMembers([]int32{int32(11362)})), NewCaseCanonicalize_DeltaSet(int32(10792), CharRanges_WithMembers([]int32{int32(11366)})), NewCaseCanonicalize_DeltaSet(int32(10795), CharRanges_WithMembers([]int32{int32(11365)}))})
	static_CaseCanonicalize_CANON_DELTA_SETS = j2gort.Ext_com_google_common_collect_ImmutableList_Of_13(NewCaseCanonicalize_DeltaSet(int32(-10743), CharRanges_WithMembers([]int32{int32(619)})), NewCaseCanonicalize_DeltaSet(int32(-10727), CharRanges_WithMembers([]int32{int32(637)})), NewCaseCanonicalize_DeltaSet(int32(-3814), CharRanges_WithMembers([]int32{int32(7549)})), NewCaseCanonicalize_DeltaSet(int32(-743), CharRanges_WithMembers([]int32{int32(181)})), NewCaseCanonicalize_DeltaSet(int32(-195), CharRanges_WithMembers([]int32{int32(384)})), NewCaseCanonicalize_DeltaSet(int32(-163), CharRanges_WithMembers([]int32{int32(410)})), NewCaseCanonicalize_DeltaSet(int32(-130), CharRanges_WithRanges([]int32{int32(414), int32(415), int32(891), int32(894)})), NewCaseCanonicalize_DeltaSet(int32(-128), CharRanges_WithRanges([]int32{int32(8056), int32(8058)})), NewCaseCanonicalize_DeltaSet(int32(-126), CharRanges_WithRanges([]int32{int32(8060), int32(8062)})), NewCaseCanonicalize_DeltaSet(int32(-121), CharRanges_WithMembers([]int32{int32(255)})), NewCaseCanonicalize_DeltaSet(int32(-112), CharRanges_WithRanges([]int32{int32(8058), int32(8060)})), NewCaseCanonicalize_DeltaSet(int32(-100), CharRanges_WithRanges([]int32{int32(8054), int32(8056)})), []any{NewCaseCanonicalize_DeltaSet(int32(-97), CharRanges_WithMembers([]int32{int32(405)})), NewCaseCanonicalize_DeltaSet(int32(-86), CharRanges_WithRanges([]int32{int32(8050), int32(8054)})), NewCaseCanonicalize_DeltaSet(int32(-84), CharRanges_WithMembers([]int32{int32(837)})), NewCaseCanonicalize_DeltaSet(int32(-74), CharRanges_WithRanges([]int32{int32(8048), int32(8050)})), NewCaseCanonicalize_DeltaSet(int32(-56), CharRanges_WithMembers([]int32{int32(447)})), NewCaseCanonicalize_DeltaSet(int32(-8), CharRanges_WithRanges([]int32{int32(7936), int32(7944), int32(7952), int32(7958), int32(7968), int32(7976), int32(7984), int32(7992), int32(8000), int32(8006), int32(8017), int32(8018), int32(8019), int32(8020), int32(8021), int32(8022), int32(8023), int32(8024), int32(8032), int32(8040), int32(8112), int32(8114), int32(8144), int32(8146), int32(8160), int32(8162)})), NewCaseCanonicalize_DeltaSet(int32(-7), CharRanges_WithMembers([]int32{int32(1010), int32(8165)})), NewCaseCanonicalize_DeltaSet(int32(1), CharRanges_WithMembers([]int32{int32(257), int32(259), int32(261), int32(263), int32(265), int32(267), int32(269), int32(271), int32(273), int32(275), int32(277), int32(279), int32(281), int32(283), int32(285), int32(287), int32(289), int32(291), int32(293), int32(295), int32(297), int32(299), int32(301), int32(303), int32(307), int32(309), int32(311), int32(314), int32(316), int32(318), int32(320), int32(322), int32(324), int32(326), int32(328), int32(331), int32(333), int32(335), int32(337), int32(339), int32(341), int32(343), int32(345), int32(347), int32(349), int32(351), int32(353), int32(355), int32(357), int32(359), int32(361), int32(363), int32(365), int32(367), int32(369), int32(371), int32(373), int32(375), int32(378), int32(380), int32(382), int32(387), int32(389), int32(392), int32(396), int32(402), int32(409), int32(417), int32(419), int32(421), int32(424), int32(429), int32(432), int32(436), int32(438), int32(441), int32(445), int32(453), int32(456), int32(459), int32(462), int32(464), int32(466), int32(468), int32(470), int32(472), int32(474), int32(476), int32(479), int32(481), int32(483), int32(485), int32(487), int32(489), int32(491), int32(493), int32(495), int32(498), int32(501), int32(505), int32(507), int32(509), int32(511), int32(513), int32(515), int32(517), int32(519), int32(521), int32(523), int32(525), int32(527), int32(529), int32(531), int32(533), int32(535), int32(537), int32(539), int32(541), int32(543), int32(547), int32(549), int32(551), int32(553), int32(555), int32(557), int32(559), int32(561), int32(563), int32(572), int32(578), int32(583), int32(585), int32(587), int32(589), int32(591), int32(985), int32(987), int32(989), int32(991), int32(993), int32(995), int32(997), int32(999), int32(1001), int32(1003), int32(1005), int32(1007), int32(1016), int32(1019), int32(1121), int32(1123), int32(1125), int32(1127), int32(1129), int32(1131), int32(1133), int32(1135), int32(1137), int32(1139), int32(1141), int32(1143), int32(1145), int32(1147), int32(1149), int32(1151), int32(1153), int32(1163), int32(1165), int32(1167), int32(1169), int32(1171), int32(1173), int32(1175), int32(1177), int32(1179), int32(1181), int32(1183), int32(1185), int32(1187), int32(1189), int32(1191), int32(1193), int32(1195), int32(1197), int32(1199), int32(1201), int32(1203), int32(1205), int32(1207), int32(1209), int32(1211), int32(1213), int32(1215), int32(1218), int32(1220), int32(1222), int32(1224), int32(1226), int32(1228), int32(1230), int32(1233), int32(1235), int32(1237), int32(1239), int32(1241), int32(1243), int32(1245), int32(1247), int32(1249), int32(1251), int32(1253), int32(1255), int32(1257), int32(1259), int32(1261), int32(1263), int32(1265), int32(1267), int32(1269), int32(1271), int32(1273), int32(1275), int32(1277), int32(1279), int32(1281), int32(1283), int32(1285), int32(1287), int32(1289), int32(1291), int32(1293), int32(1295), int32(1297), int32(1299), int32(7681), int32(7683), int32(7685), int32(7687), int32(7689), int32(7691), int32(7693), int32(7695), int32(7697), int32(7699), int32(7701), int32(7703), int32(7705), int32(7707), int32(7709), int32(7711), int32(7713), int32(7715), int32(7717), int32(7719), int32(7721), int32(7723), int32(7725), int32(7727), int32(7729), int32(7731), int32(7733), int32(7735), int32(7737), int32(7739), int32(7741), int32(7743), int32(7745), int32(7747), int32(7749), int32(7751), int32(7753), int32(7755), int32(7757), int32(7759), int32(7761), int32(7763), int32(7765), int32(7767), int32(7769), int32(7771), int32(7773), int32(7775), int32(7777), int32(7779), int32(7781), int32(7783), int32(7785), int32(7787), int32(7789), int32(7791), int32(7793), int32(7795), int32(7797), int32(7799), int32(7801), int32(7803), int32(7805), int32(7807), int32(7809), int32(7811), int32(7813), int32(7815), int32(7817), int32(7819), int32(7821), int32(7823), int32(7825), int32(7827), int32(7829), int32(7841), int32(7843), int32(7845), int32(7847), int32(7849), int32(7851), int32(7853), int32(7855), int32(7857), int32(7859), int32(7861), int32(7863), int32(7865), int32(7867), int32(7869), int32(7871), int32(7873), int32(7875), int32(7877), int32(7879), int32(7881), int32(7883), int32(7885), int32(7887), int32(7889), int32(7891), int32(7893), int32(7895), int32(7897), int32(7899), int32(7901), int32(7903), int32(7905), int32(7907), int32(7909), int32(7911), int32(7913), int32(7915), int32(7917), int32(7919), int32(7921), int32(7923), int32(7925), int32(7927), int32(7929), int32(8580), int32(11361), int32(11368), int32(11370), int32(11372), int32(11382), int32(11393), int32(11395), int32(11397), int32(11399), int32(11401), int32(11403), int32(11405), int32(11407), int32(11409), int32(11411), int32(11413), int32(11415), int32(11417), int32(11419), int32(11421), int32(11423), int32(11425), int32(11427), int32(11429), int32(11431), int32(11433), int32(11435), int32(11437), int32(11439), int32(11441), int32(11443), int32(11445), int32(11447), int32(11449), int32(11451), int32(11453), int32(11455), int32(11457), int32(11459), int32(11461), int32(11463), int32(11465), int32(11467), int32(11469), int32(11471), int32(11473), int32(11475), int32(11477), int32(11479), int32(11481), int32(11483), int32(11485), int32(11487), int32(11489), int32(11491)})), NewCaseCanonicalize_DeltaSet(int32(2), CharRanges_WithMembers([]int32{int32(454), int32(457), int32(460), int32(499)})), NewCaseCanonicalize_DeltaSet(int32(15), CharRanges_WithMembers([]int32{int32(1231)})), NewCaseCanonicalize_DeltaSet(int32(16), CharRanges_WithRanges([]int32{int32(8560), int32(8576)})), NewCaseCanonicalize_DeltaSet(int32(26), CharRanges_WithRanges([]int32{int32(9424), int32(9450)})), NewCaseCanonicalize_DeltaSet(int32(28), CharRanges_WithMembers([]int32{int32(8526)})), NewCaseCanonicalize_DeltaSet(int32(31), CharRanges_WithMembers([]int32{int32(962)})), NewCaseCanonicalize_DeltaSet(int32(32), CharRanges_WithRanges([]int32{int32(97), int32(123), int32(224), int32(247), int32(248), int32(255), int32(945), int32(962), int32(963), int32(972), int32(1072), int32(1104), int32(65345), int32(65371)})), NewCaseCanonicalize_DeltaSet(int32(37), CharRanges_WithRanges([]int32{int32(941), int32(944)})), NewCaseCanonicalize_DeltaSet(int32(38), CharRanges_WithMembers([]int32{int32(940)})), NewCaseCanonicalize_DeltaSet(int32(47), CharRanges_WithMembers([]int32{int32(981)})), NewCaseCanonicalize_DeltaSet(int32(48), CharRanges_WithRanges([]int32{int32(1377), int32(1415), int32(11312), int32(11359)})), NewCaseCanonicalize_DeltaSet(int32(54), CharRanges_WithMembers([]int32{int32(982)})), NewCaseCanonicalize_DeltaSet(int32(57), CharRanges_WithMembers([]int32{int32(977)})), NewCaseCanonicalize_DeltaSet(int32(59), CharRanges_WithMembers([]int32{int32(7835)})), NewCaseCanonicalize_DeltaSet(int32(62), CharRanges_WithMembers([]int32{int32(976)})), NewCaseCanonicalize_DeltaSet(int32(63), CharRanges_WithRanges([]int32{int32(973), int32(975)})), NewCaseCanonicalize_DeltaSet(int32(64), CharRanges_WithMembers([]int32{int32(972)})), NewCaseCanonicalize_DeltaSet(int32(69), CharRanges_WithMembers([]int32{int32(649)})), NewCaseCanonicalize_DeltaSet(int32(71), CharRanges_WithMembers([]int32{int32(652)})), NewCaseCanonicalize_DeltaSet(int32(79), CharRanges_WithMembers([]int32{int32(477)})), NewCaseCanonicalize_DeltaSet(int32(80), CharRanges_WithRanges([]int32{int32(1009), int32(1010), int32(1104), int32(1120)})), NewCaseCanonicalize_DeltaSet(int32(86), CharRanges_WithMembers([]int32{int32(1008)})), NewCaseCanonicalize_DeltaSet(int32(96), CharRanges_WithMembers([]int32{int32(1013)})), NewCaseCanonicalize_DeltaSet(int32(202), CharRanges_WithMembers([]int32{int32(601)})), NewCaseCanonicalize_DeltaSet(int32(203), CharRanges_WithMembers([]int32{int32(603)})), NewCaseCanonicalize_DeltaSet(int32(205), CharRanges_WithRanges([]int32{int32(598), int32(600), int32(608), int32(609)})), NewCaseCanonicalize_DeltaSet(int32(206), CharRanges_WithMembers([]int32{int32(596)})), NewCaseCanonicalize_DeltaSet(int32(207), CharRanges_WithMembers([]int32{int32(611)})), NewCaseCanonicalize_DeltaSet(int32(209), CharRanges_WithMembers([]int32{int32(616)})), NewCaseCanonicalize_DeltaSet(int32(210), CharRanges_WithMembers([]int32{int32(595)})), NewCaseCanonicalize_DeltaSet(int32(211), CharRanges_WithMembers([]int32{int32(617), int32(623)})), NewCaseCanonicalize_DeltaSet(int32(213), CharRanges_WithMembers([]int32{int32(626)})), NewCaseCanonicalize_DeltaSet(int32(214), CharRanges_WithMembers([]int32{int32(629)})), NewCaseCanonicalize_DeltaSet(int32(217), CharRanges_WithRanges([]int32{int32(650), int32(652)})), NewCaseCanonicalize_DeltaSet(int32(218), CharRanges_WithMembers([]int32{int32(640), int32(643), int32(648)})), NewCaseCanonicalize_DeltaSet(int32(219), CharRanges_WithMembers([]int32{int32(658)})), NewCaseCanonicalize_DeltaSet(int32(7205), CharRanges_WithMembers([]int32{int32(8126)})), NewCaseCanonicalize_DeltaSet(int32(7264), CharRanges_WithRanges([]int32{int32(11520), int32(11558)})), NewCaseCanonicalize_DeltaSet(int32(10792), CharRanges_WithMembers([]int32{int32(11366)})), NewCaseCanonicalize_DeltaSet(int32(10795), CharRanges_WithMembers([]int32{int32(11365)}))})
}

func (this_ *CaseCanonicalize) initFields_CaseCanonicalize() {
	_ = this_
}

// Java constructor: CaseCanonicalize()
func NewCaseCanonicalize() *CaseCanonicalize {
	clinit_CaseCanonicalize()
	this_ := &CaseCanonicalize{}
	InitCaseCanonicalize(this_)
	return this_
}

func InitCaseCanonicalize(this_ *CaseCanonicalize) {
	this_.initFields_CaseCanonicalize()
}

// Java static method: CaseCanonicalize.caseCanonicalize(java.lang.String s)
func CaseCanonicalize_CaseCanonicalize_1_string(s string) string {
	clinit_CaseCanonicalize()
	_ = s
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = j2gort.JStrLen(s)
		_ = n
		for ; i < n; i = i + 1 {
			var ch uint16 = j2gort.JStrCharAt(s, i)
			_ = ch
			var cu uint16 = CaseCanonicalize_CaseCanonicalize_1_uint16(ch)
			_ = cu
			if int32(cu) != int32(ch) {
				var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder_1_string(s)
				_ = sb
				sb.SetCharAt(i, cu)
				for func() int32 { old_1 := i; nw_2 := old_1 + 1; i = nw_2; return nw_2 }() < n {
					sb.SetCharAt(i, CaseCanonicalize_CaseCanonicalize_1_uint16(j2gort.JStrCharAt(s, i)))
				}
				return sb.String()
			}
		}
	}
	return s
}

// Java static method: CaseCanonicalize.caseCanonicalize(char ch)
func CaseCanonicalize_CaseCanonicalize_1_uint16(ch uint16) uint16 {
	clinit_CaseCanonicalize()
	_ = ch
	if int32(ch) < int32(128) {
		if (int32(uint16(65)) <= int32(ch)) && (int32(ch) <= int32(uint16(90))) {
			return uint16((int32(ch) | int32(32)))
		} else {
			return ch
		}
	}
	if StaticGet_CaseCanonicalize_CASE_SENSITIVE().Contains(int32(ch)) {
		{
			cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(StaticGet_CaseCanonicalize_CANON_DELTA_SETS()))
			for cur_3.HasNext() {
				rawElem_4 := cur_3.Next()
				ds := j2gort.CastPtr[CaseCanonicalize_DeltaSet](rawElem_4)
				_ = ds
				if ds.codeUnits.Contains(int32(ch)) {
					return uint16((int32(ch) - ds.delta))
				}
			}
		}
	}
	return ch
}

// Java static method: CaseCanonicalize.expandToAllMatched(com.google.javascript.jscomp.regex.CharRanges ranges)
func CaseCanonicalize_ExpandToAllMatched(ranges *CharRanges) *CharRanges {
	clinit_CaseCanonicalize()
	_ = ranges
	var caseSensitive *CharRanges = ranges.Intersection(StaticGet_CaseCanonicalize_CASE_SENSITIVE())
	_ = caseSensitive
	if caseSensitive.IsEmpty() {
		return ranges
	}
	var expanded *CharRanges = StaticGet_CharRanges_EMPTY()
	_ = expanded
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(StaticGet_CaseCanonicalize_DELTA_SETS()))
		for cur_5.HasNext() {
			rawElem_6 := cur_5.Next()
			ds := j2gort.CastPtr[CaseCanonicalize_DeltaSet](rawElem_6)
			_ = ds
			expanded = expanded.Union(caseSensitive.Intersection(ds.codeUnits).Shift(-(ds.delta)))
		}
	}
	return ranges.Union(expanded)
}

// Java static method: CaseCanonicalize.reduceToMinimum(com.google.javascript.jscomp.regex.CharRanges ranges)
func CaseCanonicalize_ReduceToMinimum(ranges *CharRanges) *CharRanges {
	clinit_CaseCanonicalize()
	_ = ranges
	var caseSensitive *CharRanges = ranges.Intersection(StaticGet_CaseCanonicalize_CASE_SENSITIVE())
	_ = caseSensitive
	if caseSensitive.IsEmpty() {
		return ranges
	}
	var expanded *CharRanges = StaticGet_CharRanges_EMPTY()
	_ = expanded
	{
		cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(StaticGet_CaseCanonicalize_CANON_DELTA_SETS()))
		for cur_7.HasNext() {
			rawElem_8 := cur_7.Next()
			ds := j2gort.CastPtr[CaseCanonicalize_DeltaSet](rawElem_8)
			_ = ds
			expanded = expanded.Union(caseSensitive.Intersection(ds.codeUnits).Shift(-(ds.delta)))
		}
	}
	expanded = expanded.Difference(StaticGet_CaseCanonicalize_UCASE_ASCII_LETTERS()).Union(expanded.Intersection(StaticGet_CaseCanonicalize_UCASE_ASCII_LETTERS()).Shift(int32(32)))
	return ranges.Difference(caseSensitive).Union(expanded)
}

// Java class: com.google.javascript.jscomp.regex.CaseCanonicalize.DeltaSet
type CaseCanonicalize_DeltaSet struct {
	codeUnits *CharRanges `gson:"codeUnits"`
	delta     int32       `gson:"delta"`
	idHash_   atomic.Int32
}

func (this_ *CaseCanonicalize_DeltaSet) BindSelf_(s any) {
	_ = s
}
func (this_ *CaseCanonicalize_DeltaSet) IsJavaNil_() bool { return this_ == nil }
func (*CaseCanonicalize_DeltaSet) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.regex.CaseCanonicalize$DeltaSet", "DeltaSet"
}
func (this_ *CaseCanonicalize_DeltaSet) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CaseCanonicalize_DeltaSet) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.regex.CaseCanonicalize.DeltaSet", "com.google.javascript.jscomp.regex.CaseCanonicalize$DeltaSet", func(v_ any) bool {
		return j2gort.InstanceOf[*CaseCanonicalize_DeltaSet](v_)
	})
}

func (this_ *CaseCanonicalize_DeltaSet) initFields_CaseCanonicalize_DeltaSet() {
	_ = this_
}

// Java constructor: DeltaSet(int delta, com.google.javascript.jscomp.regex.CharRanges codeUnits)
func NewCaseCanonicalize_DeltaSet(delta int32, codeUnits *CharRanges) *CaseCanonicalize_DeltaSet {
	this_ := &CaseCanonicalize_DeltaSet{}
	InitCaseCanonicalize_DeltaSet(this_, delta, codeUnits)
	return this_
}

func InitCaseCanonicalize_DeltaSet(this_ *CaseCanonicalize_DeltaSet, delta int32, codeUnits *CharRanges) {
	_ = delta
	_ = codeUnits
	this_.initFields_CaseCanonicalize_DeltaSet()
	this_.delta = delta
	this_.codeUnits = codeUnits
}
