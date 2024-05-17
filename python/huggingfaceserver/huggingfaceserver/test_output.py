# Copyright 2024 The KServe Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

from pytest import approx


bert_token_classification_retrun_prob_expected_output = {
    "predictions": [
        [
            [
                {
                    0: approx(10.125614166259766),
                    1: approx(-1.818465232849121),
                    2: approx(-1.3191171884536743),
                    3: approx(-1.9324347972869873),
                    4: approx(-1.4850242137908936),
                    5: approx(-1.69266676902771),
                    6: approx(-0.8981077671051025),
                    7: approx(-1.7251288890838623),
                    8: approx(0.3205742835998535),
                },
                {
                    0: approx(9.647308349609375),
                    1: approx(-2.1127445697784424),
                    2: approx(-0.8831670880317688),
                    3: approx(-2.719135284423828),
                    4: approx(-0.47473347187042236),
                    5: approx(-2.242401361465454),
                    6: approx(0.6101500988006592),
                    7: approx(-2.219207763671875),
                    8: approx(-0.6545180082321167),
                },
                {
                    0: approx(10.364564895629883),
                    1: approx(-2.240159034729004),
                    2: approx(-0.9236820936203003),
                    3: approx(-2.6233034133911133),
                    4: approx(-0.501063883304596),
                    5: approx(-1.9418429136276245),
                    6: approx(-0.04101325571537018),
                    7: approx(-2.1208934783935547),
                    8: approx(-1.2565146684646606),
                },
                {
                    0: approx(9.219614028930664),
                    1: approx(-2.1359004974365234),
                    2: approx(-0.16899853944778442),
                    3: approx(-3.0277936458587646),
                    4: approx(0.25898468494415283),
                    5: approx(-2.4426767826080322),
                    6: approx(0.4815789461135864),
                    7: approx(-2.322394371032715),
                    8: approx(-0.23837831616401672),
                },
                {
                    0: approx(11.02370834350586),
                    1: approx(-2.0757784843444824),
                    2: approx(-0.864897608757019),
                    3: approx(-2.608891010284424),
                    4: approx(-0.7810694575309753),
                    5: approx(-1.7040153741836548),
                    6: approx(-0.5265536308288574),
                    7: approx(-1.7141872644424438),
                    8: approx(-1.0127880573272705),
                },
                {
                    0: approx(9.288057327270508),
                    1: approx(-1.9795231819152832),
                    2: approx(0.30957841873168945),
                    3: approx(-2.7814087867736816),
                    4: approx(-0.4492063820362091),
                    5: approx(-1.7794853448867798),
                    6: approx(-0.31150245666503906),
                    7: approx(-2.0755019187927246),
                    8: approx(-0.8829554319381714),
                },
                {
                    0: approx(7.968825817108154),
                    1: approx(-2.0212578773498535),
                    2: approx(-0.43740740418434143),
                    3: approx(-3.6762940883636475),
                    4: approx(-0.012760266661643982),
                    5: approx(-2.3209781646728516),
                    6: approx(0.6871910095214844),
                    7: approx(-2.541992664337158),
                    8: approx(0.3260609209537506),
                },
            ]
        ],
        [
            [
                {
                    0: approx(10.125614166259766),
                    1: approx(-1.818464994430542),
                    2: approx(-1.3191170692443848),
                    3: approx(-1.9324350357055664),
                    4: approx(-1.4850239753723145),
                    5: approx(-1.6926666498184204),
                    6: approx(-0.8981078863143921),
                    7: approx(-1.7251291275024414),
                    8: approx(0.32057392597198486),
                },
                {
                    0: approx(9.647308349609375),
                    1: approx(-2.1127443313598633),
                    2: approx(-0.8831661939620972),
                    3: approx(-2.71913480758667),
                    4: approx(-0.4747332036495209),
                    5: approx(-2.2424018383026123),
                    6: approx(0.6101502180099487),
                    7: approx(-2.219207286834717),
                    8: approx(-0.6545186638832092),
                },
                {
                    0: approx(10.36456298828125),
                    1: approx(-2.2401585578918457),
                    2: approx(-0.923682451248169),
                    3: approx(-2.6233034133911133),
                    4: approx(-0.5010634660720825),
                    5: approx(-1.9418421983718872),
                    6: approx(-0.041012972593307495),
                    7: approx(-2.1208930015563965),
                    8: approx(-1.2565146684646606),
                },
                {
                    0: approx(9.219614028930664),
                    1: approx(-2.1359007358551025),
                    2: approx(-0.16899824142456055),
                    3: approx(-3.0277938842773438),
                    4: approx(0.2589845061302185),
                    5: approx(-2.442676067352295),
                    6: approx(0.4815790057182312),
                    7: approx(-2.3223938941955566),
                    8: approx(-0.23837843537330627),
                },
                {
                    0: approx(11.02370834350586),
                    1: approx(-2.075778007507324),
                    2: approx(-0.8648972511291504),
                    3: approx(-2.6088919639587402),
                    4: approx(-0.7810693979263306),
                    5: approx(-1.7040156126022339),
                    6: approx(-0.5265532732009888),
                    7: approx(-1.7141878604888916),
                    8: approx(-1.0127880573272705),
                },
                {
                    0: approx(9.288058280944824),
                    1: approx(-1.9795236587524414),
                    2: approx(0.3095783591270447),
                    3: approx(-2.7814087867736816),
                    4: approx(-0.4492063522338867),
                    5: approx(-1.7794855833053589),
                    6: approx(-0.3115025758743286),
                    7: approx(-2.0755016803741455),
                    8: approx(-0.8829552531242371),
                },
                {
                    0: approx(7.968822479248047),
                    1: approx(-2.0212578773498535),
                    2: approx(-0.4374064803123474),
                    3: approx(-3.676295280456543),
                    4: approx(-0.012759596109390259),
                    5: approx(-2.3209784030914307),
                    6: approx(0.6871915459632874),
                    7: approx(-2.5419921875),
                    8: approx(0.3260613977909088),
                },
            ]
        ],
    ]
}
